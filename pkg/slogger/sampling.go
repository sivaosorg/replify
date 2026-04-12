package slogger

import "time"

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
