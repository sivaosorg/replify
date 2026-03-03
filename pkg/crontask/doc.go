// Package crontask provides a production-grade cron and task scheduling engine
// for the replify ecosystem. It is designed to be the canonical scheduling
// sub-package for long-running Go services that require reliable, expressive,
// and observable periodic job execution.
//
// # Design Goals
//
// crontask was built with three primary goals in mind:
//
//  1. Correctness — schedules must fire at the right time, every time, across
//     DST transitions, leap years, and timezone boundaries.
//  2. Expressiveness — operators support the full standard five-field cron
//     syntax, an optional leading seconds field, semantic aliases such as
//     @daily and @weekly, step-based interval expressions, and per-job
//     jitter to spread load across a fleet.
//  3. Observability — every registered job exposes metadata (last run, next
//     run, run count, last error) and the scheduler surface accepts hook
//     interfaces for pre/post-execution callbacks, success/failure
//     notifications, and metrics instrumentation.
//
// # Architectural Philosophy
//
// The package is structured into four distinct layers, each with a single
// responsibility:
//
//   - Expression layer (expression.go, parser.go) — converts a raw string such
//     as "0 9 * * 1-5" or "@weekdays" into a typed Schedule that can compute
//     the next activation time for any reference instant.
//
//   - Job layer (job.go) — holds the function to execute, its configuration
//     (retry policy, timeout, hooks, jitter), and live runtime statistics. The
//     in-memory registry is guarded by a read/write mutex so that the scheduler
//     loop and external callers can safely inspect or mutate the job list at
//     any time.
//
//   - Execution layer (executor.go) — wraps a job invocation with timeout
//     enforcement, retry-with-backoff, context propagation, and hook dispatch.
//     Every invocation runs in its own goroutine so that a slow or stuck job
//     never delays the scheduler tick.
//
//   - Scheduler layer (scheduler.go) — owns the main goroutine, advances a
//     monotonic clock, queries each registered job's next-fire time, and
//     dispatches due jobs through the executor. The scheduler is fully
//     concurrent-safe and supports graceful shutdown via Shutdown(ctx).
//
// # Comparison with robfig/cron and gronx
//
// robfig/cron is the de-facto standard cron library for Go. crontask adopts
// its scheduler-loop design (sub-second precision, heap-ordered next-fire
// times) but extends it with a richer job model (retry, backoff, hooks,
// jitter, per-job context) and replaces its terse API with idiomatic
// option-function constructors documented in the style used throughout
// replify.
//
// gronx is primarily an expression parser and evaluator. crontask borrows its
// flexible field-parsing ideas, the @alias vocabulary, and the concept of
// validating an expression without running a scheduler. Unlike gronx, crontask
// ships a complete execution engine, so consumers do not need a separate
// library.
//
// # Quick Start
//
//	s, err := crontask.New()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	s.Start()
//
//	_, err = s.Register("0 * * * *", func(ctx context.Context) error {
//	    fmt.Println("every hour:", time.Now())
//	    return nil
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Shut down cleanly after a signal.
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	s.Shutdown(ctx)
//
// # Thread Safety
//
// All exported methods on Scheduler are safe for concurrent use from multiple
// goroutines. The job registry, scheduler loop, and executor each acquire
// their own locks independently to minimise contention.
package crontask
