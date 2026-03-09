# Getting Started with replify

**replify** is a zero-dependency Go library that standardizes and simplifies API response wrapping for RESTful services. It leverages the Decorator (Builder) Pattern to let you compose error handling, metadata, pagination, debugging, and streaming capabilities in a clean, chainable API.

This guide covers everything you need to know to go from installation to advanced production usage across the entire **replify** ecosystem.

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Standard Response Format](#standard-response-format)
5. [Building Responses — Fluent Builder API](#building-responses--fluent-builder-api)
   - [Shortcut Constructors](#shortcut-constructors)
   - [Builder Methods Reference](#builder-methods-reference)
6. [Pagination](#pagination)
7. [Metadata (meta)](#metadata-meta)
8. [Response Headers (header)](#response-headers-header)
9. [Debugging Information](#debugging-information)
10. [Error Handling](#error-handling)
11. [Serialization & Deserialization](#serialization--deserialization)
12. [Inspecting a Response](#inspecting-a-response)
13. [Normalization](#normalization)
14. [Compression & Body Streaming](#compression--body-streaming)
15. [Streaming API (StreamingWrapper)](#streaming-api-streamingwrapper)
16. [JSON Body Querying with fj](#json-body-querying-with-fj)
17. [Sub-package Ecosystem](#sub-package-ecosystem)
    - [pkg/slogger — Structured Logging](#pkgslogger--structured-logging)
    - [pkg/crontask — Task Scheduling](#pkgcrontask--task-scheduling)
    - [pkg/fj — Fast JSON Path Engine](#pkgfj--fast-json-path-engine)
    - [pkg/coll — Collection Utilities](#pkgcoll--collection-utilities)
    - [pkg/conv — Type Conversion](#pkgconv--type-conversion)
    - [pkg/randn — Random Data Generation](#pkgrandn--random-data-generation)
    - [pkg/hashy — Deterministic Hashing](#pkghashy--deterministic-hashing)
    - [pkg/match — Wildcard Pattern Matching](#pkgmatch--wildcard-pattern-matching)
    - [pkg/strutil — String Utilities](#pkgstrutil--string-utilities)
    - [pkg/truncate — String Truncation](#pkgtruncate--string-truncation)
    - [pkg/sysx — System Utilities](#pkgsysx--system-utilities)
    - [pkg/netx — Network Subnetting](#pkgnetx--network-subnetting)
    - [pkg/encoding — JSON Encoding](#pkgencoding--json-encoding)
    - [pkg/common — Reflection Utilities](#pkgcommon--reflection-utilities)
    - [pkg/ref — Pointer Utilities](#pkgref--pointer-utilities)
    - [pkg/msort — Map Sorting](#pkgmsort--map-sorting)
    - [pkg/assert — Test Assertions](#pkgassert--test-assertions)
18. [Practical Examples](#practical-examples)
19. [Best Practices](#best-practices)
20. [HTTP Status Code Reference](#http-status-code-reference)
21. [Contributing](#contributing)

---

## Overview

Building RESTful APIs typically involves repetitive boilerplate: manually constructing JSON response objects, handling errors consistently, appending request IDs and version strings, calculating pagination metadata, and conditionally injecting debug info. **replify** eliminates all of it.

### Problems it solves

| Problem | replify solution |
|---------|-----------------|
| Inconsistent response formats across endpoints | One canonical `wrapper` struct serialized to a standard JSON shape |
| Repetitive error-handling boilerplate | `WithError`, `WithErrorAck`, `AppendError` fluent methods |
| Manual metadata management | `WithRequestID`, `WithApiVersion`, `WithLocale` auto-populated by defaults |
| Scattered pagination logic | Dedicated `pagination` type with auto-calculation |
| Hard-to-debug production issues | Conditional `WithDebuggingKV` information |
| Expensive JSON body querying | Built-in `fj` path engine — no full unmarshal needed |
| Large response streaming | First-class `StreamingWrapper` with strategies and progress tracking |

### Design principles

- **Zero external dependencies** — only the Go standard library is used
- **Immutable-friendly builder pattern** — every `With*` method returns the same pointer for chaining
- **Type-safe** — Go generics used throughout the sub-packages
- **Production ready** — safe for concurrent use, nil-safe guards on every method

---

## Installation

### Requirements

- Go **1.23** or higher

### Get the module

```bash
# Latest version
go get -u github.com/sivaosorg/replify@latest

# Pin to a specific version
go get github.com/sivaosorg/replify@v0.0.1
```

### Import

```go
import "github.com/sivaosorg/replify"
```

Sub-packages are imported individually when needed:

```go
import (
    "github.com/sivaosorg/replify"
    "github.com/sivaosorg/replify/pkg/fj"
    "github.com/sivaosorg/replify/pkg/randn"
    "github.com/sivaosorg/replify/pkg/slogger"
)
```

With Go module support, `go build`, `go run`, and `go test` will automatically resolve the necessary dependencies.

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify"
)

func main() {
    // Build a successful response
    response := replify.New().
        WithStatusCode(200).
        WithMessage("User retrieved successfully").
        WithBody(map[string]string{
            "id":   "usr-123",
            "name": "Alice",
        })

    fmt.Println(response.JSONPretty())
    // → pretty-printed JSON including status_code, message, data, meta, headers
}
```

**Output:**

```json
{
  "data": {
    "id": "usr-123",
    "name": "Alice"
  },
  "headers": {
    "code": 200,
    "text": "OK"
  },
  "message": "User retrieved successfully",
  "meta": {
    "api_version": "v0.0.1",
    "locale": "en_US",
    "request_id": "d7e5ce24b796da94770911db36565bf9",
    "requested_time": "2026-01-29T10:07:05.751501+07:00"
  },
  "status_code": 200,
  "total": 0
}
```

> **Note:** `meta.request_id` and `meta.requested_time` are populated automatically by `New()` with crypto-random values and the current time. You can override them with `WithRequestID()` and `WithRequestedTime()`.

---

## Standard Response Format

Every `wrapper` serializes to the following JSON schema:

```json
{
  "status_code": 200,
  "message":     "Resource retrieved successfully",
  "path":        "/api/v1/users",
  "total":       2,
  "data": [ /* any value — object, array, scalar */ ],
  "headers": {
    "code":        200,
    "text":        "OK",
    "type":        "success",
    "description": "The request has succeeded."
  },
  "meta": {
    "request_id":     "80eafc6a1655ec5a06595d155f1e6951",
    "api_version":    "v1.0.4",
    "locale":         "en_US",
    "requested_time": "2026-02-26T17:30:28.983Z",
    "delta_value":    0.0,
    "delta_cnt":      0,
    "custom_fields": {
      "trace_id":      "80eafc6a…",
      "origin_region": "us-east-1"
    }
  },
  "pagination": {
    "page":        1,
    "per_page":    20,
    "total_items": 120,
    "total_pages": 6,
    "is_last":     false
  },
  "debug": {
    "execution_time_ms": 42,
    "query":            "SELECT * FROM users"
  }
}
```

### Field descriptions

| Field | Type | Description |
|-------|------|-------------|
| `status_code` | `int` | HTTP status code |
| `message` | `string` | Human-readable outcome message |
| `path` | `string` | Request path (optional) |
| `total` | `int` | Item count for non-paginated lists |
| `data` | `any` | Primary response payload |
| `headers.code` | `int` | Application-level status code (mirrors `status_code` by default) |
| `headers.text` | `string` | Human-readable status label (e.g., `"OK"`) |
| `headers.type` | `string` | Category string (e.g., `"success"`, `"error"`) |
| `headers.description` | `string` | Detailed status description |
| `meta.request_id` | `string` | Unique request identifier for distributed tracing |
| `meta.api_version` | `string` | API version string |
| `meta.locale` | `string` | IETF locale (e.g., `"en_US"`) |
| `meta.requested_time` | `string` | ISO 8601 request timestamp |
| `meta.delta_value` | `float64` | Magnitude of payload transformation |
| `meta.delta_cnt` | `int` | Count of transformations performed |
| `meta.custom_fields` | `object` | Arbitrary key-value metadata |
| `pagination.*` | `object` | Pagination envelope (only present when set) |
| `debug` | `object` | Development-only debugging information |

---

## Building Responses — Fluent Builder API

### Shortcut Constructors

replify provides pre-built convenience functions for the most common HTTP status codes. Each accepts a `message` string and an optional `data` payload:

```go
replify.WrapOk("OK", data)                         // 200
replify.WrapCreated("Resource created", data)       // 201
replify.WrapAccepted("Processing started", data)    // 202
replify.WrapNoContent("Deleted", nil)               // 204
replify.WrapBadRequest("Invalid input", nil)        // 400
replify.WrapUnauthorized("Please sign in", nil)     // 401
replify.WrapForbidden("Access denied", nil)         // 403
replify.WrapNotFound("Resource not found", nil)     // 404
replify.WrapConflict("Duplicate entry", nil)        // 409
replify.WrapUnprocessableEntity("Validation failed", errs) // 422
replify.WrapTooManyRequest("Rate limit exceeded", nil)     // 429
replify.WrapInternalServerError("Unexpected error", nil)   // 500
replify.WrapBadGateway("Upstream error", nil)       // 502
replify.WrapServiceUnavailable("Maintenance", nil)  // 503
replify.WrapGatewayTimeout("Upstream timeout", nil) // 504
```

Full list of available shortcut constructors:

| Function | Status |
|----------|--------|
| `WrapOk` | 200 |
| `WrapCreated` | 201 |
| `WrapAccepted` | 202 |
| `WrapNoContent` | 204 |
| `WrapProcessing` | 102 |
| `WrapBadRequest` | 400 |
| `WrapUnauthorized` | 401 |
| `WrapPaymentRequired` | 402 |
| `WrapForbidden` | 403 |
| `WrapNotFound` | 404 |
| `WrapMethodNotAllowed` | 405 |
| `WrapRequestTimeout` | 408 |
| `WrapConflict` | 409 |
| `WrapGone` | 410 |
| `WrapPreconditionFailed` | 412 |
| `WrapRequestEntityTooLarge` | 413 |
| `WrapUnsupportedMediaType` | 415 |
| `WrapUnprocessableEntity` | 422 |
| `WrapLocked` | 423 |
| `WrapUpgradeRequired` | 426 |
| `WrapTooManyRequest` | 429 |
| `WrapInternalServerError` | 500 |
| `WrapNotImplemented` | 501 |
| `WrapBadGateway` | 502 |
| `WrapServiceUnavailable` | 503 |
| `WrapGatewayTimeout` | 504 |
| `WrapHTTPVersionNotSupported` | 505 |

You can also parse a raw `map[string]any` into a wrapper:

```go
w, err := replify.WrapFrom(myMap)
```

### Builder Methods Reference

All builder methods return the same `*wrapper`, so they can be chained freely.

#### Response configuration

| Method | Description |
|--------|-------------|
| `WithStatusCode(code int)` | HTTP status code (100–599) |
| `WithBody(v any)` | Primary data payload |
| `WithJSONBody(v any) (*wrapper, error)` | Set body as serialized JSON string (returns error on marshal failure) |
| `WithMessage(msg string)` | Human-readable outcome message |
| `WithMessagef(fmt string, args ...any)` | Formatted message |
| `WithPath(v string)` | Request path |
| `WithPathf(fmt string, args ...any)` | Formatted request path |
| `WithTotal(total int)` | Total item count for non-paginated lists |
| `WithError(msg string)` | Set error from a string message |
| `WithErrorf(fmt string, args ...any)` | Formatted error message |
| `WithErrorAck(err error)` | Attach existing error with stack trace |
| `WithErrorAckf(err error, fmt string, args ...any)` | Attach error with formatted context message and stack trace |
| `AppendError(err error, msg string)` | Wrap an error with a new message |
| `AppendErrorf(err error, fmt string, args ...any)` | Wrap an error with a formatted message |
| `AppendErrorAck(err error, msg string)` | Append error with stack trace and context message |
| `BindCause()` | Promote the wrapped error cause to the wrapper's error field |

#### Metadata methods

| Method | Description |
|--------|-------------|
| `WithRequestID(v string)` | Request/trace ID |
| `WithRequestIDf(fmt string, args ...any)` | Formatted request ID |
| `RandRequestID()` | Generate a new crypto-random request ID |
| `WithApiVersion(v string)` | API version string |
| `WithApiVersionf(fmt string, args ...any)` | Formatted API version |
| `WithLocale(v string)` | Locale code (e.g., `"en_US"`, `"vi_VN"`) |
| `WithRequestedTime(v time.Time)` | Override request timestamp |
| `WithCustomFieldKV(key string, value any)` | Add a single custom metadata field |
| `WithCustomFieldKVf(key, fmt string, args ...any)` | Formatted custom metadata field |
| `WithCustomFields(values map[string]any)` | Set multiple custom metadata fields at once |
| `WithMeta(v *meta)` | Replace the entire metadata object |
| `RandDeltaValue()` | Generate a random delta value |
| `IncreaseDeltaCnt()` | Increment the transformation count |
| `DecreaseDeltaCnt()` | Decrement the transformation count |

#### Pagination methods (on wrapper)

| Method | Description |
|--------|-------------|
| `WithPagination(v *pagination)` | Attach a pagination object |
| `WithPage(v int)` | Current page number (min 1) |
| `WithPerPage(v int)` | Items per page (min 1, default 10) |
| `WithTotalItems(v int)` | Total items available |
| `WithTotalPages(v int)` | Total pages count |
| `WithIsLast(v bool)` | Whether the current page is the last |

#### Header methods

| Method | Description |
|--------|-------------|
| `WithHeader(v *header)` | Replace the entire header object |

#### Debugging methods

| Method | Description |
|--------|-------------|
| `WithDebugging(v map[string]any)` | Set the entire debug map |
| `WithDebuggingKV(key string, value any)` | Add a single debug key-value pair |
| `WithDebuggingKVf(key, fmt string, args ...any)` | Add a formatted debug value |

#### Normalization and utility methods

| Method | Description |
|--------|-------------|
| `NormAll()` | Run all normalization passes |
| `NormPaging()` | Normalize pagination (ensure defaults) |
| `NormHSC()` | Sync header code/text with status code |
| `NormMeta()` | Ensure meta fields are populated |
| `NormBody()` | Normalize body (strip empty maps, etc.) |
| `NormMessage()` | Set default message from status code if empty |
| `NormDebug()` | Remove empty debug values |
| `Clone()` | Deep-copy the wrapper |
| `Reset()` | Reset all fields to zero values |
| `CompressSafe(threshold int)` | gzip-compress body if it exceeds threshold (bytes) |
| `DecompressSafe()` | Decompress a gzip-compressed body |
| `MustHash256() (string, *wrapper)` | SHA-256 hash of the serialized response |
| `Hash256() string` | SHA-256 hash (panics on error) |
| `MustHash() (uint64, *wrapper)` | FNV-64 hash |
| `Hash() uint64` | FNV-64 hash (panics on error) |

#### Streaming setup on wrapper

| Method | Description |
|--------|-------------|
| `WithStreaming(reader io.Reader, config *StreamConfig) *StreamingWrapper` | Create a `StreamingWrapper` from this wrapper with custom config |
| `AsStreaming(reader io.Reader) *StreamingWrapper` | Create a `StreamingWrapper` with default config |

---

## Pagination

Use the dedicated `pagination` builder to construct pagination metadata independently, then attach it to the wrapper.

### Standalone builder

```go
pagination := replify.Pages().
    WithPage(3).
    WithPerPage(20).
    WithTotalItems(250).
    WithTotalPages(13).
    WithIsLast(false)

response := replify.New().
    WithStatusCode(200).
    WithBody(users).
    WithPagination(pagination).
    WithTotal(len(users))
```

### Convenience constructor

```go
// FromPages(totalItems, perPage) — auto-calculates total pages
pagination := replify.FromPages(250, 20)
```

### Pagination accessors

| Method | Returns | Description |
|--------|---------|-------------|
| `Available()` | `bool` | Non-nil guard |
| `Page()` | `int` | Current page |
| `PerPage()` | `int` | Items per page |
| `TotalPages()` | `int` | Total page count |
| `TotalItems()` | `int` | Total item count |
| `IsLast()` | `bool` | Whether this is the last page |

### Calculating pagination in handlers

```go
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
    page    := queryInt(r, "page", 1)
    perPage := queryInt(r, "per_page", 20)

    users, total, err := db.FindUsers(page, perPage)
    if err != nil {
        respondJSON(w, replify.WrapInternalServerError(err.Error(), nil))
        return
    }

    totalPages := (total + perPage - 1) / perPage

    response := replify.New().
        WithStatusCode(200).
        WithBody(users).
        WithTotal(len(users)).
        WithPagination(
            replify.Pages().
                WithPage(page).
                WithPerPage(perPage).
                WithTotalItems(total).
                WithTotalPages(totalPages).
                WithIsLast(page >= totalPages),
        ).
        WithRequestID(r.Header.Get("X-Request-ID")).
        WithPath(r.URL.Path)

    respondJSON(w, response)
}
```

---

## Metadata (meta)

The `meta` struct carries request-level context that is automatically populated with defaults:

- `request_id` — crypto-random MD5-based hex string
- `requested_time` — `time.Now()` at wrapper creation
- `api_version` — `"v0.0.1"` (override per your service)
- `locale` — `"en_US"` (override per request)

### Standalone meta builder

```go
meta := replify.Meta().
    WithApiVersion("v1.2.3").
    WithRequestID("req-abc-123").
    WithLocale("vi_VN").
    WithRequestedTime(time.Now()).
    WithCustomFieldKV("tenant_id", "acme").
    WithCustomFieldKV("region", "us-east-1")

response := replify.New().
    WithStatusCode(200).
    WithMeta(meta).
    WithBody(data)
```

### Chaining directly on wrapper

```go
response := replify.New().
    WithStatusCode(200).
    WithApiVersion("v2.0.0").
    WithLocale("fr_FR").
    WithRequestID(r.Header.Get("X-Request-ID")).
    WithCustomFieldKV("correlation_id", correlationID).
    WithBody(payload)
```

### Meta accessors

| Method | Returns | Description |
|--------|---------|-------------|
| `Available()` | `bool` | Non-nil guard |
| `ApiVersion()` | `string` | API version |
| `RequestID()` | `string` | Request/trace ID |
| `Locale()` | `string` | Locale string |
| `RequestedTime()` | `time.Time` | Request timestamp |
| `CustomFields()` | `map[string]any` | All custom metadata |
| `OnCustom(key string)` | `any` | Get one custom field |
| `CustomString(key, default)` | `string` | Typed getter with fallback |
| `CustomBool(key, default)` | `bool` | Typed getter with fallback |
| `CustomInt(key, default)` | `int` | Typed getter with fallback |
| `IsApiVersionPresent()` | `bool` | Presence check |
| `IsRequestIDPresent()` | `bool` | Presence check |
| `IsLocalePresent()` | `bool` | Presence check |
| `IsRequestedTimePresent()` | `bool` | Presence check |
| `IsCustomPresent()` | `bool` | Any custom fields exist |
| `IsCustomKeyPresent(key)` | `bool` | Specific custom field exists |
| `JSONCustomFields()` | `string` | Custom fields as compact JSON |

---

## Response Headers (header)

The `header` struct carries an application-level status envelope that mirrors or extends the HTTP status code:

```go
h := replify.Header().
    WithCode(200).
    WithText("OK").
    WithType("success").
    WithDescription("The request has succeeded and the resource has been returned.")

response := replify.New().
    WithStatusCode(200).
    WithHeader(h).
    WithBody(data)
```

### Built-in predefined headers

replify ships with pre-configured `*header` constants for every standard status:

```go
// Available via the package-level variable:
replify.Processing           // 102
replify.Ok                   // 200
replify.Created              // 201
// … and every standard HTTP status code
```

Call `NormHSC()` to automatically sync the header with whatever `WithStatusCode()` was set:

```go
response := replify.New().
    WithStatusCode(404).
    NormHSC()
// → header is automatically set to the 404 "Not Found" preset
```

### Header accessors

| Method | Returns | Description |
|--------|---------|-------------|
| `Available()` | `bool` | Non-nil guard |
| `Code()` | `int` | Application code |
| `Text()` | `string` | Status label |
| `Type()` | `string` | Category string |
| `Description()` | `string` | Detailed description |
| `IsCodePresent()` | `bool` | Presence check |
| `IsTextPresent()` | `bool` | Presence check |
| `IsTypePresent()` | `bool` | Presence check |
| `IsDescriptionPresent()` | `bool` | Presence check |
| `Respond()` | `map[string]any` | Map with only present fields |

---

## Debugging Information

Debugging fields are arbitrary key-value pairs intended for development and troubleshooting. They should be conditionally included and **never** expose secrets in production.

```go
response := replify.New().
    WithStatusCode(500).
    WithErrorAck(err).
    WithMessage("Database query failed")

// Only in development / staging:
if os.Getenv("ENV") != "production" {
    response.
        WithDebuggingKV("sql_query", query).
        WithDebuggingKVf("trace_id", "trace:%s", traceID).
        WithDebuggingKV("duration_ms", elapsed.Milliseconds())
}
```

### Accessing debug values

```go
// Raw map
debugMap := response.Debugging()

// Single value
val := response.OnDebugging("duration_ms")

// Typed getters (all accept a key and a default value)
n      := response.DebuggingInt("duration_ms", 0)
s      := response.DebuggingString("trace_id", "")
b      := response.DebuggingBool("cached", false)
f      := response.DebuggingFloat64("score", 0.0)
t      := response.DebuggingTime("created_at", time.Time{})
dur    := response.DebuggingDuration("latency", 0)

// JSON path queries inside the debug map
n2     := response.JSONDebuggingInt("nested.count", 0)
s2     := response.JSONDebuggingString("db.host", "")

// Presence checks
response.IsDebuggingPresent()
response.IsDebuggingKeyPresent("sql_query")
```

---

## Error Handling

replify provides rich error construction utilities with stack traces.

### Package-level error constructors

```go
// New error with stack trace
err := replify.NewError("something went wrong")

// Formatted error with stack trace
err := replify.NewErrorf("failed to open %s: permission denied", path)

// Annotate an existing error with a stack trace
err := replify.NewErrorAck(originalErr)

// Annotate with context message + stack trace
err := replify.NewErrorAckf(originalErr, "while processing order %s", orderID)

// Append context message to an existing error
err := replify.AppendError(originalErr, "additional context")
err := replify.AppendErrorf(originalErr, "context for order %s", orderID)
```

### Attaching errors to a wrapper

```go
// From a string
response.WithError("User not found")

// Formatted
response.WithErrorf("User %d not found", userID)

// From an existing error (preserves stack trace)
response.WithErrorAck(err)

// From error + formatted message
response.WithErrorAckf(err, "loading profile for user %d", userID)

// Wrap with context and attach to wrapper
response.AppendError(err, "database read failed")
response.AppendErrorAck(err, "during order lookup")
```

### Reading errors from a wrapper

```go
// Error message string
msg := response.Error()

// Underlying cause
cause := response.Cause()

// Presence check
hasErr := response.IsErrorPresent()

// Combined: error exists OR status is 4xx/5xx
isErr := response.IsError()
```

---

## Serialization & Deserialization

### Serializing a wrapper to JSON

```go
response := replify.New().
    WithStatusCode(200).
    WithBody(data)

// Compact JSON string
compact := response.JSON()

// Pretty-printed JSON string
pretty := response.JSONPretty()

// Raw bytes
b := response.JSONBytes()

// map[string]any representation
m := response.Respond()

// R value type (useful as return value)
r := response.Reply()

// *R pointer
rPtr := response.ReplyPtr()
```

### Parsing JSON back to a wrapper

```go
// From a JSON string
w, err := replify.UnwrapJSON(jsonStr)
if err != nil {
    log.Fatalf("parse error: %v", err)
}

// From a map
w, err := replify.WrapFrom(dataMap)
```

> `UnwrapJSON` strips JS-style comments and trailing commas before parsing, making it lenient with hand-authored JSON.

---

## Inspecting a Response

### Status inspection

```go
w.Available()            // non-nil guard
w.StatusCode()           // int, e.g. 200
w.StatusText()           // "200 (OK)"
w.IsSuccess()            // 2xx
w.IsInformational()      // 1xx
w.IsRedirection()        // 3xx
w.IsClientError()        // 4xx
w.IsServerError()        // 5xx
w.IsError()              // error exists OR 4xx/5xx
```

### Data inspection

```go
w.Body()                 // any — primary payload
w.Message()              // string
w.Total()                // int
w.Error()                // string error message
w.Cause()                // underlying error
w.Meta()                 // *meta
w.Header()               // *header
w.Pagination()           // *pagination
w.Debugging()            // map[string]any
w.DeltaValue()           // float64
w.DeltaCnt()             // int
```

### Presence checks

```go
w.IsBodyPresent()
w.IsJSONBody()           // body is a valid JSON string
w.IsHeaderPresent()
w.IsMetaPresent()
w.IsPagingPresent()
w.IsLastPage()
w.IsErrorPresent()
w.IsTotalPresent()
w.IsStatusCodePresent()
w.IsDebuggingPresent()
w.IsDebuggingKeyPresent("key")
```

### Streaming body

```go
// Stream body as chunks
ch := response.Stream()
for chunk := range ch {
    // process []byte chunk
    _ = chunk
}
```

---

## Normalization

Normalization passes auto-correct or auto-fill fields based on what is already set. Call them after building the wrapper and before serializing.

```go
response := replify.New().
    WithStatusCode(201).
    WithBody(newResource).
    NormAll()
// → NormAll() calls NormHSC + NormMeta + NormPaging + NormBody + NormMessage + NormDebug
```

Individual normalization passes:

| Method | What it does |
|--------|--------------|
| `NormHSC()` | Sets `header` fields from `status_code` if header is not explicitly set |
| `NormMeta()` | Fills missing `meta` fields (locale, version, request_id, time) |
| `NormPaging()` | Sets pagination defaults (page=1, per_page=10) if pagination exists |
| `NormBody()` | Strips empty/nil body values |
| `NormMessage()` | Derives message from status code text if message is empty |
| `NormDebug()` | Removes zero-value entries from the debug map |
| `NormAll()` | Runs all of the above in order |

---

## Compression & Body Streaming

### In-place gzip compression

When the response body is large, compress it before serialization to reduce payload size:

```go
// Compresses the body in-place if its serialized size exceeds 4 KB
response.CompressSafe(4096)

// Decompress later (e.g., on the receiving side)
response.DecompressSafe()
```

`CompressSafe` automatically injects debug metadata with original and compressed sizes so you can observe the benefit.

---

## Streaming API (StreamingWrapper)

For large datasets or long-running transfers, use `StreamingWrapper` to stream data progressively while still returning a standard replify response when the operation completes.

### Core streaming types

| Type | Description |
|------|-------------|
| `StreamConfig` | Streaming configuration (chunk size, strategy, compression, timeouts) |
| `StreamingWrapper` | Wraps `*wrapper` + `io.Reader` with progress tracking |
| `StreamProgress` | Current progress: bytes transferred, percentage, ETA, transfer rate |
| `StreamingStats` | Post-stream statistics: compression ratio, average bandwidth, failed chunks |
| `StreamChunk` | Single data chunk with checksum and metadata |
| `StreamingMetadata` | Metadata extension for streaming context |
| `BufferPool` | Reusable byte buffer pool |

### Streaming strategies

| Constant | Description |
|----------|-------------|
| `StrategyDirect` | Write directly to the writer without buffering |
| `StrategyBuffered` | Buffer data in memory before writing (default) |
| `StrategyChunked` | Process data in fixed-size chunks with concurrent workers |

### Compression types

| Constant | Description |
|----------|-------------|
| `CompressNone` | No compression (default) |
| `CompressGzip` | gzip compression |
| `CompressFlate` | DEFLATE compression |

### Creating a streaming wrapper

```go
reader := bytes.NewReader(largePayload)

// From existing wrapper
sw := response.WithStreaming(reader, &replify.StreamConfig{
    ChunkSize:           65536,           // 64 KB chunks
    Strategy:            replify.StrategyChunked,
    Compression:         replify.CompressGzip,
    UseBufferPool:       true,
    MaxConcurrentChunks: 4,
    ReadTimeout:         30 * time.Second,
    WriteTimeout:        30 * time.Second,
    ThrottleRate:        0,               // unlimited
})

// Or with default config
sw := response.AsStreaming(reader)

// Or directly
sw := replify.NewStreaming(reader, replify.NewStreamConfig())
```

### Running the stream

```go
var writer bytes.Buffer

sw := replify.New().
    WithStatusCode(200).
    WithMessage("File stream started").
    AsStreaming(reader)

// Set destination writer
sw.WithWriter(&writer)

// Register a progress callback
sw.WithCallback(func(progress *replify.StreamProgress, err error) {
    if err != nil {
        log.Printf("streaming error: %v", err)
        return
    }
    log.Printf("progress: %d%%  transferred: %d bytes  rate: %d B/s",
        progress.Percentage,
        progress.TransferredBytes,
        progress.TransferRate,
    )
})

// Start streaming
result := sw.Start(context.Background())
if !result.IsSuccess() {
    log.Printf("streaming failed: %s", result.Error())
}
```

### Buffer pool

```go
// Create a pool of 8 buffers, each 64 KB
pool := replify.NewBufferPool(65536, 8)

buf := pool.Get()   // acquire buffer
defer pool.Put(buf) // return to pool
```

---

## JSON Body Querying with fj

When the response body is a JSON string (or any struct that serializes to JSON), replify exposes the full **fj** (Fast JSON) path engine directly on the wrapper — no separate import or full unmarshal required.

### Basic queries

```go
response := replify.New().
    WithStatusCode(200).
    WithBody(map[string]any{
        "user": map[string]any{
            "name":   "Alice",
            "role":   "admin",
            "active": true,
        },
        "items": []map[string]any{
            {"id": 1, "price": 9.99},
            {"id": 2, "price": 4.50},
        },
    })

// Single path query
name := response.QueryJSONBody("user.name").String()    // "Alice"
role := response.QueryJSONBody("user.role").String()    // "admin"
n    := response.QueryJSONBody("items.#").Int()         // 2 (array length)

// Multiple paths in one serialization pass
fields := response.QueryJSONBodyMulti("user.name", "user.role")
fmt.Println(fields[0].String(), fields[1].String())   // Alice admin

// Parse once, query many times
doc := response.JSONBodyParser()
fmt.Println(doc.Get("user.name").String())
fmt.Println(doc.Get("items.0.price").Float64())

// Validate
if !response.ValidJSONBody() {
    return errors.New("body is not valid JSON")
}
```

### Aggregate helpers

```go
total := response.SumJSONBody("items.#.price")      // 14.49
min, _ := response.MinJSONBody("items.#.price")     // 4.50
max, _ := response.MaxJSONBody("items.#.price")     // 9.99
avg, _ := response.AvgJSONBody("items.#.price")     // 7.245
count  := response.CountJSONBody("items")           // 2
```

### Search and scan helpers

```go
// Full-tree substring search across all leaf values
hits := response.SearchJSONBody("admin")

// Wildcard scan
hits = response.SearchJSONBodyMatch("err*")

// Find all values under specific key names
emails := response.SearchJSONBodyByKey("email")

// Find all values under keys matching a wildcard
hits = response.SearchJSONBodyByKeyPattern("user*")

// Substring/wildcard check at a specific path
response.JSONBodyContains("user.role", "admin")
response.JSONBodyContainsMatch("user.email", "*@example.com")

// Find dot-notation path where a value first appears
path := response.FindJSONBodyPath("alice@example.com")

// All paths where value matches a wildcard
paths := response.FindJSONBodyPathsMatch("err*")
```

### Data manipulation helpers

```go
import "github.com/sivaosorg/replify/pkg/fj"

// Filter array elements
active := response.FilterJSONBody("users", func(ctx fj.Context) bool {
    return ctx.Get("active").Bool()
})

// First match
admin := response.FirstJSONBody("users", func(ctx fj.Context) bool {
    return ctx.Get("role").String() == "admin"
})

// Deduplicate (first-occurrence order preserved)
tags := response.DistinctJSONBody("tags")

// Project fields from an array of objects
rows := response.PluckJSONBody("users", "id", "email")

// Group by key
byRole := response.GroupByJSONBody("users", "role")

// Sort by field (ascending = true)
sorted := response.SortJSONBody("products", "price", true)
```

### fj path syntax

```
user.name              — field access
roles.0                — array index
roles.#                — array length
roles.#.name           — collect field from every element
roles.#(role=="admin") — first element where role == "admin"
roles.#(role=="admin")# — all elements where role == "admin"
{id,name}              — multi-select → new object
[id,name]              — multi-select → new array
name.@uppercase        — built-in transformer
name.@replace:{"target":"foo","replacement":"bar"}
..title                — recursive descent
```

### fj transformers (selected)

| Path expression | Description |
|----------------|-------------|
| `field.@uppercase` / `@lower` | Case conversion |
| `field.@trim` | Strip whitespace |
| `field.@snakecase` / `@camelcase` / `@kebabcase` | Case style conversion |
| `field.@pretty` / `@minify` | JSON formatting |
| `field.@reverse` | Reverse array/string |
| `field.@flatten` | Flatten nested arrays |
| `field.@join` | Merge array-of-objects into one object |
| `field.@keys` / `@values` | Object key/value extraction |
| `array.@filter:{"key":"active","value":true}` | Array filtering |
| `array.@pluck:name` | Extract a field from each element |
| `array.@first` / `@last` | First/last array element |
| `array.@sum` / `@min` / `@max` / `@count` | Aggregation |
| `object.@project:{"pick":["id","name"]}` | Field projection + rename |
| `object.@default:{"role":"viewer"}` | Fill missing/null fields |
| `value.@coerce:{"to":"string"}` | Type coercion |

Transformers can be piped with `|`:

```go
// Filter active users, then count them
fj.Get(json, `users.@filter:{"key":"active","value":true}|@count`).Raw()

// Pluck names and reverse the array
fj.Get(json, `users.@pluck:name|@reverse`).Raw()
```

Register custom transformers in `init()`:

```go
func init() {
    fj.AddTransformer("redact", fj.TransformerFunc(func(json, arg string) string {
        return `"[REDACTED]"`
    }))
}

// Usage
fj.Get(json, "user.password.@redact").String() // "[REDACTED]"
```

---

## Sub-package Ecosystem

All sub-packages are part of the same module and share the zero-external-dependency philosophy.

### pkg/slogger — Structured Logging

A production-grade structured logging library built entirely on the Go standard library.

**Key features:**
- Typed field constructors (no `interface{}` boxing on common types)
- JSON and text formatters with ANSI colour support
- Level-based filtering with atomic level changes at runtime
- Size- and age-based log rotation with ZIP compression
- Sliding-window rate-limiting sampler
- Side-effect hooks per log level
- Child loggers via `With()` and `Named()` for per-request scoping
- Context-aware logging — fields stored in `context.Context`
- `sync.Pool`-based entry recycling for zero allocations on the hot path

```go
import "github.com/sivaosorg/replify/pkg/slogger"

// Create a logger
log := slogger.New(
    slogger.WithLevel(slogger.LevelInfo),
    slogger.WithFormatter(slogger.NewJSONFormatter()),
)

// Structured logging
log.Info("user created",
    slogger.String("user_id", "usr-123"),
    slogger.Int("age", 30),
    slogger.Bool("admin", false),
)

// Child logger scoped to a request
reqLog := log.With(
    slogger.String("request_id", requestID),
    slogger.String("path", r.URL.Path),
)
reqLog.Info("handling request")

// Context-aware logging
ctx := slogger.WithContext(r.Context(), slogger.String("trace_id", traceID))
log.InfoContext(ctx, "processing order")

// Log rotation
rotation := slogger.NewRotation(
    slogger.WithRotationMaxSize(100), // 100 MB
    slogger.WithRotationMaxAge(7),    // 7 days
)
log = slogger.New(slogger.WithRotation(rotation))
```

### pkg/crontask — Task Scheduling

A production-grade cron and task scheduling engine with a four-layer architecture.

**Key features:**
- Standard five-field cron expressions (`* * * * *`) plus business-friendly aliases
- Timezone-correct scheduling across DST transitions
- Concurrent-safe job registry
- Retry with configurable backoff
- Rich hooks: `OnStart`, `OnComplete`, `OnError`, `OnSkip`
- Per-job metadata: next run, last error, run count, duration stats

```go
import "github.com/sivaosorg/replify/pkg/crontask"

scheduler := crontask.New()

// Add a job — runs every 5 minutes
scheduler.Add("cleanup", "*/5 * * * *", func() {
    log.Println("running cleanup")
})

// Business alias
scheduler.Add("daily-report", "@daily", func() {
    generateReport()
})

// Job with timezone
scheduler.AddWithLocation("morning-sync", "0 9 * * 1-5",
    time.LoadLocation("America/New_York"),
    func() { syncData() },
)

scheduler.Start()
defer scheduler.Stop()
```

### pkg/fj — Fast JSON Path Engine

Zero-allocation JSON path extraction — query values from raw JSON strings without a full unmarshal.

```go
import "github.com/sivaosorg/replify/pkg/fj"

json := `{"user":{"name":"Alice","scores":[95,87,92]}}`

name   := fj.Get(json, "user.name").String()          // "Alice"
top    := fj.Get(json, "user.scores.0").Int()         // 95
count  := fj.Get(json, "user.scores.#").Int()         // 3

// Zero-copy from bytes
ctx := fj.GetBytes(rawBytes, "user.name")

// Multiple paths in one pass
results := fj.GetMulti(json, "user.name", "user.scores.#")

// Parse once, query many times
doc := fj.Parse(json)
fmt.Println(doc.Get("user.name").String())
```

### pkg/coll — Collection Utilities

Type-safe, generics-based functional collection utilities.

```go
import "github.com/sivaosorg/replify/pkg/coll"

numbers := []int{1, 2, 3, 4, 5, 6}

evens   := coll.Filter(numbers, func(n int) bool { return n%2 == 0 })
squared := coll.Map(numbers, func(n int) int { return n * n })
sum     := coll.Reduce(numbers, func(acc, n int) int { return acc + n }, 0)
found   := coll.Contains(numbers, 3)
unique  := coll.Unique([]int{1, 2, 2, 3, 3, 3})
chunks  := coll.Chunk(numbers, 2)     // [[1,2],[3,4],[5,6]]

// Maps — comparable key variants
keys    := coll.KeyComp(myMap)
merged  := coll.MergeComp(map1, map2)
picked  := coll.PickComp(myMap, "a", "b")
omitted := coll.OmitComp(myMap, "secret")
inverted := coll.InvertComp(myMap)
val := coll.GetOrDefault(myMap, "key", defaultVal)

// Sets (HashSet)
set := coll.NewHashSet(1, 2, 3)
set.Add(4)
set.Contains(2) // true
inter := set.Intersection(otherSet)
union := set.Union(otherSet)
diff  := set.Difference(otherSet)

// Type-safe HashMap
hm := coll.NewHashMap[string, int]()
hm.Put("a", 1)
v := hm.Get("a")
hm.Remove("a")

// Stacks (use a slice-based approach with coll.Filter/Map)
// Deep merge maps
merged := coll.DeepMerge(dst, src)
```

### pkg/conv — Type Conversion

Type-safe, generics-based conversions between Go primitives and collections.

```go
import "github.com/sivaosorg/replify/pkg/conv"

// Generic conversion
n, err := conv.To[int]("42")
f, err := conv.To[float64]("3.14")

// Must variants (panic on error)
n := conv.MustTo[int]("42")

// With default fallback
n := conv.ToOrDefault[int]("bad", 0)

// Slice conversion
ints, err := conv.Slice[int]([]string{"1", "2", "3"})

// Specific converters
i, err  := conv.ToInt("123")
b, err  := conv.ToBool("true")
f, err  := conv.ToFloat64("3.14")
t, err  := conv.ToTime("2006-01-02", "2024-03-15")

// Struct conversion
type Dst struct { Name string; Age int }
dst, err := conv.ToStruct[Dst](srcMap)
```

### pkg/randn — Random Data Generation

Cryptographically-backed random data generation utilities.

```go
import "github.com/sivaosorg/replify/pkg/randn"

// UUID (reads from /dev/urandom — Unix only)
uuid, err := randn.UUID()          // "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

// Alphanumeric ID (16 chars)
id := randn.RandID(16)             // "aB3dE5fG7hI9jK1m"

// Crypto-secure hex ID
secureID := randn.CryptoID()       // 32-char hex string

// Time-based sortable ID
tsID := randn.TimeID()

// Random integers
n := randn.RandIntr(1, 100)        // 1–100 inclusive
n64 := randn.RandInt64r(0, 1000)

// Random floats
f := randn.RandFloat64r(0.0, 1.0)

// Random bytes
b := randn.RandBytes(32)
```

### pkg/hashy — Deterministic Hashing

Hash any Go value — structs, slices, maps, primitives — to a consistent hash.

```go
import "github.com/sivaosorg/replify/pkg/hashy"

type User struct {
    ID    int    `hash:"id"`
    Name  string `hash:"name"`
    Token string `hash:"-"` // excluded from hashing
}

u := User{ID: 1, Name: "Alice", Token: "secret"}

// FNV-64 uint64 hash
h, err := hashy.Hash(u)

// Hex string
hex, err := hashy.HashHex(u)

// SHA-256 hex digest
sha, err := hashy.Hash256(u)

// Base-32 encoded
b32, err := hashy.HashBase32(u)

// Base-16 (hex) encoded
b16, err := hashy.HashBase16(u)

// Base-10 (decimal) string
b10, err := hashy.HashBase10(u)

// With configuration options
opts := hashy.NewOptions().
    WithIgnoreZeroValue(true).
    WithSlicesAsSets(true).
    WithTagName("json").
    WithZeroNil(true).
    Build()

h, err = hashy.HashValue(u, opts)
```

### pkg/match — Wildcard Pattern Matching

Fast wildcard (`*`, `?`) pattern matching with optional complexity limits.

```go
import "github.com/sivaosorg/replify/pkg/match"

match.Match("hello world", "hello*")      // true
match.Match("user-123", "user-?*")        // true
match.Match("admin", "ad?in")             // true

// With complexity limit (prevents ReDoS-style patterns)
ok, err := match.MatchSafe("large input", "***heavy***", 1000)

// Escape literal wildcards
pattern := match.Escape("file*.txt")     // "file\*.txt"

// Pattern bounds
min, max := match.Bounds("hell?")        // 4, 5
```

### pkg/strutil — String Utilities

100+ string manipulation, validation, and formatting functions.

```go
import "github.com/sivaosorg/replify/pkg/strutil"

// Validation
strutil.IsEmpty("")                    // true
strutil.IsBlank("   ")                // true
strutil.IsNumeric("12345")            // true
strutil.IsAlpha("hello")              // true
strutil.IsAlphanumeric("hello123")    // true

// Transformation
strutil.TitleCase("hello world")      // "Hello World"
strutil.Slugify("Hello World!")       // "hello-world"
strutil.Capitalize("hello")           // "Hello"
strutil.Reverse("hello")              // "olleh"
strutil.Repeat("ab", 3)              // "ababab"
strutil.Pad("hi", 5, " ")           // "  hi "

// Inspection
strutil.CountOccurrences("hello", "l")  // 2
strutil.HasPrefix("hello", "hel")       // true
strutil.ContainsAny("hello", "aeiou")   // true

// Hashing
sha := strutil.SHA256("my secret")
```

### pkg/truncate — String Truncation

Unicode-safe string truncation with multiple strategies and omission markers.

```go
import "github.com/sivaosorg/replify/pkg/truncate"

long := "The quick brown fox jumps over the lazy dog"

// End truncation (default ellipsis)
truncate.CutEllipsis{MaxLen: 20}.Truncate(long)      // "The quick brown fox…"

// Leading truncation (omit from start)
truncate.CutEllipsisLeading{MaxLen: 20}.Truncate(long) // "…r the lazy dog"

// Middle truncation
truncate.EllipsisMiddle{MaxLen: 20}.Truncate(long)   // "The quick …azy dog"

// Fluent builder
result := truncate.New().
    WithMaxLen(30).
    WithOmission("...").
    Truncate(long)
```

### pkg/sysx — System Utilities

OS detection, runtime introspection, environment management, process control, command execution, and file-system helpers.

```go
import "github.com/sivaosorg/replify/pkg/sysx"

// OS detection
sysx.IsLinux()
sysx.IsDarwin()
sysx.IsWindows()
sysx.OSVersion()

// Runtime info
sysx.Hostname()
sysx.PID()
sysx.GoVersion()
stats := sysx.MemStats()

// Environment variables
val := sysx.Getenv("DB_HOST", "localhost")
port := sysx.GetenvInt("DB_PORT", 5432)
debug := sysx.GetenvBool("DEBUG", false)
sysx.Hasenv("SECRET_KEY")

// Command execution
result, err := sysx.Run("git", "status")
out, err := sysx.Output("ls", "-la")

// Builder API
cmd := sysx.NewCommand("go", "test", "./...").
    WithWorkDir("/app").
    WithTimeout(60 * time.Second).
    WithEnv("GO_ENV", "test")
result, err := cmd.Execute()

// File helpers
sysx.FileExists("/etc/hosts")
sysx.DirExists("/tmp")
content, err := sysx.ReadFile("/path/to/file")
err = sysx.WriteFile("/path/to/file", data)
err = sysx.WriteFileAtomic("/path/to/file", data)
```

### pkg/netx — Network Subnetting

IPv4 and IPv6 CIDR parsing, FLSM, and VLSM subnetting toolkit.

```go
import "github.com/sivaosorg/replify/pkg/netx"

// Parse a CIDR
subnet, err := netx.ParseCIDR("10.0.0.0/24")
fmt.Println(subnet.Network())    // 10.0.0.0
fmt.Println(subnet.Broadcast())  // 10.0.0.255
fmt.Println(subnet.UsableHosts()) // 254

// FLSM — split into equal subnets
subnets, err := netx.FLSM("192.168.1.0/24", 4)
// → [192.168.1.0/26, 192.168.1.64/26, ...]

// VLSM — right-size subnets for host requirements
subnets, err := netx.VLSM("10.0.0.0/24", []int{100, 50, 10})
// → [10.0.0.0/25 (126 hosts), 10.0.0.128/26 (62 hosts), 10.0.0.192/28 (14 hosts)]

// Overlap check
netx.Overlaps("10.0.0.0/24", "10.0.0.128/25") // true
```

### pkg/encoding — JSON Encoding

Wrappers around `encoding/json` with pretty-printing, minification, and colorized terminal output.

```go
import "github.com/sivaosorg/replify/pkg/encoding"

data := map[string]any{"name": "Alice", "age": 30}

// Marshal to compact JSON string
s, err := encoding.Json(data)

// Marshal to pretty JSON string
pretty, err := encoding.JsonPretty(data)

// Unmarshal
var out map[string]any
err = encoding.UnJson(jsonStr, &out)

// Colorized terminal output
colored := encoding.JsonColor(data)
fmt.Println(colored)
```

### pkg/common — Reflection Utilities

Reflection-based utilities for working with unknown types.

```go
import "github.com/sivaosorg/replify/pkg/common"

// Check if value is nil (handles interface wrapping)
common.IsNil(val)

// Dereference pointer to underlying value
v := common.Deref(ptr)

// Generic reader utilities
common.ReadAll(r)
```

### pkg/ref — Pointer Utilities

Concise helpers for creating pointers to scalar values.

```go
import "github.com/sivaosorg/replify/pkg/ref"

// Create typed pointers without declaring variables
s := ref.String("hello")       // *string
n := ref.Int(42)               // *int
b := ref.Bool(true)            // *bool
f := ref.Float64(3.14)         // *float64
t := ref.Time(time.Now())      // *time.Time
```

### pkg/msort — Map Sorting

Sort maps by key or value.

```go
import "github.com/sivaosorg/replify/pkg/msort"

m := map[string]int{"c": 3, "a": 1, "b": 2}

// Sort by key (ascending)
sorted := msort.SortKey(m)

// Sort by key (descending)
sorted = msort.SortKeyDesc(m)

// Sort by value (ascending)
sorted = msort.SortValue(m)

// Sort by value (descending)
sorted = msort.SortValueDesc(m)

// Stable sort (preserves order of equal elements)
sorted = msort.SortKeyStable(m)

// Top N items after sorting
top3 := msort.SortValue(m).Top(3)

// Convert back to map or get keys/values
top3.ToMap()
top3.Keys()
top3.Values()

// Sort map[K]time.Time by time value
times := map[string]time.Time{"a": t1, "b": t2}
byTime := msort.SortTimeValue(times)
byTimeDesc := msort.SortTimeValueDesc(times)
```

### pkg/assert — Test Assertions

Lightweight test assertion helpers built on `*testing.T`.

```go
import "github.com/sivaosorg/replify/pkg/assert"

func TestMyFunc(t *testing.T) {
    result := myFunc(42)

    assert.AssertEqual(t, result, 84)
    assert.AssertNil(t, err)
    assert.AssertNotNil(t, result)
    assert.AssertTrue(t, result > 0)
    assert.AssertFalse(t, result < 0)
}
```

> All assertion functions call `t.Helper()` so failure messages point to the call site, and they use `t.Errorf` (not `t.Fatalf`) to allow multiple failures to be reported in a single test run.

---

## Practical Examples

### Example 1 — CRUD REST API with net/http

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/sivaosorg/replify"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    id := parseID(r)
    user, err := db.FindUser(id)
    if err != nil {
        respond(w, replify.WrapNotFound("User not found", nil).
            WithRequestID(r.Header.Get("X-Request-ID")))
        return
    }
    respond(w, replify.WrapOk("User retrieved successfully", user).
        WithRequestID(r.Header.Get("X-Request-ID")))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        respond(w, replify.WrapBadRequest("Invalid request body", nil))
        return
    }
    if err := validate(user); err != nil {
        respond(w, replify.WrapUnprocessableEntity(err.Error(), nil))
        return
    }
    created, err := db.CreateUser(user)
    if err != nil {
        respond(w, replify.WrapInternalServerError("Failed to create user", nil).
            WithErrorAck(err))
        return
    }
    respond(w, replify.WrapCreated("User created successfully", created))
}

func respond(w http.ResponseWriter, r interface {
    StatusCode() int
    JSON() string
}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(r.StatusCode())
    w.Write([]byte(r.JSON()))
}
```

### Example 2 — Paginated list endpoint

```go
func ListUsers(w http.ResponseWriter, r *http.Request) {
    page    := queryInt(r, "page", 1)
    perPage := queryInt(r, "per_page", 20)

    users, total, err := db.ListUsers(page, perPage)
    if err != nil {
        respond(w, replify.WrapInternalServerError("Failed to fetch users", nil).
            WithErrorAck(err))
        return
    }

    totalPages := (total + perPage - 1) / perPage

    respond(w, replify.New().
        WithStatusCode(200).
        WithBody(users).
        WithTotal(len(users)).
        WithPagination(
            replify.Pages().
                WithPage(page).
                WithPerPage(perPage).
                WithTotalItems(total).
                WithTotalPages(totalPages).
                WithIsLast(page >= totalPages),
        ).
        WithRequestID(r.Header.Get("X-Request-ID")).
        WithPath(r.URL.Path))
}
```

### Example 3 — Error handling with stack traces

```go
func ProcessOrder(w http.ResponseWriter, r *http.Request) {
    order, err := orderService.Process(r.Context())
    resp := replify.New()

    if err != nil {
        resp.
            WithStatusCode(500).
            WithErrorAck(err).
            WithMessage("Order processing failed")

        if os.Getenv("ENV") != "production" {
            resp.
                WithDebuggingKV("order_payload", order).
                WithDebuggingKVf("trace", "goroutine:%d", goid())
        }
    } else {
        resp.
            WithStatusCode(200).
            WithBody(order).
            WithMessage("Order processed successfully")
    }

    respond(w, resp)
}
```

### Example 4 — Parsing and querying a JSON body

```go
jsonStr := `{
    "data": {
        "users": [
            {"id": 1, "name": "Alice", "role": "admin",  "active": true},
            {"id": 2, "name": "Bob",   "role": "viewer", "active": false},
            {"id": 3, "name": "Carol", "role": "editor", "active": true}
        ]
    },
    "status_code": 200
}`

w, err := replify.UnwrapJSON(jsonStr)
if err != nil {
    log.Fatal(err)
}

// Query the body
adminName := w.QueryJSONBody("data.users.#(role==\"admin\").name").String()
fmt.Println(adminName) // "Alice"

// Filter active users
active := w.FilterJSONBody("data.users", func(ctx fj.Context) bool {
    return ctx.Get("active").Bool()
})
fmt.Println(len(active)) // 2

// Collect all names
names := w.PluckJSONBody("data.users", "name")
// → [{String:"Alice"}, {String:"Bob"}, {String:"Carol"}]
```

### Example 5 — Distributed tracing across microservices

```go
// Service A — outbound request
requestID := randn.CryptoID()
response := replify.New().
    WithStatusCode(200).
    WithRequestID(requestID).
    WithApiVersion("v2.0.0").
    WithCustomFieldKV("correlation_id", correlationID).
    WithCustomFieldKV("source_service", "user-service").
    WithBody(userProfile)

// Service B — receive and propagate
inbound, _ := replify.UnwrapJSON(responseBody)
outbound := replify.New().
    WithStatusCode(200).
    WithRequestID(inbound.Meta().RequestID()).       // preserve request ID
    WithCustomFields(inbound.Meta().CustomFields()). // propagate all custom fields
    WithBody(enrichedData)
```

### Example 6 — Conditional debug info (dev vs production)

```go
func handler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    data, err := fetchData(r.Context())
    elapsed := time.Since(start)

    resp := replify.New().
        WithStatusCode(200).
        WithBody(data).
        WithRequestID(r.Header.Get("X-Request-ID"))

    if os.Getenv("LOG_LEVEL") == "debug" {
        resp.
            WithDebuggingKV("execution_time_ms", elapsed.Milliseconds()).
            WithDebuggingKV("db_queries", queryCount).
            WithDebuggingKV("cache_hit", cacheHit)
    }

    if err != nil {
        resp.
            WithStatusCode(500).
            WithErrorAck(err)
    }

    respond(w, resp)
}
```

### Example 7 — Streaming a large file

```go
func StreamFileHandler(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open("/path/to/large-file.bin")
    if err != nil {
        respond(w, replify.WrapNotFound("File not found", nil))
        return
    }
    defer file.Close()

    sw := replify.New().
        WithStatusCode(200).
        WithMessage("File stream started").
        AsStreaming(file)

    sw.WithCallback(func(p *replify.StreamProgress, err error) {
        log.Printf("%.1f%% transferred (%d bytes, %d B/s)",
            float64(p.Percentage), p.TransferredBytes, p.TransferRate)
    })

    result := sw.Start(r.Context())
    respond(w, result)
}
```

### Example 8 — Using slogger with replify

```go
import (
    "github.com/sivaosorg/replify"
    "github.com/sivaosorg/replify/pkg/slogger"
)

var log = slogger.New(
    slogger.WithLevel(slogger.LevelInfo),
    slogger.WithFormatter(slogger.NewJSONFormatter()),
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
    reqLog := log.With(
        slogger.String("request_id", r.Header.Get("X-Request-ID")),
        slogger.String("path", r.URL.Path),
    )
    reqLog.Info("received create order request")

    order, err := processOrder(r)
    if err != nil {
        reqLog.Error("order processing failed", slogger.Error(err))
        respond(w, replify.WrapInternalServerError("Failed to create order", nil).
            WithErrorAck(err).
            WithRequestID(r.Header.Get("X-Request-ID")))
        return
    }

    reqLog.Info("order created", slogger.String("order_id", order.ID))
    respond(w, replify.WrapCreated("Order created", order).
        WithRequestID(r.Header.Get("X-Request-ID")))
}
```

---

## Best Practices

### ✅ Do

1. **Always set `WithStatusCode`** — it drives header derivation, condition checks (`IsSuccess`, etc.), and downstream consumers.

2. **Use `WithErrorAck` for real errors** — preserves the Go error chain and captures a stack trace.

3. **Set `WithRequestID` from upstream headers** — enables distributed tracing across services.

4. **Use `NormAll()` before serializing** when you want guaranteed consistency (headers, meta defaults, etc.).

5. **Gate debug info behind an env check** — never expose internal state in production responses.
   ```go
   if os.Getenv("ENV") != "production" {
       resp.WithDebuggingKV("sql", query)
   }
   ```

6. **Use `JSONBodyParser()` when querying the body multiple times** — it serializes the body once and lets you re-use the parsed document.

7. **Call `Available()` before accessing a wrapper received from external code** — defensively prevents nil panics.

8. **Use `FromPages` for simple pagination** — auto-calculates `total_pages`.

### ❌ Don't

1. **Don't skip status codes.**
   ```go
   // ❌
   replify.New().WithBody(data)
   // ✅
   replify.New().WithStatusCode(200).WithBody(data)
   ```

2. **Don't store secrets in debug or custom metadata fields.**

3. **Don't use generic error messages.**
   ```go
   // ❌
   .WithError("error occurred")
   // ✅
   .WithError("failed to load user: email already registered")
   ```

4. **Don't ignore `UnwrapJSON` errors.**
   ```go
   // ❌
   w, _ := replify.UnwrapJSON(jsonStr)
   // ✅
   w, err := replify.UnwrapJSON(jsonStr)
   if err != nil { ... }
   ```

5. **Don't mutate `fj.UnsafeBytes` output** — it shares memory with the source string.

6. **Don't register custom `fj` transformers after `init()`** — the registry is not concurrency-safe for writes.

---

## HTTP Status Code Reference

### Success (2xx)

| Code | Status | Typical use |
|------|--------|-------------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | POST — resource created |
| 202 | Accepted | Async processing started |
| 204 | No Content | Successful DELETE |
| 206 | Partial Content | Range requests, video streaming |

### Redirection (3xx)

| Code | Status | Typical use |
|------|--------|-------------|
| 301 | Moved Permanently | Permanent URL change |
| 302 | Found | Temporary redirect |
| 304 | Not Modified | Cached content still valid |
| 307 | Temporary Redirect | POST redirect preserving method |
| 308 | Permanent Redirect | POST redirect preserving method |

### Client Errors (4xx)

| Code | Status | Typical use |
|------|--------|-------------|
| 400 | Bad Request | Malformed request |
| 401 | Unauthorized | Missing/invalid auth |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource absent |
| 405 | Method Not Allowed | Wrong HTTP verb |
| 408 | Request Timeout | Client too slow |
| 409 | Conflict | Duplicate or version conflict |
| 410 | Gone | Resource permanently deleted |
| 412 | Precondition Failed | ETag mismatch |
| 413 | Payload Too Large | Request body too large |
| 415 | Unsupported Media Type | Wrong content type |
| 422 | Unprocessable Entity | Validation errors |
| 423 | Locked | Resource locked |
| 426 | Upgrade Required | Protocol upgrade needed |
| 429 | Too Many Requests | Rate limit exceeded |

### Server Errors (5xx)

| Code | Status | Typical use |
|------|--------|-------------|
| 500 | Internal Server Error | Unexpected failure |
| 501 | Not Implemented | Feature unavailable |
| 502 | Bad Gateway | Upstream error |
| 503 | Service Unavailable | Maintenance/overload |
| 504 | Gateway Timeout | Upstream timeout |
| 505 | HTTP Version Not Supported | Unsupported protocol version |

---

## Contributing

```bash
# 1. Clone
git clone --depth 1 https://github.com/sivaosorg/replify.git
cd replify

# 2. Install dependencies
go mod tidy

# 3. Run tests
go test ./...

# 4. Lint (optional, if golangci-lint is installed)
golangci-lint run ./...
```

- Follow standard Go formatting (`gofmt`).
- Add tests for every new exported function or behaviour.
- Update the relevant sub-package `README.md` if you change its API.
- Submit a pull request against `main` with a clear description of the change.

---

**Issues:** [github.com/sivaosorg/replify/issues](https://github.com/sivaosorg/replify/issues)  
**Discussions:** [github.com/sivaosorg/replify/discussions](https://github.com/sivaosorg/replify/discussions)  
**License:** MIT — see [LICENSE](../LICENSE)
