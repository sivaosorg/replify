# hashy

**hashy** is a powerful, deterministic hashing library for Go that generates consistent hash values for any Go data structure. It supports primitives, structs, slices, maps, and complex nested types with configurable behavior through struct tags and options.

## Overview

The `hashy` package provides a comprehensive solution for generating deterministic hash values from Go data structures. Unlike built-in hash functions that work only on basic types, `hashy` can hash entire structs, nested data, and collections while respecting field ordering and custom hashing logic.

**Key Features:**
- 🔐 **Deterministic hashing** - identical values always produce identical hashes
- 🎯 **Deep hashing** - works with nested structs, slices, maps, and pointers
- 🏷️ **Struct tag support** - control hashing behavior via `hash:"..."` tags
- ⚙️ **Highly configurable** - customize behavior with fluent option builders
- 🔄 **Multiple output formats** - uint64, hex, base64, SHA-256, and more
- 🎨 **Custom hash interfaces** - implement `Hashable` for custom types
- ⚡ **Optimized performance** - uses FNV-1a by default with fast paths
- 🧩 **Order-independent hashing** - treat slices as sets when needed
- 🔍 **Field filtering** - selectively include/exclude fields via interfaces

**Built on FNV-1a:** By default, `hashy` uses the fast FNV-1a hash algorithm, known for excellent distribution and performance.

## Use Cases

### When to Use
- ✅ **Caching keys** - generate stable cache keys from complex objects
- ✅ **Data deduplication** - detect duplicate records in databases
- ✅ **Change detection** - track whether data has been modified
- ✅ **Distributed systems** - consistent hashing for sharding/partitioning
- ✅ **Testing** - verify object equality in unit tests
- ✅ **ETags** - generate HTTP ETags for API responses
- ✅ **Versioning** - create version fingerprints for configuration
- ✅ **Comparison** - fast equality checks for large data structures

### When Not to Use
- ❌ **Cryptographic security** - use `crypto/*` packages instead (hashy is not cryptographically secure)
- ❌ **Password hashing** - use bcrypt, argon2, or similar
- ❌ **Message authentication** - use HMAC instead
- ❌ **Digital signatures** - use proper cryptographic signing
- ❌ **When uniqueness is critical** - hash collisions are possible (though rare)

## Installation

```bash
go get github.com/sivaosorg/replify
```

Import the package in your Go code:

```go
import "github.com/sivaosorg/replify/pkg/hashy"
```

## Usage

### Basic Hashing

The simplest way to hash values:

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/hashy"
)

func main() {
    // Hash a single value
    hash, err := hashy.Hash("hello world")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hash: %d\n", hash)

    // Hash multiple values (hashed as a tuple)
    hash, err = hashy.Hash("user", 12345, true)
    fmt.Printf("Multi-value hash: %d\n", hash)

    // Hash a struct
    type User struct {
        Name string
        Age  int
    }
    user := User{Name: "Alice", Age: 30}
    hash, err = hashy.Hash(user)
    fmt.Printf("Struct hash: %d\n", hash)
}
```

### String Hash Formats

Generate hashes in various string formats:

```go
// SHA-256 hash string
hash256, _ := hashy.Hash256("data")
fmt.Println(hash256) // "a1b2c3d4..."

// Hexadecimal (16-character, zero-padded)
hexHash, _ := hashy.HashHex("data")
fmt.Println(hexHash) // "000000000a1b2c3d"

// Hexadecimal (short, no padding)
hexShort, _ := hashy.HashHexShort("data")
fmt.Println(hexShort) // "a1b2c3d"

// Base64 encoded
encoded, _ := hashy.HashEncoded("data")
fmt.Println(encoded)

// Decimal string
decimal, _ := hashy.HashBase10("data")
fmt.Println(decimal) // "12345678901234567"

// Hexadecimal string (lowercase)
hex16, _ := hashy.HashBase16("data")
fmt.Println(hex16) // "abc123def456"

// Base32 string
base32, _ := hashy.HashBase32("data")
fmt.Println(base32)
```

### Struct Tags

Control hashing behavior using struct tags:

```go
type User struct {
    ID       int    `hash:"ignore"` // Exclude from hash
    Name     string                  // Included by default
    Password string `hash:"-"`       // Same as "ignore"
    Roles    []string `hash:"set"`   // Order-independent
    Internal string  `hash:"string"` // Use fmt.Stringer if available
}

user := User{
    ID:       123,
    Name:     "Alice",
    Password: "secret",
    Roles:    []string{"admin", "user"},
}

hash, _ := hashy.Hash(user)
// ID and Password are excluded from the hash
// Roles are hashed order-independently
```

**Available Tags:**
- `hash:"ignore"` or `hash:"-"` - Skip this field
- `hash:"set"` - Treat slice as order-independent set
- `hash:"string"` - Use `fmt.Stringer` if type implements it

## Examples

### 1. Caching with Hash Keys

```go
type Product struct {
    ID          int
    Name        string
    Price       float64
    UpdatedAt   time.Time `hash:"ignore"` // Don't invalidate cache on timestamp change
}

func getCacheKey(product Product) (string, error) {
    // Generate a stable cache key
    return hashy.Hash256(product)
}

func main() {
    product := Product{
        ID:    101,
        Name:  "Laptop",
        Price: 999.99,
    }

    cacheKey, _ := getCacheKey(product)
    fmt.Println("Cache key:", cacheKey)

    // Same product = same key (even with different timestamp)
    product.UpdatedAt = time.Now()
    sameKey, _ := getCacheKey(product)
    fmt.Println("Keys match:", cacheKey == sameKey) // true
}
```

### 2. Detecting Changes

```go
type Config struct {
    Host     string
    Port     int
    Features map[string]bool
}

func hasConfigChanged(old, new Config) bool {
    oldHash, _ := hashy.Hash(old)
    newHash, _ := hashy.Hash(new)
    return oldHash != newHash
}

func main() {
    v1 := Config{Host: "localhost", Port: 8080}
    v2 := Config{Host: "localhost", Port: 8080}
    v3 := Config{Host: "localhost", Port: 9090}

    fmt.Println("v1 vs v2:", hasConfigChanged(v1, v2)) // false
    fmt.Println("v1 vs v3:", hasConfigChanged(v1, v3)) // true
}
```

### 3. Order-Independent Slice Hashing

```go
type Team struct {
    Name    string
    Members []string `hash:"set"` // Order doesn't matter
}

func main() {
    team1 := Team{
        Name:    "DevOps",
        Members: []string{"Alice", "Bob", "Charlie"},
    }

    team2 := Team{
        Name:    "DevOps",
        Members: []string{"Charlie", "Alice", "Bob"}, // Different order
    }

    hash1, _ := hashy.Hash(team1)
    hash2, _ := hashy.Hash(team2)

    fmt.Println("Hashes match:", hash1 == hash2) // true - order ignored
}
```

### 4. Custom Hash Options

```go
import "hash/fnv"

func main() {
    // Create custom options
    opts := hashy.NewOptions().
        WithTagName("json").           // Use "json" tag instead of "hash"
        WithIgnoreZeroValue(true).     // Skip zero-value fields
        WithZeroNil(true).             // Treat nil as zero value
        WithSlicesAsSets(true).        // All slices are order-independent
        WithUseStringer(true).         // Always use String() method if available
        Build()

    type User struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }

    user := User{Name: "Alice", Email: "alice@example.com", Age: 0}

    // Hash with custom options
    hash, err := hashy.Hash(user, opts)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Custom hash: %d\n", hash)
}
```

### 5. Implementing Custom Hashable

```go
type Coordinate struct {
    Lat float64
    Lon float64
}

// Implement Hashable interface for custom hashing logic
func (c Coordinate) Hash() (uint64, error) {
    // Round to 6 decimal places for hashing
    lat := int64(c.Lat * 1000000)
    lon := int64(c.Lon * 1000000)
    return hashy.Hash(lat, lon)
}

func main() {
    coord1 := Coordinate{Lat: 40.712776, Lon: -74.005974}
    coord2 := Coordinate{Lat: 40.7127760001, Lon: -74.0059740001} // Slightly different

    hash1, _ := hashy.Hash(coord1)
    hash2, _ := hashy.Hash(coord2)

    // Hashes match due to rounding in custom Hash() method
    fmt.Println("Hashes match:", hash1 == hash2)
}
```

### 6. Field Selection with Interfaces

```go
type UserProfile struct {
    Username    string
    Email       string
    LastLogin   time.Time
    PrivateData string
}

// Implement FieldSelector to control which fields are hashed
func (u UserProfile) SelectField() hashy.SelectField {
    return func(field string, value any) (bool, error) {
        // Exclude sensitive and volatile fields
        excluded := []string{"PrivateData", "LastLogin"}
        for _, f := range excluded {
            if field == f {
                return false, nil
            }
        }
        return true, nil
    }
}

func main() {
    profile := UserProfile{
        Username:    "alice",
        Email:       "alice@example.com",
        LastLogin:   time.Now(),
        PrivateData: "secret",
    }

    hash, _ := hashy.Hash(profile)
    // Only Username and Email are included in the hash
    fmt.Printf("Selective hash: %d\n", hash)
}
```

### 7. Map Entry Selection

```go
type Settings struct {
    Values map[string]string
}

// Control which map entries are hashed
func (s Settings) SelectMapEntry() hashy.SelectMapEntry {
    return func(field string, k, v any) (bool, error) {
        key, ok := k.(string)
        if !ok {
            return true, nil
        }
        // Exclude temporary settings
        return !strings.HasPrefix(key, "temp_"), nil
    }
}

func main() {
    settings := Settings{
        Values: map[string]string{
            "host":      "localhost",
            "port":      "8080",
            "temp_flag": "true", // Excluded from hash
        },
    }

    hash, _ := hashy.Hash(settings)
    fmt.Printf("Filtered map hash: %d\n", hash)
}
```

### 8. ETags for HTTP APIs

```go
type APIResponse struct {
    Data      interface{}
    Timestamp time.Time `hash:"ignore"` // Don't include in ETag
}

func generateETag(response APIResponse) (string, error) {
    hash, err := hashy.HashHex(response)
    if err != nil {
        return "", err
    }
    return `"` + hash + `"`, nil // Wrap in quotes for HTTP ETag
}

func main() {
    response := APIResponse{
        Data:      map[string]string{"status": "ok"},
        Timestamp: time.Now(),
    }

    etag, _ := generateETag(response)
    fmt.Println("ETag:", etag) // "0000000012abc34d"
}
```

## API Reference

### Core Functions

| Function | Description | Return Type |
|----------|-------------|-------------|
| `Hash(data ...any)` | Generate 64-bit hash | `uint64, error` |
| `HashValue(value any, opts)` | Hash single value with options | `uint64, error` |
| `Hash256(data ...any)` | Generate SHA-256 hash string | `string, error` |
| `HashHex(data ...any)` | 16-char hex hash (zero-padded) | `string, error` |
| `HashHexShort(data ...any)` | Hex hash (no padding) | `string, error` |
| `HashBase10(data ...any)` | Decimal string hash | `string, error` |
| `HashBase16(data ...any)` | Hexadecimal string | `string, error` |
| `HashBase32(data ...any)` | Base32 string | `string, error` |
| `HashEncoded(data ...any)` | Base64 encoded hash | `string, error` |

### Options Builder

```go
opts := hashy.NewOptions().
    WithHasher(customHasher).      // Custom hash.Hash64 implementation
    WithTagName("json").            // Use different struct tag
    WithZeroNil(true).              // Treat nil pointers as zero values
    WithIgnoreZeroValue(true).      // Skip zero-value fields
    WithSlicesAsSets(true).         // All slices are order-independent
    WithUseStringer(true).          // Use fmt.Stringer when available
    Build()
```

### Interfaces

**Hashable** - Custom hash implementation:
```go
type Hashable interface {
    Hash() (uint64, error)
}
```

**FieldSelector** - Control field inclusion:
```go
type FieldSelector interface {
    SelectField() SelectField
}

type SelectField func(field string, value any) (bool, error)
```

**MapSelector** - Control map entry inclusion:
```go
type MapSelector interface {
    SelectMapEntry() SelectMapEntry
}

type SelectMapEntry func(field string, k, v any) (bool, error)
```

### Struct Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `hash:"ignore"` | Exclude field from hash | `ID int \`hash:"ignore"\`` |
| `hash:"-"` | Same as "ignore" | `Password string \`hash:"-"\`` |
| `hash:"set"` | Order-independent slice | `Tags []string \`hash:"set"\`` |
| `hash:"string"` | Use fmt.Stringer | `Status Status \`hash:"string"\`` |

## Best Practices & Notes

### ⚠️ Common Pitfalls

1. **Hash Collisions**: While rare with 64-bit hashes, collisions are possible. Don't rely on uniqueness for security-critical operations.

2. **Non-Deterministic Types**: Be careful with:
   - `map` iteration order (handled automatically by hashy)
   - Pointer addresses (use values, not pointers)
   - Random values (exclude or use fixed seeds)

3. **Float Precision**: Floating-point values are hashed as-is. Small differences will produce different hashes:
   ```go
   hashy.Hash(0.1 + 0.2) != hashy.Hash(0.3) // May differ due to float precision
   ```

4. **Time Values**: `time.Time` includes location and monotonic clock. Strip unnecessary precision:
   ```go
   timestamp := time.Now().UTC().Truncate(time.Second)
   ```

5. **Unexported Fields**: Only exported (public) struct fields are hashed.

### 💡 Recommendations

✅ **Use struct tags** to exclude volatile fields (timestamps, IDs)

✅ **Use `hash:"set"` for unordered collections** (permissions, tags)

✅ **Implement `Hashable`** for types requiring custom logic

✅ **Use `FieldSelector`** for complex inclusion rules

✅ **Generate string hashes** for cache keys and ETags (`Hash256`, `HashHex`)

✅ **Test hash stability** across versions if persisting hashes

✅ **Document hash assumptions** in struct comments

### 🔒 Thread Safety

All hashing functions are safe for concurrent use. The FNV-1a hasher is recreated for each operation.

```go
// Safe - concurrent hashing
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(val int) {
        defer wg.Done()
        hashy.Hash(val)
    }(i)
}
wg.Wait()
```

### ⚡ Performance Tips

- **Reuse options**: Build once, pass to multiple `HashValue()` calls
- **Use appropriate hash function**: FNV-1a is fast; use custom hasher if needed
- **Minimize allocations**: Hash values directly instead of serializing first
- **Cache hash results**: For immutable data, compute once and store
- **Use uint64 format**: Faster than string conversions

### 🐛 Debugging Hash Mismatches

If two objects that should hash the same don't:

1. **Check field order**: Struct field order matters
2. **Verify zero values**: Use `WithIgnoreZeroValue(true)` consistently
3. **Inspect tags**: Ensure `hash:"set"` is on the right fields
4. **Check pointer equality**: `&x` and `&y` are different even if `x == y`
5. **Print individual field hashes**: Hash each field separately to identify differences

```go
// Debug individual fields
type User struct {
    Name string
    Age  int
}

u := User{Name: "Alice", Age: 30}
nameHash, _ := hashy.Hash(u.Name)
ageHash, _ := hashy.Hash(u.Age)
fmt.Printf("Name: %d, Age: %d\n", nameHash, ageHash)
```

### 📝 Versioning Considerations

If you plan to store or compare hashes across application versions:

- **Document hash inputs**: Clearly state which fields are included
- **Version your hash logic**: Add a version prefix if hash algorithm changes
- **Test compatibility**: Verify old data still produces expected hashes
- **Avoid breaking changes**: Don't change field order or tag behavior

```go
// Versioned hash
type VersionedData struct {
    Version int    `hash:"ignore"` // Don't hash version itself
    Data    string
}

func hashWithVersion(data VersionedData) (string, error) {
    hash, err := hashy.HashHex(data)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("v%d:%s", data.Version, hash), nil
}
```

## Error Handling

The library returns two main error types:

1. **`ErrNotStringer`**: Field has `hash:"string"` tag but doesn't implement `fmt.Stringer`
2. **General errors**: Invalid options, unsupported types

```go
hash, err := hashy.Hash(data)
if err != nil {
    var notStringer *hashy.ErrNotStringer
    if errors.As(err, &notStringer) {
        log.Printf("Field %s needs Stringer implementation", notStringer.Field)
    } else {
        log.Printf("Hash error: %v", err)
    }
}
```

## Limitations

- **Not cryptographically secure**: Use standard library `crypto/*` for security
- **No guaranteed uniqueness**: Hash collisions are theoretically possible
- **Struct field order matters**: Reordering fields changes the hash
- **Unexported fields ignored**: Only public fields are hashed
- **No cross-language compatibility**: Hashes are Go-specific

## Contributing

Contributions are welcome! Please see the main [replify repository](https://github.com/sivaosorg/replify) for contribution guidelines.

## License

This library is part of the [replify](https://github.com/sivaosorg/replify) project.
