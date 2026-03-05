package sysx

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ///////////////////////////
// Section: OS detection
// ///////////////////////////

// IsLinux reports whether the current operating system is Linux.
//
// The check is performed at runtime by comparing runtime.GOOS against the
// string "linux". The function is safe for concurrent use and allocates no
// memory.
//
// Returns:
//
//	A boolean value:
//	 - true  when the program is running on Linux;
//	 - false on any other operating system.
//
// Example:
//
//	if sysx.IsLinux() {
//	    fmt.Println("running on Linux")
//	}
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsDarwin reports whether the current operating system is macOS (Darwin).
//
// The check is performed at runtime by comparing runtime.GOOS against the
// string "darwin". The function is safe for concurrent use and allocates no
// memory.
//
// Returns:
//
//	A boolean value:
//	 - true  when the program is running on macOS/Darwin;
//	 - false on any other operating system.
//
// Example:
//
//	if sysx.IsDarwin() {
//	    fmt.Println("running on macOS")
//	}
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsWindows reports whether the current operating system is Windows.
//
// The check is performed at runtime by comparing runtime.GOOS against the
// string "windows". The function is safe for concurrent use and allocates no
// memory.
//
// Returns:
//
//	A boolean value:
//	 - true  when the program is running on Windows;
//	 - false on any other operating system.
//
// Example:
//
//	if sysx.IsWindows() {
//	    fmt.Println("running on Windows")
//	}
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// ///////////////////////////
// Section: Architecture / platform info
// ///////////////////////////

// OSName returns the operating system name as reported by the Go runtime.
//
// The value is identical to runtime.GOOS and is one of the platform strings
// defined by the Go toolchain (e.g. "linux", "darwin", "windows", "freebsd").
//
// Returns:
//
//	A non-empty string identifying the current operating system.
//
// Example:
//
//	fmt.Println(sysx.OSName()) // "linux"
func OSName() string {
	return runtime.GOOS
}

// Arch returns the CPU architecture as reported by the Go runtime.
//
// The value is identical to runtime.GOARCH and is one of the architecture
// strings defined by the Go toolchain (e.g. "amd64", "arm64", "386",
// "arm", "riscv64").
//
// Returns:
//
//	A non-empty string identifying the current CPU architecture.
//
// Example:
//
//	fmt.Println(sysx.Arch()) // "amd64"
func Arch() string {
	return runtime.GOARCH
}

// Is64Bit reports whether the program is running on a 64-bit architecture.
//
// The check inspects runtime.GOARCH for the well-known 64-bit identifiers:
// "amd64", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x",
// "riscv64", and "wasm". Any architecture not in this list is treated as
// non-64-bit.
//
// Returns:
//
//	A boolean value:
//	 - true  when the architecture is 64-bit;
//	 - false otherwise.
//
// Example:
//
//	if sysx.Is64Bit() {
//	    fmt.Println("64-bit architecture")
//	}
func Is64Bit() bool {
	switch runtime.GOARCH {
	case "amd64", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x", "riscv64", "wasm":
		return true
	}
	return false
}

// IsArm reports whether the program is running on an ARM-based architecture.
//
// Both 32-bit ARM ("arm") and 64-bit AArch64 ("arm64") targets are
// considered ARM-based.
//
// Returns:
//
//	A boolean value:
//	 - true  when the architecture is "arm" or "arm64";
//	 - false otherwise.
//
// Example:
//
//	if sysx.IsArm() {
//	    fmt.Println("ARM architecture")
//	}
func IsArm() bool {
	return runtime.GOARCH == "arm" || runtime.GOARCH == "arm64"
}

// ///////////////////////////
// Section: OS version
// ///////////////////////////

// OSVersion returns a best-effort human-readable operating system version
// string.
//
// The resolution strategy differs by platform:
//   - Linux:   reads /etc/os-release and returns the value of PRETTY_NAME=
//     if present; falls back to runtime.GOOS on any error.
//   - Darwin:  runs `sw_vers -productVersion` and returns its trimmed stdout;
//     falls back to runtime.GOOS on any error.
//   - Windows: returns a string composed of runtime.GOOS and runtime.GOARCH.
//   - Other:   returns runtime.GOOS.
//
// The function is not guaranteed to return a version number on every platform
// and environment. Callers that require precise version parsing should use
// platform-specific APIs.
//
// Returns:
//
//	A non-empty string describing the operating system version or name.
//
// Example:
//
//	fmt.Println(sysx.OSVersion()) // e.g. "Ubuntu 22.04.3 LTS"
func OSVersion() string {
	switch runtime.GOOS {
	case "linux":
		return linuxOSVersion()
	case "darwin":
		return darwinOSVersion()
	case "windows":
		return runtime.GOOS + "/" + runtime.GOARCH
	default:
		return runtime.GOOS
	}
}

// linuxOSVersion reads /etc/os-release and returns the PRETTY_NAME value.
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

// darwinOSVersion runs sw_vers to retrieve the macOS product version.
func darwinOSVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return runtime.GOOS
	}
	v := strings.TrimSpace(string(out))
	if v == "" {
		return runtime.GOOS
	}
	return v
}
