// Package sysx provides a lightweight, production-grade system utilities
// toolkit for interacting with the underlying operating system, process
// environment, runtime, network, and file system from within Go programs.
//
// All functions are designed to be safe for concurrent use and impose no
// external dependencies beyond the Go standard library.
//
// # Package Architecture
//
// Struct definitions and package-level globals live in type.go. Each file
// focuses on a single concern:
//
//   - type.go      — struct definitions (Command, CommandResult, SafeFileWriter), getters
//   - command.go   — command builder, Execute/Run/Output, and convenience functions
//   - file.go      — file existence/permission checks, read/write/atomic utilities
//   - dir.go       — directory creation, removal, and listing utilities
//   - path.go      — path manipulation helpers (base name, dir name, extension, join)
//   - io.go        — stream-oriented I/O helpers (CountLines, Head, CopyFile)
//   - net.go       — IP classification, port probing, address helpers, connectivity
//   - env.go       — environment variable helpers
//   - os.go        — OS detection and version
//   - process.go   — process existence, signalling, lookup
//   - runtime.go   — hostname, PID, UID, CPU count, memory stats
//   - utilities.go — internal helpers and UserInfo()
//   - entry.go     — SystemInfo(), IsPrivileged()
//   - resource.go  — Resource builder (NewResource + With*/From*) and lifecycle
//   - memblob.go   — MemBlob: in-memory ReadSeekCloser backing
//   - tempfile.go  — TempFile: on-disk temporary file backing with auto-cleanup
//   - spillbuf.go  — private hybrid memory→disk buffer for streaming producers
//   - mime.go      — MimeFromName helper and ErrNilResource
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
//	sysx.Getenv(key, fallback)      // env var with fallback
//	sysx.MustGetenv(key)            // panics if env var is absent or empty
//	sysx.Hasenv(key)                // true when var exists and is non-empty
//	sysx.GetenvInt(key, fallback)   // env var parsed as int
//	sysx.GetenvBool(key, fallback)  // env var parsed as bool
//	sysx.GetenvSlice(key, sep)      // env var split into a slice
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
// Builder API — configure once, execute cleanly:
//
//	sysx.NewCommand("bash").
//		WithArgs("-c", "echo hello").
//		WithTimeout(5 * time.Second).
//		WithEnv("APP_ENV=prod").
//		WithDir("/tmp").
//		Execute()       // returns *CommandResult
//
// All CommandResult fields are unexported; read them via accessors:
//
//	res.Stdout()    string
//	res.Stderr()    string
//	res.ExitCode()  int
//	res.Duration()  time.Duration
//	res.Err()       error
//	res.IsSuccess()   bool
//	res.Combined()  string
//
// Convenience functions:
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
//	sysx.FileMode(path)      // permission bits of the file
//	sysx.FileModTime(path)   // last modification time of the file
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
//	sysx.WriteFileLocked(path, data)        // per-path in-process mutex, serialize writers
//	sysx.NewSafeFileWriter(path)            // reusable mutex-protected writer for one path
//
// # Directory Utilities
//
//	sysx.CreateDir(path)               // create directory and all parents (MkdirAll, 0755)
//	sysx.RemoveDir(path)               // remove directory and all contents (RemoveAll)
//	sysx.ListDir(path)                 // []string of all entry names in a directory
//	sysx.ListDirFiles(path)            // []string of regular file names only
//	sysx.ListDirDirs(path)             // []string of subdirectory names only
//
// # Path Helpers
//
//	sysx.BaseName(path)                // last element of path (filepath.Base)
//	sysx.DirName(path)                 // directory component of path (filepath.Dir)
//	sysx.Ext(path)                     // file extension including leading dot (filepath.Ext)
//	sysx.AbsPath(path)                 // absolute representation of path (filepath.Abs)
//	sysx.JoinPath(elem...)             // join path elements (filepath.Join)
//	sysx.CleanPath(path)               // lexically clean path (filepath.Clean)
//	sysx.SplitPath(path)               // split into (dir, file) components
//
// # Stream I/O Utilities
//
//	sysx.CountLines(path)              // count newline-delimited lines in a file
//	sysx.Head(path, n)                 // first n lines of a file as []string
//	sysx.CopyFile(src, dst)            // copy a file from src to dst
//	sysx.TruncateFile(path, size)      // truncate or extend a file to size bytes
//
// # Network Utilities
//
// IP classification:
//
//	sysx.IsIPv4(ip)          // true for valid IPv4 dotted-decimal strings
//	sysx.IsIPv6(ip)          // true for valid IPv6 strings (non-IPv4)
//	sysx.IsLocalIP(ip)       // true for loopback, link-local, and RFC 1918 addresses
//
// Port probing:
//
//	sysx.IsPortOpen(host, port)      // true when TCP connect to host:port succeeds
//	sysx.IsPortAvailable(port)       // true when the port can be bound locally
//
// Network address helpers:
//
//	sysx.GetLocalIP()          // first non-loopback IPv4 from local interfaces
//	sysx.GetPublicIP()         // public IP via https://api.ipify.org (requires internet)
//	sysx.GetInterfaceIPs()     // all IPv4+IPv6 addresses across all interfaces
//
// URL and host helpers:
//
//	sysx.IsValidHost(host)       // true when host resolves via DNS or is a valid IP
//	sysx.ParseHostPort(addr)     // split "host:port" into (host, port, error)
//	sysx.IsValidURL(rawURL)      // true when rawURL has a valid scheme and host
//
// Connectivity:
//
//	sysx.PingHost(host)                              // TCP probe to port 80, 5s timeout
//	sysx.CheckTCPConn(host, port, timeout)     // TCP connect with explicit timeout
//
// # Resource: storage-agnostic exports
//
// Resource is the envelope every data-exporting workflow (reports, dumps,
// archives, attachments, backups, media generation, document rendering)
// returns instead of a raw *os.File. Consumers — Telegram bots, S3
// uploader, email services, HTTP download handlers, Discord bots,
// generic storage services, backup workers — depend exclusively on
// Resource and on the standard io interfaces.
//
// Builder API:
//
//	res := sysx.NewResource().
//	    WithName("user-report.csv").
//	    WithContentType(sysx.MimeCSV).      // optional; auto-derived from name
//	    FromBytes(payload)                   // or FromString / FromReader / FromTempFile / FromFile
//	defer res.Close()
//
// Configuration setters (all chainable):
//
//	res.WithName(string)                 // suggested filename
//	res.WithContentType(string)          // IANA media type
//	res.WithSize(int64)                  // explicit size for custom backings
//	res.WithContent(ReadSeekCloser)      // attach a custom backing
//	res.WithSpillThreshold(int64)        // FromReader memory ceiling
//	res.WithTempPattern(string)          // os.CreateTemp pattern
//	res.WithTempDir(string)              // parent dir for temp files
//	res.WithRemoveOnClose(bool)          // arm/disarm temp-file cleanup
//
// Loaders (single-input "From*" methods):
//
//	res.FromBytes([]byte)                          *Resource
//	res.FromString(string)                         *Resource
//	res.FromFile(*os.File)                         (*Resource, error)
//	res.FromTempFile(func(io.Writer) error)        (*Resource, error)
//	res.FromReader(io.Reader)                      (*Resource, error)
//
// Lifecycle and consumption:
//
//	res.Close()                          // release backing (idempotent)
//	res.Rewind()                         // seek to offset 0
//	res.CopyTo(io.Writer)                // stream payload into a sink
//	res.Drain()                          // consume and discard
//
// Built-in backings:
//
//	*sysx.MemBlob   — bytes-only, no I/O; Close is a no-op
//	*sysx.TempFile  — on-disk temporary file with auto-removal on Close
//	(private)       — hybrid memory→spill buffer used by FromReader
//
// MIME constants and helpers:
//
//	sysx.MimeCSV, MimeJSON, MimePDF, MimeGZIP, MimeOctetStream, …
//	sysx.MimeFromName("dump.tar.gz")     // returns MimeGZIP
//
// # Concurrency and Safety Notes
//
// All exported functions in this package are safe for concurrent use.
// WriteFileLocked and SafeFileWriter provide in-process serialization via
// sync.Mutex. For cross-process atomicity, use AtomicWriteFile which relies on
// the atomic rename(2) syscall on POSIX systems.
//
// # Cross-Platform Notes
//
// Path separator: use JoinPath and CleanPath instead of manual string
// concatenation to ensure correct behavior on all platforms.
//
// File permissions: mode bits set by CreateDir (0755) and write functions
// (0644) are subject to the process umask on Unix. On Windows, mode bits are
// an approximation; all files typically report as executable.
//
// Symlinks:
//
//	IsFile, IsDir, FileSize, FileMode, and FileModTime
//
// follow symbolic links. IsSymlink uses os.Lstat and does not follow links.
package sysx
