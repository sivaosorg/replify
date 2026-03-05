package sysx

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

// ///////////////////////////
// Section: Basic command execution
// ///////////////////////////

// ExecCommand runs the named program with the provided arguments and waits
// for it to complete.
//
// Both stdout and stderr are discarded. The function does not use shell
// interpolation; name must be a program name or absolute path. An error is
// returned if name is empty, if the program cannot be found, or if the
// program exits with a non-zero status.
//
// Parameters:
//   - `name`: the program name or path to execute.
//   - `args`: optional arguments to pass to the program.
//
// Returns:
//
//	An error if the command could not be started or exited non-zero, or nil on success.
//
// Example:
//
//	if err := sysx.ExecCommand("go", "build", "./..."); err != nil {
//	    log.Fatal(err)
//	}
func ExecCommand(name string, args ...string) error {
	if name == "" {
		return errors.New("sysx: command name must not be empty")
	}
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// ExecOutput runs the named program with the provided arguments, waits for it
// to complete, and returns the combined stdout and stderr output as a string.
//
// The function does not use shell interpolation. An error is returned if name
// is empty, if the program cannot be found, or if the program exits with a
// non-zero status.
//
// Parameters:
//   - `name`: the program name or path to execute.
//   - `args`: optional arguments to pass to the program.
//
// Returns:
//
//	(string, error): the combined output and nil on success, or the partial
//	output and a non-nil error on failure.
//
// Example:
//
//	out, err := sysx.ExecOutput("git", "rev-parse", "HEAD")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(strings.TrimSpace(out))
func ExecOutput(name string, args ...string) (string, error) {
	if name == "" {
		return "", errors.New("sysx: command name must not be empty")
	}
	cmd := exec.Command(name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// ///////////////////////////
// Section: Timeout variants
// ///////////////////////////

// ExecCommandWithTimeout runs the named program with the provided arguments
// and a deadline of timeout. If the command does not finish within the
// deadline the process is killed and a context deadline-exceeded error is
// returned.
//
// Parameters:
//   - `timeout`: maximum duration to wait for the command to complete.
//   - `name`:    the program name or path to execute.
//   - `args`:    optional arguments to pass to the program.
//
// Returns:
//
//	An error if the command timed out, could not be started, or exited non-zero, or nil on success.
//
// Example:
//
//	err := sysx.ExecCommandWithTimeout(5*time.Second, "ping", "-c", "1", "localhost")
func ExecCommandWithTimeout(timeout time.Duration, name string, args ...string) error {
	if name == "" {
		return errors.New("sysx: command name must not be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Run()
}

// ExecOutputWithTimeout runs the named program with the provided arguments
// and a deadline of timeout, returning the combined stdout and stderr.
//
// If the command does not finish within the deadline the process is killed and
// a context deadline-exceeded error is returned along with any output produced
// before the timeout.
//
// Parameters:
//   - `timeout`: maximum duration to wait for the command to complete.
//   - `name`:    the program name or path to execute.
//   - `args`:    optional arguments to pass to the program.
//
// Returns:
//
//	(string, error): the combined output and nil on success, or partial output
//	and a non-nil error on failure or timeout.
//
// Example:
//
//	out, err := sysx.ExecOutputWithTimeout(3*time.Second, "curl", "-s", "http://localhost")
func ExecOutputWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
	if name == "" {
		return "", errors.New("sysx: command name must not be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// ///////////////////////////
// Section: Directory-scoped execution
// ///////////////////////////

// ExecCommandInDir runs the named program with the provided arguments from
// the specified working directory and waits for it to complete.
//
// stdout and stderr are discarded. The function does not use shell
// interpolation.
//
// Parameters:
//   - `dir`:  the working directory in which to run the command.
//   - `name`: the program name or path to execute.
//   - `args`: optional arguments to pass to the program.
//
// Returns:
//
//	An error if the command could not be started or exited non-zero, or nil on success.
//
// Example:
//
//	err := sysx.ExecCommandInDir("/tmp/myproject", "go", "test", "./...")
func ExecCommandInDir(dir, name string, args ...string) error {
	if name == "" {
		return errors.New("sysx: command name must not be empty")
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}

// ExecOutputInDir runs the named program with the provided arguments from
// the specified working directory, waits for it to complete, and returns the
// combined stdout and stderr output.
//
// Parameters:
//   - `dir`:  the working directory in which to run the command.
//   - `name`: the program name or path to execute.
//   - `args`: optional arguments to pass to the program.
//
// Returns:
//
//	(string, error): the combined output and nil on success, or the partial
//	output and a non-nil error on failure.
//
// Example:
//
//	out, err := sysx.ExecOutputInDir("/tmp/myproject", "git", "status")
func ExecOutputInDir(dir, name string, args ...string) (string, error) {
	if name == "" {
		return "", errors.New("sysx: command name must not be empty")
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}
