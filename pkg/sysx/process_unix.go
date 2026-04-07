//go:build !windows

package sysx

import (
	"os"
	"syscall"
)

// processExists checks whether a process with the given PID is running on
// Unix-like systems by sending signal 0. A nil or EPERM error indicates the
// process exists.
func processExists(pid int) bool {
	err := syscall.Kill(pid, syscall.Signal(0))
	return err == nil || err == syscall.EPERM
}

// killProcess sends SIGTERM to the process identified by pid.
func killProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGTERM)
}

// killProcessForcefully sends SIGKILL to the process identified by pid.
func killProcessForcefully(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGKILL)
}
