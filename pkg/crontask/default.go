package crontask

import (
	"fmt"
	"time"
)

// New constructs and returns a new Scheduler configured by the supplied
// SchedulerOptions.
//
// The scheduler is created in a stopped state; call Start to begin processing
// jobs.
//
// Example:
//
//	s, err := crontask.New(
//	    crontask.WithLocation(time.UTC),
//	    crontask.WithSeconds(),
//	)
func New(opts ...SchedulerOption) (*Scheduler, error) {
	cfg := schedulerConfig{
		loc: time.UTC,
	}
	for _, o := range opts {
		o(&cfg)
	}
	s := &Scheduler{
		cfg:      cfg,
		registry: newRegistry(),
		exec:     &executor{onError: cfg.onError},
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
	// Pre-close doneCh so that callers blocking on it before Start get an
	// immediate return.
	close(s.doneCh)
	return s, nil
}

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
