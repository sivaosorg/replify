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

## Getting Started

### Requirements

Go version **1.19** or higher.

### Installation

```bash
go get github.com/sivaosorg/replify
```

Import the sub-package in your Go code:

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
content, err := fj.ParseJSONFile("/path/to/data.json")
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

### Working with Bytes (no-allocation)

When JSON is already in a `[]byte` slice, use `GetBytes` to avoid converting to string. To further avoid converting `ctx.Unprocessed()` back to `[]byte`, use `ctx.Index()` as a zero-allocation sub-slice:

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

var data = []byte(`{"user":{"id":"12345","roles":[{"roleId":"1","roleName":"Admin"},{"roleId":"2","roleName":"Editor"}]}}`)

func main() {
    ctx := fj.GetBytes(data, "user.roles.#.roleName")

    // Zero-allocation sub-slice of the original []byte
    var raw []byte
    if ctx.Index() > 0 {
        raw = data[ctx.Index() : ctx.Index()+len(ctx.Unprocessed())]
    } else {
        raw = []byte(ctx.Unprocessed())
    }
    fmt.Println(string(raw)) // ["Admin","Editor"]
}
```

### Existence

Use `Exists()` to distinguish a missing path from an explicit JSON `null`:

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

var data = []byte(`{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"}}}`)

func main() {
    ctx := fj.GetBytes(data, "user.name.firstName")
    if ctx.Exists() {
        fmt.Println(ctx.String()) // John
    } else {
        fmt.Println("not found")
    }
}
```

### Loop

`Foreach` iterates over array elements or object key-value pairs. Return `false` from the callback to stop early:

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

var data = []byte(`{"user":{"roles":[{"roleId":"1","roleName":"Admin"},{"roleId":"2","roleName":"Editor"}]}}`)

func main() {
    fj.GetBytes(data, "user.roles").Foreach(func(_, value fj.Context) bool {
        fmt.Println(value.Get("roleName").String())
        return true
    })
    // Admin
    // Editor
}
```

### Unmarshal

Use `Value()` to obtain a native Go type via a type assertion:

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

var data = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"}}}`)

func main() {
    // As map[string]interface{}
    name, ok := fj.GetBytes(data, "user.name").Value().(map[string]interface{})
    if ok {
        fmt.Println(name) // map[firstName:John lastName:Doe]
    }

    // As float64
    id, ok := fj.GetBytes(data, "user.id").Value().(float64)
    if ok {
        fmt.Println(id) // 12345
    }

    // Direct typed accessors are often more ergonomic
    fmt.Println(fj.GetBytes(data, "user.id").Int64()) // 12345
}
```

### Parse & Get

`ParseBytes` parses JSON once; subsequent `Get` calls navigate inside the parsed result without re-scanning the entire document:

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

var data = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"}}}`)

func main() {
    // All three are equivalent
    fmt.Println(fj.ParseBytes(data).Get("user").Get("id").Int64()) // 12345
    fmt.Println(fj.GetBytes(data, "user.id").Int64())              // 12345
    fmt.Println(fj.GetBytes(data, "user").Get("id").Int64())       // 12345
}
```

### JSON Lines (extended)

JSON Lines support uses the `..` prefix to treat a multi-line document as an array. Requires [JSON Lines](https://jsonlines.org/) format (one JSON object per line).

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

var lines = []byte(`
    {"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports"},{"permissionId":"102","permissionName":"Manage Users"}]}
    {"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content"},{"permissionId":"202","permissionName":"View Analytics"}]}
`)

func main() {
    // Count lines
    fmt.Println(fj.ParseBytes(lines).Get("..#").String()) // 2

    // Get a specific line by index
    fmt.Println(fj.GetBytes(lines, "..0").Get("roleName").String()) // Admin

    // Pluck a field from every line
    fmt.Println(fj.GetBytes(lines, "..#.roleName").String()) // ["Admin","Editor"]

    // Nested query across all lines
    fmt.Println(fj.GetBytes(lines, `..#.permissions.#(permissionId=="101").permissionName`).String())
    // ["View Reports"]
}
```

## Path Syntax Reference

A `fj` path is a sequence of elements separated by `.`. In addition to `.`, several characters carry special meaning: `|`, `#`, `@`, `\`, `*`, `!`, and `?`.

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

### Object & Array Examples

```
# Basic field access
id                              → "http://subs/base-sample-schema.json"
properties.alias.description    → "An unique identifier in a submission."
properties.alias.minLength      → 1
required                        → ["alias","taxonId","releaseDate"]
required.0                      → "alias"
required.1                      → "taxonId"
oneOf.0.required.1              → "team"

# Wildcards (* matches any sequence, ? matches one character)
anim*ls.1.name                  → "Barky"
*nimals.1.name                  → "Barky"

# Escape special characters
properties.alias\.description   → "An unique identifier in a submission."

# Array length and pluck
animals.#                       → 3
animals.#.name                  → ["Meowsy","Barky","Purrpaws"]
```

### Query Examples

Queries support `==`, `!=`, `<`, `<=`, `>`, `>=`, `%` (like), and `!%` (not like):

```
stock.#(price_2002==56.27).symbol                → "MMM"
stock.#(company=="Amazon.com").symbol             → "AMZN"
stock.#(initial_price>=10)#.symbol               → ["MMM","AMZN","CPB","DIS","DOW","XOM","F","GPS","GIS"]
stock.#(company%"D*")#.symbol                    → ["DIS","DOW"]
stock.#(company!%"D*")#.symbol                   → ["MMM","AMZN","CPB","XOM","F","GPS","GIS"]
required.#(%"*as*")#                             → ["alias","releaseDate"]
required.#(%"*as*")                              → "alias"
animals.#(foods.likes.#(%"*a*"))#.name           → ["Meowsy","Barky"]
```

### Tilde Operator

The `~` operator evaluates a value as a boolean before comparison. The most recent value that does not exist is treated as `false`:

```
~true    → interprets truthy values as true
~false   → interprets falsy and undefined values as true
~null    → interprets null and undefined values as true
~*       → interprets any defined value as true
```

```
bank.#(isActive==~true)#.name   → ["Davis Wade","Oneill Everett"]
bank.#(isActive==~false)#.name  → ["Stark Jenkins","Odonnell Rollins","Rachelle Chang","Dalton Waters"]
bank.#(eyeColor==~null)#.name   → ["Dalton Waters"]
bank.#(company==~*)#.name       → ["Stark Jenkins","Odonnell Rollins","Rachelle Chang","Davis Wade","Oneill Everett","Dalton Waters"]
```

### Dot & Pipe

The `.` is the default separator; `|` can also be used. They behave identically except after `#` inside array/query contexts:

```
bank.0.balance                         → "$1,404.23"
bank|0.balance                         → "$1,404.23"
bank.0|balance                         → "$1,404.23"
bank.#                                 → 6
bank.#(gender=="female")#|#            → 2
bank.#(gender=="female")#.name         → ["Rachelle Chang","Davis Wade"]
bank.#(gender=="female")#|name         → not-present
bank.#(gender=="female")#|0            → first female object
```

### Multi-Selectors

Comma-separated selectors inside `{...}` create a new object; inside `[...]` create a new array:

```
{version,author,type,"top":stock.#(price_2007>=10)#.symbol}
→ {"version":"1.0.0","author":"subs","type":"object","top":["MMM","AMZN","CPB","DIS","DOW","XOM","GPS","GIS"]}
```

### Literals

JSON literals are prefixed with `!` and are useful when building new objects with [Multi-Selectors](#multi-selectors):

```
{version,author,"marked":!true,"scope":!"static"}
→ {"version":"1.0.0","author":"subs","marked":true,"scope":"static"}
```

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

### Transformer Examples

```
required.1.@flip                                    → "dInoxat"
required.@reverse                                   → ["releaseDate","taxonId","alias"]
required.@reverse.0                                 → "releaseDate"
animals.@join.@minify                               → {"name":"Purrpaws","species":"cat","foods":{...}}
animals.1.@keys                                     → ["name","species","foods"]
animals.1.@values.@minify                           → ["Barky","dog",{...}]
{"id":bank.#.company,"details":bank.#(age>=10)#.eyeColor}|@group
    → [{"id":"HINWAY","details":"blue"},{"id":"NEXGENE","details":"green"},...]
{"id":bank.#.company,"details":bank.#(age>=10)#.eyeColor}|@group|#  → 6
stock.@search:#(price_2007>=50)|0.company            → "3M"
stock.@search:#(price_2007>=50)|0.company.@lowercase → "3m"
stock.0.company.@hex                                 → "334d"
stock.0.company.@bin                                 → "0011001101001101"
stock.0.description.@wc                              → 42
author|@padLeft:{"padding":"*","length":15}|@string  → "***********subs"
author|@padRight:{"padding":"*","length":15}|@string → "subs***********"
bank.0.@pretty:{"sort_keys":true}
→ {
      "address": "766 Cooke Court, Dunbar, Connecticut, 9512",
      "age": 26,
      "balance": "$1,404.23",
      ...
  }
```

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
| `ParseReader` | `ParseReader(in io.Reader) (string, error)` | Read all data from an `io.Reader` and return as a string. |
| `ParseJSONFile` | `ParseJSONFile(filepath string) (string, error)` | Read a JSON file and return its contents as a string. |
| `IsValidJSON` | `IsValidJSON(json string) bool` | Report whether a string is valid JSON. |
| `IsValidJSONBytes` | `IsValidJSONBytes(json []byte) bool` | Report whether a byte slice is valid JSON. |
| `AddTransformer` | `AddTransformer(name string, fn func(json, arg string) string)` | Register a named transformer. |

### Search Engine Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `Search` | `Search(json, keyword string) []Context` | Full-tree scan — return all scalar leaves whose string value contains `keyword`. |
| `SearchMatch` | `SearchMatch(json, pattern string) []Context` | Full-tree wildcard scan — return all scalar leaves whose string value matches `pattern` (`*`, `?`). Uses `match.Match`. |
| `SearchByKey` | `SearchByKey(json string, keys ...string) []Context` | Return all values stored under the given key name(s) at any depth (exact names). |
| `SearchByKeyPattern` | `SearchByKeyPattern(json, keyPattern string) []Context` | Return all values stored under object keys that match the wildcard `keyPattern`. Uses `match.Match`. |
| `Contains` | `Contains(json, path, target string) bool` | Report whether the value at `path` contains the substring `target`. |
| `ContainsMatch` | `ContainsMatch(json, path, pattern string) bool` | Report whether the value at `path` matches the wildcard `pattern`. Uses `match.Match`. |
| `FindPath` | `FindPath(json, value string) string` | Return the first dot-notation path at which a scalar equals `value`. |
| `FindPaths` | `FindPaths(json, value string) []string` | Return all dot-notation paths at which a scalar equals `value`. |
| `FindPathMatch` | `FindPathMatch(json, valuePattern string) string` | Return the first dot-notation path at which a scalar matches the wildcard `valuePattern`. |
| `FindPathsMatch` | `FindPathsMatch(json, valuePattern string) []string` | Return all paths at which a scalar matches the wildcard `valuePattern`. |
| `Count` | `Count(json, path string) int` | Count elements at `path` (array length or 1 for scalars; 0 when missing). |
| `Sum` | `Sum(json, path string) float64` | Sum of all numeric values at `path`. |
| `Min` | `Min(json, path string) (float64, bool)` | Minimum numeric value at `path`. |
| `Max` | `Max(json, path string) (float64, bool)` | Maximum numeric value at `path`. |
| `Avg` | `Avg(json, path string) (float64, bool)` | Arithmetic mean of numeric values at `path`. |
| `CollectFloat64` | `CollectFloat64(json, path string) []float64` | Collect numeric values at `path` using `conv.Float64` (handles string-encoded numbers). |
| `Filter` | `Filter(json, path string, fn func(Context) bool) []Context` | Keep only elements at `path` for which `fn` returns true. |
| `First` | `First(json, path string, fn func(Context) bool) Context` | First element at `path` for which `fn` returns true. |
| `Distinct` | `Distinct(json, path string) []Context` | Unique values at `path` (first-occurrence order). |
| `Pluck` | `Pluck(json, path string, fields ...string) []Context` | Extract named fields from each object in the array at `path`. |
| `GroupBy` | `GroupBy(json, path, keyField string) map[string][]Context` | Group array elements by the string value of `keyField`, using `conv.String` for key normalization. |
| `SortBy` | `SortBy(json, path, keyField string, ascending bool) []Context` | Sort array elements by a sub-field using `conv`-powered numeric and string comparison. |
| `CoerceTo` | `CoerceTo(ctx Context, into any) error` | Coerce a Context's value into any Go typed variable via `conv.Infer`. |

### `Context` Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Kind` | `Kind() Type` | Return the JSON type (`Null`, `False`, `Number`, `String`, `True`, `JSON`). |
| `Raw` | `Raw() string` | Return the raw unprocessed JSON fragment. |
| `Number` | `Number() float64` | Return the numeric value (for `Number` type). |
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
| `GetMulti` | `GetMulti(paths ...string) []Context` | Query multiple paths simultaneously. |
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

### 6. Search Engine

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/replify/pkg/fj"
)

func main() {
    json := `{
        "store": {
            "books": [
                {"id":1,"title":"The Go Programming Language","author":"Donovan","price":34.99,"genre":"tech"},
                {"id":2,"title":"Clean Code","author":"Martin","price":29.99,"genre":"tech"},
                {"id":3,"title":"Harry Potter","author":"Rowling","price":14.99,"genre":"fiction"},
                {"id":4,"title":"Dune","author":"Herbert","price":12.99,"genre":"fiction"}
            ],
            "owner": "Alice"
        },
        "ratings": [5, 3, 4, 5, 1, 4],
        "tags": ["go","json","fast","go","json"]
    }`

    // --- Full-tree substring search ---
    matches := fj.Search(json, "tech")
    fmt.Println(len(matches)) // 2

    // --- Full-tree wildcard search (match.Match) ---
    wildMatches := fj.SearchMatch(json, "D*")
    fmt.Println(len(wildMatches)) // 2: "Donovan", "Dune"

    // --- Search by exact key name ---
    authors := fj.SearchByKey(json, "author")
    for _, a := range authors {
        fmt.Println(a.String()) // Donovan, Martin, Rowling, Herbert
    }

    // --- Search by key wildcard pattern (match.Match) ---
    keyMatches := fj.SearchByKeyPattern(json, "auth*")
    fmt.Println(len(keyMatches)) // 4 (author fields)

    // --- Substring contains ---
    fmt.Println(fj.Contains(json, "store.owner", "Ali")) // true

    // --- Wildcard contains (match.Match) ---
    fmt.Println(fj.ContainsMatch(json, "store.owner", "Al*")) // true

    // --- Path discovery (exact) ---
    fmt.Println(fj.FindPath(json, "Rowling")) // store.books.2.author

    // --- Path discovery (wildcard, match.Match) ---
    fmt.Println(fj.FindPathMatch(json, "Row*")) // store.books.2.author

    // --- All paths matching a wildcard ---
    paths := fj.FindPathsMatch(json, "D*")
    fmt.Println(paths) // [store.books.0.author store.books.3.title]

    // --- Aggregate functions ---
    fmt.Println(fj.Sum(json, "ratings"))       // 22
    fmt.Println(fj.Count(json, "store.books")) // 4
    v, _ := fj.Min(json, "ratings")
    fmt.Println(v) // 1
    v, _ = fj.Max(json, "ratings")
    fmt.Println(v) // 5
    avg, _ := fj.Avg(json, "ratings")
    fmt.Printf("%.4f\n", avg) // 3.6667

    // --- Collect numbers with conv.Float64 (handles string-encoded numbers too) ---
    nums := fj.CollectFloat64(json, "ratings")
    fmt.Println(nums) // [5 3 4 5 1 4]

    // --- Filter ---
    fiction := fj.Filter(json, "store.books", func(ctx fj.Context) bool {
        return ctx.Get("genre").String() == "fiction"
    })
    fmt.Println(len(fiction)) // 2

    // --- First ---
    cheap := fj.First(json, "store.books", func(ctx fj.Context) bool {
        return ctx.Get("price").Float64() < 20
    })
    fmt.Println(cheap.Get("title").String()) // Harry Potter

    // --- Distinct ---
    unique := fj.Distinct(json, "tags")
    fmt.Println(len(unique)) // 3: go, json, fast

    // --- Pluck ---
    projected := fj.Pluck(json, "store.books", "id", "title")
    for _, p := range projected {
        fmt.Println(p.String())
    }

    // --- GroupBy (uses conv.String for key normalization) ---
    groups := fj.GroupBy(json, "store.books", "genre")
    fmt.Println(len(groups["tech"]))    // 2
    fmt.Println(len(groups["fiction"])) // 2

    // --- SortBy (uses conv.Float64/conv.String for comparison) ---
    byPrice := fj.SortBy(json, "store.books", "price", true)
    fmt.Println(byPrice[0].Get("title").String()) // Dune (cheapest)

    // --- CoerceTo (uses conv.Infer) ---
    ctx := fj.Get(json, "store.books.0.price")
    var price float64
    if err := fj.CoerceTo(ctx, &price); err == nil {
        fmt.Printf("%.2f\n", price) // 34.99
    }
}
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
