package slogger

import "time"

// ///////////////////////////////////////////////////////////////////////////
// SamplingOptions accessors
// ///////////////////////////////////////////////////////////////////////////

// First returns the number of messages always logged within the period.
//
// Returns:
//
// the first value.
func (s *SamplingOptions) First() int {
	if s == nil {
		return 0
	}
	return s.first
}

// SetFirst sets the number of messages always logged within the period.
//
// Parameters:
//   - `first`: the number of messages
func (s *SamplingOptions) SetFirst(first int) {
	if s == nil {
		return
	}
	s.first = first
}

// Period returns the sampling window duration.
//
// Returns:
//
// the period value.
func (s *SamplingOptions) Period() time.Duration {
	if s == nil {
		return 0
	}
	return s.period
}

// SetPeriod sets the sampling window duration.
//
// Parameters:
//   - `period`: the sampling window
func (s *SamplingOptions) SetPeriod(period time.Duration) {
	if s == nil {
		return
	}
	s.period = period
}

// Thereafter returns the interval for logging after First is exhausted.
//
// Returns:
//
// the thereafter value.
func (s *SamplingOptions) Thereafter() int {
	if s == nil {
		return 0
	}
	return s.thereafter
}

// SetThereafter sets the interval for logging after First is exhausted.
//
// Parameters:
//   - `thereafter`: every Nth message after First
func (s *SamplingOptions) SetThereafter(thereafter int) {
	if s == nil {
		return
	}
	s.thereafter = thereafter
}

// WithFirst sets the first count and returns the receiver for chaining.
//
// Parameters:
//   - `first`: the number of messages
//
// Returns:
//
// the receiver, for method chaining.
func (s *SamplingOptions) WithFirst(first int) *SamplingOptions {
	s.SetFirst(first)
	return s
}

// WithPeriod sets the sampling window duration and returns the receiver for chaining.
//
// Parameters:
//   - `period`: the sampling window
//
// Returns:
//
// the receiver, for method chaining.
func (s *SamplingOptions) WithPeriod(period time.Duration) *SamplingOptions {
	s.SetPeriod(period)
	return s
}

// WithThereafter sets the interval for logging after First is exhausted and returns the receiver for chaining.
//
// Parameters:
//   - `thereafter`: every Nth message after First
//
// Returns:
//
// the receiver, for method chaining.
func (s *SamplingOptions) WithThereafter(thereafter int) *SamplingOptions {
	s.SetThereafter(thereafter)
	return s
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
		b.resetAt = now.Add(s.opts.period)
	}

	b.count++
	cnt := b.count

	first := uint64(s.opts.first)
	if cnt <= first {
		return true
	}
	if s.opts.thereafter == 0 {
		return false
	}
	return (cnt-first-1)%uint64(s.opts.thereafter) == 0
}
