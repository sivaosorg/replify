package sysx

import (
	"bufio"
	"fmt"
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

// Write appends data to the file, creating it if it does not exist.
// The operation is serialised by an internal mutex, making it safe to call
// concurrently from multiple goroutines.
//
// Parameters:
//   - `data`: the bytes to append.
//
// Returns:
//
//	An error if the file could not be opened or written; nil on success.
func (w *SafeFileWriter) Write(data []byte) error {
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
// The operation is serialised by an internal mutex.
//
// Parameters:
//   - `s`: the string to append.
//
// Returns:
//
//	An error if the file could not be opened or written; nil on success.
func (w *SafeFileWriter) WriteString(s string) error {
	return w.Write([]byte(s))
}

// Overwrite replaces the entire file content with data, atomically using the
// temporary-file-and-rename pattern. The operation is serialised by an
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
	return AtomicWriteFile(w.path, data)
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

// WriteFile writes data to the file at path, creating the file if it does not
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
//	if err := sysx.WriteFile("/tmp/output.bin", payload); err != nil {
//	    log.Fatal(err)
//	}
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

// WriteFileString writes content to the file at path, creating it if it does
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
//	if err := sysx.WriteFileString("/tmp/hello.txt", "hello world\n"); err != nil {
//	    log.Fatal(err)
//	}
func WriteFileString(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// WriteFileLocked writes data to path using a per-path in-process mutex,
// ensuring that concurrent calls with the same path value are serialised. The
// file is created with permission 0644 if it does not exist, or truncated if
// it does.
//
// Note: this provides in-process synchronisation only. For cross-process
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
//	go sysx.WriteFileLocked("/tmp/shared.json", data1)
//	go sysx.WriteFileLocked("/tmp/shared.json", data2)
func WriteFileLocked(path string, data []byte) error {
	mu := getFileMutex(path)
	mu.Lock()
	defer mu.Unlock()
	return os.WriteFile(path, data, 0o644)
}

// AppendFile appends data to the file at path, creating it with permission
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
//	if err := sysx.AppendFile("/var/log/app.log", []byte("entry\n")); err != nil {
//	    log.Fatal(err)
//	}
func AppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
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
	return AppendFile(path, []byte(content))
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

// AtomicWriteFile writes data to path atomically using the
// temporary-file-and-rename pattern: data is first flushed to a temporary
// file located in the same directory as path, then the temporary file is
// renamed to path.
//
// On POSIX systems, os.Rename is atomic when the source and destination share
// the same filesystem, so readers will never observe a partial write.
//
// The temporary file is created with permission 0644. If the rename fails,
// the temporary file is cleaned up automatically.
//
// Parameters:
//   - `path`: the destination file path.
//   - `data`: the bytes to write.
//
// Returns:
//
//	An error if the temporary file could not be created, written, or renamed;
//	nil on success.
//
// Example:
//
//	if err := sysx.AtomicWriteFile("/etc/app/config.json", jsonData); err != nil {
//	    log.Fatal(err)
//	}
func AtomicWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
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
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	renamed = true
	return nil
}
