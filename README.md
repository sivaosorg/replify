# replify

**replify** is a Go library designed to simplify and standardize API response wrapping for RESTful services. It leverages the Decorator Pattern to dynamically add error handling, metadata, pagination, and other response features in a clean and human-readable format.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Overview

Building RESTful APIs often requires repetitive boilerplate code for standardizing responses. **replify** eliminates this by providing a fluent, chainable API that ensures consistent response formats across all your endpoints.

### What Problems Does It Solve?

- ❌ **Inconsistent response formats** across different endpoints
- ❌ **Repetitive error handling** boilerplate in every handler
- ❌ **Manual metadata management** (request IDs, timestamps, versions)
- ❌ **Complex pagination logic** scattered throughout the codebase
- ❌ **Debugging difficulties** in production vs development environments

### The Solution

✅ **Standardized response structure** - One format for all endpoints  
✅ **Fluent API** - Chainable methods for building responses  
✅ **Built-in pagination** - Complete pagination support out of the box  
✅ **Metadata management** - Request IDs, timestamps, API versions, locales  
✅ **Conditional debugging** - Development-only debug information  
✅ **Error handling** - Stack traces, error wrapping, contextual messages  
✅ **Type safety** - Full type safety with Go generics  
✅ **Zero dependencies** - Only uses Go standard library

## Features

### Core Capabilities

- 🎯 **Standardized JSON Format** - Consistent structure across all API responses
- 🔗 **Fluent Builder Pattern** - Chain methods to construct complex responses
- 📄 **Pagination Support** - Built-in page, per_page, total_items, total_pages, is_last
- 🔍 **Request Tracing** - Track requests with unique IDs across microservices
- 🌍 **Internationalization** - Locale support for multi-language APIs
- 🐛 **Debug Mode** - Conditional debugging information for development
- ⚡ **Error Handling** - Rich error information with stack traces
- 📊 **Metadata** - API version, custom fields, timestamps
- ✅ **Status Helpers** - IsSuccess(), IsClientError(), IsServerError()
- 🔄 **JSON Parsing** - Parse JSON strings back to wrapper objects

## Requirements

- Go version 1.23 or higher

## Installation

### Install Package

```bash
# Latest version
go get -u github.com/sivaosorg/replify@latest

# Specific version
go get github.com/sivaosorg/replify@v0.0.1
```

### Import in Code

```go
import "github.com/sivaosorg/replify"
```

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import.

## Quick Start

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify"
)

func main() {
    // Create a simple success response
    response := replify.New().
        WithStatusCode(200).
        WithMessage("User retrieved successfully").
        WithBody(map[string]string{
            "id":   "123",
            "name": "John Doe",
        })
    
    fmt.Println(response.JsonPretty())
}
```

**Output:**
```json
{
    "data": {
        "id": "123",
        "name": "John Doe"
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

## Standard Response Format

The library produces responses in this standardized format:

```json
{
  "data": "response body here",
  "status_code": 200,
  "message": "How are you? I'm good",
  "total": 1,
  "path": "/api/v1/users",
  "meta": {
    "request_id": "80eafc6a1655ec5a06595d155f1e6951",
    "api_version": "v0.0.1",
    "locale": "en_US",
    "requested_time": "2024-12-14T20:24:23.983839+07:00",
    "custom_fields": {
      "fields": "userID: 103"
    }
  },
  "pagination": {
    "page": 1000,
    "per_page": 2,
    "total_items": 120,
    "total_pages": 34,
    "is_last": true
  },
  "debug": {
    "___abc": "trace sessions_id: 4919e84fc26881e9fe790f5d07465db4",
    "refer": 1234
  }
}
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `data` | `interface{}` | The primary data payload of the response |
| `status_code` | `int` | HTTP status code for the response |
| `message` | `string` | Human-readable message providing context |
| `total` | `int` | Total number of items (used in non-paginated responses) |
| `path` | `string` | Request path for which the response is generated |
| `meta` | `object` | Metadata about the API response |
| `meta.request_id` | `string` | Unique identifier for the request, useful for debugging |
| `meta.api_version` | `string` | API version used for the request |
| `meta.locale` | `string` | Locale used for the request (e.g., "en_US") |
| `meta.requested_time` | `string` | Timestamp when the request was made (ISO 8601) |
| `meta.custom_fields` | `object` | Additional custom metadata fields |
| `pagination` | `object` | Pagination details, if applicable |
| `pagination.page` | `int` | Current page number |
| `pagination.per_page` | `int` | Number of items per page |
| `pagination.total_items` | `int` | Total number of items available |
| `pagination.total_pages` | `int` | Total number of pages |
| `pagination.is_last` | `bool` | Indicates whether this is the last page |
| `debug` | `object` | Debugging information (useful for development) |

## Usage

### 1. Creating Basic Responses

#### Success Response

```go
response := replify.New().
    WithStatusCode(200).
    WithMessage("Operation successful").
    WithBody(data)
```

#### Error Response

```go
response := replify.New().
    WithStatusCode(400).
    WithError("Invalid input: email is required").
    WithMessage("Validation failed")
```

#### Response with Metadata

```go
response := replify.New().
    WithStatusCode(200).
    WithBody(users).
    WithRequestID("req-123-456").
    WithApiVersion("v1.0.0").
    WithLocale("en_US").
    WithPath("/api/v1/users")
```

### 2. Pagination

#### Creating Pagination

```go
pagination := replify.Pages().
    WithPage(1).
    WithPerPage(20).
    WithTotalItems(150).
    WithTotalPages(8).
    WithIsLast(false)

response := replify.New().
    WithStatusCode(200).
    WithBody(users).
    WithPagination(pagination).
    WithTotal(20)
```

### 3. Debugging Information

```go
response := replify.New().
    WithStatusCode(500).
    WithError("Database connection failed").
    WithDebuggingKV("query", "SELECT * FROM users").
    WithDebuggingKV("error_code", "CONN_TIMEOUT").
    WithDebuggingKV("retry_count", 3)
```

### 4. Complete Example

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify"
    "github.com/sivaosorg/replify/pkg/randn"
)

func main() {
    // Create pagination
    p := replify.Pages().
        WithIsLast(true).
        WithPage(1000).
        WithTotalItems(120).
        WithTotalPages(34).
        WithPerPage(2)
    
    // Create response
    w := replify.New().
        WithStatusCode(200).
        WithTotal(1).
        WithMessagef("How are you? %v", "I'm good").
        WithDebuggingKV("refer", 1234).
        WithDebuggingKVf("___abc", "trace sessions_id: %v", randn.CryptoID()).
        WithBody("response body here").
        WithPath("/api/v1/users").
        WithCustomFieldKVf("fields", "userID: %v", 103).
        WithPagination(p)
    
    if !w.Available() {
        return
    }
    
    // Access response properties
    fmt.Println(w.Json())
    fmt.Println(w.StatusCode())
    fmt.Println(w.StatusText())
    fmt.Println(w.Message())
    fmt.Println(w.Body())
    fmt.Println(w.IsSuccess())
    fmt.Println(w.Respond())
    
    // Check metadata
    fmt.Println(w.Meta().IsCustomPresent())
    fmt.Println(w.Meta().IsApiVersionPresent())
    fmt.Println(w.Meta().IsRequestIDPresent())
    fmt.Println(w.Meta().IsRequestedTimePresent())
}
```

### 5. Parsing JSON to Response

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/sivaosorg/replify"
)

func main() {
    jsonStr := `{
        "data": "response body here",
        "debug": {
          "___abc": "trace sessions_id: 4919e84fc26881e9fe790f5d07465db4",
          "refer": 1234
        },
        "message": "How do you do? I'm good",
        "meta": {
          "api_version": "v0.0.1",
          "custom_fields": {
            "fields": "userID: 103"
          },
          "locale": "en_US",
          "request_id": "80eafc6a1655ec5a06595d155f1e6951",
          "requested_time": "2024-12-14T20:24:23.983839+07:00"
        },
        "pagination": {
          "is_last": true,
          "page": 1000,
          "per_page": 2,
          "total_items": 120,
          "total_pages": 34
        },
        "path": "/api/v1/users",
        "status_code": 200,
        "total": 1
    }`
    
    t := time.Now()
    w, err := replify.UnwrapJSON(jsonStr)
    diff := time.Since(t)
    
    if err != nil {
        log.Fatalf("Error parsing JSON: %v", err)
    }
    
    fmt.Printf("Exe time: %+v\n", diff.String())
    fmt.Printf("%+v\n", w.OnDebugging("___abc"))
    fmt.Printf("%+v\n", w.JsonPretty())
}
```

## Practical Examples

### Example 1: RESTful CRUD API

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

// GET /users/:id
func GetUser(w http.ResponseWriter, r *http.Request) {
    id := getIDFromPath(r)
    user, err := findUserByID(id)
    
    var response *replify.R
    if err != nil {
        response = replify.New().
            WithStatusCode(404).
            WithError(err.Error()).
            WithMessage("User not found").
            WithRequestID(r.Header.Get("X-Request-ID"))
    } else {
        response = replify.New().
            WithStatusCode(200).
            WithBody(user).
            WithMessage("User retrieved successfully").
            WithRequestID(r.Header.Get("X-Request-ID"))
    }
    
    respondJSON(w, response)
}

// POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        response := replify.New().
            WithStatusCode(400).
            WithError(err.Error()).
            WithMessage("Invalid request body")
        respondJSON(w, response)
        return
    }
    
    if err := validateUser(user); err != nil {
        response := replify.New().
            WithStatusCode(422).
            WithError(err.Error()).
            WithMessage("Validation failed")
        respondJSON(w, response)
        return
    }
    
    createdUser, err := createUser(user)
    if err != nil {
        response := replify.New().
            WithStatusCode(500).
            WithErrorAck(err).
            WithMessage("Failed to create user")
        respondJSON(w, response)
        return
    }
    
    response := replify.New().
        WithStatusCode(201).
        WithBody(createdUser).
        WithMessage("User created successfully")
    respondJSON(w, response)
}

func respondJSON(w http.ResponseWriter, response *replify.R) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(response.StatusCode())
    w.Write([]byte(response.Json()))
}
```

### Example 2: Paginated List API

```go
func ListUsers(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    page := getQueryInt(r, "page", 1)
    perPage := getQueryInt(r, "per_page", 10)
    search := r.URL.Query().Get("search")
    
    // Fetch users with pagination
    users, total, err := db.FindUsers(search, page, perPage)
    if err != nil {
        response := replify.New().
            WithStatusCode(500).
            WithErrorAck(err).
            WithMessage("Failed to fetch users").
            WithDebuggingKV("search", search).
            WithDebuggingKV("page", page)
        respondJSON(w, response)
        return
    }
    
    // Calculate pagination metadata
    totalPages := (total + perPage - 1) / perPage
    isLast := page >= totalPages
    
    pagination := replify.Pages().
        WithPage(page).
        WithPerPage(perPage).
        WithTotalItems(total).
        WithTotalPages(totalPages).
        WithIsLast(isLast)
    
    response := replify.New().
        WithStatusCode(200).
        WithBody(users).
        WithPagination(pagination).
        WithTotal(len(users)).
        WithMessage("Users retrieved successfully").
        WithPath(r.URL.Path).
        WithRequestID(r.Header.Get("X-Request-ID"))
    
    respondJSON(w, response)
}
```

### Example 3: Error Handling with Stack Traces

```go
func ProcessOrder(w http.ResponseWriter, r *http.Request) {
    order, err := processOrderLogic(r)
    
    response := replify.New()
    
    if err != nil {
        response.
            WithStatusCode(500).
            WithErrorAck(err).
            WithMessage("Order processing failed")
        
        // Add debug info in development
        if os.Getenv("ENV") == "development" {
            response.
                WithDebuggingKV("timestamp", time.Now()).
                WithDebuggingKV("stack_trace", err.Error()).
                WithDebuggingKV("order_data", order)
        }
    } else {
        response.
            WithStatusCode(200).
            WithBody(order).
            WithMessage("Order processed successfully")
    }
    
    respondJSON(w, response)
}
```

## API Reference

### Wrapper Type (R)

```go
type R struct {
    *wrapper
}
```

The `R` type is a high-level abstraction providing a simplified interface for handling API responses.

### Core Functions

| Function | Description |
|----------|-------------|
| `New() *wrapper` | Creates a new response wrapper |
| `Pages() *pagination` | Creates a new pagination object |
| `UnwrapJSON(jsonStr string) (*wrapper, error)` | Parses JSON string to wrapper |

### Configuration Methods

#### Response Configuration

| Method | Description |
|--------|-------------|
| `WithStatusCode(code int)` | Sets HTTP status code |
| `WithBody(v interface{})` | Sets response body/data |
| `WithMessage(message string)` | Sets response message |
| `WithMessagef(format string, args...)` | Sets formatted message |
| `WithError(message string)` | Sets error message |
| `WithErrorf(format string, args...)` | Sets formatted error |
| `WithErrorAck(err error)` | Sets error with stack trace |
| `AppendError(err error, message string)` | Wraps error with context |
| `AppendErrorf(err error, format string, args...)` | Wraps error with formatted context |
| `WithPath(v string)` | Sets request path |
| `WithPathf(v string, args...)` | Sets formatted request path |
| `WithTotal(total int)` | Sets total items count |

#### Metadata Methods

| Method | Description |
|--------|-------------|
| `WithRequestID(v string)` | Sets request ID |
| `WithRequestIDf(format string, args...)` | Sets formatted request ID |
| `WithApiVersion(v string)` | Sets API version |
| `WithApiVersionf(format string, args...)` | Sets formatted API version |
| `WithLocale(v string)` | Sets locale (e.g., "en_US") |
| `WithRequestedTime(v time.Time)` | Sets request timestamp |
| `WithCustomFieldKV(key string, value interface{})` | Adds custom metadata field |
| `WithCustomFieldKVf(key, format string, args...)` | Adds formatted custom field |
| `WithCustomFields(values map[string]interface{})` | Sets multiple custom fields |
| `WithMeta(v *meta)` | Sets entire metadata object |
| `WithHeader(v *header)` | Sets the header |

#### Pagination Methods

| Method | Description |
|--------|-------------|
| `WithPagination(v *pagination)` | Sets pagination object |
| `WithPage(v int)` | Sets current page number |
| `WithPerPage(v int)` | Sets items per page |
| `WithTotalItems(v int)` | Sets total items count |
| `WithTotalPages(v int)` | Sets total pages count |
| `WithIsLast(v bool)` | Sets if current page is last |

#### Debugging Methods

| Method | Description |
|--------|-------------|
| `WithDebugging(v map[string]interface{})` | Sets debug information map |
| `WithDebuggingKV(key string, value interface{})` | Adds single debug key-value |
| `WithDebuggingKVf(key, format string, args...)` | Adds formatted debug value |

### Query Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `Available()` | `bool` | Checks if wrapper is non-nil |
| `StatusCode()` | `int` | Gets HTTP status code |
| `StatusText()` | `string` | Gets status text (e.g., "OK") |
| `Body()` | `interface{}` | Gets response body |
| `Message()` | `string` | Gets response message |
| `Error()` | `string` | Gets error message |
| `Cause()` | `error` | Gets underlying error cause |
| `Total()` | `int` | Gets total items |
| `Meta()` | `*meta` | Gets metadata object |
| `Header()` | `*header` | Gets header object |
| `Pagination()` | `*pagination` | Gets pagination object |
| `Debugging()` | `map[string]interface{}` | Gets debug information |
| `OnDebugging(key string)` | `interface{}` | Gets specific debug value |

### Conditional Check Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `IsSuccess()` | `bool` | Checks if status is 2xx |
| `IsClientError()` | `bool` | Checks if status is 4xx |
| `IsServerError()` | `bool` | Checks if status is 5xx |
| `IsRedirection()` | `bool` | Checks if status is 3xx |
| `IsError()` | `bool` | Checks if error exists or status is 4xx/5xx |
| `IsErrorPresent()` | `bool` | Checks if error field exists |
| `IsBodyPresent()` | `bool` | Checks if body exists |
| `IsPagingPresent()` | `bool` | Checks if pagination exists |
| `IsMetaPresent()` | `bool` | Checks if metadata exists |
| `IsHeaderPresent()` | `bool` | Checks if header exists |
| `IsDebuggingPresent()` | `bool` | Checks if debug info exists |
| `IsDebuggingKeyPresent(key string)` | `bool` | Checks if specific debug key exists |
| `IsLastPage()` | `bool` | Checks if current page is last |
| `IsStatusCodePresent()` | `bool` | Checks if valid status code exists |
| `IsTotalPresent()` | `bool` | Checks if total count exists |

### Serialization Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `Json()` | `string` | Returns compact JSON string |
| `JsonPretty()` | `string` | Returns pretty-printed JSON |
| `Respond()` | `map[string]interface{}` | Returns map representation |
| `Reply()` | `R` | Returns R wrapper |

## HTTP Status Codes Reference

### Common API Scenarios

| **Scenario** | **HTTP Status Codes** | **Example** |
|--------------|----------------------|-------------|
| **Successful Resource Retrieval** | 200 OK, 304 Not Modified | `GET /users/123` - Returns user data |
| **Resource Creation** | 201 Created | `POST /users` - Creates a new user |
| **Asynchronous Processing** | 202 Accepted | `POST /large-file` - File upload starts |
| **Validation Errors** | 400 Bad Request | `POST /users` - Missing required field |
| **Authentication Issues** | 401 Unauthorized, 403 Forbidden | Invalid credentials or permissions |
| **Rate Limiting** | 429 Too Many Requests | Exceeded API request limits |
| **Missing Resource** | 404 Not Found | `GET /users/999` - User not found |
| **Server Failures** | 500 Internal Server Error, 503 Service Unavailable | Database failure or maintenance |
| **Version Conflicts** | 409 Conflict | Outdated version causing conflict |

### Detailed Status Codes

#### Success (2xx)

| Code | Status | Use Case |
|------|--------|----------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST (resource created) |
| 202 | Accepted | Async processing started |
| 204 | No Content | Successful DELETE |
| 206 | Partial Content | Video streaming, range requests |

#### Redirection (3xx)

| Code | Status | Use Case |
|------|--------|----------|
| 301 | Moved Permanently | Resource permanently moved |
| 302 | Found | Temporary redirect |
| 304 | Not Modified | Cached content still valid |
| 307 | Temporary Redirect | POST redirect maintaining method |
| 308 | Permanent Redirect | Permanent redirect maintaining method |

#### Client Errors (4xx)

| Code | Status | Use Case |
|------|--------|----------|
| 400 | Bad Request | Invalid request format/data |
| 401 | Unauthorized | Missing/invalid authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource conflict (duplicate) |
| 413 | Payload Too Large | Request body too large |
| 415 | Unsupported Media Type | Invalid content type |
| 422 | Unprocessable Entity | Validation errors |
| 429 | Too Many Requests | Rate limiting |

#### Server Errors (5xx)

| Code | Status | Use Case |
|------|--------|----------|
| 500 | Internal Server Error | Unexpected server error |
| 501 | Not Implemented | Feature not implemented |
| 502 | Bad Gateway | Upstream service error |
| 503 | Service Unavailable | Service down/maintenance |
| 504 | Gateway Timeout | Upstream timeout |

## Best Practices

### ✅ Do's

1. **Always set status codes**
   ```go
   response := replify.New().
       WithStatusCode(200).
       WithBody(data)
   ```

2. **Use request IDs for tracing**
   ```go
   response := replify.New().
       WithRequestID(r.Header.Get("X-Request-ID")).
       WithBody(data)
   ```

3. **Include API version**
   ```go
   response := replify.New().
       WithApiVersion("v1.0.0").
       WithBody(data)
   ```

4. **Use WithErrorAck for stack traces**
   ```go
   response := replify.New().
       WithStatusCode(500).
       WithErrorAck(err)
   ```

5. **Check response status before processing**
   ```go
   if response.IsSuccess() {
       processData(response.Body())
   }
   ```

6. **Use pagination for list endpoints**
   ```go
   pagination := replify.Pages().
       WithPage(page).
       WithPerPage(perPage).
       WithTotalItems(total)
   ```

### ❌ Don'ts

1. **Don't forget to set status codes**
   ```go
   // ❌ Bad
   response := replify.New().WithBody(data)
   
   // ✅ Good
   response := replify.New().WithStatusCode(200).WithBody(data)
   ```

2. **Don't expose sensitive debug info in production**
   ```go
   // ❌ Bad
   response := replify.New().
       WithDebuggingKV("database_password", dbPass)
   
   // ✅ Good
   if os.Getenv("ENV") == "development" {
       response.WithDebuggingKV("query", sqlQuery)
   }
   ```

3. **Don't use generic error messages**
   ```go
   // ❌ Bad
   WithError("Error occurred")
   
   // ✅ Good
   WithError("Failed to create user: email already exists")
   ```

4. **Don't ignore error checking**
   ```go
   // ❌ Bad
   wrapper, _ := replify.UnwrapJSON(jsonStr)
   
   // ✅ Good
   wrapper, err := replify.UnwrapJSON(jsonStr)
   if err != nil {
       log.Printf("Failed to parse JSON: %v", err)
   }
   ```

## Use Cases

### ✅ When to Use

- **RESTful API Development** - Standardizing API responses
- **Microservices** - Consistent responses across services
- **API Versioning** - Including version metadata
- **Error Standardization** - Consistent error formats
- **Pagination** - APIs returning paginated results
- **Multi-tenant APIs** - Including tenant/locale information
- **Request Tracing** - Tracking requests across services
- **Development Debugging** - Conditional debug information

### ❌ When Not to Use

- **GraphQL APIs** - GraphQL has its own response format
- **gRPC Services** - Protocol Buffers define the structure
- **WebSocket APIs** - Real-time bidirectional communication
- **Simple CLIs** - Overkill for command-line tools
- **Internal Services** - Where custom formats are required
- **High-Performance** - Direct JSON encoding may be faster

## Contributing

To contribute to this project, follow these steps:

1. **Clone the repository**
   ```bash
   git clone --depth 1 https://github.com/sivaosorg/replify.git
   ```

2. **Navigate to the project directory**
   ```bash
   cd replify
   ```

3. **Prepare the project environment**
   ```bash
   go mod tidy
   ```

4. **Make your changes**
   - Follow Go best practices
   - Add tests for new features
   - Update documentation

5. **Run tests**
   ```bash
   go test ./...
   ```

6. **Submit a pull request**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Sub-packages

### `fj` — Fast JSON Path Querying and Transformation

**Import path:** `github.com/sivaosorg/replify/pkg/fj`

`fj` (_Fast JSON_) provides a fast and simple way to retrieve, query, and transform values from a JSON document without unmarshalling the entire structure into Go types.

#### Quick Start

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

func main() {
    json := `{
        "user": {
            "name": "Alice",
            "roles": [{"roleName":"Admin"},{"roleName":"Editor"}],
            "age": 30,
            "active": true
        }
    }`

    fmt.Println(fj.Get(json, "user.name").String())           // Alice
    fmt.Println(fj.Get(json, "user.age").Int64())             // 30
    fmt.Println(fj.Get(json, "user.active").Bool())           // true
    fmt.Println(fj.Get(json, "user.roles.#.roleName").String()) // ["Admin","Editor"]
}
```

#### Path Syntax Reference

| Syntax | Description | Example |
|--------|-------------|---------|
| `field` | Basic object field access | `user.name` → `"Alice"` |
| `field.N` | Array index access | `roles.0.name` → first element name |
| `field.#` | Array length | `roles.#` → `2` |
| `field.#.key` | Iterate all array elements | `roles.#.roleName` → `["Admin","Editor"]` |
| `field.*` / `field.?` | Wildcard match (`*` = any chars, `?` = one char) | `anim*ls.0.name` |
| `field\.key` | Escape special chars with `\` | `properties.alias\.description` |
| `field.#(key==val)` | Query: first match | `stock.#(symbol=="MMM").price` |
| `field.#(key==val)#` | Query: all matches | `stock.#(active==true)#.name` |
| `~true` / `~false` / `~null` / `~*` | Tilde: boolean coercion in queries | `bank.#(isActive==~true)#.name` |
| `.` / `\|` | Dot and pipe separators (same except after `#`) | `a.b.c` or `a\|b\|c` |
| `{f1,f2}` | Multi-selector → new object | `{id,name}` |
| `[f1,f2]` | Multi-selector → new array | `[id,name]` |
| `!"value"` | JSON literal | `{"x":!true,"y":!"static"}` |
| `..field` | JSON Lines: query each line | `..user.name` |

#### Built-in Transformers

| Transformer | Description | Argument format (optional) |
|-------------|-------------|---------------------------|
| `@trim` | Remove leading/trailing whitespace | — |
| `@this` | Return JSON as-is (identity) | — |
| `@valid` | Return JSON only if valid, else empty | — |
| `@pretty` | Format JSON with indentation | `@pretty:{"sort_keys":true,"indent":"\t","prefix":"","width":80}` |
| `@minify` | Compact JSON, remove whitespace | — |
| `@flip` | Reverse string characters | — |
| `@reverse` | Reverse array elements or object key order | — |
| `@flatten` | Flatten nested arrays | `@flatten:{"deep":true}` |
| `@join` | Merge array of objects into one | `@join:{"preserve":true}` |
| `@keys` | Extract object keys as array | — |
| `@values` | Extract object values as array | — |
| `@string` | Encode value as JSON string | — |
| `@json` | Convert string to JSON representation | — |
| `@group` | Group array-of-object values by key | — |
| `@search` | Search all matching values at path | — |
| `@uppercase` | Convert to uppercase | — |
| `@lowercase` | Convert to lowercase | — |
| `@snakeCase` | Convert to snake_case | — |
| `@camelCase` | Convert to camelCase | — |
| `@kebabCase` | Convert to kebab-case | — |
| `@replace` | Replace first occurrence of substring | `@replace:{"target":"old","replacement":"new"}` |
| `@replaceAll` | Replace all occurrences of substring | `@replaceAll:{"target":"old","replacement":"new"}` |
| `@hex` | Encode string as hex | — |
| `@bin` | Encode string as binary | — |
| `@insertAt` | Insert string at index | `@insertAt:{"index":0,"insert":"prefix"}` |
| `@wc` | Count words in string | — |
| `@padLeft` | Pad string on the left | `@padLeft:{"padding":"*","length":20}` |
| `@padRight` | Pad string on the right | `@padRight:{"padding":"*","length":20}` |

#### Full API Reference

**Top-level functions:**

| Function | Signature | Description |
|----------|-----------|-------------|
| `Get` | `Get(json, path string) Context` | Search JSON for a dot-notation path; returns matching value. |
| `GetBytes` | `GetBytes(json []byte, path string) Context` | Same as `Get` but accepts a byte slice. |
| `Parse` | `Parse(json string) Context` | Parse a JSON string into a `Context` without path querying. |
| `ParseBytes` | `ParseBytes(json []byte) Context` | Same as `Parse` but accepts a byte slice. |
| `ParseBufio` | `ParseBufio(in io.Reader) (string, error)` | Read all data from an `io.Reader` and return as a string. |
| `ParseFilepath` | `ParseFilepath(filepath string) (string, error)` | Read JSON from a file path and return its contents as a string. |
| `IsValidJSON` | `IsValidJSON(json string) bool` | Report whether a JSON string is valid. |
| `IsValidJSONBytes` | `IsValidJSONBytes(json []byte) bool` | Report whether a JSON byte slice is valid. |
| `AddTransformer` | `AddTransformer(name string, fn func(json, arg string) string)` | Register a custom transformer by name. |

**`Context` methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `Kind` | `Kind() Type` | Return the JSON type (`Null`, `False`, `Number`, `String`, `True`, `JSON`). |
| `Unprocessed` | `Unprocessed() string` | Return the raw unprocessed JSON fragment. |
| `Numeric` | `Numeric() float64` | Return the numeric value (for `Number` type). |
| `Index` | `Index() int` | Return the byte position of this value in the original JSON. |
| `Indexes` | `Indexes() []int` | Return positions of all `#`-matched elements. |
| `String` | `String() string` | Return the string representation of the value. |
| `StringColored` | `StringColored() string` | Return the string with default ANSI color styling. |
| `WithStringColored` | `WithStringColored(style *unify4g.Style) string` | Return the string with a custom ANSI color style. |
| `Bool` | `Bool() bool` | Return the boolean value. |
| `Int64` | `Int64() int64` | Return the value as `int64`. |
| `Uint64` | `Uint64() uint64` | Return the value as `uint64`. |
| `Float64` | `Float64() float64` | Return the value as `float64`. |
| `Float32` | `Float32() float32` | Return the value as `float32`. |
| `Time` | `Time() time.Time` | Parse and return the value as `time.Time`. |
| `WithTime` | `WithTime(layout string) time.Time` | Parse time using a custom layout string. |
| `Array` | `Array() []Context` | Return the value as a slice of `Context` (for JSON arrays). |
| `IsObject` | `IsObject() bool` | Report whether the value is a JSON object. |
| `IsArray` | `IsArray() bool` | Report whether the value is a JSON array. |
| `IsBool` | `IsBool() bool` | Report whether the value is a JSON boolean. |
| `Exists` | `Exists() bool` | Report whether the path was found in the JSON. |
| `Value` | `Value() interface{}` | Return the value as a native Go type. |
| `Map` | `Map() map[string]Context` | Return the value as a `map[string]Context` (for JSON objects). |
| `Foreach` | `Foreach(iterator func(key, value Context) bool)` | Iterate over array elements or object key-value pairs. |
| `Get` | `Get(path string) Context` | Query a sub-path from this context (enables chaining). |
| `GetMul` | `GetMul(paths ...string) []Context` | Query multiple paths, returning one result per path. |
| `Path` | `Path(json string) string` | Return the dot-notation path that produced this context. |
| `Paths` | `Paths(json string) []string` | Return paths for each element in an array result. |
| `Less` | `Less(token Context, caseSensitive bool) bool` | Report whether this value is less than `token`. |

#### Custom Transformer

Register a custom transformer with `AddTransformer` to extend the built-in set:

```go
package main

import (
    "fmt"
    "strings"
    "github.com/sivaosorg/replify/pkg/fj"
)

func init() {
    // Register a "shout" transformer: appends "!!!" to the value
    fj.AddTransformer("shout", func(json, arg string) string {
        return strings.Trim(json, `"`) + "!!!"
    })
}

func main() {
    json := `{"greeting": "hello"}`
    fmt.Println(fj.Get(json, "greeting.@shout")) // hello!!!
}
```

#### JSON Color Styles

`fj` ships with named color style variables for use with `Context.WithStringColored`:

| Variable | Description |
|----------|-------------|
| `DarkStyle` | Dark tones (blue, green, yellow, magenta, red, gray) |
| `NeonStyle` | Vibrant neon-like colors |
| `PastelStyle` | Soft pastel tones |
| `HighContrastStyle` | High-contrast colors for accessibility |
| `VintageStyle` | Classic, muted vintage palette |
| `CyberpunkStyle` | Futuristic cyberpunk palette |
| `OceanStyle` | Cool ocean-inspired blues and cyans |
| `FieryStyle` | Warm fiery reds and oranges |
| `GalaxyStyle` | Deep-space purples and blues |
| `SunsetStyle` | Warm sunset oranges and pinks |
| `JungleStyle` | Lush jungle greens |
| `MonochromeStyle` | Grayscale only |
| `ForestStyle` | Earthy forest greens and browns |
| `IceStyle` | Cold icy blues and whites |
| `RetroStyle` | Retro terminal amber/green |
| `AutumnStyle` | Autumn browns, oranges, and reds |
| `GothicStyle` | Dark gothic purples and blacks |
| `VaporWaveStyle` | Aesthetic vaporwave pinks and purples |
| `VampireStyle` | Deep blood reds and blacks |
| `CarnivalStyle` | Bright carnival multicolor |
| `SteampunkStyle` | Brass and copper tones |
| `WoodlandStyle` | Natural woodland tans and greens |
| `CandyStyle` | Bright candy pastels |
| `TwilightStyle` | Dusk purples and navies |
| `EarthStyle` | Warm earth tones |
| `ElectricStyle` | Electric blues and greens |
| `WitchingHourStyle` | Dark witching-hour palette |
| `MidnightStyle` | Deep midnight navy and silver |

#### Best Practices

**✅ Do's**

- Check `ctx.Exists()` before using a value to avoid processing zero-value defaults:
  ```go
  if ctx := fj.Get(json, "user.email"); ctx.Exists() {
      sendEmail(ctx.String())
  }
  ```
- Use `GetBytes` when your JSON is already a `[]byte` to avoid an extra allocation.
- Register custom transformers once at startup (e.g., in an `init()` function).
- Use `Foreach` instead of `Array()` when you only need to process elements one by one, to avoid building an intermediate slice.
- Use `IsValidJSON` / `IsValidJSONBytes` to pre-validate untrusted input.

**❌ Don'ts**

- Don't assume a zero-value `Context` means `null` — it may mean the path was not found. Always call `Exists()`.
- Don't modify the JSON byte slice passed to `GetBytes` while a query is in progress.
- Don't register transformers with names that conflict with built-ins (e.g., `@pretty`, `@minify`).
- Don't call `Map()` on a non-object value without first checking `IsObject()`.
- Don't call `Array()` on a non-array value without first checking `IsArray()`.

---

## Related Packages

Part of the **replify** ecosystem:

- [replify](https://github.com/sivaosorg/replify) - API response wrapping library (this package)
- [fj](https://github.com/sivaosorg/replify/pkg/fj) - Fast JSON path querying and transformation
- [conv](https://github.com/sivaosorg/replify/pkg/conv) - Type conversion utilities
- [coll](https://github.com/sivaosorg/replify/pkg/coll) - Type-safe collection utilities
- [common](https://github.com/sivaosorg/replify/pkg/common) - Reflection-based utilities
- [encoding](https://github.com/sivaosorg/replify/pkg/encoding) - JSON encoding utilities
- [hashy](https://github.com/sivaosorg/replify/pkg/hashy) - Deterministic hashing
- [match](https://github.com/sivaosorg/replify/pkg/match) - Wildcard pattern matching
- [msort](https://github.com/sivaosorg/replify/pkg/msort) - Map sorting utilities
- [randn](https://github.com/sivaosorg/replify/pkg/randn) - Random data generation
- [ref](https://github.com/sivaosorg/replify/pkg/ref) - Pointer utilities
- [strutil](https://github.com/sivaosorg/replify/pkg/strutil) - String utilities
- [truncate](https://github.com/sivaosorg/replify/pkg/truncate) - String truncation utilities

## Support

- **Issues**: [GitHub Issues](https://github.com/sivaosorg/replify/issues)
- **Discussions**: [GitHub Discussions](https://github.com/sivaosorg/replify/discussions)

## Acknowledgments

Built with ❤️ for the Go community.