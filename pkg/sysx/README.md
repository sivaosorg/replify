# sysx

**sysx** is a lightweight, production-grade system utilities toolkit for Go, providing a clean and consistent API for OS detection, runtime introspection, environment variable management, process control, command execution, and file system queries — all built exclusively on the Go standard library.

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
