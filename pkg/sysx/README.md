# sysx

**sysx** is a lightweight, production-grade system utilities toolkit for Go, providing a clean and consistent API for OS detection, runtime introspection, environment variable management, process control, command execution, and file system queries — all built exclusively on the Go standard library.

## Overview

The `sysx` package eliminates the boilerplate of writing low-level system queries from scratch. It addresses common pain points like:

- **OS / Architecture detection** – know at runtime whether you are on Linux, macOS, Windows, 64-bit, or ARM
- **Runtime introspection** – read hostname, PID, UID, GID, number of CPUs, goroutine count, Go version, and memory stats in a single call
- **Environment management** – read, write, and parse environment variables with typed helpers (int, bool, slice) and sensible fallbacks
- **Process utilities** – check whether a PID is alive, send signals, look up processes by PID
- **Command execution** – a builder API plus convenience functions supporting timeout, working-directory override, environment injection, real-time streaming, async launch, and shell-style pipelines
- **File system helpers** – existence checks, type checks, permission checks, size queries, directory lookups, read/write helpers, atomic writes, and concurrency-safe writers

**Problem Solved:** Querying the operating system involves a patchwork of `os`, `os/exec`, `runtime`, `syscall`, and `os/user` calls scattered across many packages. `sysx` unifies these into a single, coherent API with uniform error handling and well-documented behaviour.

## Design Philosophy

- **Zero external dependencies** — only the Go standard library is used
- **Safe for concurrent use** — all exported functions are stateless or read-only with respect to shared state
- **Explicit errors, not silent failures** — functions return errors rather than hiding them; `Must*` variants are provided where panicking on failure is a deliberate choice
- **Platform-aware** — differences between Linux, macOS, and Windows are documented and handled gracefully
- **No shell interpolation** — command execution helpers use `os/exec` directly, never `sh -c`, to avoid injection risks
- **DRY internals** — the `Command` builder is the single implementation; all convenience functions delegate to it

## Package Architecture

| File | Responsibility |
|------|----------------|
| `doc.go` | Package-level godoc documentation |
| `os.go` | OS and architecture detection (`IsLinux`, `IsDarwin`, `IsWindows`, `OSVersion`, …) |
| `runtime.go` | Runtime information (`Hostname`, `PID`, `UID`, `GoVersion`, `MemStats`, …) |
| `env.go` | Environment variable helpers (`GetEnv`, `MustGetEnv`, `HasEnv`, typed getters, `EnvMap`) |
| `process.go` | Process utilities (`ProcessExists`, `KillProcess`, `FindProcessByPID`, …) |
| `command.go` | Command execution — `Command` builder, `CommandResult`, and all convenience functions |
| `file.go` | File system helpers — existence/type/permission checks, read/write functions, atomic and concurrency-safe writers |
| `utilities.go` | Internal helpers (`isZero`, `trimSpace`, `parseBoolString`, `splitLines`, `commandBuffer`) and exported `UserInfo()` |
| `entry.go` | Top-level convenience (`SystemInfo`, `IsPrivileged`) |

## Installation

```bash
go get github.com/sivaosorg/replify
```

Import the package:

```go
import "github.com/sivaosorg/replify/pkg/sysx"
```

**Requirements:** Go 1.24.0 or higher

## API Reference

### OS Detection (`os.go`)

```go
sysx.IsLinux()   bool   // true on Linux
sysx.IsDarwin()  bool   // true on macOS
sysx.IsWindows() bool   // true on Windows

sysx.OSName()    string // runtime.GOOS  ("linux", "darwin", "windows", …)
sysx.Arch()      string // runtime.GOARCH ("amd64", "arm64", "386", …)
sysx.Is64Bit()   bool   // true for amd64, arm64, ppc64, …
sysx.IsArm()     bool   // true for arm and arm64

sysx.OSVersion() string // best-effort OS version string
```

**OSVersion resolution:**

| Platform | Source |
|----------|--------|
| Linux | `PRETTY_NAME` field from `/etc/os-release` |
| macOS | Output of `sw_vers -productVersion` |
| Windows | `runtime.GOOS + "/" + runtime.GOARCH` |
| Other | `runtime.GOOS` |

### Runtime Information (`runtime.go`)

```go
sysx.Hostname()          (string, error) // os.Hostname()
sysx.MustHostname()      string          // panics on error
sysx.PID()               int             // os.Getpid()
sysx.PPID()              int             // os.Getppid()
sysx.UID()               int             // os.Getuid()  (-1 on Windows)
sysx.GID()               int             // os.Getgid()  (-1 on Windows)
sysx.ExecutablePath()    (string, error) // os.Executable()
sysx.MustExecutablePath() string         // panics on error
sysx.NumCPU()            int             // runtime.NumCPU()
sysx.NumGoroutine()      int             // runtime.NumGoroutine()
sysx.GoVersion()         string          // runtime.Version() — e.g. "go1.24.0"
sysx.MemStats()          runtime.MemStats
```

### Environment Utilities (`env.go`)

```go
sysx.GetEnv(key, fallback string) string          // env var or fallback
sysx.MustGetEnv(key string) string                // panics if absent/empty
sysx.HasEnv(key string) bool                      // true when set and non-empty
sysx.SetEnv(key, value string) error              // os.Setenv
sysx.UnsetEnv(key string) error                   // os.Unsetenv
sysx.GetEnvInt(key string, fallback int) int      // parsed int or fallback
sysx.GetEnvBool(key string, fallback bool) bool   // parsed bool or fallback
sysx.GetEnvSlice(key, sep string) []string        // split by sep, nil if unset
sysx.Environ() []string                           // os.Environ()
sysx.EnvMap() map[string]string                   // all env vars as map
```

**Bool string recognition (case-insensitive):**

| Truthy | Falsy |
|--------|-------|
| `1`, `true`, `yes`, `on` | `0`, `false`, `no`, `off` |

### Process Utilities (`process.go`)

```go
sysx.ProcessExists(pid int) bool                     // true when process is running
sysx.KillProcess(pid int) error                      // SIGTERM
sysx.KillProcessForcefully(pid int) error            // SIGKILL
sysx.CurrentProcessName() string                     // filepath.Base of executable
sysx.FindProcessByPID(pid int) (*os.Process, error)  // os.FindProcess
```

### Command Execution (`command.go`)

The command subsystem has two API layers that share a single implementation.

#### Builder API — `Command` and `CommandResult`

```go
// Create and configure a Command:
cmd := sysx.NewCommand("bash").
    WithArgs("-c", "echo hello").
    WithTimeout(5 * time.Second).
    WithEnv("APP_ENV=prod").
    WithDir("/tmp").
    WithStdin(os.Stdin).
    WithStdout(os.Stdout).
    WithStderr(os.Stderr)

// Execute and get a structured result:
res := cmd.Execute()
// res.Stdout   string         — captured stdout (empty when WithStdout was used)
// res.Stderr   string         — captured stderr (empty when WithStderr was used)
// res.ExitCode int            — 0 on success, -1 when undetermined
// res.Duration time.Duration  — wall-clock execution time
// res.Error    error          — nil on success

// Convenience methods on *CommandResult:
res.Success()   bool   // true when Error == nil
res.Combined()  string // Stdout + Stderr

// Shorter variants on *Command:
cmd.Run()           error          // Execute().Error
cmd.Output()        (string, error) // Execute().Combined(), Execute().Error
```

#### Convenience Functions

```go
// Structured result shortcut
sysx.RunCommand(name string, args ...string) *CommandResult

// Discard output
sysx.ExecCommand(name string, args ...string) error
sysx.ExecCommandContext(ctx context.Context, name string, args ...string) error
sysx.ExecCommandWithTimeout(timeout time.Duration, name string, args ...string) error
sysx.ExecCommandInDir(dir, name string, args ...string) error

// Capture combined stdout+stderr
sysx.ExecOutput(name string, args ...string) (string, error)
sysx.ExecOutputWithTimeout(timeout time.Duration, name string, args ...string) (string, error)
sysx.ExecOutputInDir(dir, name string, args ...string) (string, error)

// Capture stdout as lines
sysx.ExecOutputLines(name string, args ...string) ([]string, error)

// Stream output in real time to provided writers
sysx.ExecStreaming(stdout, stderr io.Writer, name string, args ...string) error

// Start asynchronously — caller calls cmd.Wait()
sysx.ExecAsync(name string, args ...string) (*exec.Cmd, error)

// Chain commands as a shell pipeline
sysx.ExecPipeline(commands ...[]string) (string, error)
```

### File System Utilities (`file.go`)

#### Existence and Type Checks

```go
sysx.FileExists(path string) bool   // any file system entry
sysx.DirExists(path string) bool    // directory only
sysx.IsFile(path string) bool       // regular file (follows symlinks)
sysx.IsDir(path string) bool        // directory (follows symlinks)
sysx.IsSymlink(path string) bool    // symbolic link (does NOT follow)
```

#### Permission Checks

```go
sysx.IsExecutable(path string) bool  // owner execute bit (0100)
sysx.IsReadable(path string) bool    // owner read bit    (0400)
sysx.IsWritable(path string) bool    // attempts O_WRONLY open
```

#### File Metadata and Directories

```go
sysx.FileSize(path string) (int64, error)
sysx.TempDir() string
sysx.HomeDir() (string, error)
sysx.MustHomeDir() string
sysx.WorkingDir() (string, error)
sysx.MustWorkingDir() string
```

#### File Reading

```go
sysx.ReadFile(path string) ([]byte, error)
sysx.ReadFileString(path string) (string, error)
sysx.ReadLines(path string) ([]string, error)
sysx.StreamLines(path string, handler func(string) error) error
```

#### File Writing

```go
sysx.WriteFile(path string, data []byte) error
sysx.WriteFileString(path string, content string) error
sysx.AppendFile(path string, data []byte) error
sysx.AppendString(path string, content string) error
sysx.WriteLines(path string, lines []string) error
```

#### Atomic and Concurrency-Safe Writes

```go
// Atomic replace via temp-file + rename (POSIX atomic)
sysx.AtomicWriteFile(path string, data []byte) error

// Per-path in-process mutex (multiple goroutines, same path)
sysx.WriteFileLocked(path string, data []byte) error

// Reusable mutex-protected writer for one path
w := sysx.NewSafeFileWriter(path).WithPerm(0o600)
w.Write(data []byte) error
w.WriteString(s string) error
w.Overwrite(data []byte) error  // atomic replace
```

### Entry-Point Helpers (`entry.go`)

```go
sysx.SystemInfo() map[string]string
// keys: "os", "arch", "hostname", "pid", "go_version", "executable", "num_cpu"

sysx.IsPrivileged() bool  // true when UID() == 0 (root on Unix)
```

### Utility Helpers (`utilities.go`)

```go
sysx.UserInfo() string  // "uid=1000 gid=1000"
```

## Usage Examples

### Basic OS / Architecture Check

```go
fmt.Println("OS:      ", sysx.OSName())
fmt.Println("Arch:    ", sysx.Arch())
fmt.Println("Version: ", sysx.OSVersion())
fmt.Println("64-bit:  ", sysx.Is64Bit())
fmt.Println("ARM:     ", sysx.IsArm())
```

### Environment Configuration

```go
host  := sysx.GetEnv("DB_HOST", "localhost")
port  := sysx.GetEnvInt("DB_PORT", 5432)
debug := sysx.GetEnvBool("DEBUG", false)
tags  := sysx.GetEnvSlice("TAGS", ",")

// Panics at startup if DATABASE_URL is not configured
dsn := sysx.MustGetEnv("DATABASE_URL")
```

### System Information Snapshot

```go
info := sysx.SystemInfo()
for k, v := range info {
    fmt.Printf("%-12s = %s\n", k, v)
}
// os           = linux
// arch         = amd64
// hostname     = myserver
// pid          = 12345
// go_version   = go1.24.0
// executable   = /usr/local/bin/myapp
// num_cpu      = 8
```

## Command Execution Utilities

### Simple execution

```go
// Fire and forget
if err := sysx.ExecCommand("go", "generate", "./..."); err != nil {
    log.Fatal(err)
}

// Capture combined output
out, err := sysx.ExecOutput("git", "log", "--oneline", "-5")
fmt.Print(out)

// Capture stdout as lines
lines, err := sysx.ExecOutputLines("git", "branch", "-a")
for _, l := range lines {
    fmt.Println(l)
}
```

### Advanced Command Execution Examples

#### Builder with full configuration

```go
res := sysx.NewCommand("bash").
    WithArgs("-c", "echo $APP_ENV && ls /tmp").
    WithTimeout(10 * time.Second).
    WithEnv("APP_ENV=production").
    WithDir("/tmp").
    Execute()

if !res.Success() {
    log.Printf("command failed (exit %d, %.2fs):\nstdout: %s\nstderr: %s",
        res.ExitCode, res.Duration.Seconds(), res.Stdout, res.Stderr)
}
```

#### Stream output in real time

```go
// Print build output line-by-line as it happens
err := sysx.ExecStreaming(os.Stdout, os.Stderr, "go", "build", "-v", "./...")
```

#### Async execution

```go
cmd, err := sysx.ExecAsync("server", "--port", "8080")
if err != nil {
    log.Fatal(err)
}
// ... do other work while server runs ...
if err := cmd.Wait(); err != nil {
    log.Println("server exited:", err)
}
```

#### Shell-style pipeline

```go
out, err := sysx.ExecPipeline(
    []string{"cat", "/var/log/syslog"},
    []string{"grep", "ERROR"},
    []string{"wc", "-l"},
)
fmt.Println("error lines:", strings.TrimSpace(out))
```

#### Context-controlled execution

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := sysx.ExecCommandContext(ctx, "go", "test", "-race", "./..."); err != nil {
    log.Fatal(err)
}
```

#### Structured result inspection

```go
res := sysx.RunCommand("golangci-lint", "run", "./...")
fmt.Printf("exit=%d duration=%v\n", res.ExitCode, res.Duration.Round(time.Millisecond))
if !res.Success() {
    fmt.Println("--- stdout ---")
    fmt.Print(res.Stdout)
    fmt.Println("--- stderr ---")
    fmt.Print(res.Stderr)
}
```

## File I/O Utilities

### Reading files

```go
// Read entire file
data, err := sysx.ReadFile("/etc/hosts")

// Read as string
content, err := sysx.ReadFileString("/etc/hostname")
fmt.Println(strings.TrimSpace(content))

// Read line-by-line into slice
lines, err := sysx.ReadLines("/var/log/app.log")
fmt.Printf("%d log entries\n", len(lines))
```

### Streaming large files

```go
// Memory-efficient: one line in memory at a time
err := sysx.StreamLines("/var/log/big.log", func(line string) error {
    if strings.Contains(line, "CRITICAL") {
        alertOnCall(line)
    }
    return nil
})
```

**Early termination example:**

```go
var found string
err := sysx.StreamLines("/etc/passwd", func(line string) error {
    if strings.HasPrefix(line, "root:") {
        found = line
        return errors.New("stop") // stop processing immediately
    }
    return nil
})
```

### Writing files

```go
// Create / truncate
sysx.WriteFile("/tmp/output.bin", data)
sysx.WriteFileString("/tmp/config.txt", configContent)

// Append
sysx.AppendFile("/var/log/app.log", []byte("new entry\n"))
sysx.AppendString("/var/log/app.log", "another entry\n")

// Write line slice (each element newline-terminated, buffered writer)
sysx.WriteLines("/tmp/list.txt", []string{"alpha", "beta", "gamma"})
```

## Concurrency-Safe File Operations

### Atomic writes — no partial reads

`AtomicWriteFile` uses the temp-file + rename pattern. On POSIX systems
`os.Rename` is atomic within the same filesystem, so readers will never
observe a partial write.

```go
// Safe for readers running concurrently:
if err := sysx.AtomicWriteFile("/etc/app/config.json", newJSON); err != nil {
    log.Fatal(err)
}
```

### Per-path in-process locking

`WriteFileLocked` serialises concurrent writes to the same path using a
package-level `sync.Map` of mutexes. No goroutine observes a half-written file
within the same process.

```go
// Safe to call from multiple goroutines targeting the same path:
go sysx.WriteFileLocked("/tmp/shared.dat", payload1)
go sysx.WriteFileLocked("/tmp/shared.dat", payload2)
```

### `SafeFileWriter` — reusable concurrent writer

Share a single `SafeFileWriter` instance across goroutines for safe append
operations. The internal mutex serialises every `Write`, `WriteString`, and
`Overwrite` call.

```go
w := sysx.NewSafeFileWriter("/var/log/app.log").WithPerm(0o600)

// Fan-out writes from multiple goroutines — no external lock required:
for i := 0; i < 100; i++ {
    go func(n int) {
        w.WriteString(fmt.Sprintf("goroutine %d\n", n))
    }(i)
}
```

Use `Overwrite` for atomic replacement through the same writer:

```go
w.Overwrite(updatedConfig) // temp-file + rename under the mutex
```

## Performance Considerations

| Operation | Implementation detail |
|-----------|----------------------|
| `ReadLines` / `StreamLines` | `bufio.Scanner` — one line in memory, large files handled efficiently |
| `WriteLines` | `bufio.Writer` — batches small writes into a single syscall |
| `commandBuffer` (internal) | `strings.Builder` — avoids `bytes.Buffer` allocations for command output |
| `AtomicWriteFile` | single `os.Rename` syscall; no extra read pass |
| `WriteFileLocked` | `sync.Map` + per-path `sync.Mutex` — O(1) contention lookup |
| `SafeFileWriter` | single mutex per instance — lowest overhead for repeated appends to one file |

**Guidelines:**
- For large files, always prefer `StreamLines` over `ReadLines` to avoid loading the whole file into memory.
- Re-use a `SafeFileWriter` instance across goroutines rather than calling `AppendFile` concurrently (which uses no lock).
- Use `AtomicWriteFile` when readers and writers run concurrently; it prevents any reader from seeing partial data.

## Best Practices

1. **Prefer `Must*` only at startup** — use error-returning variants in library code; reserve `Must*` for `main()` or `init()` where a missing resource is fatal.
2. **Check `ProcessExists` before signalling** — `os.FindProcess` on Unix always succeeds; pair it with `ProcessExists` to avoid signalling phantom PIDs.
3. **Always set a timeout in long-running services** — use `WithTimeout` or `ExecCommandWithTimeout` when calling external programs.
4. **`ExecOutput*` merges stdout+stderr** — use the `Command` builder with `WithStdout`/`WithStderr` to separate them when needed.
5. **`EnvMap` is a snapshot** — the returned map is not updated when the environment changes; call again if you need fresh values.
6. **Prefer `AtomicWriteFile` over `WriteFile` for critical data** — `WriteFile` uses `os.WriteFile` which is not atomic on all kernels.
7. **`WriteFileLocked` is in-process only** — for cross-process safety combine it with `AtomicWriteFile` or use a platform lock.

## Platform Caveats

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| `UID()` / `GID()` | numeric | numeric | always -1 |
| `PPID()` | supported | supported | always 0 |
| `ProcessExists` | signal(0) | signal(0) | FindProcess (unreliable for dead PIDs) |
| `KillProcess` | SIGTERM | SIGTERM | limited |
| `KillProcessForcefully` | SIGKILL | SIGKILL | limited |
| `IsExecutable` / `IsReadable` | mode bits | mode bits | approximation |
| `IsWritable` | open test | open test | open test |
| `OSVersion` | `/etc/os-release` | `sw_vers` | GOOS/GOARCH string |
| `IsPrivileged` | UID==0 | UID==0 | always false |
| `AtomicWriteFile` rename | atomic (same fs) | atomic (same fs) | not guaranteed atomic |
| `ExecStreaming` / `ExecPipeline` | full support | full support | partial support |

## When to Use `sysx` vs stdlib

| Task | Prefer |
|------|--------|
| Quick OS check | `sysx.IsLinux()` vs `runtime.GOOS == "linux"` — both work; sysx adds readability |
| Env var with fallback | `sysx.GetEnv` — stdlib requires a conditional around `os.LookupEnv` |
| Typed env vars | `sysx.GetEnvInt` / `sysx.GetEnvBool` — no stdlib equivalent |
| Run command and capture output | `sysx.ExecOutput` — saves wiring `cmd.Stdout`, `cmd.Stderr` |
| Command with timeout | `sysx.ExecCommandWithTimeout` — saves `context.WithTimeout` boilerplate |
| Stream command output | `sysx.ExecStreaming` — pass any `io.Writer` |
| Check if a file exists | `sysx.FileExists` — stdlib `os.Stat` + error check |
| Read file as lines | `sysx.ReadLines` / `sysx.StreamLines` — stdlib requires `bufio.Scanner` setup |
| Write file atomically | `sysx.AtomicWriteFile` — stdlib has no atomic helper |
| Concurrent file writes | `sysx.SafeFileWriter` / `WriteFileLocked` — stdlib requires manual mutex |
| Check if a process is alive | `sysx.ProcessExists` — stdlib requires `syscall.Kill(pid, 0)` and type assertion |
| Complex I/O piping | `os/exec` directly — `sysx.ExecPipeline` covers linear chains only |
| Watching the file system | `fsnotify` or stdlib — outside `sysx` scope |

## Real-World Integration Examples

### Health-Check Endpoint

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    info := sysx.SystemInfo()
    json.NewEncoder(w).Encode(map[string]any{
        "status":     "ok",
        "hostname":   info["hostname"],
        "pid":        info["pid"],
        "go_version": info["go_version"],
        "num_cpu":    info["num_cpu"],
    })
}
```

### Startup Configuration Validation

```go
func mustLoadConfig() Config {
    return Config{
        DSN:     sysx.MustGetEnv("DATABASE_URL"),
        Port:    sysx.GetEnvInt("PORT", 8080),
        Debug:   sysx.GetEnvBool("DEBUG", false),
        Origins: sysx.GetEnvSlice("CORS_ORIGINS", ","),
    }
}
```

### Conditional Platform Logic

```go
func configureLogging() {
    switch {
    case sysx.IsLinux():
        // Use journald integration
    case sysx.IsDarwin():
        // Use os_log
    default:
        // Fallback to stderr
    }
}
```

### CI Pipeline Step Runner

```go
steps := [][]string{
    {"go", "vet", "./..."},
    {"go", "build", "./..."},
    {"go", "test", "-race", "-cover", "./..."},
}
for _, step := range steps {
    res := sysx.RunCommand(step[0], step[1:]...)
    fmt.Printf("$ %s => exit=%d (%.2fs)\n",
        strings.Join(step, " "), res.ExitCode, res.Duration.Seconds())
    if !res.Success() {
        fmt.Fprintln(os.Stderr, res.Combined())
        os.Exit(1)
    }
}
```

### Concurrent Log Aggregator

```go
w := sysx.NewSafeFileWriter("/var/log/aggregated.log")
var wg sync.WaitGroup
for _, source := range logSources {
    wg.Add(1)
    go func(src string) {
        defer wg.Done()
        sysx.StreamLines(src, func(line string) error {
            return w.WriteString(line + "\n")
        })
    }(source)
}
wg.Wait()
```

### Atomic Config Reload

```go
func reloadConfig(newCfg []byte) error {
    // Validate before replacing
    if err := validateConfig(newCfg); err != nil {
        return err
    }
    // Replace atomically so readers always see a complete file
    return sysx.AtomicWriteFile("/etc/app/config.json", newCfg)
}
```


## Overview

The `sysx` package eliminates the boilerplate of writing low-level system queries from scratch. It addresses common pain points like:

- **OS / Architecture detection** – know at runtime whether you are on Linux, macOS, Windows, 64-bit, or ARM
- **Runtime introspection** – read hostname, PID, UID, GID, number of CPUs, goroutine count, Go version, and memory stats in a single call
- **Environment management** – read, write, and parse environment variables with typed helpers (int, bool, slice) and sensible fallbacks
- **Process utilities** – check whether a PID is alive, send signals, look up processes by PID
- **Command execution** – run external programs with optional timeout, working-directory override, and combined output capture
- **File system helpers** – existence checks, type checks, permission checks, size queries, and directory lookups

**Problem Solved:** Querying the operating system involves a patchwork of `os`, `os/exec`, `runtime`, `syscall`, and `os/user` calls scattered across many packages. `sysx` unifies these into a single, coherent API with uniform error handling and well-documented behaviour.

## Design Philosophy

- **Zero external dependencies** — only the Go standard library is used
- **Safe for concurrent use** — all exported functions are stateless or read-only with respect to shared state
- **Explicit errors, not silent failures** — functions return errors rather than hiding them; `Must*` variants are provided where panicking on failure is a deliberate choice
- **Platform-aware** — differences between Linux, macOS, and Windows are documented and handled gracefully
- **No shell interpolation** — command execution helpers use `os/exec` directly, never `sh -c`, to avoid injection risks

## Package Architecture

| File | Responsibility |
|------|----------------|
| `doc.go` | Package-level godoc documentation |
| `os.go` | OS and architecture detection (`IsLinux`, `IsDarwin`, `IsWindows`, `OSVersion`, …) |
| `runtime.go` | Runtime information (`Hostname`, `PID`, `UID`, `GoVersion`, `MemStats`, …) |
| `env.go` | Environment variable helpers (`GetEnv`, `MustGetEnv`, `HasEnv`, typed getters, `EnvMap`) |
| `process.go` | Process utilities (`ProcessExists`, `KillProcess`, `FindProcessByPID`, …) |
| `command.go` | Command execution (`ExecCommand`, `ExecOutput`, timeout and directory variants) |
| `file.go` | File system helpers (`FileExists`, `DirExists`, permission checks, `FileSize`, `HomeDir`, …) |
| `utilities.go` | Internal helpers (`isZero`, `trimSpace`, `parseBoolString`) and exported `UserInfo()` |
| `entry.go` | Top-level convenience (`SystemInfo`, `IsPrivileged`) |

## Installation

```bash
go get github.com/sivaosorg/replify
```

Import the package:

```go
import "github.com/sivaosorg/replify/pkg/sysx"
```

**Requirements:** Go 1.24.0 or higher

## API Reference

### OS Detection (`os.go`)

```go
sysx.IsLinux()   bool   // true on Linux
sysx.IsDarwin()  bool   // true on macOS
sysx.IsWindows() bool   // true on Windows

sysx.OSName()    string // runtime.GOOS  ("linux", "darwin", "windows", …)
sysx.Arch()      string // runtime.GOARCH ("amd64", "arm64", "386", …)
sysx.Is64Bit()   bool   // true for amd64, arm64, ppc64, …
sysx.IsArm()     bool   // true for arm and arm64

sysx.OSVersion() string // best-effort OS version string
```

**OSVersion resolution:**

| Platform | Source |
|----------|--------|
| Linux | `PRETTY_NAME` field from `/etc/os-release` |
| macOS | Output of `sw_vers -productVersion` |
| Windows | `runtime.GOOS + "/" + runtime.GOARCH` |
| Other | `runtime.GOOS` |

### Runtime Information (`runtime.go`)

```go
sysx.Hostname()          (string, error) // os.Hostname()
sysx.MustHostname()      string          // panics on error
sysx.PID()               int             // os.Getpid()
sysx.PPID()              int             // os.Getppid()
sysx.UID()               int             // os.Getuid()  (-1 on Windows)
sysx.GID()               int             // os.Getgid()  (-1 on Windows)
sysx.ExecutablePath()    (string, error) // os.Executable()
sysx.MustExecutablePath() string         // panics on error
sysx.NumCPU()            int             // runtime.NumCPU()
sysx.NumGoroutine()      int             // runtime.NumGoroutine()
sysx.GoVersion()         string          // runtime.Version() — e.g. "go1.24.0"
sysx.MemStats()          runtime.MemStats
```

### Environment Utilities (`env.go`)

```go
sysx.GetEnv(key, fallback string) string          // env var or fallback
sysx.MustGetEnv(key string) string                // panics if absent/empty
sysx.HasEnv(key string) bool                      // true when set and non-empty
sysx.SetEnv(key, value string) error              // os.Setenv
sysx.UnsetEnv(key string) error                   // os.Unsetenv
sysx.GetEnvInt(key string, fallback int) int      // parsed int or fallback
sysx.GetEnvBool(key string, fallback bool) bool   // parsed bool or fallback
sysx.GetEnvSlice(key, sep string) []string        // split by sep, nil if unset
sysx.Environ() []string                           // os.Environ()
sysx.EnvMap() map[string]string                   // all env vars as map
```

**Bool string recognition (case-insensitive):**

| Truthy | Falsy |
|--------|-------|
| `1`, `true`, `yes`, `on` | `0`, `false`, `no`, `off` |

### Process Utilities (`process.go`)

```go
sysx.ProcessExists(pid int) bool                     // true when process is running
sysx.KillProcess(pid int) error                      // SIGTERM
sysx.KillProcessForcefully(pid int) error            // SIGKILL
sysx.CurrentProcessName() string                     // filepath.Base of executable
sysx.FindProcessByPID(pid int) (*os.Process, error)  // os.FindProcess
```

### Command Execution (`command.go`)

```go
sysx.ExecCommand(name string, args ...string) error
sysx.ExecOutput(name string, args ...string) (string, error)

sysx.ExecCommandWithTimeout(timeout time.Duration, name string, args ...string) error
sysx.ExecOutputWithTimeout(timeout time.Duration, name string, args ...string) (string, error)

sysx.ExecCommandInDir(dir, name string, args ...string) error
sysx.ExecOutputInDir(dir, name string, args ...string) (string, error)
```

All `ExecOutput*` functions capture **combined** stdout and stderr. None of the helpers perform shell interpolation; `name` is always passed directly to `os/exec`.

### File System Utilities (`file.go`)

```go
// Existence
sysx.FileExists(path string) bool
sysx.DirExists(path string) bool

// Type checks (follow symlinks)
sysx.IsFile(path string) bool
sysx.IsDir(path string) bool
sysx.IsSymlink(path string) bool  // does NOT follow symlinks

// Permission checks (mode-bit based; Windows is approximate)
sysx.IsExecutable(path string) bool  // owner execute bit (0100)
sysx.IsReadable(path string) bool    // owner read bit    (0400)
sysx.IsWritable(path string) bool    // tries os.O_WRONLY open

// Metadata
sysx.FileSize(path string) (int64, error)

// Special directories
sysx.TempDir() string
sysx.HomeDir() (string, error)
sysx.MustHomeDir() string
sysx.WorkingDir() (string, error)
sysx.MustWorkingDir() string
```

### Entry-Point Helpers (`entry.go`)

```go
sysx.SystemInfo() map[string]string
// keys: "os", "arch", "hostname", "pid", "go_version", "executable", "num_cpu"

sysx.IsPrivileged() bool  // true when UID() == 0 (root on Unix)
```

### Utility Helpers (`utilities.go`)

```go
sysx.UserInfo() string  // "uid=1000 gid=1000"
```

## Usage Examples

### Basic OS / Architecture Check

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/sysx"
)

func main() {
    fmt.Println("OS:      ", sysx.OSName())
    fmt.Println("Arch:    ", sysx.Arch())
    fmt.Println("Version: ", sysx.OSVersion())
    fmt.Println("64-bit:  ", sysx.Is64Bit())
    fmt.Println("ARM:     ", sysx.IsArm())
}
```

### Environment Configuration

```go
host  := sysx.GetEnv("DB_HOST", "localhost")
port  := sysx.GetEnvInt("DB_PORT", 5432)
debug := sysx.GetEnvBool("DEBUG", false)
tags  := sysx.GetEnvSlice("TAGS", ",")

// Panics at startup if DATABASE_URL is not configured
dsn := sysx.MustGetEnv("DATABASE_URL")
```

### System Information Snapshot

```go
info := sysx.SystemInfo()
for k, v := range info {
    fmt.Printf("%-12s = %s\n", k, v)
}
// os           = linux
// arch         = amd64
// hostname     = myserver
// pid          = 12345
// go_version   = go1.24.0
// executable   = /usr/local/bin/myapp
// num_cpu      = 8
```

### Running External Commands

```go
// Fire and forget
if err := sysx.ExecCommand("go", "generate", "./..."); err != nil {
    log.Fatal(err)
}

// Capture combined output
out, err := sysx.ExecOutput("git", "log", "--oneline", "-5")
if err != nil {
    log.Fatal(err)
}
fmt.Print(out)

// With timeout
out, err = sysx.ExecOutputWithTimeout(10*time.Second, "curl", "-s", "http://localhost/health")

// In a specific directory
err = sysx.ExecCommandInDir("/opt/myapp", "make", "build")
```

### Process Management

```go
pid := sysx.PID()
fmt.Println("I am running, PID:", pid)

if sysx.ProcessExists(pid) {
    fmt.Println("confirmed: still alive")
}

// Graceful then forceful termination
if err := sysx.KillProcess(targetPID); err != nil {
    sysx.KillProcessForcefully(targetPID)
}
```

### File System Checks

```go
if !sysx.FileExists("/etc/myapp/config.yaml") {
    log.Fatal("configuration file not found")
}

if sysx.DirExists("/var/cache/myapp") {
    size, _ := sysx.FileSize("/var/cache/myapp/data.bin")
    fmt.Printf("cache size: %d bytes\n", size)
}

home := sysx.MustHomeDir()
fmt.Println("home:", home)
```

## Best Practices

1. **Prefer `Must*` only at startup** — use the error-returning variants in library code; reserve `Must*` for `main()` or `init()` where a missing resource is fatal.
2. **Check `ProcessExists` before signalling** — `os.FindProcess` on Unix always succeeds; pair it with `ProcessExists` to avoid signalling phantom PIDs.
3. **Use `ExecOutput` for diagnostics, not pipelines** — both stdout and stderr are merged; if you need them separated, use `os/exec` directly.
4. **Timeouts protect production code** — always use `ExecCommandWithTimeout` / `ExecOutputWithTimeout` when calling external programs in long-running services.
5. **`EnvMap` is a snapshot** — the returned map is not updated when the environment changes; call again if you need fresh values.

## Platform Caveats

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| `UID()` / `GID()` | numeric | numeric | always -1 |
| `PPID()` | supported | supported | always 0 |
| `ProcessExists` | signal(0) | signal(0) | FindProcess (unreliable for dead PIDs) |
| `KillProcess` | SIGTERM | SIGTERM | limited |
| `KillProcessForcefully` | SIGKILL | SIGKILL | limited |
| `IsExecutable` / `IsReadable` | mode bits | mode bits | approximation |
| `IsWritable` | open test | open test | open test |
| `OSVersion` | `/etc/os-release` | `sw_vers` | GOOS/GOARCH string |
| `IsPrivileged` | UID==0 | UID==0 | always false |

## When to Use `sysx` vs stdlib

| Task | Prefer |
|------|--------|
| Quick OS check | `sysx.IsLinux()` vs `runtime.GOOS == "linux"` — both work; sysx adds readability |
| Env var with fallback | `sysx.GetEnv` — stdlib requires a conditional around `os.LookupEnv` |
| Typed env vars | `sysx.GetEnvInt` / `sysx.GetEnvBool` — no stdlib equivalent |
| Running a command and capturing output | `sysx.ExecOutput` — stdlib requires wiring up `cmd.Stdout`, `cmd.Stderr` |
| Command with timeout | `sysx.ExecCommandWithTimeout` — saves boilerplate `context.WithTimeout` setup |
| Checking if a file exists | `sysx.FileExists` — stdlib `os.Stat` + error check |
| Checking if a process is alive | `sysx.ProcessExists` — stdlib requires `syscall.Kill(pid, 0)` and error type assertion |
| Complex I/O piping | `os/exec` directly — `sysx` does not expose stdin or fine-grained output splitting |
| Watching the file system | `fsnotify` or stdlib — outside `sysx` scope |

## Real-World Integration Examples

### Health-Check Endpoint

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    info := sysx.SystemInfo()
    json.NewEncoder(w).Encode(map[string]any{
        "status":     "ok",
        "hostname":   info["hostname"],
        "pid":        info["pid"],
        "go_version": info["go_version"],
        "num_cpu":    info["num_cpu"],
    })
}
```

### Startup Configuration Validation

```go
func mustLoadConfig() Config {
    return Config{
        DSN:     sysx.MustGetEnv("DATABASE_URL"),
        Port:    sysx.GetEnvInt("PORT", 8080),
        Debug:   sysx.GetEnvBool("DEBUG", false),
        Origins: sysx.GetEnvSlice("CORS_ORIGINS", ","),
    }
}
```

### Conditional Platform Logic

```go
func configureLogging() {
    if sysx.IsLinux() {
        // Use journald integration
    } else if sysx.IsDarwin() {
        // Use os_log
    } else {
        // Fallback to stderr
    }
}
```

### Run a Build and Capture Output

```go
out, err := sysx.ExecOutputWithTimeout(
    2*time.Minute,
    "go", "build", "-o", "/tmp/myapp", "./cmd/myapp",
)
if err != nil {
    log.Printf("build failed:\n%s\n%v", out, err)
    os.Exit(1)
}
```
