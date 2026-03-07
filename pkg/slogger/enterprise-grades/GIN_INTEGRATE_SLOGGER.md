# Integrating slogger with Gin

This guide explains how to configure **slogger** as the centralised logging
solution for a [Gin](https://github.com/gin-gonic/gin) web application,
covering global setup, request middleware, scoped service loggers, structured
request tracing, and a production-grade deployment strategy.

---

## Table of Contents

1. [Why Integrate slogger with Gin](#1-why-integrate-slogger-with-gin)
2. [Global Logger Setup](#2-global-logger-setup)
3. [Gin Middleware Logging](#3-gin-middleware-logging)
4. [Scoped Logger Per Service](#4-scoped-logger-per-service)
5. [Structured Request Logging](#5-structured-request-logging)
6. [Production Logging Strategy](#6-production-logging-strategy)

---

## 1. Why Integrate slogger with Gin

Gin ships with its own `gin.Logger()` middleware that writes plaintext lines to
`os.Stdout`. That is acceptable for toy projects, but production services need
more:

| Requirement | Gin default | slogger |
|---|---|---|
| Structured JSON fields | ✗ | ✓ |
| Per-request trace/span IDs | ✗ | ✓ |
| Context-propagated fields | ✗ | ✓ |
| Per-level file rotation | ✗ | ✓ |
| Hook-based alerting | ✗ | ✓ |
| Sampling for chatty routes | ✗ | ✓ |
| Zero external dependencies | N/A | ✓ |

Replacing Gin's default logger with slogger gives you:

- **Centralised logging** — one logger configuration for all components of the
  application, from Gin routing to database layers to background workers.
- **Request tracing** — every log line emitted during a request lifecycle
  carries the same `request_id`, `trace_id`, and `user_id` automatically,
  with no manual threading of the logger through function arguments.
- **Observability** — structured JSON output integrates directly with log
  aggregators (Loki, Elasticsearch, Datadog, Splunk) without custom parsers.

---

## 2. Global Logger Setup

### Application bootstrap

Create a single slogger instance during application initialisation and register
it as the global logger. Every package that calls `slogger.Info(...)` will use
this instance.

```go
// internal/log/log.go
package log

import (
    "os"
    "time"

    "github.com/sivaosorg/replify/pkg/slogger"
)

// Init configures the global slogger instance for the application.
// Call this once at the very start of main(), before any other package
// initialises its own logger.
func Init(env, version string) {
    var formatter slogger.Formatter
    if env == "production" {
        formatter = slogger.NewJSONFormatter().
            WithTimeKey("timestamp").
            WithLevelKey("severity")
    } else {
        formatter = slogger.NewTextFormatter(os.Stderr)
    }

    log := slogger.New(func(o *slogger.Options) {
        o.Level     = parseLevel(env)
        o.Formatter = formatter
        o.Output    = os.Stdout
        o.Fields    = []slogger.Field{
            slogger.String("service", "my-api"),
            slogger.String("version", version),
            slogger.String("env",     env),
        }
    })

    slogger.SetGlobalLogger(log)
}

func parseLevel(env string) slogger.Level {
    if env == "development" {
        return slogger.DebugLevel
    }
    return slogger.InfoLevel
}
```

### Wiring into main

```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"

    applog "myapp/internal/log"
    "myapp/internal/router"
)

func main() {
    applog.Init(os.Getenv("ENV"), version)

    r := gin.New() // Do NOT use gin.Default() — it adds Gin's own logger middleware.
    r.Use(router.SloggerMiddleware())
    router.Register(r)

    _ = r.Run(":8080")
}
```

> **Important:** use `gin.New()` instead of `gin.Default()` to prevent Gin's
> built-in logger middleware from duplicating log output.

---

## 3. Gin Middleware Logging

The middleware captures the five key observability signals for every HTTP
request: path, method, status code, client IP, and latency.

```go
// internal/router/middleware.go
package router

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sivaosorg/replify/pkg/slogger"
)

// SloggerMiddleware returns a Gin middleware that logs every request using
// slogger's global logger. It captures:
//   - request_id  (from X-Request-ID header or generated)
//   - method, path, query
//   - status code
//   - latency
//   - client IP
//   - any errors set on the Gin context
func SloggerMiddleware() gin.HandlerFunc {
    log := slogger.GetGlobalLogger().Named("http")

    return func(c *gin.Context) {
        start := time.Now()

        // Resolve or generate the request ID.
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        c.Header("X-Request-ID", requestID)

        // Build a request-scoped logger and inject fields into the context.
        reqLog := log.With(
            slogger.String("request_id", requestID),
            slogger.String("method",     c.Request.Method),
            slogger.String("path",       c.FullPath()),
            slogger.String("client_ip",  c.ClientIP()),
        )
        ctx := slogger.WithContextFields(c.Request.Context(),
            slogger.String("request_id", requestID),
        )
        c.Request = c.Request.WithContext(ctx)

        // Store the scoped logger in the Gin context for handler access.
        c.Set("logger", reqLog)

        reqLog.Debug("request started")

        // Process the request.
        c.Next()

        // Log the outcome.
        elapsed := time.Since(start)
        status  := c.Writer.Status()
        fields  := []slogger.Field{
            slogger.Int("status",       status),
            slogger.Duration("latency", elapsed),
            slogger.Int("bytes",        c.Writer.Size()),
        }

        if len(c.Errors) > 0 {
            fields = append(fields, slogger.String("errors", c.Errors.String()))
        }

        switch {
        case status >= http.StatusInternalServerError:
            reqLog.Error("request completed", fields...)
        case status >= http.StatusBadRequest:
            reqLog.Warn("request completed", fields...)
        default:
            reqLog.Info("request completed", fields...)
        }
    }
}

func generateRequestID() string {
    // Replace with a UUID library or crypto/rand implementation.
    return "req-" + time.Now().Format("20060102150405.000000000")
}
```

### Accessing the request logger in handlers

```go
func GetUserHandler(c *gin.Context) {
    logVal, _ := c.Get("logger")
    log, ok := logVal.(*slogger.Logger)
    if !ok {
        log = slogger.GetGlobalLogger()
    }

    userID := c.Param("id")
    log.Debug("fetching user", slogger.String("user_id", userID))

    user, err := userService.Find(c.Request.Context(), userID)
    if err != nil {
        log.Error("user not found",
            slogger.String("user_id", userID),
            slogger.Err(err),
        )
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    log.Info("user fetched", slogger.String("user_id", userID))
    c.JSON(http.StatusOK, user)
}
```

---

## 4. Scoped Logger Per Service

Each service layer creates a **named child logger** from the global logger.
Named loggers:

- Include the service name in every log line.
- Inherit the global formatter, output, and hooks.
- Can be given their own minimum level if needed.

```go
// internal/service/user_service.go
package service

import (
    "context"

    "github.com/sivaosorg/replify/pkg/slogger"
)

// UserService handles user business logic with structured logging.
type UserService struct {
    log  *slogger.Logger
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{
        log:  slogger.GetGlobalLogger().Named("user-service"),
        repo: repo,
    }
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
    s.log.WithContext(ctx).Info("creating user",
        slogger.String("email", req.Email),
    )

    user, err := s.repo.Create(ctx, req)
    if err != nil {
        s.log.WithContext(ctx).Error("failed to create user",
            slogger.String("email", req.Email),
            slogger.Err(err),
        )
        return nil, err
    }

    s.log.WithContext(ctx).Info("user created",
        slogger.String("user_id", user.ID),
    )
    return user, nil
}
```

**Text output example:**
```
2026-01-15T10:00:00Z INFO  [user-service] creating user email=alice@example.com request_id=req-001 trace_id=abc123
```

**JSON output example:**
```json
{"ts":"2026-01-15T10:00:00Z","level":"INFO","name":"user-service","msg":"creating user","email":"alice@example.com","request_id":"req-001","trace_id":"abc123"}
```

---

## 5. Structured Request Logging

Carry observability identifiers across the entire request lifecycle using
`context.Context`. This eliminates the need to pass a logger as a function
argument throughout your call stack.

### Injecting trace and span IDs

If your service participates in distributed tracing (OpenTelemetry, Jaeger,
Zipkin), extract the trace context from the incoming request and inject it into
the slogger context:

```go
func SloggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract from OpenTelemetry span (if present).
        span := trace.SpanFromContext(c.Request.Context())
        traceID := span.SpanContext().TraceID().String()
        spanID  := span.SpanContext().SpanID().String()

        ctx := slogger.WithContextFields(c.Request.Context(),
            slogger.String("trace_id",  traceID),
            slogger.String("span_id",   spanID),
            slogger.String("request_id", c.GetHeader("X-Request-ID")),
        )
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### Fields available in every log line

After the middleware runs, every `log.WithContext(ctx)` call anywhere in the
request's call stack will automatically include:

| Field | Example | Source |
|---|---|---|
| `request_id` | `"req-abc123"` | X-Request-ID header |
| `trace_id` | `"4bf92f3577b34da6..."` | OpenTelemetry span |
| `span_id` | `"00f067aa0ba902b7"` | OpenTelemetry span |
| `user_id` | `"u-42"` | Set after authentication |

Add the user ID after authentication:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, err := validateToken(c.GetHeader("Authorization"))
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Add user_id to the logging context.
        ctx := slogger.WithContextFields(c.Request.Context(),
            slogger.String("user_id", claims.UserID),
        )
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

---

## 6. Production Logging Strategy

### Log rotation for long-running servers

```go
log := slogger.New(
    slogger.WithRotation(slogger.RotationOptions{
        Dir:      "/var/log/my-api",
        MaxBytes: 200 * 1024 * 1024, // 200 MiB per level
        MaxAge:   12 * time.Hour,
        Compress: true,
    }),
    func(o *slogger.Options) {
        o.Level     = slogger.InfoLevel
        o.Formatter = slogger.NewJSONFormatter()
        o.Output    = os.Stdout // container stdout for Fluentd/Fluent Bit
        o.Fields    = []slogger.Field{
            slogger.String("service", "my-api"),
        }
    },
)
slogger.SetGlobalLogger(log)
```

The archive structure produced:
```
/var/log/my-api/
├── info.log              (active)
├── error.log             (active)
└── archived/
    └── 2026-01-15/
        ├── 20260115060000_info.zip
        └── 20260115060000_error.zip
```

### Log aggregation

For containerised deployments:

- Write JSON to **stdout** — the container runtime captures it.
- Use **Fluentd** or **Fluent Bit** as a DaemonSet to tail container logs and
  forward to Elasticsearch/Loki/Datadog.
- Each JSON field becomes a queryable index in Kibana or Grafana.

For VM-based deployments:

- Write JSON to rotating files under `/var/log/<service>/`.
- Use **Filebeat** to ship the files to Elasticsearch.
- Alternatively, use **Promtail** to tail and push to Loki.

### Alerting hook

```go
type PagerDutyHook struct {
    client *pagerduty.Client
}

func (h *PagerDutyHook) Levels() []slogger.Level {
    return []slogger.Level{slogger.ErrorLevel, slogger.FatalLevel}
}

func (h *PagerDutyHook) Fire(e *slogger.Entry) error {
    // Fire asynchronously to avoid blocking the log call.
    go func() {
        _ = h.client.CreateIncident(pagerduty.CreateIncidentOptions{
            Title:   e.Message(),
            Urgency: "high",
        })
    }()
    return nil
}

log.AddHook(&PagerDutyHook{client: pdClient})
```

### Sampling high-traffic routes

```go
// Suppress health-check noise: log first 1 per minute, drop the rest.
healthLog := slogger.New(func(o *slogger.Options) {
    *o = *baseOpts
    o.SamplingOpts = &slogger.SamplingOptions{
        First:      1,
        Period:     time.Minute,
        Thereafter: 0,
    }
})
```
