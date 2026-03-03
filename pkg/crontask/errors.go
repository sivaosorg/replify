package crontask

import (
	"fmt"
)

// Error implements the error interface.
//
// Example:
//
//	for _, run := range runs {
//		if errors.Is(run.Error, crontask.ErrJobTimeout) {
//			// Handle timeout specifically
//		}
//	}
func (e *ExpressionError) Error() string {
	if e.Field >= 0 {
		return fmt.Sprintf("crontask: invalid expression %q (field %d): %s", e.Expression, e.Field, e.Reason)
	}
	return fmt.Sprintf("crontask: invalid expression %q: %s", e.Expression, e.Reason)
}

// Unwrap returns ErrInvalidExpression so that callers can test for the
// sentinel value with errors.Is.
//
// Example:
//
//	for _, run := range runs {
//		if errors.Is(run.Error, crontask.ErrJobTimeout) {
//			// Handle timeout specifically
//		}
//	}
func (e *ExpressionError) Unwrap() error {
	return ErrInvalidExpression
}

// Error implements the error interface.
//
// Example:
//
//	for _, run := range runs {
//		if errors.Is(run.Error, crontask.ErrJobTimeout) {
//			// Handle timeout specifically
//		}
//	}
func (e *JobError) Error() string {
	return fmt.Sprintf("crontask: job %q attempt %d failed: %v", e.JobID, e.Attempt, e.Err)
}

// Unwrap returns the underlying error.
//
// Example:
//
//	for _, run := range runs {
//		if errors.Is(run.Error, crontask.ErrJobTimeout) {
//			// Handle timeout specifically
//		}
//	}
func (e *JobError) Unwrap() error {
	return e.Err
}
