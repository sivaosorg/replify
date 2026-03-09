package sysx

import (
	"path/filepath"
)

// BaseName returns the last element of path. Trailing path separators are
// removed before extracting the final element. If path is empty, BaseName
// returns ".".
//
// It wraps filepath.Base and is safe for concurrent use.
//
// Parameters:
//   - `path`: the file system path.
//
// Returns:
//
//	A string containing the base name component of path.
//
// Example:
//
//	sysx.BaseName("/etc/hosts")   // "hosts"
//	sysx.BaseName("/usr/bin/")    // "bin"
//	sysx.BaseName("")             // "."
func BaseName(path string) string {
	return filepath.Base(path)
}

// DirName returns all but the last element of path, typically the path's
// directory. The returned path does not end in a separator unless it is the
// root directory. If path is empty, DirName returns ".".
//
// It wraps filepath.Dir and is safe for concurrent use.
//
// Parameters:
//   - `path`: the file system path.
//
// Returns:
//
//	A string containing the directory component of path.
//
// Example:
//
//	sysx.DirName("/etc/hosts")   // "/etc"
//	sysx.DirName("/usr/bin/git") // "/usr/bin"
//	sysx.DirName("file.txt")     // "."
func DirName(path string) string {
	return filepath.Dir(path)
}

// Ext returns the file name extension used by path. The extension is the
// suffix beginning at the final dot in the last element of path; it is empty
// if there is no dot.
//
// It wraps filepath.Ext and is safe for concurrent use.
//
// Parameters:
//   - `path`: the file system path.
//
// Returns:
//
//	A string containing the extension (including the leading dot), or an
//	empty string when path has no extension.
//
// Example:
//
//	sysx.Ext("archive.tar.gz") // ".gz"
//	sysx.Ext("/etc/hosts")     // ""
//	sysx.Ext("README.md")      // ".md"
func Ext(path string) string {
	return filepath.Ext(path)
}

// AbsPath returns an absolute representation of path. If path is not already
// absolute, it is joined with the current working directory to make it
// absolute. The result is cleaned via filepath.Clean.
//
// It wraps filepath.Abs and is safe for concurrent use.
//
// Parameters:
//   - `path`: the file system path to convert.
//
// Returns:
//
//	(string, error): the absolute path and nil on success, or an empty string
//	and a non-nil error if the current working directory cannot be determined.
//
// Example:
//
//	abs, err := sysx.AbsPath("relative/dir")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(abs) // e.g. "/home/user/relative/dir"
func AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

// JoinPath joins any number of path elements into a single path, separating
// them with the OS-specific path separator. Empty elements are ignored. The
// result is cleaned via filepath.Clean. An empty result is returned when all
// elements are empty.
//
// It wraps filepath.Join and is safe for concurrent use.
//
// Parameters:
//   - `elem`: zero or more path components to join.
//
// Returns:
//
//	A string containing the joined and cleaned path.
//
// Example:
//
//	sysx.JoinPath("/usr", "local", "bin") // "/usr/local/bin"
//	sysx.JoinPath("a", "b", "c")          // "a/b/c" (on Unix)
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// CleanPath returns the shortest path name equivalent to path by purely
// lexical processing. It applies the same rules as filepath.Clean: it
// eliminates redundant separators, "." and ".." elements.
//
// It wraps filepath.Clean and is safe for concurrent use.
//
// Parameters:
//   - `path`: the file system path to clean.
//
// Returns:
//
//	A string containing the cleaned path. If the argument is empty, CleanPath
//	returns ".".
//
// Example:
//
//	sysx.CleanPath("/usr//local/./bin/../lib") // "/usr/local/lib"
//	sysx.CleanPath("")                          // "."
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// SplitPath splits path immediately following the final separator, separating
// it into a directory and file name component. If there is no separator in
// path, SplitPath returns an empty dir and file set to path. The returned
// values have the property that path = dir + file.
//
// It wraps filepath.Split and is safe for concurrent use.
//
// Parameters:
//   - `path`: the file system path to split.
//
// Returns:
//
//	(dir, file string): the directory (including trailing separator) and the
//	base file name.
//
// Example:
//
//	dir, file := sysx.SplitPath("/usr/local/bin/git")
//	// dir  = "/usr/local/bin/"
//	// file = "git"
func SplitPath(path string) (dir, file string) {
	return filepath.Split(path)
}
