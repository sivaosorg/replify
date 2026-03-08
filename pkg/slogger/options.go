package slogger

import "io"

// WithLevel is a functional option that sets the minimum log level.
//
// Parameters:
//   - `level`: the minimum log level
//
// Returns:
//
// a functional option that sets Level on the logger's Options.
func WithLevel(level Level) func(*Options) {
	return func(o *Options) {
		o.Level = level
	}
}

// WithFormatter is a functional option that sets the formatter.
//
// Parameters:
//   - `formatter`: the formatter to use
//
// Returns:
//
// a functional option that sets Formatter on the logger's Options.
func WithFormatter(formatter Formatter) func(*Options) {
	return func(o *Options) {
		o.Formatter = formatter
	}
}

// WithOutput is a functional option that sets the output writer.
//
// Parameters:
//   - `output`: the output writer to use
//
// Returns:
//
// a functional option that sets Output on the logger's Options.
func WithOutput(output io.Writer) func(*Options) {
	return func(o *Options) {
		o.Output = output
	}
}

// WithCaller is a functional option that enables caller reporting.
//
// Parameters:
//   - `enable`: whether to enable caller reporting
//
// Returns:
//
// a functional option that sets CallerReporter on the logger's Options.
func WithCaller(enable bool) func(*Options) {
	return func(o *Options) {
		o.CallerReporter = enable
	}
}

// WithCallerSkip is a functional option that sets the caller skip count.
//
// Parameters:
//   - `skip`: the number of stack frames to skip
//
// Returns:
//
// a functional option that sets CallerSkip on the logger's Options.
func WithCallerSkip(skip int) func(*Options) {
	return func(o *Options) {
		o.CallerSkip = skip
	}
}

// WithFields is a functional option that adds fields to the logger.
//
// Parameters:
//   - `fields`: the fields to add
//
// Returns:
//
// a functional option that adds Fields to the logger's Options.
func WithFields(fields ...Field) func(*Options) {
	return func(o *Options) {
		o.Fields = append(o.Fields, fields...)
	}
}

// WithName is a functional option that sets the logger name.
//
// Parameters:
//   - `name`: the name to set
//
// Returns:
//
// a functional option that sets Name on the logger's Options.
func WithName(name string) func(*Options) {
	return func(o *Options) {
		o.Name = name
	}
}

// WithSamplingOpts is a functional option that sets the sampling options.
//
// Parameters:
//   - `opts`: the sampling options to set
//
// Returns:
//
// a functional option that sets SamplingOpts on the logger's Options.
func WithSamplingOpts(opts *SamplingOptions) func(*Options) {
	return func(o *Options) {
		o.SamplingOpts = opts
	}
}

// WithRotation is a functional option that enables per-level file rotation.
// Pass a configured RotationOptions to control directory, size limits, age
// limits, and compression behaviour.
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
