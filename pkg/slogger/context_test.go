package slogger

import (
	"context"
	"testing"
)

// =============================================================================
// FieldsFromContext Tests
// =============================================================================

func TestFieldsFromContext(t *testing.T) {
	t.Parallel()

	t.Run("returns fields from context", func(t *testing.T) {
		t.Parallel()
		fields := []Field{String("a", "1"), Int("b", 2)}
		ctx := context.WithValue(context.Background(), contextKey{}, fields)

		got := FieldsFromContext(ctx)
		assertLen(t, got, 2)
		assertEqual(t, "a", got[0].Key())
	})

	t.Run("returns nil for nil context", func(t *testing.T) {
		t.Parallel()
		got := FieldsFromContext(nil)
		assertNil(t, got)
	})

	t.Run("returns nil for context without fields", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		got := FieldsFromContext(ctx)
		assertNil(t, got)
	})

	t.Run("returns nil for wrong type", func(t *testing.T) {
		t.Parallel()
		ctx := context.WithValue(context.Background(), contextKey{}, "not fields")
		got := FieldsFromContext(ctx)
		assertNil(t, got)
	})
}

// =============================================================================
// WithContextFields Tests
// =============================================================================

func TestWithContextFields(t *testing.T) {
	t.Parallel()

	t.Run("adds fields to context", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = WithContextFields(ctx, String("key", "value"))

		fields := FieldsFromContext(ctx)
		assertLen(t, fields, 1)
		assertEqual(t, "key", fields[0].Key())
	})

	t.Run("merges with existing fields", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = WithContextFields(ctx, String("a", "1"))
		ctx = WithContextFields(ctx, String("b", "2"))

		fields := FieldsFromContext(ctx)
		assertLen(t, fields, 2)
		assertEqual(t, "a", fields[0].Key())
		assertEqual(t, "b", fields[1].Key())
	})

	t.Run("returns same context for empty fields", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		newCtx := WithContextFields(ctx)

		assertEqual(t, ctx, newCtx)
	})

	t.Run("adds multiple fields at once", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = WithContextFields(ctx, String("a", "1"), Int("b", 2), Bool("c", true))

		fields := FieldsFromContext(ctx)
		assertLen(t, fields, 3)
	})

	t.Run("preserves order", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = WithContextFields(ctx, String("first", "1"))
		ctx = WithContextFields(ctx, String("second", "2"))
		ctx = WithContextFields(ctx, String("third", "3"))

		fields := FieldsFromContext(ctx)
		assertLen(t, fields, 3)
		assertEqual(t, "first", fields[0].Key())
		assertEqual(t, "second", fields[1].Key())
		assertEqual(t, "third", fields[2].Key())
	})
}

// =============================================================================
// Context Integration Tests
// =============================================================================

func TestContextFieldsIntegration(t *testing.T) {
	t.Parallel()

	t.Run("context fields appear in logs", func(t *testing.T) {
		t.Parallel()
		// This is more of an integration test
		ctx := context.Background()
		ctx = WithContextFields(ctx, String("request_id", "abc123"))
		ctx = WithContextFields(ctx, String("user_id", "user456"))

		fields := FieldsFromContext(ctx)
		assertLen(t, fields, 2)
		assertEqual(t, "abc123", fields[0].Value())
		assertEqual(t, "user456", fields[1].Value())
	})
}

// =============================================================================
// GlobalWithContextFields Tests
// =============================================================================

func TestGlobalWithContextFields(t *testing.T) {
	t.Parallel()

	t.Run("works same as WithContextFields", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = GlobalWithContextFields(ctx, String("key", "value"))

		fields := FieldsFromContext(ctx)
		assertLen(t, fields, 1)
		assertEqual(t, "key", fields[0].Key())
	})
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkFieldsFromContext(b *testing.B) {
	fields := []Field{String("a", "1"), Int("b", 2), Bool("c", true)}
	ctx := context.WithValue(context.Background(), contextKey{}, fields)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FieldsFromContext(ctx)
	}
}

func BenchmarkWithContextFields(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithContextFields(ctx, String("key", "value"))
	}
}

func BenchmarkWithContextFields_Merge(b *testing.B) {
	ctx := context.Background()
	ctx = WithContextFields(ctx, String("existing", "value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithContextFields(ctx, String("new", "value"))
	}
}
