package sysx

import (
	"os"
)

// CreateDir creates the named directory, along with any necessary parent
// directories, using permission 0755.
//
// It is equivalent to os.MkdirAll(path, 0755) and is idempotent: calling it
// on a directory that already exists is a no-op.
//
// Parameters:
//   - `path`: the directory path to create.
//
// Returns:
//
//	An error if the directory could not be created; nil on success.
//
// Example:
//
//	if err := sysx.CreateDir("/tmp/app/logs"); err != nil {
//	    log.Fatal(err)
//	}
func CreateDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// RemoveDir removes the directory at path together with all of its contents.
//
// It wraps os.RemoveAll. Calling RemoveDir on a path that does not exist is
// not an error; the function returns nil in that case.
//
// Parameters:
//   - `path`: the directory path to remove.
//
// Returns:
//
//	An error if the directory could not be removed; nil on success or if the
//	path does not exist.
//
// Example:
//
//	if err := sysx.RemoveDir("/tmp/app/cache"); err != nil {
//	    log.Fatal(err)
//	}
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// ListDir returns the names of all entries (files, directories, symlinks) in
// the directory at path, in the order returned by os.ReadDir (lexicographic
// order). Only entry names are returned; use filepath.Join to build full paths.
//
// Parameters:
//   - `path`: the directory path to read.
//
// Returns:
//
//	([]string, error): the entry names and nil on success, or nil and a
//	non-nil error if the directory cannot be read.
//
// Example:
//
//	names, err := sysx.ListDir("/etc")
//	for _, name := range names {
//	    fmt.Println(name)
//	}
func ListDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

// ListDirFiles returns only the names of regular files in the directory at
// path. Directories, symbolic links, and other non-regular entries are
// excluded. Entries are returned in lexicographic order.
//
// Parameters:
//   - `path`: the directory path to read.
//
// Returns:
//
//	([]string, error): the file names and nil on success, or nil and a
//	non-nil error if the directory cannot be read.
//
// Example:
//
//	files, err := sysx.ListDirFiles("/var/log")
//	for _, f := range files {
//	    fmt.Println(f)
//	}
func ListDirFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Type().IsRegular() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// ListDirDirs returns only the names of subdirectories inside the directory
// at path. Regular files, symbolic links, and other non-directory entries are
// excluded. Entries are returned in lexicographic order.
//
// Parameters:
//   - `path`: the directory path to read.
//
// Returns:
//
//	([]string, error): the subdirectory names and nil on success, or nil and
//	a non-nil error if the directory cannot be read.
//
// Example:
//
//	dirs, err := sysx.ListDirDirs("/home")
//	for _, d := range dirs {
//	    fmt.Println(d)
//	}
func ListDirDirs(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
