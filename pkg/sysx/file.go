package sysx

import (
	"os"
	"os/user"
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
