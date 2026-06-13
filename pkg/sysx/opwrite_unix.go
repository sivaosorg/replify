//go:build !windows

package sysx

import "os"

// renameReplace renames src to dst, replacing dst atomically if it exists.
// On POSIX systems (Linux, macOS), os.Rename already provides this guarantee.
func renameReplace(src, dst string) error {
	return os.Rename(src, dst)
}
