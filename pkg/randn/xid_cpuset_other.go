//go:build !linux

package randn

// pidContainerOffset returns 0 on non-Linux platforms where
// /proc/self/cpuset is not available.
func pidContainerOffset() int {
	return 0
}
