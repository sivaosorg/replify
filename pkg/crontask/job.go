package crontask

import (
	"context"
	"sync"
	"time"
)

// JobFunc is the function signature for a scheduled job. The context passed
// to JobFunc is derived from the job's base context (or context.Background
// when none is configured) and may carry a deadline when WithTimeout is set.
// Jobs should honour context cancellation for clean shutdown.
type JobFunc func(ctx context.Context) error

// JobInfo is an immutable snapshot of a registered job's metadata and
// runtime statistics. It is returned by Jobs() and used for introspection.
type JobInfo struct {
	// ID is the unique identifier of the job, either supplied by the caller
	// via WithJobID or generated automatically at registration time.
	ID string

	// Name is an optional human-readable label set via WithJobName.
	Name string

	// Expression is the raw cron expression string as supplied to Register.
	Expression string

	// NextRun is the next scheduled activation time in the scheduler's
	// timezone. The zero value means the schedule has no future activations.
	NextRun time.Time

	// LastRun is the time the job was most recently dispatched. The zero
	// value means the job has never been executed.
	LastRun time.Time

	// LastErr is the error returned by the most recent execution, or nil if
	// the last execution succeeded or the job has never been run.
	LastErr error

	// RunCount is the total number of times the job has been dispatched
	// (across all retries within a single schedule activation, only the
	// initial dispatch is counted).
	RunCount int64
}

// entry is the internal mutable state for a single registered job. Access
// must be protected by the owning registry's mutex.
type entry struct {
	id         string
	name       string
	expression string
	schedule   Schedule
	fn         JobFunc
	cfg        jobConfig

	mu       sync.Mutex
	nextRun  time.Time
	lastRun  time.Time
	lastErr  error
	runCount int64
}

// snapshot returns an immutable JobInfo that reflects the current state of the
// entry. The caller must hold at least a read lock on the entry's own mutex OR
// be the sole owner.
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

// registry is the concurrent-safe store of all registered job entries.
type registry struct {
	mu      sync.RWMutex
	entries map[string]*entry
}

// newRegistry allocates an initialised registry.
func newRegistry() *registry {
	return &registry{entries: make(map[string]*entry)}
}

// add inserts or replaces an entry in the registry. If an entry with the same
// ID already exists it is overwritten without error.
func (r *registry) add(e *entry) {
	r.mu.Lock()
	r.entries[e.id] = e
	r.mu.Unlock()
}

// remove deletes the entry with the given id. It returns ErrJobNotFound when
// the id is not present.
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
func (r *registry) get(id string) (*entry, bool) {
	r.mu.RLock()
	e, ok := r.entries[id]
	r.mu.RUnlock()
	return e, ok
}

// list returns a slice of all entries in an unspecified order.
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

// updateNextRun recomputes the next activation time for e using the
// scheduler's reference time as the base.
func updateNextRun(e *entry, now time.Time) {
	e.mu.Lock()
	e.nextRun = e.schedule.Next(now)
	e.mu.Unlock()
}

// recordResult stores the outcome of a single execution into the entry's
// mutable state.
func recordResult(e *entry, t time.Time, err error) {
	e.mu.Lock()
	e.lastRun = t
	e.lastErr = err
	e.runCount++
	e.mu.Unlock()
}
