package slogger

import (
	"errors"
	"sync"
	"testing"
)

// =============================================================================
// Hooks Constructor Tests
// =============================================================================

func TestNewHooks(t *testing.T) {
	t.Parallel()

	hooks := NewHooks()
	assertNotNil(t, hooks)
	assertNotNil(t, hooks.hooks)
}

// =============================================================================
// Hooks.Add Tests
// =============================================================================

func TestHooks_Add(t *testing.T) {
	t.Parallel()

	t.Run("adds hook for single level", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook)

		assertEqual(t, 1, hooks.Len(InfoLevel))
		assertEqual(t, 0, hooks.Len(DebugLevel))
	})

	t.Run("adds hook for multiple levels", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook := &testHookImpl{levels: []Level{InfoLevel, WarnLevel, ErrorLevel}}
		hooks.Add(hook)

		assertEqual(t, 1, hooks.Len(InfoLevel))
		assertEqual(t, 1, hooks.Len(WarnLevel))
		assertEqual(t, 1, hooks.Len(ErrorLevel))
		assertEqual(t, 0, hooks.Len(DebugLevel))
	})

	t.Run("adds multiple hooks", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook1 := &testHookImpl{levels: []Level{InfoLevel}}
		hook2 := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook1)
		hooks.Add(hook2)

		assertEqual(t, 2, hooks.Len(InfoLevel))
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hooks *Hooks
		hook := &testHookImpl{levels: []Level{InfoLevel}}
		assertNotPanics(t, func() {
			hooks.Add(hook)
		})
	})
}

// =============================================================================
// Hooks.Fire Tests
// =============================================================================

func TestHooks_Fire(t *testing.T) {
	t.Parallel()

	t.Run("fires hook", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook)

		entry := &Entry{level: InfoLevel, message: "test"}
		err := hooks.Fire(InfoLevel, entry)

		assertNoError(t, err)
		assertEqual(t, 1, hook.fireCount)
	})

	t.Run("fires multiple hooks", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook1 := &testHookImpl{levels: []Level{InfoLevel}}
		hook2 := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook1)
		hooks.Add(hook2)

		entry := &Entry{level: InfoLevel, message: "test"}
		err := hooks.Fire(InfoLevel, entry)

		assertNoError(t, err)
		assertEqual(t, 1, hook1.fireCount)
		assertEqual(t, 1, hook2.fireCount)
	})

	t.Run("returns first error", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		expectedErr := errors.New("hook error")
		hook := &testHookImpl{levels: []Level{InfoLevel}, err: expectedErr}
		hooks.Add(hook)

		entry := &Entry{level: InfoLevel, message: "test"}
		err := hooks.Fire(InfoLevel, entry)

		assertError(t, err)
		assertEqual(t, expectedErr, err)
	})

	t.Run("continues after error", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook1 := &testHookImpl{levels: []Level{InfoLevel}, err: errors.New("error")}
		hook2 := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook1)
		hooks.Add(hook2)

		entry := &Entry{level: InfoLevel, message: "test"}
		_ = hooks.Fire(InfoLevel, entry)

		// Both hooks should fire
		assertEqual(t, 1, hook1.fireCount)
		assertEqual(t, 1, hook2.fireCount)
	})

	t.Run("no hooks for level", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook := &testHookImpl{levels: []Level{ErrorLevel}}
		hooks.Add(hook)

		entry := &Entry{level: InfoLevel, message: "test"}
		err := hooks.Fire(InfoLevel, entry)

		assertNoError(t, err)
		assertEqual(t, 0, hook.fireCount)
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hooks *Hooks
		entry := &Entry{level: InfoLevel, message: "test"}
		assertNotPanics(t, func() {
			err := hooks.Fire(InfoLevel, entry)
			assertNoError(t, err)
		})
	})
}

// =============================================================================
// Hooks.Len Tests
// =============================================================================

func TestHooks_Len(t *testing.T) {
	t.Parallel()

	t.Run("returns count", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hooks.Add(&testHookImpl{levels: []Level{InfoLevel}})
		hooks.Add(&testHookImpl{levels: []Level{InfoLevel}})
		hooks.Add(&testHookImpl{levels: []Level{InfoLevel}})

		assertEqual(t, 3, hooks.Len(InfoLevel))
	})

	t.Run("zero for no hooks", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		assertEqual(t, 0, hooks.Len(InfoLevel))
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hooks *Hooks
		assertEqual(t, 0, hooks.Len(InfoLevel))
	})
}

// =============================================================================
// Hooks.HooksFor Tests
// =============================================================================

func TestHooks_HooksFor(t *testing.T) {
	t.Parallel()

	t.Run("returns hooks for level", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook1 := &testHookImpl{levels: []Level{InfoLevel}}
		hook2 := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook1)
		hooks.Add(hook2)

		got := hooks.HooksFor(InfoLevel)
		assertLen(t, got, 2)
	})

	t.Run("returns copy", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		hook := &testHookImpl{levels: []Level{InfoLevel}}
		hooks.Add(hook)

		got := hooks.HooksFor(InfoLevel)
		got[0] = nil

		// Original should be unchanged
		assertLen(t, hooks.HooksFor(InfoLevel), 1)
		assertNotNil(t, hooks.HooksFor(InfoLevel)[0])
	})

	t.Run("nil for no hooks", func(t *testing.T) {
		t.Parallel()
		hooks := NewHooks()
		assertNil(t, hooks.HooksFor(InfoLevel))
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hooks *Hooks
		assertNil(t, hooks.HooksFor(InfoLevel))
	})
}

// =============================================================================
// Hooks Concurrent Access Tests
// =============================================================================

func TestHooks_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	hooks := NewHooks()
	var wg sync.WaitGroup

	// Concurrent adds
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hooks.Add(&testHookImpl{levels: []Level{InfoLevel}})
		}()
	}

	// Concurrent fires
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = hooks.Fire(InfoLevel, &Entry{level: InfoLevel})
		}()
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = hooks.Len(InfoLevel)
			_ = hooks.HooksFor(InfoLevel)
		}()
	}

	wg.Wait()

	// Should have 100 hooks
	assertEqual(t, 100, hooks.Len(InfoLevel))
}

// =============================================================================
// testHookImpl helper
// =============================================================================

type testHookImpl struct {
	mu        sync.Mutex
	levels    []Level
	err       error
	fireCount int
}

func (h *testHookImpl) Levels() []Level {
	return h.levels
}

func (h *testHookImpl) Fire(e *Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fireCount++
	return h.err
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkHooks_Add(b *testing.B) {
	hooks := NewHooks()
	hook := &testHookImpl{levels: []Level{InfoLevel}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hooks.Add(hook)
	}
}

func BenchmarkHooks_Fire(b *testing.B) {
	hooks := NewHooks()
	hooks.Add(&testHookImpl{levels: []Level{InfoLevel}})
	entry := &Entry{level: InfoLevel}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hooks.Fire(InfoLevel, entry)
	}
}

func BenchmarkHooks_Len(b *testing.B) {
	hooks := NewHooks()
	hooks.Add(&testHookImpl{levels: []Level{InfoLevel}})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hooks.Len(InfoLevel)
	}
}
