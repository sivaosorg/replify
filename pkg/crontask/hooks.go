package crontask

import (
	"context"
	"time"
)

// CustomHooks is a struct that implements the Hooks interface with custom
// implementations of the methods you care about. It embeds NoopHooks to
// provide default no-op implementations for the remaining methods.
//
// Example:
//
//	type MyHooks struct {
//		crontask.NoopHooks
//	}
//
//	func (h *MyHooks) OnFailure(_ context.Context, id string, _ time.Duration, err error) {
//		log.Printf("ALERT: job %s failed: %v", id, err)
//	}

// OnStart implements Hooks. It does nothing.
func (NoopHooks) OnStart(_ context.Context, _ string) {}

// OnSuccess implements Hooks. It does nothing.
func (NoopHooks) OnSuccess(_ context.Context, _ string, _ time.Duration) {}

// OnFailure implements Hooks. It does nothing.
func (NoopHooks) OnFailure(_ context.Context, _ string, _ time.Duration, _ error) {}

// OnComplete implements Hooks. It does nothing.
func (NoopHooks) OnComplete(_ context.Context, _ string, _ time.Duration) {}
