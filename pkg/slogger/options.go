package slogger

import (
	"io"
	"os"
)

// Options configures a Logger at construction time.
// Pass functional options to New to override the defaults.
type Options struct {
	// Level is the minimum level that will be logged.
	Level Level
	// Formatter serialises entries; defaults to TextFormatter on os.Stderr.
	Formatter Formatter
	// Output is the destination writer; defaults to os.Stderr.
	Output io.Writer
	// CallerReporter enables automatic source-location capture.
	CallerReporter bool
	// CallerSkip adds extra skip frames for library wrappers.
	CallerSkip int
	// Fields are attached to every entry emitted by the logger.
	Fields []Field
	// Name is a dot-separated logger identifier prepended to output.
	Name string
	// SamplingOpts, when non-nil, enables per-message rate limiting.
	SamplingOpts *SamplingOptions
}

// defaultOptions returns production-ready defaults:
// InfoLevel, TextFormatter writing to os.Stderr, no caller capture.
//
// Returns:
//
// an *Options pre-populated with sensible production settings.
func defaultOptions() *Options {
	out := os.Stderr
	return &Options{
		Level:     InfoLevel,
		Output:    out,
		Formatter: NewTextFormatter(out),
	}
}
