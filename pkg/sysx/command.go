package sysx

import (
"context"
"errors"
"fmt"
"io"
"os"
"os/exec"
"time"
)

// ///////////////////////////
// Section: Command builder methods
// ///////////////////////////

// NewCommand creates a new Command for the program identified by name.
//
// name may be an absolute path, a relative path, or a bare program name that
// is resolved through the process $PATH. No validation is performed until
// Execute, Run, or Output is called.
//
// Parameters:
//   - `name`: the program name or path to execute.
//
// Returns:
//
//A pointer to a new Command ready for configuration.
//
// Example:
//
//result := sysx.NewCommand("git").
//    WithArgs("rev-parse", "HEAD").
//    WithDir("/path/to/repo").
//    Execute()
//if result.Success() {
//    fmt.Println(strings.TrimSpace(result.Stdout()))
//}
func NewCommand(name string) *Command {
return &Command{name: name}
}

// WithArgs sets the positional arguments passed to the program.
// Calling WithArgs replaces any previously set arguments.
//
// Parameters:
//   - `args`: the command-line arguments.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithArgs(args ...string) *Command {
c.args = args
return c
}

// WithDir sets the working directory for the command.
// An empty string leaves the working directory unchanged
// (the child inherits the calling process's directory).
//
// Parameters:
//   - `dir`: the working directory path.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithDir(dir string) *Command {
c.dir = dir
return c
}

// WithEnv appends one or more environment variable bindings in "KEY=VALUE"
// form to the command's environment. Bindings are merged on top of the
// calling process environment; later values for the same key shadow earlier ones.
// Calling WithEnv multiple times accumulates bindings.
//
// Parameters:
//   - `env`: one or more "KEY=VALUE" environment variable bindings.
//
// Returns:
//
//The receiver, enabling method chaining.
//
// Example:
//
//sysx.NewCommand("go").
//    WithArgs("build", "./...").
//    WithEnv("GOOS=linux", "GOARCH=amd64").
//    Execute()
func (c *Command) WithEnv(env ...string) *Command {
c.env = append(c.env, env...)
return c
}

// WithTimeout sets a maximum duration for the command. If the command does not
// finish within the deadline the process is killed and Execute returns a
// CommandResult whose Err() wraps context.DeadlineExceeded.
//
// A zero or negative duration means no timeout is applied.
//
// Parameters:
//   - `d`: the maximum duration to wait.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithTimeout(d time.Duration) *Command {
c.timeout = d
return c
}

// WithContext attaches an existing context to the command, enabling external
// cancellation and deadline propagation. If both WithContext and WithTimeout
// are used, the more restrictive deadline wins.
//
// Parameters:
//   - `ctx`: the context to attach.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithContext(ctx context.Context) *Command {
c.ctx = ctx
return c
}

// WithStdin sets the reader that supplies the command's standard input.
//
// Parameters:
//   - `r`: the io.Reader to use as stdin.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithStdin(r io.Reader) *Command {
c.stdin = r
return c
}

// WithStdout sets the writer to which the command's standard output is
// forwarded in real time. When a custom writer is provided,
// CommandResult.Stdout() will be empty because data is not buffered internally.
//
// Parameters:
//   - `w`: the io.Writer to receive stdout.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithStdout(w io.Writer) *Command {
c.stdout = w
return c
}

// WithStderr sets the writer to which the command's standard error is
// forwarded in real time. When a custom writer is provided,
// CommandResult.Stderr() will be empty because data is not buffered internally.
//
// Parameters:
//   - `w`: the io.Writer to receive stderr.
//
// Returns:
//
//The receiver, enabling method chaining.
func (c *Command) WithStderr(w io.Writer) *Command {
c.stderr = w
return c
}

// Execute runs the command and returns a CommandResult containing captured
// stdout, stderr, exit code, wall-clock duration, and any error.
//
// If WithStdout or WithStderr were provided, the corresponding fields in
// CommandResult are empty because the data was streamed to those writers.
//
// Execute is safe to call multiple times; each call spawns a new process.
//
// Returns:
//
//A non-nil *CommandResult describing the outcome of the command.
//
// Example:
//
//res := sysx.NewCommand("bash").
//    WithArgs("-c", "echo hello").
//    WithTimeout(5 * time.Second).
//    WithEnv("APP_ENV=prod").
//    WithDir("/tmp").
//    Execute()
//fmt.Printf("exit=%d stdout=%q duration=%v\n", res.ExitCode(), res.Stdout(), res.Duration())
func (c *Command) Execute() *CommandResult {
if c.name == "" {
return &CommandResult{
err:      errors.New("sysx: command name must not be empty"),
exitCode: -1,
}
}
cmd, cancel := c.buildCmd()
defer cancel()

var outBuf, errBuf commandBuffer
if c.stdout != nil {
cmd.Stdout = c.stdout
} else {
cmd.Stdout = &outBuf
}
if c.stderr != nil {
cmd.Stderr = c.stderr
} else {
cmd.Stderr = &errBuf
}

start := time.Now()
runErr := cmd.Run()
dur := time.Since(start)

res := &CommandResult{duration: dur}
if c.stdout == nil {
res.stdout = outBuf.String()
}
if c.stderr == nil {
res.stderr = errBuf.String()
}
if runErr != nil {
res.err = runErr
var exitErr *exec.ExitError
if errors.As(runErr, &exitErr) {
res.exitCode = exitErr.ExitCode()
} else {
res.exitCode = -1
}
}
return res
}

// Run runs the command, discarding all output, and returns only the error.
// It is a convenience method equivalent to calling Execute and inspecting
// the Err() result.
//
// Returns:
//
//An error if the command could not be started or exited non-zero; nil on success.
func (c *Command) Run() error {
return c.Execute().Err()
}

// Output runs the command and returns the combined stdout+stderr as a string
// along with any error.
//
// Returns:
//
//(string, error): the combined output and nil on success, or combined
//partial output and a non-nil error on failure.
func (c *Command) Output() (string, error) {
r := c.Execute()
return r.Combined(), r.Err()
}

// buildCmd constructs the underlying *exec.Cmd from the Command fields.
// The returned CancelFunc must always be called to release context resources.
func (c *Command) buildCmd() (*exec.Cmd, context.CancelFunc) {
base := c.ctx
if base == nil {
base = context.Background()
}
var cancel context.CancelFunc
if c.timeout > 0 {
base, cancel = context.WithTimeout(base, c.timeout)
} else {
cancel = func() {}
}
cmd := exec.CommandContext(base, c.name, c.args...)
if c.dir != "" {
cmd.Dir = c.dir
}
if len(c.env) > 0 {
cmd.Env = append(os.Environ(), c.env...)
}
if c.stdin != nil {
cmd.Stdin = c.stdin
}
return cmd, cancel
}

// ///////////////////////////
// Section: Top-level convenience functions
// ///////////////////////////

// RunCommand creates a Command for name, runs it with the provided args, and
// returns the structured CommandResult.
//
// It is the single-call equivalent of:
//
//NewCommand(name).WithArgs(args...).Execute()
//
// Parameters:
//   - `name`: the program name or path.
//   - `args`: optional arguments.
//
// Returns:
//
//A non-nil *CommandResult.
//
// Example:
//
//res := sysx.RunCommand("git", "status")
//if !res.Success() {
//    log.Printf("exit %d: %s", res.ExitCode(), res.Stderr())
//}
func RunCommand(name string, args ...string) *CommandResult {
return NewCommand(name).WithArgs(args...).Execute()
}

// ExecCommand runs the named program with the provided arguments and waits
// for it to complete. Both stdout and stderr are discarded.
//
// The function does not use shell interpolation; name must be a program name
// or absolute path. An error is returned if name is empty, if the program
// cannot be found, or if the program exits with a non-zero status.
//
// Parameters:
//   - `name`: the program name or path to execute.
//   - `args`: optional arguments to pass to the program.
//
// Returns:
//
//An error if the command could not be started or exited non-zero, or nil on success.
//
// Example:
//
//if err := sysx.ExecCommand("go", "build", "./..."); err != nil {
//    log.Fatal(err)
//}
func ExecCommand(name string, args ...string) error {
if name == "" {
return errors.New("sysx: command name must not be empty")
}
return NewCommand(name).WithArgs(args...).Run()
}

// ExecCommandContext runs the named program under the provided context.
// The command is cancelled when ctx is cancelled or its deadline expires.
//
// Parameters:
//   - `ctx`:  the context controlling cancellation and deadline.
//   - `name`: the program name or path.
//   - `args`: optional arguments.
//
// Returns:
//
//An error if the command failed, was cancelled, or timed out; nil on success.
//
// Example:
//
//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//defer cancel()
//if err := sysx.ExecCommandContext(ctx, "go", "test", "./..."); err != nil {
//    log.Fatal(err)
//}
func ExecCommandContext(ctx context.Context, name string, args ...string) error {
if name == "" {
return errors.New("sysx: command name must not be empty")
}
return NewCommand(name).WithArgs(args...).WithContext(ctx).Run()
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
//(string, error): the combined output and nil on success, or the partial
//output and a non-nil error on failure.
//
// Example:
//
//out, err := sysx.ExecOutput("git", "rev-parse", "HEAD")
//if err != nil {
//    log.Fatal(err)
//}
//fmt.Println(strings.TrimSpace(out))
func ExecOutput(name string, args ...string) (string, error) {
if name == "" {
return "", errors.New("sysx: command name must not be empty")
}
r := NewCommand(name).WithArgs(args...).Execute()
return r.Combined(), r.Err()
}

// ExecOutputLines runs the named program and returns its stdout split into
// individual lines. Line endings are stripped; empty lines are preserved.
// Stderr is captured but discarded on success. On failure, the error from the
// command is returned.
//
// Parameters:
//   - `name`: the program name or path.
//   - `args`: optional arguments.
//
// Returns:
//
//([]string, error): lines of stdout and nil on success, or nil and a
//non-nil error on failure.
//
// Example:
//
//lines, err := sysx.ExecOutputLines("git", "log", "--oneline", "-5")
//for _, l := range lines {
//    fmt.Println(l)
//}
func ExecOutputLines(name string, args ...string) ([]string, error) {
if name == "" {
return nil, errors.New("sysx: command name must not be empty")
}
r := NewCommand(name).WithArgs(args...).Execute()
if r.Err() != nil {
return nil, r.Err()
}
return splitLines(r.Stdout()), nil
}

// ExecStreaming runs the named program, forwarding its stdout and stderr to
// the provided writers in real time as the command executes. Either writer
// may be nil to discard the corresponding stream.
//
// Parameters:
//   - `stdout`: writer receiving standard output; nil to discard.
//   - `stderr`: writer receiving standard error; nil to discard.
//   - `name`:   the program name or path.
//   - `args`:   optional arguments.
//
// Returns:
//
//An error if the command failed or could not be started; nil on success.
//
// Example:
//
//err := sysx.ExecStreaming(os.Stdout, os.Stderr, "go", "build", "./...")
func ExecStreaming(stdout, stderr io.Writer, name string, args ...string) error {
if name == "" {
return errors.New("sysx: command name must not be empty")
}
c := NewCommand(name).WithArgs(args...)
if stdout != nil {
c = c.WithStdout(stdout)
}
if stderr != nil {
c = c.WithStderr(stderr)
}
return c.Run()
}

// ExecAsync starts the named program asynchronously and returns the underlying
// *exec.Cmd without waiting for it to finish. The caller is responsible for
// calling cmd.Wait() to release associated resources and obtain the exit status.
//
// Parameters:
//   - `name`: the program name or path.
//   - `args`: optional arguments.
//
// Returns:
//
//(*exec.Cmd, error): the started command handle and nil on success, or nil
//and a non-nil error if the command could not be started.
//
// Example:
//
//cmd, err := sysx.ExecAsync("long-running-server", "--port", "8080")
//if err != nil {
//    log.Fatal(err)
//}
//// ... do other work ...
//cmd.Wait()
func ExecAsync(name string, args ...string) (*exec.Cmd, error) {
if name == "" {
return nil, errors.New("sysx: command name must not be empty")
}
cmd := exec.Command(name, args...)
if err := cmd.Start(); err != nil {
return nil, err
}
return cmd, nil
}

// ExecPipeline executes a sequence of commands as a shell pipeline: the
// standard output of each command is connected to the standard input of the
// next. The stdout of the final command is returned as a string.
//
// Each element of commands is a []string where the first element is the
// program name and the remaining elements are its arguments. An error is
// returned if any command in the pipeline fails to start or exits non-zero.
//
// Parameters:
//   - `commands`: one or more [name, arg...] command descriptors.
//
// Returns:
//
//(string, error): the stdout of the last command and nil on success, or
//partial output and a non-nil error on failure.
//
// Example:
//
//out, err := sysx.ExecPipeline(
//    []string{"cat", "/etc/passwd"},
//    []string{"grep", "root"},
//    []string{"cut", "-d:", "-f1"},
//)
func ExecPipeline(commands ...[]string) (string, error) {
if len(commands) == 0 {
return "", errors.New("sysx: pipeline requires at least one command")
}
cmds := make([]*exec.Cmd, len(commands))
for i, args := range commands {
if len(args) == 0 {
return "", fmt.Errorf("sysx: pipeline command at index %d has no name", i)
}
cmds[i] = exec.Command(args[0], args[1:]...)
}
// Wire stdout of each command to stdin of the next.
for i := 0; i < len(cmds)-1; i++ {
pipe, err := cmds[i].StdoutPipe()
if err != nil {
return "", fmt.Errorf("sysx: pipeline pipe error at index %d: %w", i, err)
}
cmds[i+1].Stdin = pipe
}
var outBuf commandBuffer
cmds[len(cmds)-1].Stdout = &outBuf

// Start all commands in order.
for i, cmd := range cmds {
if err := cmd.Start(); err != nil {
return "", fmt.Errorf("sysx: pipeline start error at index %d: %w", i, err)
}
}
// Wait for all commands in order so pipes drain correctly.
for i, cmd := range cmds {
if err := cmd.Wait(); err != nil {
return outBuf.String(), fmt.Errorf("sysx: pipeline wait error at index %d: %w", i, err)
}
}
return outBuf.String(), nil
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
//An error if the command timed out, could not be started, or exited non-zero, or nil on success.
//
// Example:
//
//err := sysx.ExecCommandWithTimeout(5*time.Second, "ping", "-c", "1", "localhost")
func ExecCommandWithTimeout(timeout time.Duration, name string, args ...string) error {
if name == "" {
return errors.New("sysx: command name must not be empty")
}
return NewCommand(name).WithArgs(args...).WithTimeout(timeout).Run()
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
//(string, error): the combined output and nil on success, or partial output
//and a non-nil error on failure or timeout.
//
// Example:
//
//out, err := sysx.ExecOutputWithTimeout(3*time.Second, "curl", "-s", "http://localhost")
func ExecOutputWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
if name == "" {
return "", errors.New("sysx: command name must not be empty")
}
r := NewCommand(name).WithArgs(args...).WithTimeout(timeout).Execute()
return r.Combined(), r.Err()
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
//An error if the command could not be started or exited non-zero, or nil on success.
//
// Example:
//
//err := sysx.ExecCommandInDir("/tmp/myproject", "go", "test", "./...")
func ExecCommandInDir(dir, name string, args ...string) error {
if name == "" {
return errors.New("sysx: command name must not be empty")
}
return NewCommand(name).WithArgs(args...).WithDir(dir).Run()
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
//(string, error): the combined output and nil on success, or the partial
//output and a non-nil error on failure.
//
// Example:
//
//out, err := sysx.ExecOutputInDir("/tmp/myproject", "git", "status")
func ExecOutputInDir(dir, name string, args ...string) (string, error) {
if name == "" {
return "", errors.New("sysx: command name must not be empty")
}
r := NewCommand(name).WithArgs(args...).WithDir(dir).Execute()
return r.Combined(), r.Err()
}
