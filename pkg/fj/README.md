# fj

**fj** (_Fast JSON_) is a Go package that provides a fast and simple way to retrieve, query, and transform values from a JSON document without unmarshalling the entire structure into Go types.

## Overview

The `fj` package uses a dot-notation path syntax that supports wildcards, array indexing, conditional queries, multi-selectors, pipe operators, and a rich set of built-in transformers. Custom transformers can be registered at runtime without modifying the core library.

**Key Features:**
- ⚡ **Zero-allocation hot paths** — `Get` and `GetBytes` minimize heap pressure
- 🔍 **Rich path syntax** — wildcards, array indexing, queries, multi-selectors, literals
- 🔧 **Built-in transformers** — 28+ transformers including `@pretty`, `@minify`, `@reverse`, `@flatten`, `@join`, `@snakeCase`, `@camelCase`, and more
- 🧩 **Extensible** — register custom transformers via `AddTransformer`
- 🎨 **JSON color output** — 27 named color styles for terminal display
- 🔒 **Concurrency-safe** — all public API functions are safe for concurrent goroutine use
- 📁 **File and stream support** — parse JSON from files and `io.Reader` sources

## Use Cases

### When to Use
- ✅ Extracting specific fields from large JSON documents without full unmarshalling
- ✅ Building lightweight JSON pipelines with composable transformers
- ✅ Quickly querying deeply nested or dynamic JSON structures
- ✅ Colorizing JSON output for terminal applications
- ✅ Validating JSON at the edge before heavy processing

### When Not to Use
- ❌ When you need full struct binding (use `encoding/json` or `json-iterator`)
- ❌ When you need to _write_ or _modify_ JSON (use `encoding/json`)
- ❌ When JSON schema validation is required (use a dedicated schema library)

## Installation

```bash
go get github.com/sivaosorg/replify
```

Import the package in your Go code:

```go
import "github.com/sivaosorg/replify/pkg/fj"
```

## Usage

### Basic Field Access

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
            "age": 30,
            "active": true,
            "roles": ["Admin", "Editor"]
        }
    }`

    fmt.Println(fj.Get(json, "user.name").String())    // Alice
    fmt.Println(fj.Get(json, "user.age").Int64())      // 30
    fmt.Println(fj.Get(json, "user.active").Bool())    // true
    fmt.Println(fj.Get(json, "user.roles.0").String()) // Admin
    fmt.Println(fj.Get(json, "user.roles.#").Int64())  // 2 (array length)
}
```

### Array Iteration

```go
json := `{"roles":[{"name":"Admin"},{"name":"Editor"}]}`

// Collect all names
ctx := fj.Get(json, "roles.#.name")
fmt.Println(ctx.String()) // ["Admin","Editor"]

// Iterate with Foreach
fj.Get(json, "roles").Foreach(func(_, val fj.Context) bool {
    fmt.Println(val.Get("name").String())
    return true
})
```

### Queries

```go
json := `{"items":[
    {"name":"apple","price":1.2},
    {"name":"banana","price":0.8},
    {"name":"cherry","price":3.5}
]}`

// First item with price > 1.0
fmt.Println(fj.Get(json, `items.#(price>1.0).name`).String()) // apple

// All items with price > 1.0
fmt.Println(fj.Get(json, `items.#(price>1.0)#.name`).String()) // ["apple","cherry"]
```

### Multi-Selectors

```go
json := `{"id":1,"name":"Alice","email":"alice@example.com","age":30}`

// Build a new object with selected fields
fmt.Println(fj.Get(json, `{id,name}`).String())
// {"id":1,"name":"Alice"}

// Build a new array with selected fields
fmt.Println(fj.Get(json, `[id,name]`).String())
// [1,"Alice"]
```

### Transformers

```go
json := `{"message":"  Hello, World!  "}`

// Trim whitespace
fmt.Println(fj.Get(json, "message.@trim"))

// Convert to uppercase
fmt.Println(fj.Get(json, "message.@uppercase"))

// Pretty-print JSON
data := `{"b":2,"a":1}`
fmt.Println(fj.Get(data, "@pretty"))

// Pretty-print with sorted keys
fmt.Println(fj.Get(data, `@pretty:{"sort_keys":true}`))

// Minify
fmt.Println(fj.Get(data, "@minify"))
```

### Working with Bytes

```go
data := []byte(`{"status":"ok","code":200}`)
ctx := fj.GetBytes(data, "code")
fmt.Println(ctx.Int64()) // 200
```

### Reading from a File

```go
content, err := fj.ParseFilepath("/path/to/data.json")
if err != nil {
    log.Fatal(err)
}
ctx := fj.Get(content, "user.name")
fmt.Println(ctx.String())
```

### Validation

```go
fmt.Println(fj.IsValidJSON(`{"key":"value"}`)) // true
fmt.Println(fj.IsValidJSON(`{bad json}`))       // false
```

### Chaining

```go
json := `{"a":{"b":{"c":"deep"}}}`
ctx := fj.Get(json, "a").Get("b").Get("c")
fmt.Println(ctx.String()) // deep
```

### JSON Lines

```go
lines := `
{"name":"Alice","age":30}
{"name":"Bob","age":25}
{"name":"Carol","age":35}
`
fmt.Println(fj.Get(lines, "..name").String())
// ["Alice","Bob","Carol"]
```

## Path Syntax Reference

| Syntax | Description | Example |
|--------|-------------|---------|
| `field` | Object field access | `user.name` |
| `field.N` | Array index access (0-based) | `roles.0` |
| `field.#` | Array length | `roles.#` → `3` |
| `field.#.key` | Pluck field from all elements | `roles.#.name` → `["Admin","Editor"]` |
| `field.*` | Wildcard — matches any characters | `us*.name` |
| `field.?` | Wildcard — matches exactly one char | `us?.name` |
| `field\.key` | Escape special character `.` | `a\.b` |
| `field.#(k==v)` | Query: first match | `#(active==true).name` |
| `field.#(k==v)#` | Query: all matches | `#(price>1)#.name` |
| `k%"pat"` | Query: like (wildcard match) | `#(name%"A*")` |
| `k!%"pat"` | Query: not like | `#(name!%"A*")` |
| `~true` / `~false` / `~null` / `~*` | Tilde — boolean coercion in queries | `#(active==~true)` |
| `a.b` / `a\|b` | Dot and pipe separators | `user.name` or `user\|name` |
| `{f1,f2}` | Multi-selector → new object | `{id,name}` |
| `[f1,f2]` | Multi-selector → new array | `[id,name]` |
| `!"value"` | JSON literal in multi-selector | `{"active":!true}` |
| `..field` | JSON Lines — query each line | `..user.name` |

## Built-in Transformers

Transformers are path components prefixed with `@`. They can be chained with the pipe (`|`) operator.

| Transformer | Description | Argument format (optional) |
|-------------|-------------|---------------------------|
| `@trim` | Remove leading/trailing whitespace | — |
| `@this` | Return the value as-is (identity) | — |
| `@valid` | Return value only if it is valid JSON; else empty string | — |
| `@pretty` | Format JSON with indentation | `@pretty:{"sort_keys":true,"indent":"  ","prefix":"","width":80}` |
| `@minify` | Compact JSON, remove all whitespace | — |
| `@flip` | Reverse string characters | — |
| `@reverse` | Reverse array elements or object key order | — |
| `@flatten` | Flatten nested arrays (shallow by default) | `@flatten:{"deep":true}` |
| `@join` | Merge array of objects into one object | `@join:{"preserve":true}` |
| `@keys` | Extract object keys as a JSON array | — |
| `@values` | Extract object values as a JSON array | — |
| `@string` | Encode the value as a JSON string | — |
| `@json` | Convert a string to its JSON representation | — |
| `@group` | Group array-of-object values by key | — |
| `@search` | Search all values at a given sub-path | `@search:author` |
| `@uppercase` | Convert to uppercase | — |
| `@lowercase` | Convert to lowercase | — |
| `@snakeCase` | Convert to snake_case (spaces → underscores, lowercase) | — |
| `@camelCase` | Convert to camelCase | — |
| `@kebabCase` | Convert to kebab-case (spaces → hyphens, lowercase) | — |
| `@replace` | Replace the **first** occurrence of a substring | `@replace:{"target":"old","replacement":"new"}` |
| `@replaceAll` | Replace **all** occurrences of a substring | `@replaceAll:{"target":"old","replacement":"new"}` |
| `@hex` | Encode string as hexadecimal | — |
| `@bin` | Encode string as binary | — |
| `@insertAt` | Insert a string at a given byte index | `@insertAt:{"index":5,"insert":"XYZ"}` |
| `@wc` | Count words in the string | — |
| `@padLeft` | Pad on the left to a target length | `@padLeft:{"padding":"*","length":10}` |
| `@padRight` | Pad on the right to a target length | `@padRight:{"padding":"*","length":10}` |

## Custom Transformers

Register a custom transformer with `AddTransformer`. Registrations are thread-safe and take effect immediately for all subsequent calls.

```go
package main

import (
    "fmt"
    "strings"
    "github.com/sivaosorg/replify/pkg/fj"
)

func init() {
    // @shout appends "!!!" to the string value
    fj.AddTransformer("shout", func(json, arg string) string {
        return strings.Trim(json, `"`) + "!!!"
    })

    // @repeat repeats the value N times (arg is the count as a string)
    fj.AddTransformer("repeat", func(json, arg string) string {
        n := fj.Parse(arg).Int64()
        if n <= 0 {
            return json
        }
        v := strings.Trim(json, `"`)
        return strings.Repeat(v, int(n))
    })
}

func main() {
    json := `{"greeting":"hello"}`

    fmt.Println(fj.Get(json, "greeting.@shout"))
    // hello!!!

    fmt.Println(fj.Get(json, "greeting.@repeat:3"))
    // hellohellohello
}
```

## JSON Color Styles

Use `Context.StringColored()` for default ANSI coloring or `Context.WithStringColored(style)` for a custom style.

```go
ctx := fj.Get(json, "user")
fmt.Println(ctx.StringColored())                         // default style
fmt.Println(ctx.WithStringColored(fj.NeonStyle))         // neon style
fmt.Println(ctx.WithStringColored(fj.DarkStyle))         // dark style
```

Available named style variables:

| Variable | Description |
|----------|-------------|
| `DarkStyle` | Dark tones — navy blue, dark green, amber, dark magenta |
| `NeonStyle` | Vibrant neon — bright cyan, lime, yellow |
| `PastelStyle` | Soft pastel tones |
| `HighContrastStyle` | High-contrast for accessibility |
| `VintageStyle` | Muted vintage palette |
| `CyberpunkStyle` | Futuristic cyberpunk neons |
| `OceanStyle` | Cool ocean blues and cyans |
| `FieryStyle` | Warm reds and oranges |
| `GalaxyStyle` | Deep-space purples and blues |
| `SunsetStyle` | Warm oranges and pinks |
| `JungleStyle` | Lush jungle greens |
| `MonochromeStyle` | Grayscale only |
| `ForestStyle` | Earthy forest greens and browns |
| `IceStyle` | Cold icy blues and whites |
| `RetroStyle` | Retro terminal amber/green |
| `AutumnStyle` | Browns, oranges, and reds |
| `GothicStyle` | Dark purples and blacks |
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

## API Reference

### Top-level Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `Get` | `Get(json, path string) Context` | Search JSON for a dot-notation path; return first match. |
| `GetBytes` | `GetBytes(json []byte, path string) Context` | Same as `Get` but accepts a byte slice. |
| `Parse` | `Parse(json string) Context` | Parse a JSON string into a `Context` without path querying. |
| `ParseBytes` | `ParseBytes(json []byte) Context` | Same as `Parse` but accepts a byte slice. |
| `ParseBufio` | `ParseBufio(in io.Reader) (string, error)` | Read all data from an `io.Reader` and return as a string. |
| `ParseFilepath` | `ParseFilepath(filepath string) (string, error)` | Read a JSON file and return its contents as a string. |
| `IsValidJSON` | `IsValidJSON(json string) bool` | Report whether a string is valid JSON. |
| `IsValidJSONBytes` | `IsValidJSONBytes(json []byte) bool` | Report whether a byte slice is valid JSON. |
| `AddTransformer` | `AddTransformer(name string, fn func(json, arg string) string)` | Register a named transformer. |

### `Context` Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Kind` | `Kind() Type` | Return the JSON type (`Null`, `False`, `Number`, `String`, `True`, `JSON`). |
| `Unprocessed` | `Unprocessed() string` | Return the raw unprocessed JSON fragment. |
| `Numeric` | `Numeric() float64` | Return the numeric value (for `Number` type). |
| `Index` | `Index() int` | Return the byte offset of this value in the original JSON. |
| `Indexes` | `Indexes() []int` | Return positions of all `#`-matched elements. |
| `String` | `String() string` | Return a string representation of the value. |
| `StringColored` | `StringColored() string` | Return the string with default ANSI color styling. |
| `WithStringColored` | `WithStringColored(style *unify4g.Style) string` | Return the string with a custom ANSI color style. |
| `Bool` | `Bool() bool` | Return the boolean value. |
| `Int64` | `Int64() int64` | Return the value as `int64`. |
| `Uint64` | `Uint64() uint64` | Return the value as `uint64`. |
| `Float64` | `Float64() float64` | Return the value as `float64`. |
| `Float32` | `Float32() float32` | Return the value as `float32`. |
| `Time` | `Time() time.Time` | Parse the value as `time.Time` using RFC 3339. |
| `WithTime` | `WithTime(layout string) time.Time` | Parse the value as `time.Time` with a custom layout. |
| `Array` | `Array() []Context` | Return all array elements as a `[]Context`. |
| `IsObject` | `IsObject() bool` | Report whether the value is a JSON object. |
| `IsArray` | `IsArray() bool` | Report whether the value is a JSON array. |
| `IsBool` | `IsBool() bool` | Report whether the value is a JSON boolean. |
| `Exists` | `Exists() bool` | Report whether the path was found in the JSON. |
| `Value` | `Value() interface{}` | Return the value as a native Go type (`map`, `[]interface{}`, etc.). |
| `Map` | `Map() map[string]Context` | Return the value as a `map[string]Context` (for JSON objects). |
| `Foreach` | `Foreach(iterator func(key, value Context) bool)` | Iterate over array elements or object key-value pairs. |
| `Get` | `Get(path string) Context` | Query a sub-path (enables chaining). |
| `GetMul` | `GetMul(paths ...string) []Context` | Query multiple paths simultaneously. |
| `Path` | `Path(json string) string` | Return the dot-notation path that produced this context. |
| `Paths` | `Paths(json string) []string` | Return paths for each element in an array result. |
| `Less` | `Less(token Context, caseSensitive bool) bool` | Report whether this value is less than `token`. |

### JSON Type Constants

| Constant | Description |
|----------|-------------|
| `Null` | JSON `null` |
| `False` | JSON `false` |
| `Number` | JSON number |
| `String` | JSON string |
| `True` | JSON `true` |
| `JSON` | JSON object or array |

## Examples

### 1. Nested Object Access

```go
json := `{"store":{"book":{"title":"Go Programming","price":29.99}}}`
title := fj.Get(json, "store.book.title").String()
price := fj.Get(json, "store.book.price").Float64()
fmt.Printf("%s costs $%.2f\n", title, price) // Go Programming costs $29.99
```

### 2. Array Queries

```go
json := `{"users":[
    {"name":"Alice","role":"admin","active":true},
    {"name":"Bob","role":"user","active":false},
    {"name":"Carol","role":"admin","active":true}
]}`

// All active admin names
admins := fj.Get(json, `users.#(role=="admin")#.name`)
fmt.Println(admins.String()) // ["Alice","Carol"]
```

### 3. Transformer Pipeline

```go
json := `{"tags":["Go","json","fast"]}`
result := fj.Get(json, "tags.@reverse")
fmt.Println(result.String()) // ["fast","json","Go"]

// Chain transformers using pipe
reversed := fj.Get(json, "tags.@reverse|0")
fmt.Println(reversed.String()) // "fast"
```

### 4. Multi-Selector with Literals

```go
json := `{"version":"1.0","author":"Alice","items":[1,2,3]}`
result := fj.Get(json, `{version,author,"count":items.#,"published":!true}`)
fmt.Println(result.String())
// {"version":"1.0","author":"Alice","count":3,"published":true}
```

### 5. Struct Unmarshalling via Value()

```go
json := `{"name":"Alice","scores":[95,87,91]}`
m := fj.Get(json, "").Map()
name := m["name"].String()
scores := m["scores"].Array()
fmt.Println(name)                           // Alice
fmt.Println(len(scores), scores[0].Int64()) // 3  95
```

## Best Practices

### ✅ Do's

1. **Check `Exists()` before using a value** to distinguish a missing path from a JSON `null`:

   ```go
   if ctx := fj.Get(json, "user.email"); ctx.Exists() {
       sendEmail(ctx.String())
   }
   ```

2. **Use `GetBytes` when JSON is already `[]byte`** to avoid an extra allocation:

   ```go
   ctx := fj.GetBytes(rawBytes, "user.id")
   ```

3. **Register custom transformers once at startup** (e.g., in `init()`) to avoid race conditions:

   ```go
   func init() {
       fj.AddTransformer("myTransform", myFn)
   }
   ```

4. **Use `Foreach` instead of `Array()`** when you only need to process elements one at a time:

   ```go
   fj.Get(json, "items").Foreach(func(_, val fj.Context) bool {
       process(val)
       return true
   })
   ```

5. **Validate untrusted input** with `IsValidJSON` before heavy processing:

   ```go
   if !fj.IsValidJSON(input) {
       return errors.New("invalid JSON")
   }
   ```

### ❌ Don'ts

1. **Don't assume a zero-value `Context` means JSON `null`** — it may mean the path was not found:

   ```go
   // ❌ Bad
   val := fj.Get(json, "missing.field").String() // returns "" even if not null

   // ✅ Good
   if ctx := fj.Get(json, "missing.field"); ctx.Exists() {
       val := ctx.String()
   }
   ```

2. **Don't mutate the `[]byte` passed to `GetBytes`** while a query is in progress.

3. **Don't use `@replace` when you intend to replace all occurrences** — use `@replaceAll` instead:

   ```go
   // @replace replaces only the FIRST occurrence
   fj.Get(`"foo foo foo"`, `@replace:{"target":"foo","replacement":"bar"}`)
   // → bar foo foo

   // @replaceAll replaces ALL occurrences
   fj.Get(`"foo foo foo"`, `@replaceAll:{"target":"foo","replacement":"bar"}`)
   // → bar bar bar
   ```

4. **Don't call `Map()` on a non-object** without first checking `IsObject()`.

5. **Don't call `Array()` on a non-array** without first checking `IsArray()`.

## Thread Safety

All public API functions (`Get`, `GetBytes`, `Parse*`, `IsValid*`, `AddTransformer`) are safe for concurrent use by multiple goroutines. The transformer registry uses a `sync.RWMutex` for safe concurrent reads and writes.

## Contributing

To contribute to this project, follow these steps:

1. Clone the repository:
   ```bash
   git clone --depth 1 https://github.com/sivaosorg/replify.git
   ```

2. Navigate to the project directory:
   ```bash
   cd replify
   ```

3. Prepare the project environment:
   ```bash
   go mod tidy
   ```

4. Run the tests:
   ```bash
   go test ./pkg/fj/...
   ```

5. Submit a pull request.

## License

This project is licensed under the MIT License — see the [LICENSE](../../LICENSE) file for details.
