//go:build windows

package sysx

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	processQueryLimitedInformation = 0x1000
	stillActive                    = 259 // STILL_ACTIVE value returned by GetExitCodeProcess
)

var (
	modkernel32Windows     = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess        = modkernel32Windows.NewProc("OpenProcess")
	procCloseHandle        = modkernel32Windows.NewProc("CloseHandle")
	procGetExitCodeProcess = modkernel32Windows.NewProc("GetExitCodeProcess")
)

// processExists checks whether a process with the given PID is running on
// Windows. Uses OpenProcess with PROCESS_QUERY_LIMITED_INFORMATION to verify
// the PID refers to a living process. Returns false for any PID that cannot
// be opened or that has already exited.
//
// Note: os.FindProcess always returns a non-nil handle on Windows regardless
// of whether the PID is still alive, so we use the Windows API directly here.
func processExists(pid int) bool {
	handle, _, _ := procOpenProcess.Call(
		uintptr(processQueryLimitedInformation),
		0,
		uintptr(pid),
	)
	if handle == 0 {
		return false
	}
	defer procCloseHandle.Call(handle)

	var exitCode uint32
	ret, _, _ := procGetExitCodeProcess.Call(handle, uintptr(unsafe.Pointer(&exitCode)))
	if ret == 0 {
		return false
	}
	// STILL_ACTIVE (259) means the process is still running.
	return exitCode == stillActive
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
