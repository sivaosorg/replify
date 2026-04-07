package sysx

import (
	"os"
	"path/filepath"
)

// ProcessExists reports whether a process with the given PID is currently
// running.
//
// On Unix-like systems the check is performed by sending signal 0 to the
// process via syscall.Kill. A nil error indicates the process exists and the
// caller has permission to signal it. A syscall.EPERM error also indicates
// the process exists (but the caller lacks permission to signal it).
//
// On Windows, os.FindProcess never returns an error for a non-existent PID;
// the function therefore returns true for any positive PID on Windows. This
// is a known platform limitation.
//
// Parameters:
//   - `pid`: the process identifier to check.
//
// Returns:
//
//	A boolean value:
//	 - true  when the process appears to be running;
//	 - false when the PID is invalid (≤ 0) or the process does not exist.
//
// Example:
//
//	if sysx.ProcessExists(os.Getpid()) {
//	    fmt.Println("current process is running")
//	}
func ProcessExists(pid int) bool {
	if pid <= 0 {
		return false
	}
	return processExists(pid)
}

// KillProcess sends SIGTERM to the process identified by pid.
//
// SIGTERM is the standard termination signal that gives the process an
// opportunity to clean up before exiting. Use KillProcessForcefully when
// an immediate kill is required.
//
// Parameters:
//   - `pid`: the process identifier of the target process.
//
// Returns:
//
//	An error if the process could not be found or signalled, or nil on success.
//
// Example:
//
//	if err := sysx.KillProcess(pid); err != nil {
//	    log.Printf("failed to SIGTERM process %d: %v", pid, err)
//	}
func KillProcess(pid int) error {
	return killProcess(pid)
}

// KillProcessForcefully sends SIGKILL to the process identified by pid.
//
// SIGKILL cannot be caught or ignored by the target process; it is
// terminated immediately. Prefer KillProcess (SIGTERM) when a graceful
// shutdown is possible.
//
// Parameters:
//   - `pid`: the process identifier of the target process.
//
// Returns:
//
//	An error if the process could not be found or signalled, or nil on success.
//
// Example:
//
//	if err := sysx.KillProcessForcefully(pid); err != nil {
//	    log.Printf("failed to SIGKILL process %d: %v", pid, err)
//	}
func KillProcessForcefully(pid int) error {
	return killProcessForcefully(pid)
}

// CurrentProcessName returns the base name of the executable that started
// the current process.
//
// The value is derived from os.Executable followed by filepath.Base. If the
// executable path cannot be determined, an empty string is returned.
//
// Returns:
//
//	A string containing the base filename of the current executable.
//
// Example:
//
//	fmt.Println(sysx.CurrentProcessName()) // e.g. "myapp"
func CurrentProcessName() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Base(exe)
}

// FindProcessByPID returns the os.Process associated with the given PID.
//
// On Unix-like systems, os.FindProcess always succeeds for any integer PID.
// Use ProcessExists to verify that the process is actually running before
// operating on the returned value.
//
// Parameters:
//   - `pid`: the process identifier to look up.
//
// Returns:
//
//	(*os.Process, error): the process handle and nil on success, or nil and a
//	non-nil error if the lookup fails.
//
// Example:
//
//	proc, err := sysx.FindProcessByPID(os.Getpid())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(proc.Pid)
func FindProcessByPID(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}
