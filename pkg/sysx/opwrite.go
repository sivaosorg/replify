package sysx

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// NewSafeFileWriter creates a SafeFileWriter targeting path with the default
// file permission of 0644.
//
// Parameters:
//   - `path`: the file path to write to.
//
// Returns:
//
//	A pointer to a new SafeFileWriter.
//
// Example:
//
//	w := sysx.NewSafeFileWriter("/var/log/app.log")
//	go w.WriteString("line from goroutine 1\n")
//	go w.WriteString("line from goroutine 2\n")
func NewSafeFileWriter(path string) *SafeFileWriter {
	return &SafeFileWriter{path: path, perm: 0o644}
}

// WithPerm overrides the file permission used when creating the file.
// The default is 0644.
//
// Parameters:
//   - `perm`: the os.FileMode to apply on file creation.
//
// Returns:
//
//	The receiver, enabling method chaining.
func (w *SafeFileWriter) WithPerm(perm os.FileMode) *SafeFileWriter {
	w.perm = perm
	return w
}

// WriteBytes appends data to the file, creating it if it does not exist.
// The operation is serialized by an internal mutex, making it safe to call
// concurrently from multiple goroutines.
//
// Parameters:
//   - `data`: the bytes to append.
//
// Returns:
//
//	An error if the file could not be opened or written; nil on success.
func (w *SafeFileWriter) WriteBytes(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	f, err := os.OpenFile(w.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, w.perm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// WriteString appends s to the file, creating it if it does not exist.
// The operation is serialized by an internal mutex.
//
// Parameters:
//   - `s`: the string to append.
//
// Returns:
//
//	An error if the file could not be opened or written; nil on success.
func (w *SafeFileWriter) WriteString(s string) error {
	return w.WriteBytes([]byte(s))
}

// Overwrite replaces the entire file content with data, atomically using the
// temporary-file-and-rename pattern. The operation is serialized by an
// internal mutex.
//
// Parameters:
//   - `data`: the bytes to write.
//
// Returns:
//
//	An error if the write or rename failed; nil on success.
func (w *SafeFileWriter) Overwrite(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return AtomicWriteBytes(w.path, data)
}

// Path returns the file path targeted by this SafeFileWriter.
//
// Returns:
//
//	A string containing the file path.
func (w *SafeFileWriter) Path() string {
	return w.path
}

// Perm returns the file permission used when creating the file.
//
// Returns:
//
//	An os.FileMode representing the file permission.
func (w *SafeFileWriter) Perm() os.FileMode {
	return w.perm
}

// WriteBytes writes data to the file at path, creating the file if it does not
// exist or truncating it if it does. The file is created with permission 0644.
//
// Parameters:
//   - `path`: the destination file path.
//   - `data`: the bytes to write.
//
// Returns:
//
//	An error if the file could not be created or written; nil on success.
//
// Example:
//
//	if err := sysx.WriteBytes("/tmp/output.bin", payload); err != nil {
//	    log.Fatal(err)
//	}
func WriteBytes(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

// WriteBytesWithLocked writes data to path using a per-path in-process mutex,
// ensuring that concurrent calls with the same path value are serialized. The
// file is created with permission 0644 if it does not exist, or truncated if
// it does.
//
// Note: this provides in-process synchronization only. For cross-process
// safety, use AtomicWriteFile or a platform-specific file locking mechanism.
//
// Parameters:
//   - `path`: the destination file path.
//   - `data`: the bytes to write.
//
// Returns:
//
//	An error if the write failed; nil on success.
//
// Example:
//
//	// Safe to call concurrently from multiple goroutines targeting the same path.
//	go sysx.WriteBytesWithLocked("/tmp/shared.json", data1)
//	go sysx.WriteBytesWithLocked("/tmp/shared.json", data2)
func WriteBytesWithLocked(path string, data []byte) error {
	mu := getFileMutex(path)
	mu.Lock()
	defer mu.Unlock()
	return os.WriteFile(path, data, 0o644)
}

// WriteString writes content to the file at path, creating it if it does
// not exist or truncating it if it does. The file is created with permission 0644.
//
// Parameters:
//   - `path`:    the destination file path.
//   - `content`: the string to write.
//
// Returns:
//
//	An error if the file could not be created or written; nil on success.
//
// Example:
//
//	if err := sysx.WriteString("/tmp/hello.txt", "hello world\n"); err != nil {
//	    log.Fatal(err)
//	}
func WriteString(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// WriteLines writes each element of lines to the file at path as a separate
// line terminated by "\n", creating or truncating the file. A buffered writer
// is used for efficiency.
//
// Parameters:
//   - `path`:  the destination file path.
//   - `lines`: the slice of strings to write.
//
// Returns:
//
//	An error if the file could not be created or written; nil on success.
//
// Example:
//
//	if err := sysx.WriteLines("/tmp/list.txt", []string{"alpha", "beta", "gamma"}); err != nil {
//	    log.Fatal(err)
//	}
func WriteLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return w.Flush()
}

// AtomicWriteBytes writes data to path atomically using the
// temporary-file-and-rename pattern: data is first flushed to a temporary
// file in the same directory as path, then renamed to path.
//
// Atomicity guarantees by platform:
//   - Linux / macOS (native FS): rename(2) atomically replaces the
//     destination — readers never see a partial write.
//   - Windows: MoveFileEx(MOVEFILE_REPLACE_EXISTING) provides the same
//     guarantee on the same volume and works even when the destination
//     already exists.
//   - exFAT / FAT32 / SMB mounts: rename returns EEXIST when the destination
//     already exists; a remove-then-rename fallback is used. This sacrifices
//     strict atomicity on those filesystems, but prevents the EEXIST failure.
//
// Security note (CWE-732): os.CreateTemp creates the staging file with mode
// 0600 (owner-only). After the rename the final file is explicitly chmod'd to
// 0644. Callers writing sensitive data should call os.Chmod afterward.
//
// Parameters:
//   - `path`: the destination file path; parent directories are created automatically.
//   - `data`: the bytes to write.
//
// Returns:
//
//	An error if any step fails; nil on success.
//
// Example:
//
//	if err := sysx.AtomicWriteBytes("/etc/app/config.json", jsonData); err != nil {
//	    log.Fatal(err)
//	}
func AtomicWriteBytes(path string, data []byte) error {
	dir := filepath.Dir(path)
	// Ensure the parent directory exists (idempotent, works on all OS).
	if err := CreateDir(dir); err != nil {
		return fmt.Errorf("sysx: AtomicWriteBytes: mkdir %q: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".tmp_atomic_*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	renamed := false
	defer func() {
		if !renamed {
			os.Remove(tmpPath)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// renameReplace is OS-specific:
	//   - POSIX (Linux/macOS): delegates to os.Rename, which is atomic.
	//   - Windows: uses MoveFileEx(MOVEFILE_REPLACE_EXISTING), which is
	//     atomic on the same volume and handles existing destinations.
	if err := renameReplace(tmpPath, path); err != nil {
		// Last-resort fallback for non-POSIX filesystems (exFAT, FAT32,
		// some SMB mounts) where even renameReplace may fail with EEXIST.
		if os.IsExist(err) {
			if delErr := os.Remove(path); delErr != nil && !os.IsNotExist(delErr) {
				return fmt.Errorf("sysx: AtomicWriteBytes: remove existing %q: %w", path, delErr)
			}
			if err = os.Rename(tmpPath, path); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	renamed = true
	// Security fix (CWE-732): os.CreateTemp creates with 0600; set 0644 so
	// group/other read access is predictable after the rename.
	if err := os.Chmod(path, 0o644); err != nil {
		return err
	}
	return nil
}

// AppendBytes appends data to the file at path, creating it with permission
// 0644 if it does not exist.
//
// Parameters:
//   - `path`: the destination file path.
//   - `data`: the bytes to append.
//
// Returns:
//
//	An error if the file could not be opened or written; nil on success.
//
// Example:
//
//	if err := sysx.AppendBytes("/var/log/app.log", []byte("entry\n")); err != nil {
//	    log.Fatal(err)
//	}
func AppendBytes(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendOrWriteBytes writes data to path using the following strategy:
//
//   - If the file does not exist or is empty, data is written atomically via
//     the temp-file-and-rename pattern (same as [AtomicWriteBytes]).
//   - If the file already exists and contains data, sep followed by data is
//     appended, so the caller can choose a suitable record separator (e.g. "\n"
//     for JSON-Lines / NDJSON, "\n---\n" for human-readable multi-entry logs).
//
// Parent directories are created automatically (idempotent).
//
// Parameters:
//   - `path`: the destination file path.
//   - `data`: the bytes to write or append.
//   - `sep`:  separator prepended to data only when appending to a non-empty
//     file; pass nil or an empty slice to append without a separator.
//
// Returns:
//
//	An error if any step fails; nil on success.
//
// Example:
//
//	// Append JSON objects separated by a newline (NDJSON / JSON-Lines):
//	if err := sysx.AppendOrWriteBytes("/var/log/app/responses.jsonl", entry, []byte("\n")); err != nil {
//	    log.Fatal(err)
//	}
func AppendOrWriteBytes(path string, data, sep []byte) error {
	dir := filepath.Dir(path)
	if err := CreateDir(dir); err != nil {
		return fmt.Errorf("sysx: AppendOrWriteBytes: mkdir %q: %w", dir, err)
	}
	info, err := os.Stat(path)
	if err == nil && info.Size() > 0 {
		// File exists and is non-empty — open for append.
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("sysx: AppendOrWriteBytes: open %q: %w", path, err)
		}
		defer f.Close()
		if len(sep) > 0 {
			if _, err := f.Write(sep); err != nil {
				return fmt.Errorf("sysx: AppendOrWriteBytes: write sep to %q: %w", path, err)
			}
		}
		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("sysx: AppendOrWriteBytes: write data to %q: %w", path, err)
		}
		return nil
	}
	// File does not exist or is empty — write atomically.
	return AtomicWriteBytes(path, data)
}

// AppendString appends content to the file at path, creating it with
// permission 0644 if it does not exist.
//
// Parameters:
//   - `path`:    the destination file path.
//   - `content`: the string to append.
//
// Returns:
//
//	An error if the file could not be opened or written; nil on success.
//
// Example:
//
//	if err := sysx.AppendString("/var/log/app.log", "entry\n"); err != nil {
//	    log.Fatal(err)
//	}
func AppendString(path string, content string) error {
	return AppendBytes(path, []byte(content))
}

// AppendOrCopyFrom writes data from src to path using the same
// append-or-create strategy as [AppendOrWriteBytes], but accepts an
// [io.Reader] so arbitrarily large payloads are never fully buffered in
// memory.
//
//   - If the file does not exist or is empty, src is written atomically via
//     the temp-file-and-rename pattern so readers never observe a partial
//     write.
//   - If the file already exists and contains data, sep followed by the bytes
//     from src is appended. The OS-level O_APPEND flag ensures concurrent
//     appenders from separate processes do not interleave partial writes.
//
// Parent directories are created automatically (idempotent).
//
// Parameters:
//   - `path`: the destination file path.
//   - `src`:  an [io.Reader] providing the payload; consumed exactly once.
//   - `sep`:  separator prepended to src only when appending to a non-empty
//     file; pass nil or an empty slice to append without a separator.
//
// Returns:
//
//	An error if any step fails; nil on success.
//
// Example:
//
//	// Stream-append JSON objects separated by a newline (NDJSON / JSON-Lines):
//	if err := sysx.AppendOrCopyFrom("/var/log/app/bodies.jsonl", src, []byte("\n")); err != nil {
//	    log.Fatal(err)
//	}
func AppendOrCopyFrom(path string, src io.Reader, sep []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("sysx: AppendOrCopyFrom: mkdir %q: %w", dir, err)
	}
	info, err := os.Stat(path)
	if err == nil && info.Size() > 0 {
		// File exists and is non-empty — open for append.
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("sysx: AppendOrCopyFrom: open %q: %w", path, err)
		}
		defer f.Close()
		if len(sep) > 0 {
			if _, err := f.Write(sep); err != nil {
				return fmt.Errorf("sysx: AppendOrCopyFrom: write sep to %q: %w", path, err)
			}
		}
		if _, err := io.Copy(f, src); err != nil {
			return fmt.Errorf("sysx: AppendOrCopyFrom: copy to %q: %w", path, err)
		}
		return nil
	}
	// File does not exist or is empty — write atomically via temp-rename.
	tmp, err := os.CreateTemp(dir, ".tmp_atomic_*")
	if err != nil {
		return fmt.Errorf("sysx: AppendOrCopyFrom: create temp in %q: %w", dir, err)
	}
	tmpPath := tmp.Name()
	_, copyErr := io.Copy(tmp, src)
	closeErr := tmp.Close()
	if copyErr != nil || closeErr != nil {
		os.Remove(tmpPath)
		if copyErr != nil {
			return fmt.Errorf("sysx: AppendOrCopyFrom: write temp: %w", copyErr)
		}
		return fmt.Errorf("sysx: AppendOrCopyFrom: close temp: %w", closeErr)
	}
	if err := renameReplace(tmpPath, path); err != nil {
		if os.IsExist(err) {
			os.Remove(path)
			if err2 := os.Rename(tmpPath, path); err2 != nil {
				os.Remove(tmpPath)
				return fmt.Errorf("sysx: AppendOrCopyFrom: rename %q→%q: %w", tmpPath, path, err2)
			}
		} else {
			os.Remove(tmpPath)
			return fmt.Errorf("sysx: AppendOrCopyFrom: rename %q→%q: %w", tmpPath, path, err)
		}
	}
	_ = os.Chmod(path, 0o644)
	return nil
}
