# translit — architecture, performance, and design notes

A from-scratch rewrite of [`github.com/mozillazg/go-unidecode`](https://github.com/mozillazg/go-unidecode),
built around one goal: transliterate Unicode text to ASCII as fast as
possible, with zero heap allocations, predictable latency, and full
concurrency safety — while producing **byte-for-byte identical output**
to the original library.

The transliteration _data_ (which character maps to which ASCII
replacement) is unchanged from the original project. What's rewritten is
everything about how that data is stored and walked.

---

## 1. Summary of results

Measured on this machine (`Intel(R) Xeon(R) Processor @ 2.10GHz`, `GOMAXPROCS=1`,
Go 1.22, `-benchtime=1s`). Full raw output in `compare/` — run
`go test ./compare/... -bench . -benchmem` to reproduce.

| Category (running text, `Append` zero-alloc path) | Original ns/op | translit ns/op |   Speedup | Original allocs/op | translit allocs/op |
| ------------------------------------------------- | -------------: | -------------: | --------: | -----------------: | -----------------: |
| ASCII                                             |          772.8 |           77.5 | **10.0x** |                  5 |              **0** |
| Latin-extended (Czech/diacritics)                 |          865.4 |          281.5 |  **3.1x** |                  5 |              **0** |
| Vietnamese                                        |          868.8 |          252.5 |  **3.4x** |                  5 |              **0** |
| Chinese                                           |          703.2 |          266.9 |  **2.6x** |                  5 |              **0** |
| Japanese                                          |          745.7 |          313.0 |  **2.4x** |                  5 |              **0** |
| Korean                                            |          774.1 |          252.2 |  **3.1x** |                  5 |              **0** |
| Mixed scripts                                     |          959.9 |          295.7 |  **3.2x** |                  5 |              **0** |
| 1 KB input                                        |           7384 |           2563 |  **2.9x** |                  8 |              **0** |
| 10 KB input                                       |          76434 |          31530 |  **2.4x** |                 15 |              **0** |
| 100 KB input                                      |         747018 |         247176 |  **3.0x** |                 23 |              **0** |

Even the convenience `Unidecode(s) string` wrapper — which necessarily
allocates once, to build the returned string — beats the original by
1.5x-3.2x while using 2 allocations instead of 5-23, because it sizes
its one allocation exactly via `SizeHint` instead of growing a
`strings.Builder` repeatedly.

The `Append`/`AppendBytes` hot path is **0 B/op, 0 allocs/op** for every
category, confirmed both by `testing.B`'s allocation counter and by an
explicit `testing.AllocsPerRun` regression test (`TestZeroAllocations`)
that fails the build if this ever regresses.

---

## 2. Why the original is slower

Reading `unidecode.go` in the original project:

```go
func unidecode(s string) string {
	var ret strings.Builder
	for _, r := range s {
		if r < unicode.MaxASCII {
			ret.WriteRune(r)
			continue
		}
		if r > 0xeffff {
			continue
		}
		section := r >> 8
		position := r % 256
		if tb, ok := table.Tables[section]; ok {
			if len(tb) > int(position) {
				ret.WriteString(tb[position])
			}
		}
	}
	return ret.String()
}
```

Three things dominate its cost:

1. **`map[rune][]string` lookup on every non-ASCII rune.** Go maps are
   hash tables: computing a hash, probing buckets, and following a
   pointer to the value are all inherent per-lookup costs — and that
   pointer chases to a `[]string`, which is itself a slice of string
   headers (pointer+length pairs), each pointing to a separately
   allocated backing array. That's two to three pointer-chasing,
   cache-unfriendly indirections _per rune_, before any bytes are even
   copied.
2. **`strings.Builder` growth.** Builder starts empty and grows by
   doubling as content is written. For a 100 KB input this means
   repeated reallocation-and-copy of ever-larger buffers — visible
   directly in the allocation counts above (23 allocations for 100 KB,
   growing roughly logarithmically with input size).
3. **No fast path distinguishing "plain ASCII" from "everything else."**
   Every rune, including every ASCII letter, goes through the same
   `for range` rune-decode step (cheap, but not free) and the same
   comparison chain.

None of this is a criticism of the original — it's a small, correct,
readable ~25-line function, which is exactly the right tradeoff for a
project that doesn't need to sit in anyone's hot path. This rewrite
makes the opposite tradeoff: more code, more generated data, in
exchange for eliminating all of the above.

---

## 3. Data layout

### 3.1 The original layout

```go
var Tables = map[rune][]string{}
Tables[0x000] = x000  // []string of up to 256 entries
Tables[0x001] = x001
...
```

Each `x0NN` is a `[]string`. In memory this is: a map (hash table) whose
values are slice headers, each pointing to an independently-allocated
array of **string headers**, each of _those_ pointing to an
independently-allocated byte array holding the actual replacement text.
For "北" (U+5317), resolving to the replacement text means: hash `0x53`,
probe the map, follow the slice header, index to position `0x17`, follow
that string header, land on the bytes "Bei ". Four dependent pointer
dereferences, each a potential cache miss, scattered across the heap in
whatever order the Go allocator happened to place ~48,000 tiny objects
during `init()`.

### 3.2 This project's layout

Generated by `gen/main.go` (see `tables_gen.go`, `go generate`-produced,
**do not hand-edit**) into four flat, contiguous, immutable arrays plus
one shared string:

```go
const data = "...157,244 bytes of replacement text, concatenated..."

var sectionStart [503]int32   // section -> index into entryOff/entryLen, or -1
var sectionLen   [503]uint16  // section -> number of valid positions (0..255)
var entryOff     [48597]uint32 // entry -> byte offset into data
var entryLen     [48597]uint8  // entry -> byte length within data
```

A lookup for rune `r` is:

```go
section := r >> 8
start   := sectionStart[section]   // array index — O(1), no hashing
pos     := r & 0xff
idx     := start + pos             // array index — O(1)
off, ln := entryOff[idx], entryLen[idx]
replacement := data[off : off+ln]  // slice of one shared, read-only string
```

Four array index operations (each a single bounds-checked memory read —
the Go compiler proves these are in-range in the common case and elides
redundant checks) and one string slice (a pointer+length computation,
**not a copy**). No hashing, no map probing, no pointer chasing through
scattered heap objects, no interface dispatch, no reflection.

### 3.3 Why this shape specifically

- **Two-level (section, position) split, not a flat `[0x110000]` array.**
  A single flat array over all of Plane 0-14 would be ~1.1M entries,
  the overwhelming majority unused (real text clusters into a handful
  of scripts), wasting memory and — more importantly — cache: touching
  one Chinese character would pull in a cache line mostly full of
  entries for codepoints that will never be looked up in this
  process's lifetime. The two-level split means `sectionStart` /
  `sectionLen` are small (503 × (4+2) bytes ≈ 3 KB — comfortably
  L1-resident) and the `entryOff`/`entryLen` region touched for any
  one piece of text is exactly the sections that text actually uses.
- **Offset+length into a shared byte string, not `[48597]string`.**
  A Go `string` header is 16 bytes (pointer + length) even before you
  reach the actual bytes. Storing 48,597 separate string headers would
  be 777 KB of headers alone, most pointing into the _same_ underlying
  `data` constant (since it's one contiguous literal) — pure overhead.
  `uint32` offset + `uint8` length is 5 bytes per entry (243 KB total),
  and reconstructing the string is one slice expression, not a
  dereference.
- **`uint8` for length, not `int`.** The generator asserts every
  replacement is ≤255 bytes (the longest actual entry is 13 bytes), so
  a full byte is already generous headroom; using `int` would triple
  this array's size for no benefit.
- **Section-major layout, not rune-major.** Because sections are laid
  out contiguously in ascending order and a script's codepoints are
  themselves contiguous within Unicode (Latin, Cyrillic, Hiragana,
  CJK blocks, etc.), transliterating a run of same-script text touches
  a short, sequential stretch of `entryOff`/`entryLen` — good spatial
  locality, prefetcher-friendly.

---

## 4. The ASCII fast path

```go
c := src[i]
if c < asciiFastPathMax {   // 0x7f, matching the original's unicode.MaxASCII cutoff
    dst = append(dst, c)
    i++
    continue
}
```

This is the single most important optimization for realistic text,
since most human-language text — even non-Latin scripts mixed with
punctuation, digits, and spaces — is majority ASCII by byte count. A
single unsigned comparison against a constant, no function call, no
UTF-8 decode, no table touch, directly into the destination slice. The
`ASCII` benchmark category (10x over the original) isolates exactly this
path.

Note the boundary is deliberately `0x7f` (127), not `utf8.RuneSelf`
(128): the original library's check is `r < unicode.MaxASCII`, and
`unicode.MaxASCII == '\u007f'`. That means DEL (0x7F) is _not_
fast-pathed by the original — it goes through the table, whose entry
for DEL is empty, so DEL is silently dropped. This rewrite reproduces
that exact cutoff (`asciiFastPathMax = 0x7f`) rather than the more
"obvious" 128, specifically so the two libraries agree on every single
byte value, including this one edge case. (Found and pinned down by the
exhaustive `TestGoldenEveryRune` test — see §7.)

Non-ASCII bytes fall through to `utf8.DecodeRuneInString`/`DecodeRune`
(the standard library's decoder, as required — this project doesn't
reimplement UTF-8 decoding), which the runtime inlines and which is
itself branch-light for well-formed input.

---

## 5. Zero allocations, caller-owned buffers

```go
func Append(dst []byte, src string) []byte
func AppendBytes(dst []byte, src []byte) []byte
```

Both follow Go's standard `append`-style convention: the destination
buffer is supplied by the caller, and grown (via the builtin `append`,
which the library never bypasses with manual `unsafe` buffer tricks) only
if it runs out of spare capacity. If the caller pre-sizes `dst` — which
`SizeHint` makes exact and cheap to do — **no allocation happens at
all**, for any input, of any size, in any script mix. This is verified
by:

- `go test -bench . -benchmem` showing `0 B/op 0 allocs/op` for every
  `BenchmarkFast_Append_*` case (§1).
- `TestZeroAllocations` / `TestZeroAllocationsModeKeep`, which use
  `testing.AllocsPerRun` to assert exactly `0` allocations per call
  across ASCII, Latin, Vietnamese, Chinese, Japanese, Korean, and mixed
  input, and fail the test suite (not just a benchmark run someone
  might not look at) if this ever regresses.

No `strings.Builder`, no `bytes.Buffer`, no `fmt`, no intermediate
`[]rune` or `[]byte(string)` conversion anywhere in the hot path. The
only place this project allocates at all is the `Unidecode(s) string`
convenience wrapper, which is required to allocate exactly once because
Go strings are immutable and must own their backing array — and even
that allocation is sized exactly via `SizeHint`, so it never
over-allocates or reallocates.

`SizeHint` itself performs the identical table walk `Append` does
(same O(1)-per-rune lookups) but only accumulates a length, writing
nothing — so it's cheap enough to call on every request without being a
second pass that meaningfully changes the algorithm's cost, and it lets
a caller achieve a **guaranteed** zero-allocation, zero-realloc call:

```go
buf := make([]byte, 0, unidecode.SizeHint(src))
buf = unidecode.Append(buf, src)
```

---

## 6. Complexity analysis

- **Per-rune lookup: O(1).** Every non-ASCII codepoint costs exactly
  four fixed-size array reads (`sectionStart`, `sectionLen`,
  `entryOff`, `entryLen`) plus one bounds-checked slice of the shared
  `data` string — no loop, no probing, no data-dependent iteration
  count. ASCII bytes cost one comparison and one append. Both are O(1)
  independent of input size or table size.
- **Overall: O(n)** in input bytes. `Append`/`AppendBytes` make one
  linear pass over `src`, decoding each rune once (`utf8.DecodeRune*`
  is itself O(1) amortized — a UTF-8 sequence is at most 4 bytes) and
  performing O(1) work per rune. There is no nested iteration over the
  table (`sectionLen`/`maxSection` bound checks are O(1) reads, not
  scans), no recursion, and no backtracking.
- **No data-dependent branching that affects asymptotic cost.** The
  `Mode` check (`ModeSkip` vs `ModeKeep`) is a single comparison against
  a `uint8`, not a virtual/interface dispatch — the compiler can (and
  does, since `Mode` has only two values) treat it as a simple branch.
  There's no `switch` over rune ranges, no per-script special-casing in
  the algorithm itself; the _data_ (which section a rune falls in)
  determines the _replacement text_, not the _code path_ taken to find
  it.

---

## 7. Correctness: exhaustive golden testing

`compare/golden_test.go`'s `TestGoldenEveryRune` iterates **every legal
Go rune value** (`0` through `utf8.MaxRune`, skipping the UTF-16
surrogate range, which isn't representable as a Go rune) — 1,114,112
codepoints — transliterates each one with both the original library and
this one, and asserts byte-for-byte identical output. This is stronger
than spot-checking known tricky scripts: it's a complete proof that the
two implementations agree on the _entire_ input domain, not just the
authors' guesses about what might differ.

This exhaustive test is what caught the DEL (`0x7f`) boundary
discrepancy described in §4 — the kind of one-character-out-of-a-million
edge case that hand-picked test cases would very plausibly miss, and
exactly the class of bug that a rewrite claiming "identical behavior"
needs to rule out completely rather than probabilistically.

`TestGoldenRealWorldText` and `TestGoldenInvalidUTF8` supplement this
with realistic multi-script sentences and malformed-UTF-8 inputs, for
readable regression failures if something does drift (the exhaustive
test reports _which_ rune diverged, but a paragraph-level test is easier
to reason about when debugging).

---

## 8. Thread safety

Every piece of package state (`data`, `sectionStart`, `sectionLen`,
`entryOff`, `entryLen`, and the `maxSection` constant) is:

- Declared as a package-level `const`/array **literal** — populated
  entirely at compile time by the generator, with no `init()` function
  performing runtime construction.
- Never written to after program start. There are no setter functions,
  no `sync.Once`-guarded lazy build, no mutation of any kind.

Because of this, every exported function is trivially safe for
unlimited concurrent use: concurrent goroutines calling `Append` are
only ever _reading_ shared immutable memory and _writing_ to their own,
separately-owned `dst` slices — there is no shared mutable state for a
race to occur on. `go test -race` passes across the full suite,
including `TestConcurrentAppend` (64 goroutines × 200 iterations against
shared input strings, comparing every result against a precomputed
expected value) and `TestConcurrentAppendSharedDstIsolated` (32
goroutines each repeatedly re-deriving their own output and checking it
never drifts).

No `sync.Mutex`, `sync.RWMutex`, `sync.Map`, `sync.Once`, or any other
synchronization primitive appears anywhere in the library — none is
needed.

---

## 9. Cache-locality analysis

- **`sectionStart` + `sectionLen`**: 503 × 4 bytes + 503 × 2 bytes ≈
  3 KB total. Comfortably fits in a modern L1 data cache (typically
  32-48 KB) alongside everything else a hot loop touches, so after the
  first access these are effectively always resident.
- **`entryOff` + `entryLen`** for a given section: contiguous runs of
  `sectionLen[section]` entries (up to 256 × (4+1) = 1,280 bytes per
  section). Transliterating a run of same-script text — the overwhelmingly
  common case in real documents — touches one or a few such contiguous
  regions repeatedly, which the CPU prefetcher picks up on directly
  after the first couple of accesses.
- **`data`**: 157 KB total, but any single document only ever touches
  the slices its specific characters reference — and because entries
  for a section are emitted in the _same order_ as they're laid out in
  `entryOff` (ascending position within ascending section), the
  replacement text for adjacent codepoints in a script tends to sit
  physically near each other in `data` too, compounding the locality
  benefit from the offset table.
- **Compare to the original's per-rune cost**: a map lookup touches the
  map's bucket array (a hash-dependent, effectively random location),
  then a `[]string` header on a separately-allocated slice, then a
  `string` header, then finally the bytes — four essentially unrelated
  regions of the heap, with no relationship between "characters that
  appear near each other in text" and "memory locations that end up
  near each other," because Go's map bucket placement is
  hash-order, not insertion- or key-order.

---

## 10. Explanation of every optimization, one by one

| Optimization                                                                                | Where                                         | Why it helps                                                                                                                                                                                                                                                                                                                                                               |
| ------------------------------------------------------------------------------------------- | --------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ASCII fast path (single comparison, direct append)                                          | `Append`/`AppendBytes`/`SizeHint`/`runeSize`  | Skips UTF-8 decode and table lookup entirely for the majority byte class in real text.                                                                                                                                                                                                                                                                                     |
| Flat `sectionStart`/`sectionLen` two-level index                                            | `tables_gen.go` (generated)                   | Replaces a hash lookup with two array reads; O(1) with no hashing, no probing, no chained pointer dereference through map internals.                                                                                                                                                                                                                                       |
| Offset+length into one shared `data` string, not `[]string`                                 | `tables_gen.go` (generated)                   | Removes 48,597 separately-allocated string headers/backing arrays in favor of 5 bytes/entry into one contiguous, read-only buffer — far fewer distinct cache lines touched, far less memory overhead.                                                                                                                                                                      |
| `uint8` length, `uint32` offset (not `int`)                                                 | `tables_gen.go` (generated)                   | Every real replacement is ≤13 bytes; sizing fields to their true range roughly halves the offset/length table's footprint vs. native `int`.                                                                                                                                                                                                                                |
| Caller-supplied `dst []byte`, `append`-style growth                                         | `Append`/`AppendBytes`                        | Lets a caller with a known-size or reusable buffer achieve true zero allocations; matches the idiomatic Go convention (`append`, `strconv.AppendInt`, etc.) instead of inventing a new one.                                                                                                                                                                                |
| `SizeHint` as an exact O(n) pre-pass                                                        | `SizeHint`/`sizeHintMode`                     | Lets a caller guarantee zero reallocation with a single `make` up front, without over-allocating.                                                                                                                                                                                                                                                                          |
| No `map`, no `interface{}`/`any`, no reflection anywhere in the hot path                    | whole package                                 | Removes hashing cost, removes indirect/virtual call overhead, and lets the Go compiler devirtualize and potentially inline every call in the loop.                                                                                                                                                                                                                         |
| No recursion                                                                                | whole package                                 | Every function is a straight loop; no call-stack growth, no risk of stack-depth-dependent latency spikes on pathological input.                                                                                                                                                                                                                                            |
| Immutable package-level state only, no `init()` mutation, no `sync.*`                       | `tables_gen.go` (generated)                   | Makes every exported function trivially safe for unlimited concurrent callers with zero synchronization cost — reads of immutable memory never need a lock.                                                                                                                                                                                                                |
| Manual `utf8.DecodeRune`/`DecodeRuneInString` loop instead of `for range`                   | `AppendMode`/`AppendBytesMode`/`sizeHintMode` | `for range` over a string already uses the same decoder internally, but a manual loop lets the ASCII fast path run _before_ any decode call (checking the raw byte first), and lets malformed-byte handling be explicit (needed to match the original's exact drop-one-byte-and-continue behavior, and to support `ModeKeep`'s "preserve the raw invalid byte" semantics). |
| Two explicit `Mode`s (`ModeSkip`/`ModeKeep`) via a `uint8` constant, not a closure/callback | `unidecode.go`                                | Satisfies the "unknown characters pass through or are skipped, depending on configuration" requirement without introducing interface dispatch or a function-pointer indirection in the hot loop.                                                                                                                                                                           |
| Generator-time data preparation (`go generate`), zero runtime parsing                       | `gen/main.go` → `tables_gen.go`               | All of the above structure is _built once, at generation time_, and committed as plain Go array literals — there is no CSV/JSON parsing, no `init()`-time table construction, and no runtime cost at all associated with "loading" the data; it's simply linked into the binary as `.rodata`.                                                                              |

---

## 11. API reference

```go
// Core, zero-allocation (given sufficient dst capacity) API.
func Append(dst []byte, src string) []byte
func AppendBytes(dst []byte, src []byte) []byte

// Same, with explicit control over unmapped-codepoint behavior.
type Mode uint8
const (
	ModeSkip Mode = iota // drop unmapped codepoints (matches the original library)
	ModeKeep              // pass unmapped codepoints through unchanged
)
func AppendMode(dst []byte, src string, mode Mode) []byte
func AppendBytesMode(dst []byte, src []byte, mode Mode) []byte

// Sizing and validation helpers.
func SizeHint(src string) int
func SizeHintMode(src string, mode Mode) int
func ValidUTF8(src string) bool

// Convenience wrapper matching the original library's signature.
// Allocates exactly once (sized exactly via SizeHint).
func Unidecode(s string) string
```

---

## 12. Repository layout

```
translit/
├── go.mod                # zero external dependencies
├── doc.go                # package doc + //go:generate directive
├── unidecode.go          # the entire runtime implementation (~260 lines)
├── tables_gen.go         # generated, DO NOT EDIT — immutable lookup tables
├── unidecode_test.go     # unit tests, edge cases, zero-allocation regression tests
├── fuzz_test.go          # FuzzAppend, FuzzValidUTF8
├── property_test.go      # property-based tests + concurrency tests
├── gen/
│   ├── go.mod            # isolated module: depends on the original library
│   └── main.go           # table generator (`go generate ./...` from the module root)
└── compare/
    ├── go.mod              # isolated module: depends on both libraries
    ├── golden_test.go      # exhaustive + real-world correctness vs. the original
    └── benchmark_test.go   # head-to-head benchmarks vs. the original
```

The generator and comparison/benchmark code live in their own Go
modules (linked back to the main module via a `replace` directive) so
that the main `unidecode` module itself has **zero external
dependencies** — `go.mod` for the library proper lists nothing beyond
the Go standard library.

## 13. Regenerating the tables

```sh
go generate ./...
```

This re-derives `tables_gen.go` from the authoritative data (currently
vendored from `github.com/mozillazg/go-unidecode`'s generated table
package — see `gen/main.go`'s doc comment for exactly what it reads and
how it flattens it). Regeneration is deterministic: running it twice in
a row produces a byte-identical file both times.
