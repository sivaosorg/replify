package crontask

import (
	"time"
)

// snapshot returns an immutable JobInfo that reflects the current state of the
// entry. The caller must hold at least a read lock on the entry's own mutex OR
// be the sole owner.
//
// Example:
//
//	info := entry.snapshot()
func (e *entry) snapshot() JobInfo {
	e.mu.Lock()
	defer e.mu.Unlock()
	return JobInfo{
		ID:         e.id,
		Name:       e.name,
		Expression: e.expression,
		NextRun:    e.nextRun,
		LastRun:    e.lastRun,
		LastErr:    e.lastErr,
		RunCount:   e.runCount,
	}
}

// add inserts or replaces an entry in the registry. If an entry with the same
// ID already exists it is overwritten without error.
//
// Example:
//
//	r.add(entry)
func (r *registry) add(e *entry) {
	r.mu.Lock()
	r.entries[e.id] = e
	r.mu.Unlock()
}

// remove deletes the entry with the given id. It returns ErrJobNotFound when
// the id is not present.
//
// Example:
//
//	r.remove("my-job")
func (r *registry) remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[id]; !ok {
		return ErrJobNotFound
	}
	delete(r.entries, id)
	return nil
}

// get returns the entry for the given id and a boolean indicating whether
// the entry was found.
//
// Example:
//
//	e, ok := r.get("my-job")
func (r *registry) get(id string) (*entry, bool) {
	r.mu.RLock()
	e, ok := r.entries[id]
	r.mu.RUnlock()
	return e, ok
}

// list returns a slice of all entries in an unspecified order.
//
// Example:
//
//	entries := r.list()
func (r *registry) list() []*entry {
	r.mu.RLock()
	out := make([]*entry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	r.mu.RUnlock()
	return out
}

// nextDue scans all entries and returns those whose nextRun is at or before
// the given reference time, together with a duration until the closest
// upcoming activation among all remaining entries. When all entries have a
// zero nextRun (schedules exhausted), remaining is set to a large value
// (maxSleep) so the loop simply parks.
//
// Example:
//
//	due, remaining := r.nextDue(time.Now())
func (r *registry) nextDue(now time.Time) (due []*entry, remaining time.Duration) {
	const maxSleep = time.Hour

	remaining = maxSleep
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.entries {
		e.mu.Lock()
		nr := e.nextRun
		e.mu.Unlock()

		if nr.IsZero() {
			continue
		}
		if !nr.After(now) {
			due = append(due, e)
		} else {
			if d := nr.Sub(now); d < remaining {
				remaining = d
			}
		}
	}
	return
}
