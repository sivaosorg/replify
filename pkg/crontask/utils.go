package crontask

import (
	"fmt"
	"time"
)

// Expression is a parsed cron expression that exposes schedule utilities
// without requiring a running Scheduler instance. Obtain one via MustParse.
//
// Expression is immutable and safe for concurrent use.
type Expression struct {
	raw      string
	schedule Schedule
}

// Raw returns the original expression string as passed to MustParse.
func (e Expression) Raw() string { return e.raw }

// Next returns the first activation time strictly after from. It returns the
// zero time when no future activation exists (e.g. the schedule is exhausted).
func (e Expression) Next(from time.Time) time.Time {
	return e.schedule.Next(from)
}

// NextN returns the next n activation times starting from from. If fewer than
// n future activations exist, the slice is shorter than n.
func (e Expression) NextN(from time.Time, n int) []time.Time {
	if n <= 0 {
		return nil
	}
	out := make([]time.Time, 0, n)
	cur := from
	for len(out) < n {
		next := e.schedule.Next(cur)
		if next.IsZero() {
			break
		}
		out = append(out, next)
		cur = next
	}
	return out
}

// IsDue reports whether the expression is due at the given time using
// second-level granularity. See the package-level IsDue for the exact
// semantics.
func (e Expression) IsDue(at time.Time) bool {
	return isDue(e.schedule, at)
}

// MustParse is like Parse but panics instead of returning an error when the
// expression is invalid. It is intended for use in package-level variable
// initializers where the expression is a compile-time constant.
//
// Example:
//
//	var reportSchedule = crontask.MustParse("0 9 * * 1-5")
func MustParse(expr string) Expression {
	s, err := Parse(expr)
	if err != nil {
		panic(fmt.Sprintf("crontask.MustParse: %v", err))
	}
	return Expression{raw: expr, schedule: s}
}

// IsValidCronExpr reports whether expr is a syntactically and semantically
// valid cron expression recognised by this package. It is equivalent to
// Validate(expr) == nil.
//
// Thread-safe; does not require a Scheduler instance.
func IsValidCronExpr(expr string) bool {
	return Validate(expr) == nil
}

// ValidateCronExpr validates expr and returns a descriptive error when it is
// invalid. The returned error, if non-nil, is an *ExpressionError that also
// satisfies errors.Is(err, ErrInvalidExpression).
//
// ValidateCronExpr is a convenience wrapper around Validate and is provided
// for callers who prefer the longer, self-documenting name.
//
// Thread-safe; does not require a Scheduler instance.
func ValidateCronExpr(expr string) error {
	return Validate(expr)
}

// IsDue reports whether the cron expression expr was due at the given time at.
//
// The check uses second-level granularity: IsDue returns true when the first
// activation of the schedule after (at − 1 second) is at or before at. For
// standard five-field expressions this means the function returns true exactly
// at the start of a matching minute; for six-field or @every expressions it
// returns true at the matching second.
//
// IsDue returns false for invalid expressions without panicking.
//
// Thread-safe; does not require a Scheduler instance.
//
// Example:
//
//	now := time.Now().Truncate(time.Minute)
//	if crontask.IsDue("0 9 * * 1-5", now) {
//	    sendDailyReport()
//	}
func IsDue(expr string, at time.Time) bool {
	sched, err := Parse(expr)
	if err != nil {
		return false
	}
	return isDue(sched, at)
}

// NextRun returns the first activation time for expr after the reference time
// from. It is a convenience wrapper around Parse and Schedule.Next that avoids
// the need to create a Scheduler or call Parse manually.
//
// Thread-safe; does not require a Scheduler instance.
//
// Example:
//
//	next, err := crontask.NextRun("0 9 * * 1-5", time.Now())
func NextRun(expr string, from time.Time) (time.Time, error) {
	sched, err := Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return sched.Next(from), nil
}

// NextRuns returns the next n activation times for expr starting from from.
// It is a convenience wrapper around Parse and Schedule.Next.
//
// When n ≤ 0, NextRuns returns a nil slice without an error. When the schedule
// has fewer than n future activations, the returned slice is shorter than n.
//
// Thread-safe; does not require a Scheduler instance.
//
// Example:
//
//	runs, err := crontask.NextRuns("0 9 * * 1-5", time.Now(), 5)
func NextRuns(expr string, from time.Time, n int) ([]time.Time, error) {
	sched, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, nil
	}
	out := make([]time.Time, 0, n)
	cur := from
	for len(out) < n {
		next := sched.Next(cur)
		if next.IsZero() {
			break
		}
		out = append(out, next)
		cur = next
	}
	return out, nil
}

// isDue is the internal implementation shared by the package-level IsDue and
// Expression.IsDue. It returns true when the first activation of sched after
// (at − 1 second) is at or before at.
func isDue(sched Schedule, at time.Time) bool {
	at = at.Truncate(time.Second)
	next := sched.Next(at.Add(-time.Second))
	return !next.IsZero() && !next.After(at)
}
