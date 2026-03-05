// Package sysx provides a lightweight, production-grade system utilities
// toolkit for interacting with the underlying operating system, process
// environment, and runtime from within Go programs.
//
// All functions are designed to be safe for concurrent use and impose no
// external dependencies beyond the Go standard library.
//
// # OS Detection
//
//	sysx.IsLinux()    // true when running on Linux
//	sysx.IsDarwin()   // true when running on macOS/Darwin
//	sysx.IsWindows()  // true when running on Windows
//	sysx.OSName()     // returns runtime.GOOS (e.g. "linux", "darwin", "windows")
//	sysx.Arch()       // returns runtime.GOARCH (e.g. "amd64", "arm64")
//	sysx.OSVersion()  // best-effort human-readable OS version string
//
// # Runtime Information
//
//	sysx.Hostname()       // os.Hostname(), returns (string, error)
//	sysx.PID()            // current process identifier
//	sysx.NumCPU()         // number of logical CPUs
//	sysx.NumGoroutine()   // number of active goroutines
//	sysx.GoVersion()      // Go runtime version string (e.g. "go1.24.0")
//	sysx.MemStats()       // snapshot of runtime memory statistics
//
// # Environment Utilities
//
//	sysx.GetEnv(key, fallback)      // env var with fallback
//	sysx.MustGetEnv(key)            // panics if env var is absent or empty
//	sysx.HasEnv(key)                // true when var exists and is non-empty
//	sysx.GetEnvInt(key, fallback)   // env var parsed as int
//	sysx.GetEnvBool(key, fallback)  // env var parsed as bool
//	sysx.GetEnvSlice(key, sep)      // env var split into a slice
//	sysx.EnvMap()                   // all env vars as map[string]string
//
// # Process Utilities
//
//	sysx.ProcessExists(pid)         // true when process with pid is running
//	sysx.KillProcess(pid)           // send SIGTERM to process
//	sysx.KillProcessForcefully(pid) // send SIGKILL to process
//	sysx.CurrentProcessName()       // base name of the current executable
//	sysx.FindProcessByPID(pid)      // *os.Process for the given pid
//
// # Command Execution
//
//	sysx.ExecCommand(name, args...)                            // run command
//	sysx.ExecOutput(name, args...)                             // run and capture combined output
//	sysx.ExecCommandWithTimeout(timeout, name, args...)        // run with deadline
//	sysx.ExecOutputWithTimeout(timeout, name, args...)         // run with deadline, capture output
//	sysx.ExecCommandInDir(dir, name, args...)                  // run in directory
//	sysx.ExecOutputInDir(dir, name, args...)                   // run in directory, capture output
//
// # File System Utilities
//
//	sysx.FileExists(path)    // true when a file exists at path
//	sysx.DirExists(path)     // true when a directory exists at path
//	sysx.IsFile(path)        // true when path is a regular file
//	sysx.IsDir(path)         // true when path is a directory
//	sysx.IsSymlink(path)     // true when path is a symbolic link
//	sysx.IsExecutable(path)  // true when file is executable by owner
//	sysx.IsReadable(path)    // true when file is readable by owner
//	sysx.IsWritable(path)    // true when file can be opened for writing
//	sysx.FileSize(path)      // size of file in bytes
//	sysx.HomeDir()           // user's home directory
//	sysx.WorkingDir()        // current working directory
//
// All functions in this package are safe for concurrent use.
package sysx
