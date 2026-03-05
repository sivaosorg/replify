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
// The command subsystem offers two API layers.
//
// Builder API – configure once, execute cleanly:
//
//	sysx.NewCommand("bash").
//	    WithArgs("-c", "echo hello").
//	    WithTimeout(5 * time.Second).
//	    WithEnv("APP_ENV=prod").
//	    WithDir("/tmp").
//	    Execute()       // returns *CommandResult with Stdout, Stderr, ExitCode, Duration
//
// Convenience functions – one-liners for common patterns:
//
//	sysx.RunCommand(name, args...)                             // structured *CommandResult
//	sysx.ExecCommand(name, args...)                            // run, discard output
//	sysx.ExecCommandContext(ctx, name, args...)                // run under context
//	sysx.ExecOutput(name, args...)                             // run, capture combined output
//	sysx.ExecOutputLines(name, args...)                        // run, capture stdout as []string
//	sysx.ExecStreaming(stdout, stderr, name, args...)           // stream output in real time
//	sysx.ExecAsync(name, args...)                              // start without waiting
//	sysx.ExecPipeline([]string{...}, []string{...})            // shell-style pipe chain
//	sysx.ExecCommandWithTimeout(timeout, name, args...)        // run with deadline
//	sysx.ExecOutputWithTimeout(timeout, name, args...)         // run with deadline, capture
//	sysx.ExecCommandInDir(dir, name, args...)                  // run in directory
//	sysx.ExecOutputInDir(dir, name, args...)                   // run in directory, capture
//
// # File System Utilities
//
// Existence and type checks:
//
//	sysx.FileExists(path)    // true when a file system entry exists at path
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
// File reading:
//
//	sysx.ReadFile(path)                     // []byte contents
//	sysx.ReadFileString(path)               // string contents
//	sysx.ReadLines(path)                    // []string, one element per line
//	sysx.StreamLines(path, handler)         // line-by-line callback, memory-efficient
//
// File writing:
//
//	sysx.WriteFile(path, data)              // create/truncate, write bytes
//	sysx.WriteFileString(path, content)     // create/truncate, write string
//	sysx.AppendFile(path, data)             // append bytes
//	sysx.AppendString(path, content)        // append string
//	sysx.WriteLines(path, lines)            // write slice as newline-terminated lines
//
// Concurrency-safe and atomic writes:
//
//	sysx.AtomicWriteFile(path, data)        // temp-file + rename, prevents partial reads
//	sysx.WriteFileLocked(path, data)        // per-path in-process mutex, serialises writers
//	sysx.NewSafeFileWriter(path)            // reusable mutex-protected writer for one path
//
// All functions in this package are safe for concurrent use.
package sysx
