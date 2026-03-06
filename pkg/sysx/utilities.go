package sysx

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/sivaosorg/replify/pkg/conv"
	"github.com/sivaosorg/replify/pkg/strutil"
)

// setenv sets the environment variable named by key to value.
//
// It delegates directly to os.Setenv and propagates any error.
//
// Parameters:
//   - `key`:   the name of the environment variable to set.
//   - `value`: the value to assign.
//
// Returns:
//
//	An error if the variable could not be set, or nil on success.
//
// Example:
//
//	if err := sysx.setenv("LOG_LEVEL", "debug"); err != nil {
//	    log.Fatal(err)
//	}
func setenv(key, value string) error {
	if strutil.IsEmpty(key) {
		return errors.New("sysx: key must not be empty")
	}
	return os.Setenv(key, value)
}

// linuxOSVersion retrieves the Linux distribution version from the standard
// /etc/os-release file.
//
// The function specifically looks for the "PRETTY_NAME" key, which provides a
// human-readable name for the distribution. If the file is inaccessible, or
// if the key is not found, it returns runtime.GOOS as a fallback.
//
// Returns:
//
//	A string identifying the Linux distribution and version, or "linux".
func linuxOSVersion() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return runtime.GOOS
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			val := strings.TrimPrefix(line, "PRETTY_NAME=")
			val = strings.Trim(val, `"`)
			if val != "" {
				return val
			}
		}
	}
	return runtime.GOOS
}

// darwinOSVersion retrieves the macOS product version using the system's
// sw_vers utility.
//
// It executes "sw_vers -productVersion" and returns the resulting string
// after trimming whitespace. If the command execution fails or returns no
// content, it returns runtime.GOOS as a fallback.
//
// Returns:
//
//	A string identifying the macOS product version (e.g. "14.2.1"), or "darwin".
func darwinOSVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return runtime.GOOS
	}
	v := strings.TrimSpace(string(out))
	if strutil.IsEmpty(v) {
		return runtime.GOOS
	}
	return v
}

// getFileMutex retrieves the sync.Mutex associated with the given file path.
//
// It performs an atomic load-or-store operation on a global map of mutexes.
// If no mutex exists for the path, a new one is created and returned.
//
// Returns:
//
//	A non-nil *sync.Mutex unique to the provided path.
func getFileMutex(path string) *sync.Mutex {
	v, _ := fileMutexes.LoadOrStore(path, &sync.Mutex{})
	return v.(*sync.Mutex)
}

// parseBool parses a lowercase, trimmed string as a boolean.
//
// Recognised true  values: "1", "true", "yes", "on"
// Recognised false values: "0", "false", "no", "off"
//
// Returns (value, true) when the string is recognised, or (false, false) when
// it is not.
func parseBool(s string) (bool, bool) {
	if strutil.IsEmpty(s) {
		return false, false
	}
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "on":
		return true, true
	case "off":
		return false, true
	}
	v, err := conv.Bool(s)
	if err != nil {
		return false, false
	}
	return v, true
}

// splitLines splits s into individual lines by "\n", stripping any trailing
// "\r" (to handle "\r\n" line endings). A single trailing newline is consumed
// so that "a\nb\n" returns ["a","b"] rather than ["a","b",""]. An empty input
// returns nil.
func splitLines(s string) []string {
	if strutil.IsEmpty(s) {
		return nil
	}
	s = strings.TrimRight(s, "\n")
	parts := strings.Split(s, "\n")
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = strings.TrimRight(p, "\r")
	}
	return result
}
