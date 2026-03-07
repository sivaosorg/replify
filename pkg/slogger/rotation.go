package slogger

import (
"fmt"
"os"
"path/filepath"
"strings"
"time"

"github.com/sivaosorg/replify/pkg/sysx"
)

// Levels implements Hook by returning the set of log levels this hook handles.
//
// Returns:
//
// the slice of Level values for which this hook's Fire method will be called.
func (h *LevelWriterHook) Levels() []Level { return h.levels }

// Fire implements Hook by serialising the entry with the configured Formatter
// and writing the result to the appropriate level-specific file.
//
// Parameters:
//   - `e`: the log entry to write
//
// Returns:
//
// an error if formatting or writing fails; nil on success.
func (h *LevelWriterHook) Fire(e *Entry) error {
data, err := h.formatter.Format(e)
if err != nil {
return err
}
lvl := e.GetLevel()
if _, err = h.writer.WriteLevel(lvl, data); err != nil {
return fmt.Errorf("slogger: LevelWriterHook write failed for level %s: %w", lvl, err)
}
return nil
}

// NewLevelFileWriter creates a LevelFileWriter with the given options.
// The logs directory and per-level files are created automatically if they
// do not already exist.
//
// Parameters:
//   - `opts`: rotation configuration including directory, size limits, and compression
//
// Returns:
//
// a ready-to-use *LevelFileWriter and any initialisation error.
func newLevelFileWriter(opts RotationOptions) (*LevelFileWriter, error) {
if opts.Dir == "" {
opts.Dir = defaultLogDir
}
if opts.MaxBytes <= 0 {
opts.MaxBytes = defaultMaxBytes
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
// If no exact file exists for level, the nearest coarser-grained file is used
// (Trace → Debug, Fatal/Panic → Error).
//
// Parameters:
//   - `level`: the severity level determining which file receives p
//   - `p`: the bytes to write
//
// Returns:
//
// the number of bytes written and any write error.
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

// Write implements io.Writer by routing p to the InfoLevel file.
// This satisfies io.Writer so that a LevelFileWriter can be used as a generic
// writer destination.
//
// Parameters:
//   - `p`: the bytes to write
//
// Returns:
//
// the number of bytes written and any write error.
func (lfw *LevelFileWriter) Write(p []byte) (int, error) {
return lfw.WriteLevel(InfoLevel, p)
}

// Close closes all open file handles managed by this writer.
// After Close, the LevelFileWriter must not be used.
//
// Returns:
//
// the first error encountered while closing files, or nil if all succeed.
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

// Rotate forces immediate rotation of all level files, regardless of their
// current size or age. Useful for external rotation signals (e.g. SIGHUP).
//
// Returns:
//
// the first rotation error encountered, or nil if all files rotated successfully.
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
dateDir := filepath.Join(rf.dir, defaultArchivedDir, now.Format(archiveDateFormat))
if !sysx.DirExists(dateDir) {
if err := os.MkdirAll(dateDir, 0755); err != nil {
return fmt.Errorf("slogger: cannot create archive dir %q: %w", dateDir, err)
}
}

stamp := now.Format(archiveStampFormat)
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
