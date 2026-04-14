//go:build windows

package sysx

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	lockfileExclusiveLock   = 2
	lockfileFailImmediately = 1
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = modkernel32.NewProc("LockFileEx")
	procUnlockFileEx = modkernel32.NewProc("UnlockFileEx")
)

// lockFile applies an advisory lock to the given file handle.
// On Windows, it uses LockFileEx to provide exclusive or shared access
// for the entire file.
func lockFile(f *os.File, isWrite bool) error {
	var flags uint32 = 0
	if isWrite {
		flags = lockfileExclusiveLock
	}

	// We use 0xffffffff for both Low and High parts to mean "infinity"
	// (entire file) or just very large.
	ol := new(syscall.Overlapped)
	r1, _, err := procLockFileEx.Call(
		f.Fd(),
		uintptr(flags),
		0,
		uintptr(0xffffffff),
		uintptr(0xffffffff),
		uintptr(unsafe.Pointer(ol)),
	)
	if r1 == 0 {
		// err is always a syscall.Errno after a Call; compare directly
		// instead of relying on the locale-dependent Error() string.
		if err != syscall.Errno(0) {
			return err
		}
		return syscall.EINVAL
	}
	return nil
}

// unlockFile releases the advisory lock on the given file handle.
func unlockFile(f *os.File) error {
	ol := new(syscall.Overlapped)
	r1, _, err := procUnlockFileEx.Call(
		f.Fd(),
		0,
		uintptr(0xffffffff),
		uintptr(0xffffffff),
		uintptr(unsafe.Pointer(ol)),
	)
	if r1 == 0 {
		if err != syscall.Errno(0) {
			return err
		}
		return syscall.EINVAL
	}
	return nil
}
