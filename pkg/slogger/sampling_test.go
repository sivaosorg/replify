package slogger

import (
	"testing"
	"time"
)

// =============================================================================
// SamplingOptions Constructor Tests
// =============================================================================

func TestNewSamplingOptions(t *testing.T) {
	t.Parallel()

	opts := NewSamplingOptions()
	assertNotNil(t, opts)
}

// =============================================================================
// SamplingOptions Accessor Tests
// =============================================================================

func TestSamplingOptions_First(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		opts := NewSamplingOptions()
		opts.SetFirst(10)
		assertEqual(t, 10, opts.First())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *SamplingOptions
		assertEqual(t, 0, opts.First())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *SamplingOptions
		assertNotPanics(t, func() {
			opts.SetFirst(10)
		})
	})
}

func TestSamplingOptions_Period(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		opts := NewSamplingOptions()
		opts.SetPeriod(time.Second)
		assertEqual(t, time.Second, opts.Period())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *SamplingOptions
		assertEqual(t, time.Duration(0), opts.Period())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *SamplingOptions
		assertNotPanics(t, func() {
			opts.SetPeriod(time.Second)
		})
	})
}

func TestSamplingOptions_Thereafter(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		opts := NewSamplingOptions()
		opts.SetThereafter(5)
		assertEqual(t, 5, opts.Thereafter())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *SamplingOptions
		assertEqual(t, 0, opts.Thereafter())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *SamplingOptions
		assertNotPanics(t, func() {
			opts.SetThereafter(5)
		})
	})
}

// =============================================================================
// SamplingOptions Fluent API Tests
// =============================================================================

func TestSamplingOptions_WithFirst(t *testing.T) {
	t.Parallel()

	opts := NewSamplingOptions().WithFirst(10)
	assertEqual(t, 10, opts.First())
}

func TestSamplingOptions_WithPeriod(t *testing.T) {
	t.Parallel()

	opts := NewSamplingOptions().WithPeriod(time.Second)
	assertEqual(t, time.Second, opts.Period())
}

func TestSamplingOptions_WithThereafter(t *testing.T) {
	t.Parallel()

	opts := NewSamplingOptions().WithThereafter(5)
	assertEqual(t, 5, opts.Thereafter())
}

func TestSamplingOptions_FluentChaining(t *testing.T) {
	t.Parallel()

	opts := NewSamplingOptions().
		WithFirst(10).
		WithPeriod(time.Second).
		WithThereafter(5)

	assertEqual(t, 10, opts.First())
	assertEqual(t, time.Second, opts.Period())
	assertEqual(t, 5, opts.Thereafter())
}

// =============================================================================
// sampler.allow Tests
// =============================================================================

func TestSampler_allow(t *testing.T) {
	t.Parallel()

	t.Run("allows first N messages", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      3,
			period:     time.Hour, // Long period so no reset
			thereafter: 0,
		})

		assertTrue(t, s.allow("test"))
		assertTrue(t, s.allow("test"))
		assertTrue(t, s.allow("test"))
		assertFalse(t, s.allow("test")) // 4th should be blocked
	})

	t.Run("thereafter sampling", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      2,
			period:     time.Hour, // Long period so no reset
			thereafter: 3,         // Every 3rd after first 2
		})

		assertTrue(t, s.allow("test"))  // 1st - first
		assertTrue(t, s.allow("test"))  // 2nd - first
		assertTrue(t, s.allow("test"))  // 3rd - (3-2-1) % 3 == 0
		assertFalse(t, s.allow("test")) // 4th - (4-2-1) % 3 != 0
		assertFalse(t, s.allow("test")) // 5th - (5-2-1) % 3 != 0
		assertTrue(t, s.allow("test"))  // 6th - (6-2-1) % 3 == 0
	})

	t.Run("different messages independent", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      1,
			period:     time.Hour,
			thereafter: 0,
		})

		assertTrue(t, s.allow("msg1"))
		assertTrue(t, s.allow("msg2"))
		assertTrue(t, s.allow("msg3"))
		assertFalse(t, s.allow("msg1")) // Second occurrence
		assertFalse(t, s.allow("msg2")) // Second occurrence
	})

	t.Run("resets after period", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      1,
			period:     10 * time.Millisecond,
			thereafter: 0,
		})

		assertTrue(t, s.allow("test"))
		assertFalse(t, s.allow("test"))

		// Wait for period to elapse
		time.Sleep(20 * time.Millisecond)

		assertTrue(t, s.allow("test")) // Should be allowed again
	})

	t.Run("thereafter zero blocks all after first", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      2,
			period:     time.Hour,
			thereafter: 0, // Block all after first
		})

		assertTrue(t, s.allow("test"))
		assertTrue(t, s.allow("test"))
		assertFalse(t, s.allow("test"))
		assertFalse(t, s.allow("test"))
		assertFalse(t, s.allow("test"))
	})
}

// =============================================================================
// newSampler Tests
// =============================================================================

func TestNewSampler(t *testing.T) {
	t.Parallel()

	opts := SamplingOptions{
		first:      5,
		period:     time.Minute,
		thereafter: 10,
	}
	s := newSampler(opts)
	assertNotNil(t, s)
	assertEqual(t, opts, s.opts)
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestSampler_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("first zero", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      0,
			period:     time.Hour,
			thereafter: 2,
		})

		// First message should be allowed if thereafter is set
		assertTrue(t, s.allow("test"))  // (1-0-1) % 2 == 0
		assertFalse(t, s.allow("test")) // (2-0-1) % 2 != 0
		assertTrue(t, s.allow("test"))  // (3-0-1) % 2 == 0
	})

	t.Run("empty message", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      1,
			period:     time.Hour,
			thereafter: 0,
		})

		assertTrue(t, s.allow(""))
		assertFalse(t, s.allow(""))
	})

	t.Run("long message", func(t *testing.T) {
		t.Parallel()
		s := newSampler(SamplingOptions{
			first:      1,
			period:     time.Hour,
			thereafter: 0,
		})

		longMsg := string(make([]byte, 10000))
		assertTrue(t, s.allow(longMsg))
		assertFalse(t, s.allow(longMsg))
	})
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkSamplingOptions_First(b *testing.B) {
	opts := NewSamplingOptions()
	opts.SetFirst(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opts.First()
	}
}

func BenchmarkSampler_allow_allowed(b *testing.B) {
	s := newSampler(SamplingOptions{
		first:      1000000,
		period:     time.Hour,
		thereafter: 0,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.allow("test")
	}
}

func BenchmarkSampler_allow_blocked(b *testing.B) {
	s := newSampler(SamplingOptions{
		first:      1,
		period:     time.Hour,
		thereafter: 0,
	})
	s.allow("test") // First message
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.allow("test")
	}
}
