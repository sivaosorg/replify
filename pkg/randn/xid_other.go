//go:build !darwin && !linux && !freebsd && !windows

package randn

import "errors"

// readPlatformMachineID is a fallback implementation that returns an error.
// It is used on platforms where no specific host ID retrieval logic is implemented.
//
// Returns:
//   - An error indicating that the host ID cannot be retrieved on this platform.
func readPlatformMachineID() (string, error) {
	return "", errors.New("cannot find host id on this platform")
}
