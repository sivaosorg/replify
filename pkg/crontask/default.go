package crontask

import "fmt"

// newRegistry allocates an initialised registry.
func newRegistry() *registry {
	return &registry{entries: make(map[string]*entry)}
}

// newExpressionError constructs an ExpressionError for the given expression
// and field index with a formatted reason string.
func newExpressionError(expr string, field int, format string, args ...any) error {
	return &ExpressionError{
		Expression: expr,
		Field:      field,
		Reason:     fmt.Sprintf(format, args...),
	}
}
