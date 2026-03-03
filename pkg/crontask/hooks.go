package crontask

import (
	"context"
	"time"
)

// Hooks is the interface that callers may implement to observe the lifecycle
// of a job invocation. All methods are optional — embed NoopHooks to satisfy
// the interface without implementing every method.
//
// Hook methods are called synchronously within the executor goroutine. Hooks
// must not block for long periods; spawn a goroutine if you need to perform
// expensive work (e.g. remote metric writes) without slowing the executor.
type Hooks interface {
	// OnStart is called immediately before the job function is invoked,
	// after jitter has been applied and after the execution context has been
	// derived. The jobID parameter identifies the job being dispatched.
	OnStart(ctx context.Context, jobID string)

	// OnSuccess is called when the job function returns nil. The duration
	// parameter is the wall-clock time of the invocation, excluding jitter.
	OnSuccess(ctx context.Context, jobID string, duration time.Duration)

	// OnFailure is called when the job function returns a non-nil error after
	// all retry attempts are exhausted. err is the final error.
	OnFailure(ctx context.Context, jobID string, duration time.Duration, err error)

	// OnComplete is called after OnSuccess or OnFailure and regardless of the
	// outcome. It is useful for releasing resources that were acquired in
	// OnStart.
	OnComplete(ctx context.Context, jobID string, duration time.Duration)
}

// NoopHooks is a zero-value implementation of Hooks whose methods all do
// nothing. Embed NoopHooks into your own struct to selectively override only
// the methods you care about.
//
// Example:
//
//	type MyHooks struct {
//	    crontask.NoopHooks
//	}
//
//	func (h *MyHooks) OnFailure(_ context.Context, id string, _ time.Duration, err error) {
//	    log.Printf("ALERT: job %s failed: %v", id, err)
//	}
type NoopHooks struct{}

// OnStart implements Hooks. It does nothing.
func (NoopHooks) OnStart(_ context.Context, _ string) {}

// OnSuccess implements Hooks. It does nothing.
func (NoopHooks) OnSuccess(_ context.Context, _ string, _ time.Duration) {}

// OnFailure implements Hooks. It does nothing.
func (NoopHooks) OnFailure(_ context.Context, _ string, _ time.Duration, _ error) {}

// OnComplete implements Hooks. It does nothing.
func (NoopHooks) OnComplete(_ context.Context, _ string, _ time.Duration) {}
