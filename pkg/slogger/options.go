package slogger

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
