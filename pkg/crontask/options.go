package crontask

import (
	"context"
	"time"
)

// WithLocation sets the default timezone for the scheduler. Jobs that do not
// carry their own timezone specifier will have their next-run times computed
// relative to loc.
//
// Example:
//
//	tz, _ := time.LoadLocation("Europe/Paris")
//	s, _ := crontask.New(crontask.WithLocation(tz))
func WithLocation(loc *time.Location) SchedulerOption {
	return func(c *schedulerConfig) {
		if loc != nil {
			c.loc = loc
		}
	}
}

// WithSeconds enables the six-field cron format where the first field
// represents seconds. When this option is set, Parse is called in six-field
// mode and the scheduler ticks at one-second granularity.
//
// Example:
//
//	s, _ := crontask.New(crontask.WithSeconds())
func WithSeconds() SchedulerOption {
	return func(c *schedulerConfig) {
		c.withSecs = true
	}
}

// WithErrorHandler registers a callback that is invoked synchronously by the
// executor whenever a job returns a non-nil error (after all retries are
// exhausted). The callback receives the job ID and the final error.
//
// Example:
//
//	s, _ := crontask.New(crontask.WithErrorHandler(func(id string, err error) {
//	    log.Printf("job %s failed: %v", id, err)
//	}))
func WithErrorHandler(fn func(id string, err error)) SchedulerOption {
	return func(c *schedulerConfig) {
		if fn != nil {
			c.onError = fn
		}
	}
}

// WithJobID sets an explicit, caller-supplied identifier for the job. If not
// provided, a random UUID-like identifier is generated automatically.
//
// Example:
//
//	s.Register("@hourly", fn, crontask.WithJobID("price-refresh"))
func WithJobID(id string) JobOption {
	return func(c *jobConfig) {
		c.id = id
	}
}

// WithJobName attaches a human-readable display name to the job. The name
// appears in the JobInfo returned by Jobs() and is useful for dashboards.
//
// Example:
//
//	s.Register("@daily", fn, crontask.WithJobName("Daily Report"))
func WithJobName(name string) JobOption {
	return func(c *jobConfig) {
		c.name = name
	}
}

// WithMaxRetries configures the number of times the executor will retry a
// failing job before recording a final error. A value of 0 (the default)
// means the job is attempted exactly once with no retries.
//
// Example:
//
//	s.Register("0 * * * *", fn, crontask.WithMaxRetries(3))
func WithMaxRetries(n int) JobOption {
	return func(c *jobConfig) {
		if n >= 0 {
			c.maxRetries = n
		}
	}
}

// WithBackoff sets the BackoffPolicy used between retry attempts. The default
// policy applies no delay between retries. Use ExponentialBackoff or
// ConstantBackoff for more controlled retry behaviour.
//
// Example:
//
//	s.Register("@hourly", fn,
//	    crontask.WithMaxRetries(3),
//	    crontask.WithBackoff(crontask.ExponentialBackoff(time.Second)),
//	)
func WithBackoff(p BackoffPolicy) JobOption {
	return func(c *jobConfig) {
		if p != nil {
			c.backoff = p
		}
	}
}

// WithTimeout sets a per-invocation execution deadline for the job. If a
// single execution does not complete within the specified duration, the
// context passed to the job function is cancelled and ErrJobTimeout is
// recorded.
//
// Example:
//
//	s.Register("@minutely", fn, crontask.WithTimeout(10*time.Second))
func WithTimeout(d time.Duration) JobOption {
	return func(c *jobConfig) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// WithJitter adds a random delay in the range [0, max) before each job
// execution. Jitter is useful in distributed systems where many nodes share
// the same schedule and simultaneous load spikes are undesirable.
//
// Example:
//
//	s.Register("@hourly", fn, crontask.WithJitter(30*time.Second))
func WithJitter(max time.Duration) JobOption {
	return func(c *jobConfig) {
		if max > 0 {
			c.jitter = max
		}
	}
}

// WithHooks attaches one or more Hooks implementations to the job. When
// multiple hooks are supplied they are composed into a chain that dispatches
// each hook method to all members in order. Nil values are silently ignored.
//
// Backward-compatible: existing callers that supply a single Hooks value
// continue to work without modification.
//
// Example:
//
//	s.Register("@daily", fn, crontask.WithHooks(myHooks))
//	s.Register("@daily", fn, crontask.WithHooks(crontask.LoggingHook(), crontask.MetricsHook()))
func WithHooks(hooks ...Hooks) JobOption {
	return func(c *jobConfig) {
		active := make([]Hooks, 0, len(hooks))
		for _, h := range hooks {
			if h != nil {
				active = append(active, h)
			}
		}
		if len(active) == 0 {
			return
		}
		if len(active) == 1 {
			c.hooks = active[0]
			return
		}
		c.hooks = ChainHooks(active...)
	}
}

// WithContext associates a base context with the job. The executor derives a
// child context from this base for each invocation, allowing per-job
// cancellation or value propagation.
//
// Example:
//
//	s.Register("@daily", fn, crontask.WithContext(reqCtx))
func WithContext(ctx context.Context) JobOption {
	return func(c *jobConfig) {
		if ctx != nil {
			c.ctx = ctx
		}
	}
}

// WithSchedulerHooks registers one or more default Hooks that are applied to
// every job that does not supply its own per-job hooks via WithHooks. This
// provides a convenient way to enable logging, metrics, or other observability
// for all jobs in one place at scheduler construction time.
//
// When multiple hooks are supplied they are composed into a chain (equivalent
// to calling ChainHooks). Nil values are silently ignored.
//
// Per-job hooks set via WithHooks always take precedence over the scheduler-
// level default. To combine both, call ChainHooks explicitly:
//
//	s.Register("@daily", fn, crontask.WithHooks(
//	    crontask.ChainHooks(schedulerDefaultHooks, myJobSpecificHook),
//	))
//
// Example:
//
//	s, _ := crontask.New(
//	    crontask.WithSchedulerHooks(
//	        crontask.LoggingHook(),
//	        crontask.MetricsHook(),
//	    ),
//	)
func WithSchedulerHooks(hooks ...Hooks) SchedulerOption {
	return func(c *schedulerConfig) {
		active := make([]Hooks, 0, len(hooks))
		for _, h := range hooks {
			if h != nil {
				active = append(active, h)
			}
		}
		if len(active) == 0 {
			return
		}
		c.defaultHooks = ChainHooks(active...)
	}
}

// ConstantBackoff returns a BackoffPolicy that waits the same fixed delay
// between every retry attempt.
//
// Example:
//
//	crontask.ConstantBackoff(5 * time.Second)
func ConstantBackoff(delay time.Duration) BackoffPolicy {
	return func(_ int) time.Duration { return delay }
}

// ExponentialBackoff returns a BackoffPolicy that doubles the base delay on
// each successive attempt (base, 2×base, 4×base, …).
//
// Example:
//
//	crontask.ExponentialBackoff(time.Second)
func ExponentialBackoff(base time.Duration) BackoffPolicy {
	return func(attempt int) time.Duration {
		d := base
		for i := 1; i < attempt; i++ {
			d *= 2
		}
		return d
	}
}
