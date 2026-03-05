package sysx

import (
	"fmt"
	"strconv"
)

// ///////////////////////////
// Section: System-wide information
// ///////////////////////////

// SystemInfo returns a map of key/value strings summarising the most
// commonly needed runtime and operating system attributes of the current
// process.
//
// The map always contains the following keys:
//   - "os"         – operating system name (runtime.GOOS)
//   - "arch"       – CPU architecture (runtime.GOARCH)
//   - "hostname"   – machine hostname; empty string on lookup failure
//   - "pid"        – current process identifier as a decimal string
//   - "go_version" – Go runtime version (e.g. "go1.24.0")
//   - "executable" – path of the current executable; empty on failure
//   - "num_cpu"    – number of logical CPUs as a decimal string
//
// Returns:
//
//	A non-nil map[string]string containing the system information entries.
//
// Example:
//
//	info := sysx.SystemInfo()
//	for k, v := range info {
//	    fmt.Printf("%s = %s\n", k, v)
//	}
func SystemInfo() map[string]string {
	host, _ := Hostname()
	exe, _ := ExecutablePath()
	return map[string]string{
		"os":         OSName(),
		"arch":       Arch(),
		"hostname":   host,
		"pid":        strconv.Itoa(PID()),
		"go_version": GoVersion(),
		"executable": exe,
		"num_cpu":    fmt.Sprintf("%d", NumCPU()),
	}
}

// ///////////////////////////
// Section: Privilege detection
// ///////////////////////////

// IsPrivileged reports whether the current process is running with root
// (super-user) privileges.
//
// On Unix-like systems this is determined by checking whether the effective
// user identifier is 0. On Windows, UID() returns -1 and this function
// always returns false; use the Windows API for accurate privilege checks.
//
// Returns:
//
//	A boolean value:
//	 - true  when UID() == 0 (running as root on Unix);
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsPrivileged() {
//	    fmt.Println("running as root")
//	}
func IsPrivileged() bool {
	return UID() == 0
}
