package slogger

import "context"

// WithContextFields returns a new context that carries the provided fields.
// Any fields already stored in ctx are preserved; the new fields are appended
// after the existing ones.
//
// Parameters:
//   - `ctx`: the parent context
//   - `fields`: the fields to attach to the new context
//
// Returns:
//
// a derived context containing the merged field slice.
func WithContextFields(ctx context.Context, fields ...Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}
	existing := FieldsFromContext(ctx)
	merged := make([]Field, 0, len(existing)+len(fields))
	merged = append(merged, existing...)
	merged = append(merged, fields...)
	return context.WithValue(ctx, contextKey{}, merged)
}

// FieldsFromContext extracts the log fields stored in ctx.
//
// Parameters:
//   - `ctx`: a context that may carry log fields
//
// Returns:
//
// the []Field stored in ctx, or nil when no fields are present.
func FieldsFromContext(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}
	v, _ := ctx.Value(contextKey{}).([]Field)
	return v
}
