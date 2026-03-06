package slogger

import "os"

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

// WithRotation is a functional option that enables per-level file rotation.
//
// Parameters:
//   - `opts`: the rotation configuration
//
// Returns:
//
// a functional option that sets RotationOpts on the logger's Options.
func WithRotation(opts RotationOptions) func(*Options) {
return func(o *Options) {
o.RotationOpts = &opts
}
}
