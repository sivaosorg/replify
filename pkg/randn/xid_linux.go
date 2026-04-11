//go:build linux || freebsd

package randn

import "os"

// readPlatformMachineID attempts to read the host's machine ID from the filesystem.
// It first tries to read /etc/machine-id, and if that fails or is empty, it falls back
// to reading /sys/class/dmi/id/product_uuid.
//
// Returns:
//   - The machine ID as a string (trimmed of whitespace).
//   - An error if both attempts fail.
func readPlatformMachineID() (string, error) {
	b, err := os.ReadFile("/etc/machine-id")
	if err != nil || len(b) == 0 {
		b, err = os.ReadFile("/sys/class/dmi/id/product_uuid")
	}
	return string(b), err
}
