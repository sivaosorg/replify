//go:build darwin

package randn

import (
	"errors"
	"os/exec"
	"strings"
)

// readPlatformMachineID retrieves the host's platform UUID on Darwin (macOS).
// It executes the 'ioreg' command to query the IOPlatformExpertDevice class and
// parses the output for the 'IOPlatformUUID' property.
//
// Returns:
//   - A string containing the platform's UUID in lowercase.
//   - An error if 'ioreg' is not found, if the command fails, or if the UUID cannot be parsed.
func readPlatformMachineID() (string, error) {
	ioreg, err := exec.LookPath("ioreg")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(ioreg, "-rd1", "-c", "IOPlatformExpertDevice")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.SplitAfter(line, `" = "`)
			if len(parts) == 2 {
				uuid := strings.TrimRight(parts[1], `"`)
				return strings.ToLower(uuid), nil
			}
		}
	}

	return "", errors.New("cannot find host id")
}
