//go:build !windows

package sysx

import (
	"os"
	"syscall"
)

// lockFile applies an advisory lock to the given file handle.
// If isWrite is true, it applies an exclusive lock (LOCK_EX);
// otherwise, it applies a shared lock (LOCK_SH).
func lockFile(f *os.File, isWrite bool) error {
	how := syscall.LOCK_SH
	if isWrite {
		how = syscall.LOCK_EX
	}
	// LOCK_NB is not used here to ensure we wait for the lock.
	// If the user wants non-blocking, we'd need another API.
	return syscall.Flock(int(f.Fd()), how)
}

// unlockFile releases the advisory lock on the given file handle.
func unlockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
