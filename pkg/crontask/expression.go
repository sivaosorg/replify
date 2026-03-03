package crontask

import (
	"time"
)

// Raw returns the original expression string as passed to MustParse.
func (e Expression) Raw() string { return e.raw }

// Next returns the first activation time strictly after from. It returns the
// zero time when no future activation exists (e.g. the schedule is exhausted).
//
// Example:
//
//	next := expr.Next(time.Now())
func (e Expression) Next(from time.Time) time.Time {
	return e.schedule.Next(from)
}

// NextN returns the next n activation times starting from from. If fewer than
// n future activations exist, the slice is shorter than n.
//
// Example:
//
//	next := expr.NextN(time.Now(), 5)
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
//
// Example:
//
//	now := time.Now().Truncate(time.Minute)
//	if expr.IsDue(now) {
//		sendDailyReport()
//	}
func (e Expression) IsDue(at time.Time) bool {
	return isDue(e.schedule, at)
}

// Next returns the earliest time after t at which the schedule would activate.
// The receiver's location is applied before field matching; if no activation
// exists within the next four years, the zero time is returned.
func (s *cronSchedule) Next(t time.Time) time.Time {
	// Normalise to the schedule's timezone.
	t = t.In(s.loc)

	// Advance by one second (or one minute in five-field mode) to ensure we
	// return a time strictly after t.
	if len(s.second) > 0 {
		t = t.Add(time.Second)
	} else {
		t = t.Add(time.Minute).Truncate(time.Minute)
	}

	// Search forward up to ~4 years to find the next matching instant.
	deadline := t.Add(4 * 365 * 24 * time.Hour)

WRAP:
	for t.Before(deadline) {
		// Month check.
		if !s.month[t.Month()] {
			// Advance to the first day of the next valid month.
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, s.loc)
			continue WRAP
		}
		// Day-of-month and day-of-week check.
		if !s.dayOfMonth[t.Day()] || !s.dayOfWeek[t.Weekday()] {
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, s.loc)
			continue WRAP
		}
		// Hour check.
		if !s.hour[t.Hour()] {
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, s.loc)
			continue WRAP
		}
		// Minute check.
		if !s.minute[t.Minute()] {
			t = t.Add(time.Minute).Truncate(time.Minute)
			continue WRAP
		}
		// Second check (six-field mode only).
		if len(s.second) > 0 && !s.second[t.Second()] {
			t = t.Add(time.Second).Truncate(time.Second)
			continue WRAP
		}
		return t
	}
	return time.Time{}
}

// Next returns the earliest multiple of the interval that is strictly after t.
func (s *intervalSchedule) Next(t time.Time) time.Time {
	return t.Add(s.interval - time.Duration(t.UnixNano())%s.interval)
}

// lookupAlias performs a concurrent-safe lookup in aliasMap using the
// lower-cased alias name (including the "@" prefix).
func lookupAlias(name string) (string, bool) {
	aliasMapMu.RLock()
	v, ok := aliasMap[name]
	aliasMapMu.RUnlock()
	return v, ok
}
