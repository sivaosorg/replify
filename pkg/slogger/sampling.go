package slogger

import (
	"sync"
	"time"
)

// SamplingOptions configures per-message rate limiting for a Logger.
// The first First messages per Period are always logged.
// After that, every Thereafter-th message is logged (0 means drop all remaining).
type SamplingOptions struct {
	// First is the number of identical messages always logged within Period.
	First int
	// Period is the window after which the counter resets.
	Period time.Duration
	// Thereafter logs every Nth message after the First are exhausted.
	// Zero means drop all subsequent messages.
	Thereafter int
}

type samplerBucket struct {
	mu      sync.Mutex
	count   uint64
	resetAt time.Time
}

type sampler struct {
	opts    SamplingOptions
	buckets sync.Map // string (message) -> *samplerBucket
}

// newSampler creates a sampler with the given options.
//
// Parameters:
//   - `opts`: the sampling configuration
//
// Returns:
//
// a ready-to-use *sampler.
func newSampler(opts SamplingOptions) *sampler {
	return &sampler{opts: opts}
}

// allow reports whether this invocation of msg should be logged.
//
// Parameters:
//   - `msg`: the log message used as the bucket key
//
// Returns:
//
// true when the message should be emitted.
func (s *sampler) allow(msg string) bool {
	v, _ := s.buckets.LoadOrStore(msg, &samplerBucket{})
	b := v.(*samplerBucket)

	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if now.After(b.resetAt) {
		b.count = 0
		b.resetAt = now.Add(s.opts.Period)
	}

	b.count++
	cnt := b.count

	first := uint64(s.opts.First)
	if cnt <= first {
		return true
	}
	if s.opts.Thereafter == 0 {
		return false
	}
	return (cnt-first-1)%uint64(s.opts.Thereafter) == 0
}
