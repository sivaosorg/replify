package sysx

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

// ///////////////////////////
// Section: Existence checks
// ///////////////////////////

// FileExists reports whether a file exists at the given path.
//
// The function returns true for any path that exists in the file system,
// including directories and symbolic links. Use IsFile to restrict the check
// to regular files.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when any file system entry exists at path;
//	 - false when the path does not exist or the existence cannot be determined.
//
// Example:
//
//	if sysx.FileExists("/etc/hosts") {
//	    fmt.Println("hosts file found")
//	}
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists reports whether a directory exists at the given path.
//
// The function returns true only when the path exists and refers to a
// directory (not a regular file or other entry).
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when a directory exists at path;
//	 - false otherwise.
//
// Example:
//
//	if sysx.DirExists("/tmp") {
//	    fmt.Println("tmp dir exists")
//	}
func DirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// IsFile reports whether the given path exists and is a regular file.
//
// Symbolic links are followed; if the link target is a regular file this
// function returns true.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when path exists and is a regular file;
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsFile("/etc/passwd") {
//	    fmt.Println("regular file")
//	}
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

// IsDir reports whether the given path exists and is a directory.
//
// Symbolic links are followed; if the link target is a directory this
// function returns true.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when path exists and is a directory;
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsDir("/home") {
//	    fmt.Println("is a directory")
//	}
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// IsSymlink reports whether the given path is a symbolic link.
//
// Unlike IsFile and IsDir, this function does NOT follow symbolic links; it
// inspects the link entry itself using os.Lstat.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when path exists and is a symbolic link;
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsSymlink("/usr/bin/python") {
//	    fmt.Println("python is a symlink")
//	}
func IsSymlink(path string) bool {
	fi, err := os.Lstat(path)
	return err == nil && fi.Mode()&os.ModeSymlink != 0
}

// ///////////////////////////
// Section: Permission checks
// ///////////////////////////

// IsExecutable reports whether the file at the given path is executable by
// its owner.
//
// The check is based on file mode bits (0100). On Windows, mode bits are an
// approximation; all files typically report as executable.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when the file exists and its owner-execute bit is set;
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsExecutable("/usr/bin/git") {
//	    fmt.Println("git is executable")
//	}
func IsExecutable(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode()&0100 != 0
}

// IsReadable reports whether the file at the given path is readable by its
// owner.
//
// The check is based on file mode bits (0400). On Windows, mode bits are an
// approximation.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when the file exists and its owner-read bit is set;
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsReadable("/etc/hosts") {
//	    fmt.Println("readable")
//	}
func IsReadable(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode()&0400 != 0
}

// IsWritable reports whether the file at the given path can be opened for
// writing by the current process.
//
// The check attempts to open the file with os.O_WRONLY. This approach
// respects ACLs and is more accurate than mode-bit inspection alone.
//
// Parameters:
//   - `path`: the file system path to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when the file exists and can be opened for writing;
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsWritable("/tmp/out.log") {
//	    fmt.Println("writable")
//	}
func IsWritable(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

// ///////////////////////////
// Section: File metadata
// ///////////////////////////

// FileSize returns the size of the file at the given path in bytes.
//
// Parameters:
//   - `path`: the file system path to inspect.
//
// Returns:
//
//	(int64, error): the file size in bytes and nil on success, or 0 and a
//	non-nil error if the file does not exist or cannot be stat'd.
//
// Example:
//
//	size, err := sysx.FileSize("/var/log/syslog")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%d bytes\n", size)
func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// ///////////////////////////
// Section: Special directories
// ///////////////////////////

// TempDir returns the default directory to use for temporary files.
//
// It delegates directly to os.TempDir. On Unix, it returns $TMPDIR if set,
// else /tmp. On Windows it returns %TMP%, %TEMP%, or %USERPROFILE%.
//
// Returns:
//
//	A non-empty string containing the temporary directory path.
//
// Example:
//
//	fmt.Println(sysx.TempDir()) // "/tmp"
func TempDir() string {
	return os.TempDir()
}

// HomeDir returns the current user's home directory.
//
// The lookup is performed via os/user.Current(). On most systems this
// respects the HOME environment variable.
//
// Returns:
//
//	(string, error): the home directory path and nil on success, or an empty
//	string and a non-nil error if the current user cannot be determined.
//
// Example:
//
//	home, err := sysx.HomeDir()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(home)
func HomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

// MustHomeDir returns the current user's home directory and panics if the
// lookup fails.
//
// Returns:
//
//	A non-empty string containing the home directory path.
//
// Example:
//
//	fmt.Println(sysx.MustHomeDir())
func MustHomeDir() string {
	h, err := HomeDir()
	if err != nil {
		panic("sysx: cannot retrieve home directory: " + err.Error())
	}
	return h
}

// WorkingDir returns the current working directory of the process.
//
// It delegates directly to os.Getwd.
//
// Returns:
//
//	(string, error): the working directory path and nil on success, or an
//	empty string and a non-nil error on failure.
//
// Example:
//
//	wd, err := sysx.WorkingDir()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(wd)
func WorkingDir() (string, error) {
	return os.Getwd()
}

// MustWorkingDir returns the current working directory of the process and
// panics if the lookup fails.
//
// Returns:
//
//	A non-empty string containing the current working directory path.
//
// Example:
//
//	fmt.Println(sysx.MustWorkingDir())
func MustWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic("sysx: cannot retrieve working directory: " + err.Error())
	}
	return wd
}

// ///////////////////////////
// Section: File reading
// ///////////////////////////

// ReadFile reads the entire contents of the file at path and returns them as
// a byte slice.
//
// Parameters:
//   - `path`: the file system path to read.
//
// Returns:
//
//	([]byte, error): the file contents and nil on success, or nil and a
//	non-nil error if the file does not exist or cannot be read.
//
// Example:
//
//	data, err := sysx.ReadFile("/etc/hosts")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%d bytes read\n", len(data))
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileString reads the entire contents of the file at path and returns
// them as a string.
//
// Parameters:
//   - `path`: the file system path to read.
//
// Returns:
//
//	(string, error): the file contents and nil on success, or an empty string
//	and a non-nil error if the file does not exist or cannot be read.
//
// Example:
//
//	content, err := sysx.ReadFileString("/etc/hostname")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(strings.TrimSpace(content))
func ReadFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadLines reads the file at path and returns its contents as a slice of
// strings, one element per line. Line endings ("\n" and "\r\n") are stripped
// from each element. An empty file returns a non-nil empty slice.
//
// The file is read with a buffered scanner, making it efficient even for
// large files as long as no single line exceeds bufio.MaxScanTokenSize
// (64 KiB by default).
//
// Parameters:
//   - `path`: the file system path to read.
//
// Returns:
//
//	([]string, error): lines of the file and nil on success, or a partial
//	result and a non-nil error on failure.
//
// Example:
//
//	lines, err := sysx.ReadLines("/var/log/app.log")
//	for _, l := range lines {
//	    fmt.Println(l)
//	}
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines := make([]string, 0, 64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// StreamLines opens the file at path and calls handler for each line in order.
// Processing stops immediately and the handler's error is returned when handler
// returns a non-nil value. A bufio.Scanner error encountered during reading is
// also returned. Line endings are stripped before handler is invoked.
//
// StreamLines is designed for memory-efficient processing of large files: only
// one line is held in memory at a time.
//
// Parameters:
//   - `path`:    the file system path to read.
//   - `handler`: the function called for each line; return a non-nil error to stop.
//
// Returns:
//
//	An error if the file could not be opened, the scanner failed, or handler
//	returned a non-nil error; nil on success.
//
// Example:
//
//	count := 0
//	err := sysx.StreamLines("/var/log/access.log", func(line string) error {
//	    if strings.Contains(line, "ERROR") {
//	        count++
//	    }
//	    return nil
//	})
func StreamLines(path string, handler func(string) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := handler(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// ///////////////////////////
// Section: File writing
// ///////////////////////////

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

// ///////////////////////////
// Section: Atomic and concurrency-safe writes
// ///////////////////////////

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

// ///////////////////////////
// Section: SafeFileWriter
// ///////////////////////////

// SafeFileWriter provides concurrency-safe append and overwrite operations
// targeting a single file path. A single SafeFileWriter instance can be shared
// across goroutines; all write operations are serialised by an internal mutex.
//
// Create a SafeFileWriter with NewSafeFileWriter and optionally adjust the
// file permission with WithPerm before sharing across goroutines.
type SafeFileWriter struct {
	mu   sync.Mutex
	path string
	perm os.FileMode
}

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

// ///////////////////////////
// Section: Path-level in-process locking
// ///////////////////////////

// fileMutexes is the package-level registry of per-path in-process mutexes.
var fileMutexes sync.Map // map[string]*sync.Mutex

// getFileMutex returns the mutex associated with path, creating one if necessary.
func getFileMutex(path string) *sync.Mutex {
	v, _ := fileMutexes.LoadOrStore(path, &sync.Mutex{})
	return v.(*sync.Mutex)
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
