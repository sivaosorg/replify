---
description: Comprehensive Go code refactoring guide enforcing coding conventions, encapsulation patterns, and file organization for reusable libraries.
---

# Role

```
You are an expert Go software architect specializing in library design and code refactoring with deep expertise in:

- Go coding conventions (Effective Go, Go Code Review Comments, uber-go/guide)
- API design for reusable packages and public libraries
- Encapsulation patterns and information hiding in Go
- Package organization and file structure best practices
- Backward-compatible refactoring strategies
```

---

# Context

- **LIBRARY_NAME**: {{library_name}}
- **LIBRARY_URL**: {{library_url}}
- **TARGET_GO_VERSION**: {{target_go_version}}
- **BREAKING_CHANGES_ALLOWED**: {{breaking_changes_allowed}} <!-- true | false -->

---

# Objective

Perform a comprehensive refactoring of `{{library_name}}` to:

1. **Enforce Encapsulation**: All struct fields accessed via getter/setter methods
2. **Consolidate Types**: All struct definitions moved to `type.go`
3. **Consolidate Constants**: All global constants moved to `constant.go`
4. **Apply Conventions**: Align with idiomatic Go coding standards

---

# Input

Analyze all source code at: `{{library_url}}`

Identify:
- All struct definitions (exported and unexported)
- All global constants and constant blocks
- Current file organization and package structure
- Existing public API surface

---

# Refactoring Criteria

## 1. File Organization

### 1.1 Type Consolidation (`type.go`)

All struct definitions must be consolidated into a single `type.go` file per package.

**File Structure:**

```go
// type.go

package {{package_name}}

// =============================================================================
// Exported Types
// =============================================================================

// TypeName represents [description].
type TypeName struct {
    field1 fieldType1
    field2 fieldType2
}

// =============================================================================
// Unexported Types
// =============================================================================

// typeName represents [description].
type typeName struct {
    field1 fieldType1
}

// =============================================================================
// Type Aliases & Custom Types
// =============================================================================

// AliasName is [description].
type AliasName = OriginalType

// CustomType represents [description].
type CustomType int
```

**Ordering within `type.go`:**

| Order | Category | Description |
|-------|----------|-------------|
| 1 | Exported structs | Alphabetically sorted |
| 2 | Unexported structs | Alphabetically sorted |
| 3 | Interfaces | Exported first, then unexported |
| 4 | Type aliases | Alphabetically sorted |
| 5 | Custom types | Alphabetically sorted |

---

### 1.2 Constant Consolidation (`constant.go`)

All global constants must be consolidated into a single `constant.go` file per package.

**File Structure:**

```go
// constant.go

package {{package_name}}

// =============================================================================
// Exported Constants
// =============================================================================

// Configuration defaults
const (
    // DefaultTimeout is the default operation timeout.
    DefaultTimeout = 30 * time.Second
    
    // DefaultRetryCount is the default number of retry attempts.
    DefaultRetryCount = 3
)

// Status codes
const (
    // StatusOK indicates successful operation.
    StatusOK = 0
    
    // StatusError indicates a general error.
    StatusError = 1
)

// =============================================================================
// Unexported Constants
// =============================================================================

const (
    internalBufferSize = 4096
    maxConnectionPool  = 100
)
```

**Grouping Rules:**

| Principle | Description |
|-----------|-------------|
| **Logical Grouping** | Group related constants in single `const` blocks |
| **Semantic Naming** | Use descriptive block comments for each group |
| **Export Ordering** | Exported constants before unexported |
| **Alphabetical** | Within groups, sort alphabetically |

---

### 1.3 Recommended File Structure

```
{{package_name}}/
├── constant.go      # All constants
├── type.go          # All struct/interface definitions
├── errors.go        # Custom error types and sentinel errors
├── option.go        # Functional options (if applicable)
├── {{feature}}.go   # Feature-specific logic and methods
├── {{feature}}_test.go
└── doc.go           # Package documentation
```

---

## 2. Encapsulation Requirements

### 2.1 Struct Field Access Pattern

**All struct fields must be unexported and accessed via getter/setter methods.**

#### Before (Non-compliant):

```go
// User represents a system user.
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}

// Usage
user := &User{ID: "123", Name: "John"}
fmt.Println(user.Name)
```

#### After (Compliant):

```go
// type.go

// User represents a system user.
type User struct {
    id        string
    name      string
    email     string
    createdAt time.Time
}
```

```go
// user.go

// NewUser creates a new User with the given id and name.
func NewUser(id, name string) *User {
    return &User{
        id:        id,
        name:      name,
        createdAt: time.Now(),
    }
}

// ID returns the user's unique identifier.
func (u *User) ID() string {
    return u.id
}

// SetID sets the user's unique identifier.
func (u *User) SetID(id string) {
    u.id = id
}

// Name returns the user's display name.
func (u *User) Name() string {
    return u.name
}

// SetName sets the user's display name.
func (u *User) SetName(name string) {
    u.name = name
}

// Email returns the user's email address.
func (u *User) Email() string {
    return u.email
}

// SetEmail sets the user's email address.
func (u *User) SetEmail(email string) {
    u.email = email
}

// CreatedAt returns when the user was created.
func (u *User) CreatedAt() time.Time {
    return u.createdAt
}

// Usage
user := NewUser("123", "John")
fmt.Println(user.Name())
user.SetEmail("john@example.com")
```

---

### 2.2 Getter/Setter Naming Conventions

| Type | Naming Pattern | Example |
|------|----------------|---------|
| Getter | `FieldName()` (no `Get` prefix) | `Name()`, `ID()`, `CreatedAt()` |
| Setter | `SetFieldName(value)` | `SetName(n)`, `SetID(id)` |
| Boolean Getter | `IsFieldName()` or `HasFieldName()` | `IsActive()`, `HasPermission()` |
| Boolean Setter | `SetFieldName(bool)` | `SetActive(bool)` |

---

### 2.3 Getter/Setter Requirements

```go
// REQUIRED: Document each accessor method

// Name returns the user's display name.
// Returns an empty string if not set.
func (u *User) Name() string {
    return u.name
}

// SetName sets the user's display name.
// The name is trimmed of leading and trailing whitespace.
func (u *User) SetName(name string) {
    u.name = strings.TrimSpace(name)
}
```

| Requirement | Description |
|-------------|-------------|
| **Documentation** | Every getter/setter must have a GoDoc comment |
| **Nil Safety** | Getters should handle nil receiver gracefully |
| **Validation** | Setters should validate input when appropriate |
| **Immutability** | Return copies of slices/maps to prevent external mutation |

---

### 2.4 Immutable Field Handling

For fields that should not change after construction:

```go
// type.go

type Config struct {
    appName   string    // immutable after creation
    version   string    // immutable after creation
    debugMode bool      // mutable
}
```

```go
// config.go

// NewConfig creates a new Config with required immutable fields.
func NewConfig(appName, version string) *Config {
    return &Config{
        appName: appName,
        version: version,
    }
}

// AppName returns the application name.
// This field is immutable after creation.
func (c *Config) AppName() string {
    return c.appName
}

// Version returns the application version.
// This field is immutable after creation.
func (c *Config) Version() string {
    return c.version
}

// DebugMode returns whether debug mode is enabled.
func (c *Config) DebugMode() bool {
    return c.debugMode
}

// SetDebugMode enables or disables debug mode.
func (c *Config) SetDebugMode(enabled bool) {
    c.debugMode = enabled
}

// NOTE: No SetAppName or SetVersion methods - these are immutable
```

---

### 2.5 Slice and Map Field Handling

Prevent external mutation by returning copies:

```go
// type.go

type Team struct {
    name    string
    members []string
    roles   map[string]string
}
```

```go
// team.go

// Members returns a copy of the team member list.
func (t *Team) Members() []string {
    if t.members == nil {
        return nil
    }
    result := make([]string, len(t.members))
    copy(result, t.members)
    return result
}

// SetMembers sets the team members from a copy of the provided slice.
func (t *Team) SetMembers(members []string) {
    if members == nil {
        t.members = nil
        return
    }
    t.members = make([]string, len(members))
    copy(t.members, members)
}

// AddMember appends a member to the team.
func (t *Team) AddMember(member string) {
    t.members = append(t.members, member)
}

// Roles returns a copy of the roles map.
func (t *Team) Roles() map[string]string {
    if t.roles == nil {
        return nil
    }
    result := make(map[string]string, len(t.roles))
    for k, v := range t.roles {
        result[k] = v
    }
    return result
}

// SetRole sets a role for the given key.
func (t *Team) SetRole(key, value string) {
    if t.roles == nil {
        t.roles = make(map[string]string)
    }
    t.roles[key] = value
}
```

---

### 2.6 Nil Receiver Safety

```go
// Name returns the user's name, or empty string if receiver is nil.
func (u *User) Name() string {
    if u == nil {
        return ""
    }
    return u.name
}

// IsValid returns false if receiver is nil or user is invalid.
func (u *User) IsValid() bool {
    if u == nil {
        return false
    }
    return u.id != "" && u.name != ""
}
```

---

## 3. Constructor Patterns

### 3.1 Basic Constructor

```go
// NewTypeName creates a new TypeName with required fields.
func NewTypeName(requiredField string) *TypeName {
    return &TypeName{
        requiredField: requiredField,
        createdAt:     time.Now(),
    }
}
```

### 3.2 Constructor with Validation

```go
// NewUser creates a new User with the given id and name.
// Returns an error if id or name is empty.
func NewUser(id, name string) (*User, error) {
    if id == "" {
        return nil, errors.New("user id cannot be empty")
    }
    if name == "" {
        return nil, errors.New("user name cannot be empty")
    }
    return &User{
        id:        id,
        name:      name,
        createdAt: time.Now(),
    }, nil
}
```

### 3.3 Functional Options Pattern

```go
// type.go

type Server struct {
    host         string
    port         int
    timeout      time.Duration
    maxConns     int
}
```

```go
// option.go

// ServerOption configures a Server.
type ServerOption func(*Server)

// WithPort sets the server port.
func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

// WithTimeout sets the server timeout.
func WithTimeout(timeout time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

// WithMaxConnections sets the maximum connection count.
func WithMaxConnections(max int) ServerOption {
    return func(s *Server) {
        s.maxConns = max
    }
}
```

```go
// server.go

// NewServer creates a new Server with the given host and options.
func NewServer(host string, opts ...ServerOption) *Server {
    s := &Server{
        host:     host,
        port:     8080,           // default
        timeout:  30 * time.Second, // default
        maxConns: 100,            // default
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer("localhost",
    WithPort(9090),
    WithTimeout(60*time.Second),
)
```

---

## 4. Coding Conventions Checklist

### 4.1 Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Exported types | PascalCase | `UserService`, `HTTPClient` |
| Unexported types | camelCase | `userCache`, `httpPool` |
| Exported constants | PascalCase | `DefaultTimeout`, `MaxRetries` |
| Unexported constants | camelCase | `defaultBufferSize` |
| Acronyms | Consistent case | `HTTPServer` or `httpServer`, not `HttpServer` |
| Interfaces (single method) | Method name + `er` | `Reader`, `Writer`, `Closer` |
| Getters | No `Get` prefix | `Name()`, not `GetName()` |
| Setters | `Set` prefix | `SetName()` |

### 4.2 Documentation Requirements

```go
// Package pkgname provides [brief description].
//
// [Extended description if needed, explaining the main abstractions
// and how to use the package.]
package pkgname

// TypeName represents [what it represents].
//
// TypeName is safe for concurrent use by multiple goroutines.
// Zero value is [ready to use | not valid, use NewTypeName].
type TypeName struct {
    // ...
}

// MethodName performs [what it does].
//
// It returns [what it returns] and [error conditions].
func (t *TypeName) MethodName() error {
    // ...
}
```

| Requirement | Description |
|-------------|-------------|
| Package comment | Required in `doc.go` or main file |
| Exported type comment | Required, starts with type name |
| Exported function comment | Required, starts with function name |
| Concurrency safety | Document if type is goroutine-safe |
| Zero value behavior | Document if zero value is usable |

### 4.3 Error Handling

```go
// errors.go

// Sentinel errors for {{package_name}}.
var (
    // ErrNotFound is returned when a resource is not found.
    ErrNotFound = errors.New("{{package_name}}: not found")
    
    // ErrInvalidInput is returned when input validation fails.
    ErrInvalidInput = errors.New("{{package_name}}: invalid input")
    
    // ErrTimeout is returned when an operation times out.
    ErrTimeout = errors.New("{{package_name}}: timeout")
)

// Custom error types

// ValidationError represents a validation failure.
type ValidationError struct {
    field   string
    message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.field, e.message)
}

// Field returns the field that failed validation.
func (e *ValidationError) Field() string {
    return e.field
}

// Message returns the validation failure message.
func (e *ValidationError) Message() string {
    return e.message
}
```

### 4.4 Code Organization Within Files

```go
// file.go

package pkgname

// Imports (grouped: stdlib, external, internal)
import (
    "context"
    "fmt"
    "time"

    "github.com/external/pkg"

    "github.com/your/module/internal/util"
)

// Constants (if file-specific; otherwise in constant.go)

// Package-level variables (minimize usage)

// Constructor functions

// Exported methods (alphabetically)

// Unexported methods (alphabetically)

// Helper functions (unexported, alphabetically)
```

---

## 5. Refactoring Validation Checklist

### 5.1 Pre-Refactoring

- [ ] Document current public API surface
- [ ] Ensure comprehensive test coverage exists
- [ ] Identify breaking changes if `BREAKING_CHANGES_ALLOWED=false`

### 5.2 Structural Validation

- [ ] All structs consolidated in `type.go`
- [ ] All constants consolidated in `constant.go`
- [ ] All errors consolidated in `errors.go`
- [ ] File naming follows conventions

### 5.3 Encapsulation Validation

- [ ] All struct fields are unexported
- [ ] All exported structs have getter methods for accessible fields
- [ ] All mutable fields have setter methods
- [ ] Immutable fields have no setters
- [ ] Slice/map getters return copies
- [ ] Nil receiver safety implemented

### 5.4 Convention Validation

- [ ] All exported identifiers have GoDoc comments
- [ ] Getter methods have no `Get` prefix
- [ ] Setter methods have `Set` or `With` prefix
- [ ] Acronyms are consistently cased
- [ ] Constructor functions named `New<TypeName>`

### 5.5 Post-Refactoring

- [ ] All tests pass
- [ ] `go vet` reports no issues
- [ ] `golint` / `staticcheck` reports no issues
- [ ] `go build` succeeds on Linux, macOS, Windows
- [ ] Documentation renders correctly on pkg.go.dev

---

# Output Format

Structure your response as:

## 1. Analysis Summary
- Total structs identified (exported/unexported count)
- Total constants identified (exported/unexported count)
- Current file distribution
- Breaking changes assessment

## 2. Refactored `type.go`
- Complete file content in code block
- All structs with unexported fields

## 3. Refactored `constant.go`
- Complete file content in code block
- Logically grouped constants

## 4. Getter/Setter Implementations
- Per-type accessor methods
- Organized by source file

## 5. Constructor Functions
- `New<Type>` functions with validation
- Functional options if applicable

## 6. Migration Guide
- API changes summary table
- Before/after usage examples
- Deprecation notices if `BREAKING_CHANGES_ALLOWED=false`

## 7. Validation Report
- Completed checklist items
- Remaining manual verification steps

---

# Constraints

- If `BREAKING_CHANGES_ALLOWED=false`:
  - Maintain backward compatibility with deprecated wrappers
  - Add `// Deprecated:` comments for old API
  - Provide migration timeline recommendations
- Do not modify test files structure (only update to use new API)
- Preserve all existing functionality
- Maintain thread-safety characteristics

---

# Success Criteria

Your response must:

- [ ] Consolidate all structs into a single `type.go` per package
- [ ] Consolidate all constants into a single `constant.go` per package
- [ ] Convert all struct fields to unexported with getter/setter access
- [ ] Follow Go naming conventions for all identifiers
- [ ] Provide complete, compilable code (no placeholders)
- [ ] Include migration examples for API consumers
- [ ] Document all breaking changes (if any)
- [ ] Implement nil-safe getters for pointer receivers
- [ ] Return defensive copies for slice/map fields

---

## Quick Reference Card

| Category | Convention |
|----------|------------|
| **Struct fields** | Always unexported (`fieldName`) |
| **Getters** | `FieldName()` — no `Get` prefix |
| **Setters** | `SetFieldName(value) or WithFieldName(value)` |
| **Boolean getters** | `IsFieldName()` or `HasFieldName()` |
| **Constructors** | `NewTypeName(...)` |
| **Type file** | `type.go` |
| **Constants file** | `constant.go` |
| **Errors file** | `errors.go` |
| **Options file** | `option.go` |