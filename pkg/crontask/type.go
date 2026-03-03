package crontask

import (
	"context"
	"time"
)

// Alias is a short-hand name for a commonly used cron schedule, inspired by
// both Vixie-cron and gronx. Aliases are resolved by the parser before
// field-level parsing occurs.
//
// The following aliases are recognised out of the box:
//
//	@yearly   (or @annually)  — "0 0 1 1 *"
//	@monthly                  — "0 0 1 * *"
//	@weekly                   — "0 0 * * 0"
//	@daily    (or @midnight)  — "0 0 * * *"
//	@hourly                   — "0 * * * *"
//	@minutely                 — "* * * * *"
//	@weekdays                 — "0 0 * * 1-5"
//	@weekends                 — "0 0 * * 0,6"
//
// Business-oriented aliases:
//
//	@businessDaily   — "0 9 * * 1-5"       (09:00 on weekdays)
//	@businessHourly  — "0 9-17 * * 1-5"    (top of each hour, 09–17, weekdays)
//	@quarterly       — "0 0 1 1,4,7,10 *"  (midnight, first day of each quarter)
//	@semiMonthly     — "0 0 1,15 * *"      (midnight, 1st and 15th of each month)
//	@workhours       — "* 9-17 * * 1-5"    (every minute 09:00–17:59, weekdays)
//	@marketOpen      — "30 9 * * 1-5"      (09:30, weekdays)
//	@marketClose     — "0 16 * * 1-5"      (16:00, weekdays)
//
// Additional aliases can be registered at runtime with RegisterAlias.
//
// When using the six-field (seconds-first) format, the above expansions gain
// a leading "0" second field automatically.
type Alias = string

// Schedule is the interface implemented by any type that can compute the next
// activation time for a scheduled job. The single method Next receives the
// current (or reference) time and returns the earliest future time at which
// the job should run next.
//
// Implementing Schedule allows users to inject custom scheduling logic — for
// example, a schedule driven by an external calendar API — without forking the
// package.
type Schedule interface {
	// Next returns the next activation time after the given reference time t.
	// If no further activation exists (e.g. a once-only schedule in the past),
	// Next returns the zero time.Time.
	Next(t time.Time) time.Time
}

// Expression is a parsed cron expression that exposes schedule utilities
// without requiring a running Scheduler instance. Obtain one via MustParse.
//
// Expression is immutable and safe for concurrent use.
type Expression struct {
	raw      string
	schedule Schedule
}

// Hooks is the interface that callers may implement to observe the lifecycle
// of a job invocation. All methods are optional — embed NoopHooks to satisfy
// the interface without implementing every method.
//
// Hook methods are called synchronously within the executor goroutine. Hooks
// must not block for long periods; spawn a goroutine if you need to perform
// expensive work (e.g. remote metric writes) without slowing the executor.
type Hooks interface {
	// OnStart is called immediately before the job function is invoked,
	// after jitter has been applied and after the execution context has been
	// derived. The jobID parameter identifies the job being dispatched.
	OnStart(ctx context.Context, jobID string)

	// OnSuccess is called when the job function returns nil. The duration
	// parameter is the wall-clock time of the invocation, excluding jitter.
	OnSuccess(ctx context.Context, jobID string, duration time.Duration)

	// OnFailure is called when the job function returns a non-nil error after
	// all retry attempts are exhausted. err is the final error.
	OnFailure(ctx context.Context, jobID string, duration time.Duration, err error)

	// OnComplete is called after OnSuccess or OnFailure and regardless of the
	// outcome. It is useful for releasing resources that were acquired in
	// OnStart.
	OnComplete(ctx context.Context, jobID string, duration time.Duration)
}

// NoopHooks is a zero-value implementation of Hooks whose methods all do
// nothing. Embed NoopHooks into your own struct to selectively override only
// the methods you care about.
//
// Example:
//
//	type MyHooks struct {
//	    crontask.NoopHooks
//	}
//
//	func (h *MyHooks) OnFailure(_ context.Context, id string, _ time.Duration, err error) {
//	    log.Printf("ALERT: job %s failed: %v", id, err)
//	}
type NoopHooks struct{}

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

// fieldSpec describes the valid range for a single cron field,
// used by the parser to validate field values.
type fieldSpec struct {
	min  int
	max  int
	name string
}

// cronSchedule is the internal, parsed representation of a standard cron
// expression. It implements the Schedule interface.
type cronSchedule struct {
	second     []bool // [0..59] — only populated in six-field mode
	minute     []bool // [0..59]
	hour       []bool // [0..23]
	dayOfMonth []bool // [1..31]
	month      []bool // [1..12]
	dayOfWeek  []bool // [0..6] (0 = Sunday)
	loc        *time.Location
}

// intervalSchedule fires every fixed Duration starting from the first tick
// after the reference time. It implements Schedule.
type intervalSchedule struct {
	interval time.Duration
}

// executor is responsible for running a single job entry within its own
// goroutine. It enforces timeouts, applies jitter, runs retry loops, and
// dispatches hook calls.
type executor struct {
	onError func(id string, err error)
}
