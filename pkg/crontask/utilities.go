package crontask

import "time"

// isDue is the internal implementation shared by the package-level IsDue and
// Expression.IsDue. It returns true when the first activation of sched after
// (at − 1 second) is at or before at.
func isDue(sched Schedule, at time.Time) bool {
	at = at.Truncate(time.Second)
	next := sched.Next(at.Add(-time.Second))
	return !next.IsZero() && !next.After(at)
}
