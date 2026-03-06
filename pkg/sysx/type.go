package sysx

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

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

// commandBuffer is a minimal zero-allocation-friendly byte accumulator used
// internally by Execute to capture stdout and stderr. It satisfies io.Writer.
type commandBuffer struct {
	buf strings.Builder
}

// fileMutexes is the package-level registry of per-path in-process mutexes used
// by WriteFileLocked to serialise concurrent writes to the same file path.
// It is populated lazily by getFileMutex and is safe for concurrent access.
var fileMutexes sync.Map // map[string]*sync.Mutex
