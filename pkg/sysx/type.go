package sysx

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ///////////////////////////
// Section: Command types
// ///////////////////////////

// Command holds the fully resolved configuration for a single external process.
// Use NewCommand to create a Command, then chain With* methods to configure it
// before calling Execute, Run, or Output.
//
// Command is not safe for concurrent mutation; build a new value per goroutine
// when the same logical command must be launched from multiple goroutines.
type Command struct {
	name    string
	args    []string
	dir     string
	env     []string
	timeout time.Duration // zero means no timeout
	ctx     context.Context
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

// Name returns the program name or path configured on the Command.
//
// Returns:
//
//	A string containing the program name or path.
//
// Example:
//
//	cmd := sysx.NewCommand("git")
//	fmt.Println(cmd.Name()) // "git"
func (c *Command) Name() string { return c.name }

// Args returns the positional arguments configured on the Command.
//
// Returns:
//
//	A slice of strings containing the command-line arguments.
//
// Example:
//
//	cmd := sysx.NewCommand("git").WithArgs("log", "--oneline")
//	fmt.Println(cmd.Args()) // ["log", "--oneline"]
func (c *Command) Args() []string { return c.args }

// Dir returns the working directory configured on the Command.
// An empty string means the child process inherits the caller's directory.
//
// Returns:
//
//	A string containing the working directory path, or empty if not set.
func (c *Command) Dir() string { return c.dir }

// Env returns the extra environment variable bindings ("KEY=VALUE") that
// will be merged on top of the calling process environment when the command
// is executed.
//
// Returns:
//
//	A slice of "KEY=VALUE" strings; nil if no extra bindings were added.
func (c *Command) Env() []string { return c.env }

// Timeout returns the maximum execution duration configured on the Command.
// A zero duration means no timeout is applied.
//
// Returns:
//
//	A time.Duration; zero if no timeout was set.
func (c *Command) Timeout() time.Duration { return c.timeout }

// ///////////////////////////
// Section: CommandResult type
// ///////////////////////////

// CommandResult holds the structured outcome of a completed command execution.
// All fields are always populated; zero values indicate the field was not
// applicable (e.g. ExitCode() == 0 means success, Duration() == 0 on immediate failure).
//
// Use the accessor methods to read individual result fields.
type CommandResult struct {
	stdout   string
	stderr   string
	exitCode int
	duration time.Duration
	err      error
}

// Stdout returns the captured standard output of the command.
// Returns an empty string when a custom io.Writer was provided via WithStdout.
//
// Returns:
//
//	A string containing the captured stdout of the command.
func (r *CommandResult) Stdout() string { return r.stdout }

// Stderr returns the captured standard error of the command.
// Returns an empty string when a custom io.Writer was provided via WithStderr.
//
// Returns:
//
//	A string containing the captured stderr of the command.
func (r *CommandResult) Stderr() string { return r.stderr }

// ExitCode returns the process exit code; 0 indicates success.
// -1 indicates that the exit code could not be determined
// (e.g. the process was killed by a signal or a context was cancelled).
//
// Returns:
//
//	An int representing the process exit code.
func (r *CommandResult) ExitCode() int { return r.exitCode }

// Duration returns the wall-clock time spent waiting for the command to complete.
//
// Returns:
//
//	A time.Duration representing the execution time.
func (r *CommandResult) Duration() time.Duration { return r.duration }

// Err returns the error from command execution.
// A non-nil value indicates the command could not be started or exited non-zero.
//
// Returns:
//
//	An error describing the failure, or nil on success.
func (r *CommandResult) Err() error { return r.err }

// Success reports whether the command completed without error.
//
// Returns:
//
//	true when Err() is nil; false otherwise.
func (r *CommandResult) Success() bool { return r.err == nil }

// Combined returns the concatenation of Stdout followed by Stderr.
//
// Returns:
//
//	A string containing the combined output of the command.
func (r *CommandResult) Combined() string { return r.stdout + r.stderr }

// ///////////////////////////
// Section: SafeFileWriter type
// ///////////////////////////

// SafeFileWriter provides concurrency-safe append and overwrite operations
// targeting a single file path. A single SafeFileWriter instance can be shared
// across goroutines; all write operations are serialised by an internal mutex.
//
// Create a SafeFileWriter with NewSafeFileWriter and optionally adjust the
// file permission with WithPerm before sharing across goroutines.
type SafeFileWriter struct {
	mu   sync.Mutex
	path string
	perm os.FileMode
}

// Path returns the file path targeted by this SafeFileWriter.
//
// Returns:
//
//	A string containing the file path.
func (w *SafeFileWriter) Path() string { return w.path }

// Perm returns the file permission used when creating the file.
//
// Returns:
//
//	An os.FileMode representing the file permission.
func (w *SafeFileWriter) Perm() os.FileMode { return w.perm }

// ///////////////////////////
// Section: Internal I/O helper
// ///////////////////////////

// commandBuffer is a minimal zero-allocation-friendly byte accumulator used
// internally by Execute to capture stdout and stderr. It satisfies io.Writer.
type commandBuffer struct {
	buf strings.Builder
}

// Write appends p to the buffer.
func (b *commandBuffer) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

// String returns the accumulated content as a string.
func (b *commandBuffer) String() string {
	return b.buf.String()
}

// ///////////////////////////
// Section: Package-level global variables
// ///////////////////////////

// fileMutexes is the package-level registry of per-path in-process mutexes used
// by WriteFileLocked to serialise concurrent writes to the same file path.
// It is populated lazily by getFileMutex and is safe for concurrent access.
var fileMutexes sync.Map // map[string]*sync.Mutex
