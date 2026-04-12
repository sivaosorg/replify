package slogger

import (
	"testing"
	"time"
)

// =============================================================================
// RotationOptions Constructor Tests
// =============================================================================

func TestNewRotationOptions(t *testing.T) {
	t.Parallel()

	opts := NewRotationOptions()
	assertNotNil(t, opts)
}

// =============================================================================
// RotationOptions Accessor Tests
// =============================================================================

func TestRotationOptions_Directory(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		opts := NewRotationOptions()
		opts.SetDirectory("/var/log/app")
		assertEqual(t, "/var/log/app", opts.Directory())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertEqual(t, "", opts.Directory())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertNotPanics(t, func() {
			opts.SetDirectory("/var/log")
		})
	})
}

func TestRotationOptions_MaxBytes(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		opts := NewRotationOptions()
		opts.SetMaxBytes(1024 * 1024 * 10)
		assertEqual(t, int64(10485760), opts.MaxBytes())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertEqual(t, int64(0), opts.MaxBytes())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertNotPanics(t, func() {
			opts.SetMaxBytes(1024)
		})
	})
}

func TestRotationOptions_MaxAge(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		opts := NewRotationOptions()
		age := 24 * time.Hour
		opts.SetMaxAge(age)
		assertEqual(t, age, opts.MaxAge())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertEqual(t, time.Duration(0), opts.MaxAge())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertNotPanics(t, func() {
			opts.SetMaxAge(time.Hour)
		})
	})
}

func TestRotationOptions_IsCompress(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		opts := NewRotationOptions()
		opts.SetCompress(true)
		assertTrue(t, opts.IsCompress())
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		opts := NewRotationOptions()
		opts.SetCompress(false)
		assertFalse(t, opts.IsCompress())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertFalse(t, opts.IsCompress())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var opts *RotationOptions
		assertNotPanics(t, func() {
			opts.SetCompress(true)
		})
	})
}

// =============================================================================
// RotationOptions Fluent API Tests
// =============================================================================

func TestRotationOptions_WithDirectory(t *testing.T) {
	t.Parallel()

	opts := NewRotationOptions().WithDirectory("/var/log/app")
	assertEqual(t, "/var/log/app", opts.Directory())
}

func TestRotationOptions_WithMaxBytes(t *testing.T) {
	t.Parallel()

	opts := NewRotationOptions().WithMaxBytes(1024 * 1024)
	assertEqual(t, int64(1048576), opts.MaxBytes())
}

func TestRotationOptions_WithMaxAge(t *testing.T) {
	t.Parallel()

	opts := NewRotationOptions().WithMaxAge(24 * time.Hour)
	assertEqual(t, 24*time.Hour, opts.MaxAge())
}

func TestRotationOptions_WithCompress(t *testing.T) {
	t.Parallel()

	opts := NewRotationOptions().WithCompress(true)
	assertTrue(t, opts.IsCompress())
}

func TestRotationOptions_FluentChaining(t *testing.T) {
	t.Parallel()

	opts := NewRotationOptions().
		WithDirectory("/var/log").
		WithMaxBytes(10 * 1024 * 1024).
		WithMaxAge(7 * 24 * time.Hour).
		WithCompress(true)

	assertEqual(t, "/var/log", opts.Directory())
	assertEqual(t, int64(10485760), opts.MaxBytes())
	assertEqual(t, 7*24*time.Hour, opts.MaxAge())
	assertTrue(t, opts.IsCompress())
}

// =============================================================================
// LevelFileWriter Tests
// =============================================================================

func TestLevelFileWriter_Options(t *testing.T) {
	t.Parallel()

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var lfw *LevelFileWriter
		opts := lfw.Options()
		assertEqual(t, RotationOptions{}, opts)
	})
}

// =============================================================================
// LevelWriterHook Tests
// =============================================================================

func TestLevelWriterHook_Writer(t *testing.T) {
	t.Parallel()

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hook *LevelWriterHook
		assertNil(t, hook.Writer())
	})
}

func TestLevelWriterHook_Formatter(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		formatter := NewJSONFormatter()
		hook := &LevelWriterHook{formatter: formatter}
		assertEqual(t, formatter, hook.Formatter())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hook *LevelWriterHook
		assertNil(t, hook.Formatter())
	})
}

func TestLevelWriterHook_SetFormatter(t *testing.T) {
	t.Parallel()

	t.Run("sets formatter", func(t *testing.T) {
		t.Parallel()
		hook := &LevelWriterHook{}
		formatter := NewJSONFormatter()
		hook.SetFormatter(formatter)
		assertEqual(t, formatter, hook.Formatter())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var hook *LevelWriterHook
		assertNotPanics(t, func() {
			hook.SetFormatter(NewJSONFormatter())
		})
	})
}

func TestLevelWriterHook_Levels(t *testing.T) {
	t.Parallel()

	hook := &LevelWriterHook{
		levels: []Level{InfoLevel, WarnLevel, ErrorLevel},
	}
	levels := hook.Levels()
	assertLen(t, levels, 3)
}

// =============================================================================
// NewLevelWriterHook Tests
// =============================================================================

func TestNewLevelWriterHook(t *testing.T) {
	t.Parallel()

	t.Run("with specific levels", func(t *testing.T) {
		t.Parallel()
		formatter := NewTextFormatter(nil)
		hook := NewLevelWriterHook(nil, formatter, InfoLevel, WarnLevel)
		assertLen(t, hook.Levels(), 2)
		assertEqual(t, formatter, hook.Formatter())
	})

	t.Run("with no levels uses all", func(t *testing.T) {
		t.Parallel()
		formatter := NewTextFormatter(nil)
		hook := NewLevelWriterHook(nil, formatter)
		assertLen(t, hook.Levels(), 7) // All levels
	})
}

// =============================================================================
// levelFileName Tests
// =============================================================================

func TestLevelFileName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level Level
		want  string
	}{
		{DebugLevel, "debug.log"},
		{InfoLevel, "info.log"},
		{WarnLevel, "warn.log"},
		{ErrorLevel, "error.log"},
		{FatalLevel, "error.log"},
		{PanicLevel, "error.log"},
		{TraceLevel, "error.log"}, // Falls to default
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.level.String(), func(t *testing.T) {
			t.Parallel()
			got := levelFileName(tt.level)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// rotatingFile Tests
// =============================================================================

func TestRotatingFile_shouldRotate(t *testing.T) {
	t.Parallel()

	t.Run("size exceeded", func(t *testing.T) {
		t.Parallel()
		rf := &rotatingFile{
			size:     900,
			maxBytes: 1000,
		}
		assertTrue(t, rf.shouldRotate(200))
	})

	t.Run("size not exceeded", func(t *testing.T) {
		t.Parallel()
		rf := &rotatingFile{
			size:     500,
			maxBytes: 1000,
		}
		assertFalse(t, rf.shouldRotate(200))
	})

	t.Run("age exceeded", func(t *testing.T) {
		t.Parallel()
		rf := &rotatingFile{
			size:     0,
			maxBytes: 1000,
			maxAge:   time.Hour,
			openedAt: time.Now().Add(-2 * time.Hour),
		}
		assertTrue(t, rf.shouldRotate(100))
	})

	t.Run("age not exceeded", func(t *testing.T) {
		t.Parallel()
		rf := &rotatingFile{
			size:     0,
			maxBytes: 1000,
			maxAge:   time.Hour,
			openedAt: time.Now(),
		}
		assertFalse(t, rf.shouldRotate(100))
	})

	t.Run("maxBytes disabled", func(t *testing.T) {
		t.Parallel()
		rf := &rotatingFile{
			size:     1000000,
			maxBytes: 0, // Disabled
		}
		assertFalse(t, rf.shouldRotate(100))
	})

	t.Run("maxAge disabled", func(t *testing.T) {
		t.Parallel()
		rf := &rotatingFile{
			maxAge:   0, // Disabled
			openedAt: time.Now().Add(-100 * time.Hour),
		}
		assertFalse(t, rf.shouldRotate(0))
	})
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkRotationOptions_Directory(b *testing.B) {
	opts := NewRotationOptions()
	opts.SetDirectory("/var/log/app")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opts.Directory()
	}
}

func BenchmarkLevelFileName(b *testing.B) {
	levels := []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = levelFileName(levels[i%len(levels)])
	}
}
