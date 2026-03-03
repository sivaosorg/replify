package crontask

import (
	"context"
	"math/rand"
	"time"
)

// executor is responsible for running a single job entry within its own
// goroutine. It enforces timeouts, applies jitter, runs retry loops, and
// dispatches hook calls.
type executor struct {
	onError func(id string, err error)
}

// dispatch launches the job entry in a new goroutine. The caller (the
// scheduler loop) never blocks on dispatch.
func (ex *executor) dispatch(e *entry, scheduledAt time.Time) {
	go ex.run(e, scheduledAt)
}

// run is the goroutine body for a single job invocation. It applies jitter,
// derives the execution context, runs the retry loop, calls hooks, and
// records the result.
func (ex *executor) run(e *entry, scheduledAt time.Time) {
	// Determine the base context. Fall back to Background when none is
	// configured.
	base := e.cfg.ctx
	if base == nil {
		base = context.Background()
	}

	// Apply jitter if configured.
	if e.cfg.jitter > 0 {
		delay := time.Duration(rand.Int63n(int64(e.cfg.jitter))) //nolint:gosec
		timer := time.NewTimer(delay)
		select {
		case <-base.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}

	// Derive the execution context (with optional timeout).
	ctx := base
	var cancel context.CancelFunc
	if e.cfg.timeout > 0 {
		ctx, cancel = context.WithTimeout(base, e.cfg.timeout)
		defer cancel()
	}

	// Dispatch hook.
	if e.cfg.hooks != nil {
		e.cfg.hooks.OnStart(ctx, e.id)
	}

	start := time.Now()
	var finalErr error

	// Retry loop.
	for attempt := 1; attempt <= e.cfg.maxRetries+1; attempt++ {
		// Check if the context is already done before each attempt.
		select {
		case <-ctx.Done():
			finalErr = ctx.Err()
			break
		default:
		}
		if finalErr != nil {
			break
		}

		err := e.fn(ctx)
		if err == nil {
			finalErr = nil
			break
		}
		finalErr = &JobError{JobID: e.id, Attempt: attempt, Err: err}

		if attempt <= e.cfg.maxRetries {
			// Compute backoff delay.
			if e.cfg.backoff != nil {
				delay := e.cfg.backoff(attempt)
				if delay > 0 {
					timer := time.NewTimer(delay)
					select {
					case <-ctx.Done():
						timer.Stop()
						finalErr = ctx.Err()
						goto done
					case <-timer.C:
					}
				}
			}
		}
	}

done:
	elapsed := time.Since(start)

	if finalErr == nil {
		if e.cfg.hooks != nil {
			e.cfg.hooks.OnSuccess(ctx, e.id, elapsed)
		}
	} else {
		if e.cfg.hooks != nil {
			e.cfg.hooks.OnFailure(ctx, e.id, elapsed, finalErr)
		}
		if ex.onError != nil {
			ex.onError(e.id, finalErr)
		}
	}

	if e.cfg.hooks != nil {
		e.cfg.hooks.OnComplete(ctx, e.id, elapsed)
	}

	recordResult(e, scheduledAt, finalErr)
}
