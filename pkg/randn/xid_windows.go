//go:build windows

package randn

import (
	"fmt"
	"syscall"
	"unsafe"
)

// readPlatformMachineID retrieves the host's machine GUID on Windows.
// It reads the 'MachineGuid' value from the HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography registry key.
//
// Returns:
//   - A string containing the machine GUID (expected format: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx").
//   - An error if the registry key cannot be opened, the value cannot be read, or the GUID has an unexpected format.
func readPlatformMachineID() (string, error) {
	var h syscall.Handle

	regKeyCryptoPtr, err := syscall.UTF16PtrFromString(`SOFTWARE\Microsoft\Cryptography`)
	if err != nil {
		return "", fmt.Errorf(`error reading registry key "SOFTWARE\Microsoft\Cryptography": %w`, err)
	}

	err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, regKeyCryptoPtr, 0, syscall.KEY_READ|syscall.KEY_WOW64_64KEY, &h)
	if err != nil {
		return "", err
	}
	defer func() { _ = syscall.RegCloseKey(h) }()

	const syscallRegBufLen = 74
	const uuidLen = 36

	var regBuf [syscallRegBufLen]uint16
	bufLen := uint32(syscallRegBufLen)
	var valType uint32

	mGuidPtr, err := syscall.UTF16PtrFromString(`MachineGuid`)
	if err != nil {
		return "", fmt.Errorf("error reading machine GUID: %w", err)
	}

	err = syscall.RegQueryValueEx(h, mGuidPtr, nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
	if err != nil {
		return "", fmt.Errorf("error reading MachineGuid value: %w", err)
	}

	hostID := syscall.UTF16ToString(regBuf[:])
	hostIDLen := len(hostID)
	if hostIDLen != uuidLen {
		return "", fmt.Errorf("randn: unexpected MachineGuid length %d: %q", hostIDLen, hostID)
	}

	return hostID, nil
}
