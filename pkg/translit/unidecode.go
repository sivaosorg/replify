package translit

import "unicode/utf8"

// Mode controls how [Append], [AppendBytes], [AppendMode], [AppendBytesMode],
// [SizeHint], and [SizeHintMode] handle a codepoint that has no entry in
// the transliteration table — either because it falls outside the table's
// covered range (above U+EFFFF), or because its specific slot within a
// covered section was never assigned a replacement.
type Mode uint8

const (
	// ModeSkip drops any codepoint for which no transliteration exists.
	// It is the zero value and the default used by [Append], [AppendBytes],
	// [SizeHint], and [Unidecode]. This matches the behavior of the original
	// github.com/mozillazg/go-unidecode library.
	ModeSkip Mode = iota

	// ModeKeep passes any codepoint for which no transliteration exists
	// through to the output unchanged, re-encoded as valid UTF-8. For a
	// malformed input byte that cannot be decoded as a rune (utf8.RuneError
	// with size ≤ 1), the raw byte is emitted as-is.
	ModeKeep
)

const (
	// maxRune is the highest codepoint the transliteration table covers.
	// This mirrors the original library's r > 0xeffff cutoff.
	maxRune = 0xeffff

	// asciiFastPathMax mirrors the original library's `r < unicode.MaxASCII`
	// check (unicode.MaxASCII == '\u007f', i.e. 127). Bytes below this value
	// are guaranteed identity-mapped and go through the fast path with no
	// table lookup. Note this deliberately excludes 0x7f (DEL) itself: the
	// original library does not fast-path DEL, it runs DEL through the
	// table like any other codepoint -- and section 0's table entry for DEL
	// is empty, so DEL is silently dropped in ModeSkip. Matching that exact
	// cutoff (rather than the more obvious utf8.RuneSelf == 128) keeps this
	// implementation byte-for-byte identical to the original across the
	// entire codepoint space.
	asciiFastPathMax = 0x7f
)

// Append transliterates src into 7-bit ASCII and appends the result to
// dst, returning the extended buffer. Codepoints with no transliteration
// are silently dropped ([ModeSkip]); use [AppendMode] to keep them.
//
// The function follows Go's standard append convention: dst is grown only
// if it lacks sufficient spare capacity, and the returned slice may share
// dst's backing array. When dst is pre-sized via [SizeHint], no allocation
// occurs for any input, regardless of length or script mix:
//
//	buf := make([]byte, 0, translit.SizeHint(src))
//	buf = translit.Append(buf, src)
//
// DEL (U+007F) is dropped, matching the original library's
// unicode.MaxASCII fast-path boundary. All other ASCII bytes (U+0000–U+007E)
// pass through unchanged.
//
// Append is safe to call concurrently from any number of goroutines; all
// package-level lookup state is immutable after package initialization.
func Append(dst []byte, src string) []byte {
	return AppendMode(dst, src, ModeSkip)
}

// AppendBytes is [Append] for a []byte source. It avoids a string
// conversion when the caller already holds the input as a byte slice,
// and shares the same zero-allocation and concurrency guarantees.
// Codepoints with no transliteration are dropped; use [AppendBytesMode]
// to keep them instead.
func AppendBytes(dst []byte, src []byte) []byte {
	return AppendBytesMode(dst, src, ModeSkip)
}

// AppendMode transliterates src into 7-bit ASCII and appends the result
// to dst, using mode to determine the treatment of codepoints that have
// no transliteration entry. It returns the extended buffer.
//
// The mode argument governs three distinct codepoint classes:
//
//   - Codepoints outside the covered range (above U+EFFFF): dropped by
//     [ModeSkip]; re-encoded as UTF-8 and appended by [ModeKeep].
//   - Codepoints within a covered section whose slot was never assigned
//     a replacement: same as above.
//   - Codepoints whose table entry is explicitly empty (combining marks
//     and similar): always dropped in both modes, consistent with the
//     original library's behavior of appending an empty string.
//
// AppendMode carries the same zero-allocation and concurrency guarantees
// as [Append]. For the common [ModeSkip] case, prefer the shorter [Append].
func AppendMode(dst []byte, src string, mode Mode) []byte {
	i := 0
	n := len(src)
	for i < n {
		c := src[i]

		// ASCII fast path: a single comparison identifies the
		// overwhelming majority of bytes in typical text, and they
		// are appended directly with no table lookup at all.
		if c < asciiFastPathMax {
			dst = append(dst, c)
			i++
			continue
		}

		r, size := utf8.DecodeRuneInString(src[i:])
		if r == utf8.RuneError && size <= 1 {
			if mode == ModeKeep {
				dst = append(dst, c)
			}
			i++
			continue
		}

		dst = appendRune(dst, r, mode)
		i += size
	}
	return dst
}

// AppendBytesMode transliterates src into 7-bit ASCII and appends the
// result to dst, using mode to determine the treatment of codepoints that
// have no transliteration entry. It returns the extended buffer.
//
// It is the []byte-source counterpart of [AppendMode]; the two functions
// are identical in behavior except for the source type. See [AppendMode]
// for a full description of mode semantics.
//
// AppendBytesMode carries the same zero-allocation and concurrency
// guarantees as [Append].
func AppendBytesMode(dst []byte, src []byte, mode Mode) []byte {
	i := 0
	n := len(src)
	for i < n {
		c := src[i]

		if c < asciiFastPathMax {
			dst = append(dst, c)
			i++
			continue
		}

		r, size := utf8.DecodeRune(src[i:])
		if r == utf8.RuneError && size <= 1 {
			if mode == ModeKeep {
				dst = append(dst, c)
			}
			i++
			continue
		}

		dst = appendRune(dst, r, mode)
		i += size
	}
	return dst
}

// SizeHint returns the exact number of bytes [Append](nil, src) would
// produce — equivalently, len([AppendMode](nil, src, [ModeSkip])). It
// performs the same O(1)-per-rune table walk as [Append] but writes
// nothing, enabling a single guaranteed-zero-allocation call:
//
//	buf := make([]byte, 0, translit.SizeHint(src))
//	buf = translit.Append(buf, src)
//
// SizeHint is safe to call concurrently and performs no allocation.
func SizeHint(src string) int {
	return sizeHintMode(src, ModeSkip)
}

// SizeHintMode returns the exact number of bytes [AppendMode](nil, src,
// mode) would produce. It is [SizeHint] for callers using a non-default
// mode, and enables a guaranteed zero-allocation call to [AppendMode]
// or [AppendBytesMode]:
//
//	buf := make([]byte, 0, translit.SizeHintMode(src, translit.ModeKeep))
//	buf = translit.AppendMode(buf, src, translit.ModeKeep)
func SizeHintMode(src string, mode Mode) int {
	return sizeHintMode(src, mode)
}

// ValidUTF8 reports whether src is valid UTF-8. It is an allocation-free
// wrapper around [utf8.ValidString], provided so callers can validate input
// without an additional import. Passing invalid UTF-8 to [Append] is safe —
// malformed bytes are handled gracefully per-mode — but callers that need
// to surface encoding errors before transliterating can do so here.
func ValidUTF8(src string) bool {
	return utf8.ValidString(src)
}

// Unidecode transliterates s and returns the result as a new string.
// It is a convenience wrapper that matches the original
// github.com/mozillazg/go-unidecode API, suitable for call sites that
// deal in strings rather than byte slices.
//
// Unlike [Append], Unidecode must allocate: a Go string's backing array
// is immutable and caller-owned. It performs exactly one allocation,
// sized precisely via [SizeHint] so it never over-allocates or
// reallocates internally. For high-throughput paths where a []byte result
// is acceptable, prefer [Append] with a reused buffer.
func Unidecode(s string) string {
	buf := make([]byte, 0, SizeHint(s))
	buf = Append(buf, s)
	return string(buf)
}

// sizeHintMode is the shared implementation of [SizeHint] and [SizeHintMode].
// It mirrors [AppendMode]'s loop exactly — same ASCII fast path, same table
// walk, same mode semantics — but accumulates a byte count rather than
// writing to a buffer. Keeping the logic parallel ensures [SizeHint] never
// drifts out of sync with [Append].
func sizeHintMode(src string, mode Mode) int {
	total := 0
	i := 0
	n := len(src)
	for i < n {
		c := src[i]
		if c < asciiFastPathMax {
			total++
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(src[i:])
		if r == utf8.RuneError && size <= 1 {
			if mode == ModeKeep {
				total++
			}
			i++
			continue
		}
		total += runeSize(r, mode)
		i += size
	}
	return total
}

// appendRune resolves the transliteration for a single non-ASCII rune
// and appends it to dst. The table walk is unconditionally O(1): two
// reads from the section-level arrays (sectionStart, sectionLen), two
// from the entry-level arrays (entryOff, entryLen), and one slice of
// the shared data constant — no hashing, no pointer chasing, no
// allocation, no interface dispatch.
func appendRune(dst []byte, r rune, mode Mode) []byte {
	if r < 0 || r > maxRune {
		if mode == ModeKeep {
			return utf8.AppendRune(dst, r)
		}
		return dst
	}

	section := int(uint32(r) >> 8)
	if section > maxSection {
		if mode == ModeKeep {
			return utf8.AppendRune(dst, r)
		}
		return dst
	}

	start := sectionStart[section]
	if start < 0 {
		if mode == ModeKeep {
			return utf8.AppendRune(dst, r)
		}
		return dst
	}

	pos := uint32(r) & 0xff
	if pos >= uint32(sectionLen[section]) {
		if mode == ModeKeep {
			return utf8.AppendRune(dst, r)
		}
		return dst
	}

	idx := uint32(start) + pos
	off := entryOff[idx]
	ln := entryLen[idx]
	if ln == 0 {
		// A present-but-empty entry (a combining mark or similar with
		// no ASCII equivalent) is always dropped, in both modes -- this
		// matches the original library, which would append "" here.
		return dst
	}
	return append(dst, data[off:off+uint32(ln)]...)
}

// runeSize reports how many output bytes [appendRune] would produce for
// r under the given mode, without writing anything. It mirrors
// [appendRune]'s full lookup logic so that [SizeHint] remains exact
// regardless of which code path [appendRune] would take.
func runeSize(r rune, mode Mode) int {
	if r < asciiFastPathMax {
		return 1
	}
	if r < 0 || r > maxRune {
		if mode == ModeKeep {
			return utf8.RuneLen(r)
		}
		return 0
	}
	section := int(uint32(r) >> 8)
	if section > maxSection {
		if mode == ModeKeep {
			return utf8.RuneLen(r)
		}
		return 0
	}
	start := sectionStart[section]
	if start < 0 {
		if mode == ModeKeep {
			return utf8.RuneLen(r)
		}
		return 0
	}
	pos := uint32(r) & 0xff
	if pos >= uint32(sectionLen[section]) {
		if mode == ModeKeep {
			return utf8.RuneLen(r)
		}
		return 0
	}
	idx := uint32(start) + pos
	return int(entryLen[idx])
}
