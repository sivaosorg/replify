//go:build windows

package sysx

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const movefileReplaceExisting = 0x1

// renameReplace renames src to dst, replacing dst if it already exists.
// On Windows, os.Rename cannot atomically replace an existing file; this
// implementation calls MoveFileEx with MOVEFILE_REPLACE_EXISTING which is
// the idiomatic Windows equivalent of POSIX rename(2).
func renameReplace(src, dst string) error {
	srcUTF16, err := syscall.UTF16PtrFromString(src)
	if err != nil {
		return fmt.Errorf("sysx: renameReplace: encode src %q: %w", src, err)
	}
	dstUTF16, err := syscall.UTF16PtrFromString(dst)
	if err != nil {
		return fmt.Errorf("sysx: renameReplace: encode dst %q: %w", dst, err)
	}
	r, _, e := procMoveFileExW.Call(
		uintptr(unsafe.Pointer(srcUTF16)),
		uintptr(unsafe.Pointer(dstUTF16)),
		movefileReplaceExisting,
	)
	if r == 0 {
		if e != nil && e != syscall.Errno(0) {
			return &os.PathError{Op: "MoveFileEx", Path: dst, Err: e}
		}
		return fmt.Errorf("sysx: renameReplace: MoveFileEx %q -> %q failed", src, dst)
	}
	return nil
}
