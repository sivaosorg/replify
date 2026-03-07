package slogger

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
