package crontask

import (
	"errors"
	"fmt"
)

// Sentinel errors returned by the crontask package. Callers can test for
// these values with errors.Is when they need to handle a specific condition.
var (
	// ErrInvalidExpression is returned by Parse and Validate when the provided
	// cron expression is syntactically or semantically invalid.
	ErrInvalidExpression = errors.New("crontask: invalid cron expression")

	// ErrJobNotFound is returned by Remove and similar methods when the given
	// job ID does not exist in the registry.
	ErrJobNotFound = errors.New("crontask: job not found")

	// ErrSchedulerRunning is returned by Start when the scheduler is already
	// in a running state.
	ErrSchedulerRunning = errors.New("crontask: scheduler is already running")

	// ErrSchedulerStopped is returned by Register and similar methods when the
	// caller attempts to mutate a scheduler that has been permanently shut down.
	ErrSchedulerStopped = errors.New("crontask: scheduler has been stopped")

	// ErrJobTimeout is wrapped into the error returned by the executor when a
	// job's execution deadline is exceeded.
	ErrJobTimeout = errors.New("crontask: job execution timed out")

	// ErrMaxRetriesExceeded is returned when a job exhausts its configured
	// retry budget without succeeding.
	ErrMaxRetriesExceeded = errors.New("crontask: maximum retries exceeded")
)

// ExpressionError describes a parse or validation error for a specific cron
// expression. It implements the error interface and can be unwrapped to
// ErrInvalidExpression for sentinel matching.
type ExpressionError struct {
	// Expression is the raw string that triggered the error.
	Expression string

	// Field is the zero-based index of the field within the expression that
	// caused the error, or -1 when the error is not field-specific.
	Field int

	// Reason is a human-readable description of why the expression is invalid.
	Reason string
}

// Error implements the error interface.
func (e *ExpressionError) Error() string {
	if e.Field >= 0 {
		return fmt.Sprintf("crontask: invalid expression %q (field %d): %s", e.Expression, e.Field, e.Reason)
	}
	return fmt.Sprintf("crontask: invalid expression %q: %s", e.Expression, e.Reason)
}

// Unwrap returns ErrInvalidExpression so that callers can test for the
// sentinel value with errors.Is.
func (e *ExpressionError) Unwrap() error {
	return ErrInvalidExpression
}

// newExpressionError constructs an ExpressionError for the given expression
// and field index with a formatted reason string.
func newExpressionError(expr string, field int, format string, args ...any) error {
	return &ExpressionError{
		Expression: expr,
		Field:      field,
		Reason:     fmt.Sprintf(format, args...),
	}
}

// JobError wraps a job execution error with the job ID and attempt number so
// that hook implementations and callers can correlate errors with their origin.
type JobError struct {
	// JobID is the identifier of the job that failed.
	JobID string

	// Attempt is the one-based attempt number (1 = first try, 2 = first retry,
	// etc.).
	Attempt int

	// Err is the underlying error returned by the job function.
	Err error
}

// Error implements the error interface.
func (e *JobError) Error() string {
	return fmt.Sprintf("crontask: job %q attempt %d failed: %v", e.JobID, e.Attempt, e.Err)
}

// Unwrap returns the underlying error.
func (e *JobError) Unwrap() error {
	return e.Err
}
