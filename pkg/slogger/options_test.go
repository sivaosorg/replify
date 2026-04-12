package slogger

import (
	"bytes"
	"io"
	"testing"
	"time"
)

// =============================================================================
// Functional Options Tests
// =============================================================================

func TestWithLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level Level
	}{
		{name: "trace", level: TraceLevel},
		{name: "debug", level: DebugLevel},
		{name: "info", level: InfoLevel},
		{name: "warn", level: WarnLevel},
		{name: "error", level: ErrorLevel},
		{name: "fatal", level: FatalLevel},
		{name: "panic", level: PanicLevel},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opt := WithLevel(tt.level)
			o := &Options{}
			opt(o)
			assertEqual(t, tt.level, o.level)
		})
	}
}

func TestWithFormatter(t *testing.T) {
	t.Parallel()

	t.Run("json formatter", func(t *testing.T) {
		t.Parallel()
		formatter := NewJSONFormatter()
		opt := WithFormatter(formatter)
		o := &Options{}
		opt(o)
		assertEqual(t, formatter, o.formatter)
	})

	t.Run("text formatter", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		formatter := NewTextFormatter(&buf)
		opt := WithFormatter(formatter)
		o := &Options{}
		opt(o)
		assertEqual(t, formatter, o.formatter)
	})

	t.Run("nil formatter", func(t *testing.T) {
		t.Parallel()
		opt := WithFormatter(nil)
		o := &Options{}
		opt(o)
		assertNil(t, o.formatter)
	})
}

func TestWithOutput(t *testing.T) {
	t.Parallel()

	t.Run("buffer output", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		opt := WithOutput(&buf)
		o := &Options{}
		opt(o)
		assertEqual(t, &buf, o.output)
	})

	t.Run("nil output", func(t *testing.T) {
		t.Parallel()
		opt := WithOutput(nil)
		o := &Options{}
		opt(o)
		assertNil(t, o.output)
	})
}

func TestWithCaller(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		opt := WithCaller(true)
		o := &Options{}
		opt(o)
		assertTrue(t, o.caller)
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		opt := WithCaller(false)
		o := &Options{}
		o.caller = true // Set to true first
		opt(o)
		assertFalse(t, o.caller)
	})
}

func TestWithCallerSkip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		skip int
	}{
		{name: "zero skip", skip: 0},
		{name: "positive skip", skip: 3},
		{name: "large skip", skip: 100},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opt := WithCallerSkip(tt.skip)
			o := &Options{}
			opt(o)
			assertEqual(t, tt.skip, o.callerSkip)
		})
	}
}

func TestWithFields(t *testing.T) {
	t.Parallel()

	t.Run("single field", func(t *testing.T) {
		t.Parallel()
		opt := WithFields(String("key", "value"))
		o := &Options{}
		opt(o)
		assertLen(t, o.fields, 1)
		assertEqual(t, "key", o.fields[0].Key())
	})

	t.Run("multiple fields", func(t *testing.T) {
		t.Parallel()
		opt := WithFields(String("a", "1"), Int("b", 2), Bool("c", true))
		o := &Options{}
		opt(o)
		assertLen(t, o.fields, 3)
	})

	t.Run("no fields", func(t *testing.T) {
		t.Parallel()
		opt := WithFields()
		o := &Options{}
		opt(o)
		assertEmpty(t, o.fields)
	})

	t.Run("append to existing", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.fields = []Field{String("existing", "value")}
		opt := WithFields(String("new", "value"))
		opt(o)
		assertLen(t, o.fields, 2)
	})
}

func TestWithName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{name: "simple name", input: "logger"},
		{name: "dotted name", input: "app.module.component"},
		{name: "empty name", input: ""},
		{name: "unicode name", input: "日本語ロガー"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opt := WithName(tt.input)
			o := &Options{}
			opt(o)
			assertEqual(t, tt.input, o.name)
		})
	}
}

func TestWithSamplingOpts(t *testing.T) {
	t.Parallel()

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		sampOpts := &SamplingOptions{first: 10, period: time.Second, thereafter: 5}
		opt := WithSamplingOpts(sampOpts)
		o := &Options{}
		opt(o)
		assertEqual(t, sampOpts, o.samplingOpts)
	})

	t.Run("nil options", func(t *testing.T) {
		t.Parallel()
		opt := WithSamplingOpts(nil)
		o := &Options{}
		opt(o)
		assertNil(t, o.samplingOpts)
	})
}

func TestWithRotation(t *testing.T) {
	t.Parallel()

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		rotOpts := &RotationOptions{dir: "/tmp/logs", maxBytes: 1024}
		opt := WithRotation(rotOpts)
		o := &Options{}
		opt(o)
		assertEqual(t, rotOpts, o.rotationOpts)
	})

	t.Run("nil options", func(t *testing.T) {
		t.Parallel()
		opt := WithRotation(nil)
		o := &Options{}
		opt(o)
		assertNil(t, o.rotationOpts)
	})
}

// =============================================================================
// Options Accessor Tests
// =============================================================================

func TestOptions_Level(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetLevel(DebugLevel)
		assertEqual(t, DebugLevel, o.Level())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertEqual(t, InfoLevel, o.Level())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetLevel(DebugLevel)
		})
	})
}

func TestOptions_Formatter(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		f := NewJSONFormatter()
		o.SetFormatter(f)
		assertEqual(t, f, o.Formatter())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNil(t, o.Formatter())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetFormatter(NewJSONFormatter())
		})
	})
}

func TestOptions_Output(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		var buf bytes.Buffer
		o.SetOutput(&buf)
		assertEqual(t, &buf, o.Output())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNil(t, o.Output())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetOutput(io.Discard)
		})
	})
}

func TestOptions_IsCaller(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetCaller(true)
		assertTrue(t, o.IsCaller())
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetCaller(false)
		assertFalse(t, o.IsCaller())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertFalse(t, o.IsCaller())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetCaller(true)
		})
	})
}

func TestOptions_CallerSkip(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetCallerSkip(5)
		assertEqual(t, 5, o.CallerSkip())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertEqual(t, 0, o.CallerSkip())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetCallerSkip(5)
		})
	})
}

func TestOptions_Fields(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		fields := []Field{String("a", "1"), Int("b", 2)}
		o.SetFields(fields)
		got := o.Fields()
		assertLen(t, got, 2)
		assertEqual(t, "a", got[0].Key())
	})

	t.Run("returns copy", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetFields([]Field{String("a", "1")})
		got := o.Fields()
		got[0] = String("modified", "value")
		// Original should be unchanged
		assertEqual(t, "a", o.Fields()[0].Key())
	})

	t.Run("set nil", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetFields([]Field{String("a", "1")})
		o.SetFields(nil)
		assertNil(t, o.Fields())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNil(t, o.Fields())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetFields([]Field{String("a", "1")})
		})
	})
}

func TestOptions_AddFields(t *testing.T) {
	t.Parallel()

	t.Run("add to empty", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.AddFields(String("a", "1"), Int("b", 2))
		assertLen(t, o.Fields(), 2)
	})

	t.Run("add to existing", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetFields([]Field{String("existing", "value")})
		o.AddFields(String("new", "value"))
		assertLen(t, o.Fields(), 2)
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.AddFields(String("a", "1"))
		})
	})
}

func TestOptions_Name(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetName("mylogger")
		assertEqual(t, "mylogger", o.Name())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertEqual(t, "", o.Name())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetName("test")
		})
	})
}

func TestOptions_SamplingOpts(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		sampOpts := &SamplingOptions{first: 5}
		o.SetSamplingOpts(sampOpts)
		assertEqual(t, sampOpts, o.SamplingOpts())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNil(t, o.SamplingOpts())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetSamplingOpts(&SamplingOptions{})
		})
	})
}

func TestOptions_RotationOpts(t *testing.T) {
	t.Parallel()

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		rotOpts := &RotationOptions{dir: "/tmp/logs"}
		o.SetRotationOpts(rotOpts)
		assertEqual(t, rotOpts, o.RotationOpts())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNil(t, o.RotationOpts())
	})

	t.Run("nil receiver set", func(t *testing.T) {
		t.Parallel()
		var o *Options
		assertNotPanics(t, func() {
			o.SetRotationOpts(&RotationOptions{})
		})
	})
}

// =============================================================================
// Options Edge Cases
// =============================================================================

func TestOptions_ZeroValue(t *testing.T) {
	t.Parallel()

	t.Run("zero value options", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		assertEqual(t, Level(0), o.Level())
		assertNil(t, o.Formatter())
		assertNil(t, o.Output())
		assertFalse(t, o.IsCaller())
		assertEqual(t, 0, o.CallerSkip())
		assertNil(t, o.Fields())
		assertEqual(t, "", o.Name())
		assertNil(t, o.SamplingOpts())
		assertNil(t, o.RotationOpts())
	})
}

func TestOptions_ChainedSetters(t *testing.T) {
	t.Parallel()

	t.Run("multiple level changes", func(t *testing.T) {
		t.Parallel()
		o := &Options{}
		o.SetLevel(TraceLevel)
		o.SetLevel(DebugLevel)
		o.SetLevel(InfoLevel)
		assertEqual(t, InfoLevel, o.Level())
	})
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkWithLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := WithLevel(InfoLevel)
		o := &Options{}
		opt(o)
	}
}

func BenchmarkOptions_SetLevel(b *testing.B) {
	o := &Options{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.SetLevel(InfoLevel)
	}
}

func BenchmarkOptions_Level(b *testing.B) {
	o := &Options{level: InfoLevel}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = o.Level()
	}
}
