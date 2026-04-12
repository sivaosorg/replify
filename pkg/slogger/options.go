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
		o.level = level
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
		o.formatter = formatter
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
		o.output = output
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
		o.caller = enable
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
		o.callerSkip = skip
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
		o.AddFields(fields...)
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
		o.name = name
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
		o.samplingOpts = opts
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
func WithRotation(opts *RotationOptions) func(*Options) {
	return func(o *Options) {
		o.rotationOpts = opts
	}
}

// ///////////////////////////////////////////////////////////////////////////
// Options accessors
// ///////////////////////////////////////////////////////////////////////////

// Level returns the minimum log level.
//
// Returns:
//
// the Level configured in Options.
func (o *Options) Level() Level {
	if o == nil {
		return InfoLevel
	}
	return o.level
}

// SetLevel sets the minimum log level.
//
// Parameters:
//   - `level`: the minimum log level
func (o *Options) SetLevel(level Level) {
	if o == nil {
		return
	}
	o.level = level
}

// Formatter returns the formatter.
//
// Returns:
//
// the Formatter configured in Options.
func (o *Options) Formatter() Formatter {
	if o == nil {
		return nil
	}
	return o.formatter
}

// SetFormatter sets the formatter.
//
// Parameters:
//   - `formatter`: the formatter to use
func (o *Options) SetFormatter(formatter Formatter) {
	if o == nil {
		return
	}
	o.formatter = formatter
}

// Output returns the output writer.
//
// Returns:
//
// the io.Writer configured in Options.
func (o *Options) Output() io.Writer {
	if o == nil {
		return nil
	}
	return o.output
}

// SetOutput sets the output writer.
//
// Parameters:
//   - `output`: the output writer
func (o *Options) SetOutput(output io.Writer) {
	if o == nil {
		return
	}
	o.output = output
}

// IsCaller returns whether caller reporting is enabled.
//
// Returns:
//
// true if caller reporting is enabled.
func (o *Options) IsCaller() bool {
	if o == nil {
		return false
	}
	return o.caller
}

// SetCaller enables or disables caller reporting.
//
// Parameters:
//   - `enable`: whether to enable caller reporting
func (o *Options) SetCaller(enable bool) {
	if o == nil {
		return
	}
	o.caller = enable
}

// CallerSkip returns the caller skip count.
//
// Returns:
//
// the number of stack frames to skip.
func (o *Options) CallerSkip() int {
	if o == nil {
		return 0
	}
	return o.callerSkip
}

// SetCallerSkip sets the caller skip count.
//
// Parameters:
//   - `skip`: the number of stack frames to skip
func (o *Options) SetCallerSkip(skip int) {
	if o == nil {
		return
	}
	o.callerSkip = skip
}

// Fields returns a copy of the fields.
//
// Returns:
//
// a copy of the []Field slice configured in Options.
func (o *Options) Fields() []Field {
	if o == nil || o.fields == nil {
		return nil
	}
	result := make([]Field, len(o.fields))
	copy(result, o.fields)
	return result
}

// SetFields sets the fields.
//
// Parameters:
//   - `fields`: the fields to set
func (o *Options) SetFields(fields []Field) {
	if o == nil {
		return
	}
	if fields == nil {
		o.fields = nil
		return
	}
	o.fields = make([]Field, len(fields))
	copy(o.fields, fields)
}

// AddFields appends fields to the existing field list.
//
// Parameters:
//   - `fields`: the fields to append
func (o *Options) AddFields(fields ...Field) {
	if o == nil {
		return
	}
	o.fields = append(o.fields, fields...)
}

// Name returns the logger name.
//
// Returns:
//
// the name configured in Options.
func (o *Options) Name() string {
	if o == nil {
		return ""
	}
	return o.name
}

// SetName sets the logger name.
//
// Parameters:
//   - `name`: the name to set
func (o *Options) SetName(name string) {
	if o == nil {
		return
	}
	o.name = name
}

// SamplingOpts returns the sampling options.
//
// Returns:
//
// the *SamplingOptions configured in Options.
func (o *Options) SamplingOpts() *SamplingOptions {
	if o == nil {
		return nil
	}
	return o.samplingOpts
}

// SetSamplingOpts sets the sampling options.
//
// Parameters:
//   - `opts`: the sampling options
func (o *Options) SetSamplingOpts(opts *SamplingOptions) {
	if o == nil {
		return
	}
	o.samplingOpts = opts
}

// RotationOpts returns the rotation options.
//
// Returns:
//
// the *RotationOptions configured in Options.
func (o *Options) RotationOpts() *RotationOptions {
	if o == nil {
		return nil
	}
	return o.rotationOpts
}

// SetRotationOpts sets the rotation options.
//
// Parameters:
//   - `opts`: the rotation options
func (o *Options) SetRotationOpts(opts *RotationOptions) {
	if o == nil {
		return
	}
	o.rotationOpts = opts
}
