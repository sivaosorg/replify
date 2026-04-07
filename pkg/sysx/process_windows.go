//go:build windows

package sysx

import (
	"os"
)

// processExists checks whether a process with the given PID is running on
// Windows. os.FindProcess always succeeds on Windows (it does not verify that
// the PID actually refers to a running process), so this is a best-effort check.
func processExists(pid int) bool {
	p, err := os.FindProcess(pid)
	return err == nil && p != nil
}

// killProcess sends an interrupt signal to the process identified by pid.
// On Windows, SIGTERM does not exist; os.Interrupt (CTRL_C_EVENT) is used as
// the closest equivalent for a graceful shutdown request.
func killProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(os.Interrupt)
}

// killProcessForcefully terminates the process identified by pid immediately.
// On Windows this calls TerminateProcess internally via p.Kill().
func killProcessForcefully(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}
