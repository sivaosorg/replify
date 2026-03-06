package slogger

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sivaosorg/replify/pkg/sysx"
)

// RotationOptions configures the rotating file writer.
type RotationOptions struct {
	// Dir is the base log directory. Defaults to "logs".
	Dir string
	// MaxBytes is the maximum file size before rotation. Defaults to 10 MB.
	MaxBytes int64
	// MaxAge is the maximum age of a log file before rotation. Zero means no age-based rotation.
	MaxAge time.Duration
	// Compress controls whether rotated files are zipped. Defaults to true.
	Compress bool
}

// LevelFileWriter writes log entries to separate files per log level.
// It supports automatic rotation and ZIP compression of archived logs.
type LevelFileWriter struct {
	mu      sync.Mutex
	opts    RotationOptions
	writers map[Level]*rotatingFile
}

type rotatingFile struct {
	mu       sync.Mutex
	path     string
	file     *os.File
	size     int64
	maxBytes int64
	openedAt time.Time
	maxAge   time.Duration
	compress bool
	dir      string
	level    Level
}

// LevelWriterHook is a Hook that writes log entries to level-specific files.
// Use it with AddHook to enable automatic per-level file logging.
type LevelWriterHook struct {
	writer    *LevelFileWriter
	formatter Formatter
	levels    []Level
}

// NewLevelWriterHook creates a LevelWriterHook that writes to lfw using formatter.
// If levels is empty, all levels are enabled.
func NewLevelWriterHook(lfw *LevelFileWriter, formatter Formatter, levels ...Level) *LevelWriterHook {
	if len(levels) == 0 {
		levels = []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	}
	return &LevelWriterHook{
		writer:    lfw,
		formatter: formatter,
		levels:    levels,
	}
}

// Levels implements Hook.
func (h *LevelWriterHook) Levels() []Level { return h.levels }

// Fire implements Hook.
func (h *LevelWriterHook) Fire(e *Entry) error {
	data, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	if _, err = h.writer.WriteLevel(e.level, data); err != nil {
		return fmt.Errorf("slogger: LevelWriterHook write failed for level %s: %w", e.level, err)
	}
	return nil
}

// NewLevelFileWriter creates a LevelFileWriter with the given options.
func NewLevelFileWriter(opts RotationOptions) (*LevelFileWriter, error) {
	if opts.Dir == "" {
		opts.Dir = "logs"
	}
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = 10 * 1024 * 1024 // 10 MB
	}

	if !sysx.DirExists(opts.Dir) {
		if err := os.MkdirAll(opts.Dir, 0755); err != nil {
			return nil, fmt.Errorf("slogger: cannot create log directory %q: %w", opts.Dir, err)
		}
	}

	w := &LevelFileWriter{
		opts:    opts,
		writers: make(map[Level]*rotatingFile),
	}

	// Four files are created: debug, info, warn, and error.
	// Trace routes to debug; Fatal and Panic route to error (see WriteLevel).
	for _, lvl := range []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel} {
		rf, err := newRotatingFile(opts, lvl)
		if err != nil {
			_ = w.Close()
			return nil, err
		}
		w.writers[lvl] = rf
	}
	return w, nil
}

// WriteLevel routes p to the file for the given level.
func (lfw *LevelFileWriter) WriteLevel(level Level, p []byte) (int, error) {
	lfw.mu.Lock()
	rf, ok := lfw.writers[level]
	if !ok {
		switch {
		case level <= DebugLevel:
			rf = lfw.writers[DebugLevel]
		case level <= InfoLevel:
			rf = lfw.writers[InfoLevel]
		case level <= WarnLevel:
			rf = lfw.writers[WarnLevel]
		default:
			rf = lfw.writers[ErrorLevel]
		}
	}
	lfw.mu.Unlock()
	if rf == nil {
		return 0, nil
	}
	return rf.write(p)
}

// Write implements io.Writer by routing to the InfoLevel file.
func (lfw *LevelFileWriter) Write(p []byte) (int, error) {
	return lfw.WriteLevel(InfoLevel, p)
}

// Close closes all open file handles.
func (lfw *LevelFileWriter) Close() error {
	lfw.mu.Lock()
	defer lfw.mu.Unlock()
	var first error
	for _, rf := range lfw.writers {
		if err := rf.close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// Rotate forces rotation of all level files.
func (lfw *LevelFileWriter) Rotate() error {
	lfw.mu.Lock()
	defer lfw.mu.Unlock()
	var first error
	for _, rf := range lfw.writers {
		if err := rf.rotate(); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func newRotatingFile(opts RotationOptions, level Level) (*rotatingFile, error) {
	rf := &rotatingFile{
		path:     filepath.Join(opts.Dir, levelFileName(level)),
		maxBytes: opts.MaxBytes,
		maxAge:   opts.MaxAge,
		compress: opts.Compress,
		dir:      opts.Dir,
		level:    level,
	}
	if err := rf.open(); err != nil {
		return nil, err
	}
	return rf, nil
}

func levelFileName(level Level) string {
	switch level {
	case DebugLevel:
		return "debug.log"
	case InfoLevel:
		return "info.log"
	case WarnLevel:
		return "warn.log"
	default:
		return "error.log"
	}
}

func (rf *rotatingFile) open() error {
	f, err := os.OpenFile(rf.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("slogger: cannot open log file %q: %w", rf.path, err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return fmt.Errorf("slogger: cannot stat log file %q: %w", rf.path, err)
	}
	rf.file = f
	rf.size = info.Size()
	rf.openedAt = time.Now()
	return nil
}

func (rf *rotatingFile) write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.needsRotation(int64(len(p))) {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.file.Write(p)
	rf.size += int64(n)
	return n, err
}

func (rf *rotatingFile) needsRotation(incoming int64) bool {
	if rf.maxBytes > 0 && rf.size+incoming > rf.maxBytes {
		return true
	}
	if rf.maxAge > 0 && time.Since(rf.openedAt) > rf.maxAge {
		return true
	}
	return false
}

func (rf *rotatingFile) rotate() error {
	if rf.file != nil {
		if err := rf.file.Close(); err != nil {
			return fmt.Errorf("slogger: cannot close log file for rotation: %w", err)
		}
		rf.file = nil
	}

	now := time.Now()
	dateDir := filepath.Join(rf.dir, "archived", now.Format("2006-01-02"))
	if !sysx.DirExists(dateDir) {
		if err := os.MkdirAll(dateDir, 0755); err != nil {
			return fmt.Errorf("slogger: cannot create archive dir %q: %w", dateDir, err)
		}
	}

	stamp := now.Format("20060102150405")
	levelName := strings.ToLower(rf.level.String())

	if rf.compress {
		zipPath := filepath.Join(dateDir, fmt.Sprintf("%s_%s.zip", stamp, levelName))
		if err := compressToZip(rf.path, zipPath); err != nil {
			return fmt.Errorf("slogger: archive compression failed: %w", err)
		}
		_ = os.Remove(rf.path)
	} else {
		archivePath := filepath.Join(dateDir, fmt.Sprintf("%s_%s.log", stamp, levelName))
		if err := os.Rename(rf.path, archivePath); err != nil {
			return fmt.Errorf("slogger: cannot archive log file: %w", err)
		}
	}

	return rf.open()
}

func compressToZip(srcPath, zipPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	w, err := zw.Create(filepath.Base(srcPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, src)
	return err
}

func (rf *rotatingFile) close() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.file != nil {
		err := rf.file.Close()
		rf.file = nil
		return err
	}
	return nil
}
