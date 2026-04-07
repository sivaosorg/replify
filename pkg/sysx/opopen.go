package sysx

import (
	"fmt"
	"os"
)

// OpenFile opens the named file with the specified flags (e.g. CWA, RO, RW)
// and permission bits (e.g. 0644).
//
// By using the FileOpenFlags type, this function provides IDE auto-completion
// and type safety, guiding developers to use the predefined convenience constants.
//
// Parameters:
//   - `path`:  the file system path to open.
//   - `flags`: the type-safe FileOpenFlags combination.
//   - `perm`:  the os.FileMode to apply if the file is created.
//
// Returns:
//
//	(*os.File, error): a handle to the open file and nil on success, or
//	nil and a non-nil error if the file cannot be opened.
//
// Example:
//
//	// IDEs will suggest CWA, CWT, RO, etc. when typing the second argument.
//	f, err := sysx.OpenFile("app.log", sysx.CWA, 0644)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer f.Close()
func OpenFile(path string, flags FileOpenFlags, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, int(flags), perm)
}

// AtomicOpenFile opens a file with strong consistency and isolation guarantees.
// It uses a combination of in-process mutexes (via getFileMutex) and
// cross-process advisory locking (flock/LockFileEx) to ensure that the file
// handle is protected from concurrent access.
//
// The lock is automatically released when the returned f.Close() is called.
//
// Parameters:
//   - `path`:  the file system path to open.
//   - `flags`: the type-safe FileOpenFlags combination.
//   - `perm`:  the os.FileMode to apply if the file is created.
//
// Returns:
//
//	(*os.File, error): a handle to the isolated file and nil on success.
//
// Example:
//
//	f, err := sysx.AtomicOpenFile("config.json", sysx.RW, 0644)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer f.Close() // releases the lock
func AtomicOpenFile(path string, flags FileOpenFlags, perm os.FileMode) (*os.File, error) {
	// 1. In-process synchronisation
	mu := getFileMutex(path)
	mu.Lock()
	defer mu.Unlock()

	// 2. Open the file
	f, err := os.OpenFile(path, int(flags), perm)
	if err != nil {
		return nil, err
	}

	// 3. Advisory locking
	// Determine mode: RO uses shared, anything else uses exclusive
	isWrite := (int(flags) & (os.O_WRONLY | os.O_RDWR)) != 0
	if err := lockFile(f, isWrite); err != nil {
		f.Close()
		return nil, fmt.Errorf("sysx: failed to acquire advisory lock on %s: %w", path, err)
	}

	return f, nil
}

// String returns a human-readable label for the flag combination.
//
// Parameters:
//   - `f`: the FileOpenFlags combination.
//
// Returns:
//
//	A string containing the label.
//
// Example:
//
//	fmt.Println(sysx.CWA.String()) // "CWA (CREATE|WRONLY|APPEND)"
func (f FileOpenFlags) String() string {
	switch f {
	case CWA:
		return "CWA (CREATE|WRONLY|APPEND)"
	case CWT:
		return "CWT (CREATE|WRONLY|TRUNC)"
	case CW:
		return "CW (CREATE|WRONLY)"
	case AW:
		return "AW (APPEND|WRONLY)"
	case TW:
		return "TW (TRUNC|WRONLY)"
	case RO:
		return "RO (RDONLY)"
	case RW:
		return "RW (RDWR)"
	case CRW:
		return "CRW (RDWR|CREATE)"
	case CRWT:
		return "CRWT (RDWR|CREATE|TRUNC)"
	case ARW:
		return "ARW (RDWR|APPEND)"
	case CRWA:
		return "CRWA (RDWR|CREATE|APPEND)"
	case TARW:
		return "TARW (RDWR|APPEND|TRUNC)"
	case TRW:
		return "TRW (RDWR|TRUNC)"
	default:
		return "Unknown (Custom Flags)"
	}
}

// IsValid reports whether the flag combination is one of the predefined constants.
//
// Parameters:
//   - `f`: the FileOpenFlags combination.
//
// Returns:
//
//	bool: true if the flag combination is valid, false otherwise.
//
// Example:
//
//	fmt.Println(sysx.CWA.IsValid()) // true
//	fmt.Println(sysx.FileOpenFlags(999).IsValid()) // false
func (f FileOpenFlags) IsValid() bool {
	switch f {
	case CWA, CWT, CW, AW, TW, RO, RW, CRW, CRWT, ARW, CRWA, TARW, TRW:
		return true
	default:
		return false
	}
}
