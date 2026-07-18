// package translit transliterates Unicode text into plain 7-bit ASCII,
// e.g. Append(nil, "kožušček") produces "kozuscek".
//
// This is a from-scratch, performance-oriented reimplementation of
// github.com/mozillazg/go-unidecode. The transliteration data (which
// characters map to which ASCII replacement) is unchanged -- it is the
// same authoritative table, sourced from that project's generated Go
// source (itself derived from the Python "unidecode" package). What is
// different is everything about how that data is stored and how the
// hot path walks it:
//
//   - The caller supplies the destination buffer (Append/AppendBytes),
//     so the library performs zero heap allocations of its own during
//     transliteration. There is no strings.Builder, no bytes.Buffer, no
//     intermediate []rune or []byte conversion of the input.
//   - Every codepoint lookup is a fixed number of array index operations
//     into immutable, compile-time-generated arrays -- never a map, never
//     a slice-of-slices, never reflection or an interface call.
//   - ASCII bytes (the overwhelming majority of most real-world text)
//     are detected with a single comparison and appended directly, with
//     no table lookup at all.
//   - All package-level state is immutable (const/array literals with no
//     init-time mutation), so every exported function is safe to call
//     from any number of goroutines concurrently with no locking of any
//     kind.
//
// See ARCHITECTURE.md in the repository for a full write-up of the data
// layout, the complexity/cache analysis, and a comparison against the
// original implementation.
package translit

//go:generate sh -c "cd gen && go run ."
