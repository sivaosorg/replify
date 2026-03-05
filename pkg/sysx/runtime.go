package sysx

import (
	"os"
	"runtime"
)

// ///////////////////////////
// Section: Host and process identity
// ///////////////////////////

// Hostname returns the host name reported by the kernel.
//
// It delegates directly to os.Hostname and propagates any error returned by
// the operating system.
//
// Returns:
//
//	(string, error): the hostname and nil on success, or an empty string and
//	a non-nil error on failure.
//
// Example:
//
//	name, err := sysx.Hostname()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(name)
func Hostname() (string, error) {
	return os.Hostname()
}

// MustHostname returns the host name reported by the kernel and panics if the
// lookup fails.
//
// Use this function only in contexts where a hostname lookup failure is
// considered a programmer error or an unrecoverable condition (e.g. during
// program initialisation).
//
// Returns:
//
//	A non-empty string containing the current hostname.
//
// Example:
//
//	host := sysx.MustHostname()
//	fmt.Println(host)
func MustHostname() string {
	h, err := os.Hostname()
	if err != nil {
		panic("sysx: cannot retrieve hostname: " + err.Error())
	}
	return h
}

// PID returns the process identifier of the current process.
//
// It delegates directly to os.Getpid and always returns a positive integer
// on well-behaved platforms.
//
// Returns:
//
//	An int representing the current process identifier.
//
// Example:
//
//	fmt.Println("PID:", sysx.PID())
func PID() int {
	return os.Getpid()
}

// PPID returns the process identifier of the parent of the current process.
//
// It delegates directly to os.Getppid. On Windows, this always returns 0.
//
// Returns:
//
//	An int representing the parent process identifier.
//
// Example:
//
//	fmt.Println("PPID:", sysx.PPID())
func PPID() int {
	return os.Getppid()
}

// UID returns the numeric user identifier of the calling process.
//
// On Windows, it always returns -1 (os.Getuid is unsupported there).
//
// Returns:
//
//	An int representing the current user identifier.
//
// Example:
//
//	fmt.Println("UID:", sysx.UID())
func UID() int {
	return os.Getuid()
}

// GID returns the numeric group identifier of the calling process.
//
// On Windows, it always returns -1 (os.Getgid is unsupported there).
//
// Returns:
//
//	An int representing the current group identifier.
//
// Example:
//
//	fmt.Println("GID:", sysx.GID())
func GID() int {
	return os.Getgid()
}

// ///////////////////////////
// Section: Executable path
// ///////////////////////////

// ExecutablePath returns the path name for the executable that started the
// current process.
//
// It delegates to os.Executable. The path is not guaranteed to be an absolute
// path on all platforms; callers requiring a canonical path should call
// filepath.EvalSymlinks on the result.
//
// Returns:
//
//	(string, error): the path and nil on success, or an empty string and a
//	non-nil error on failure.
//
// Example:
//
//	path, err := sysx.ExecutablePath()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(path)
func ExecutablePath() (string, error) {
	return os.Executable()
}

// MustExecutablePath returns the path name for the executable that started
// the current process and panics if the lookup fails.
//
// Returns:
//
//	A non-empty string containing the executable path.
//
// Example:
//
//	fmt.Println(sysx.MustExecutablePath())
func MustExecutablePath() string {
	p, err := os.Executable()
	if err != nil {
		panic("sysx: cannot retrieve executable path: " + err.Error())
	}
	return p
}

// ///////////////////////////
// Section: Runtime statistics
// ///////////////////////////

// NumCPU returns the number of logical CPUs usable by the current process.
//
// The value may be less than the total number of CPUs on the machine if the
// process has been constrained (e.g. via cgroups or GOMAXPROCS).
//
// Returns:
//
//	A positive int representing the number of logical CPUs.
//
// Example:
//
//	fmt.Println("CPUs:", sysx.NumCPU())
func NumCPU() int {
	return runtime.NumCPU()
}

// NumGoroutine returns the number of goroutines that currently exist in the
// current program.
//
// The count includes all goroutines that have been created and not yet
// terminated, regardless of whether they are running or blocked.
//
// Returns:
//
//	A positive int representing the current goroutine count.
//
// Example:
//
//	fmt.Println("goroutines:", sysx.NumGoroutine())
func NumGoroutine() int {
	return runtime.NumGoroutine()
}

// GoVersion returns the Go runtime version string of the binary.
//
// The value is identical to runtime.Version() and has the form "go1.X.Y"
// (e.g. "go1.24.0").
//
// Returns:
//
//	A non-empty string identifying the Go version used to build the binary.
//
// Example:
//
//	fmt.Println(sysx.GoVersion()) // "go1.24.0"
func GoVersion() string {
	return runtime.Version()
}

// MemStats returns a snapshot of the Go runtime memory allocator statistics
// for the current process.
//
// The snapshot is obtained by calling runtime.ReadMemStats, which stops the
// world briefly on older Go versions. On Go 1.16+, the stop-the-world pause
// is significantly shorter.
//
// Returns:
//
//	A runtime.MemStats value populated with the current memory statistics.
//
// Example:
//
//	stats := sysx.MemStats()
//	fmt.Printf("heap alloc: %d bytes\n", stats.HeapAlloc)
func MemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}
