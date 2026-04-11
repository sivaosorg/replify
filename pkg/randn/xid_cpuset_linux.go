//go:build linux

package randn

import (
	"hash/crc32"
	"os"
)

// pidContainerOffset returns an XOR offset derived from /proc/self/cpuset.
// On Linux, this helps disambiguate XIDs generated within different
// cgroups or containers that may share the same host PID.
// Returns 0 if the file cannot be read or contains no useful data.
func pidContainerOffset() int {
	b, err := os.ReadFile("/proc/self/cpuset")
	if err != nil || len(b) <= 1 {
		return 0
	}
	return int(crc32.ChecksumIEEE(b))
}
