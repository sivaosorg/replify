package translit

import (
	"math/rand"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"
)

// TestPropertyConcatenation checks that transliteration distributes over
// concatenation: Append(a+b) == Append(a) followed by Append(b). This
// must hold because the algorithm is a stateless, per-rune map -- there
// is no cross-rune context that could make it fail, and this test
// guards against ever accidentally introducing any.
func TestPropertyConcatenation(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	pool := []string{
		"a", "Z", "0", " ", "!", "kožušček", "北京", "日本語", "Москва",
		"Ελλάδα", "café", "naïve", "한국어", "", "\x80", "\xff",
	}
	for i := 0; i < 500; i++ {
		a := pool[rng.Intn(len(pool))]
		b := pool[rng.Intn(len(pool))]
		got := string(Append(nil, a+b))
		want := string(Append(nil, a)) + string(Append(nil, b))
		if got != want {
			t.Fatalf("Append(%q+%q) = %q, want %q (concatenation property violated)", a, b, got, want)
		}
	}
}

// TestPropertyASCIISubsetIsIdentity checks that any string consisting
// solely of ASCII bytes passes through completely unchanged.
func TestPropertyASCIISubsetIsIdentity(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 200; i++ {
		n := rng.Intn(64)
		var b strings.Builder
		for j := 0; j < n; j++ {
			b.WriteByte(byte(rng.Intn(127))) // [0,126]: excludes DEL (0x7f), which the table drops
		}
		s := b.String()
		got := string(Append(nil, s))
		if got != s {
			t.Fatalf("Append(%q) = %q, want identity for pure-ASCII input", s, got)
		}
	}
}

// TestPropertyOutputNeverLongerThanKeepEncoding checks that ModeSkip
// output is never longer than the ModeKeep output for the same input:
// skipping can only omit bytes relative to keeping, never add them,
// since the transliteration tables are shared between modes and only
// the unmapped-codepoint fallback differs.
func TestPropertyOutputNeverLongerThanKeepEncoding(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	pool := []rune{'a', 'Z', '0', ' ', 'ž', 'č', '北', '京', 'Ж', 0xf0000, 0x10ffff}
	for i := 0; i < 200; i++ {
		n := rng.Intn(20)
		var b strings.Builder
		for j := 0; j < n; j++ {
			b.WriteRune(pool[rng.Intn(len(pool))])
		}
		s := b.String()
		skip := Append(nil, s)
		keep := AppendMode(nil, s, ModeKeep)
		if len(skip) > len(keep) {
			t.Fatalf("ModeSkip output longer than ModeKeep for %q: %d > %d", s, len(skip), len(keep))
		}
	}
}

// TestConcurrentAppend exercises Append from many goroutines
// simultaneously against a shared set of inputs, verifying no data race
// (run with -race) and that every goroutine sees consistent output --
// i.e. the package's global lookup tables are safe for concurrent read
// access with no locking.
func TestConcurrentAppend(t *testing.T) {
	inputs := []string{
		"kožušček", "北京市", "日本語のテスト", "Москва", "Ελληνικά",
		"café naïve", "한국어 테스트 문자열", "plain ascii text",
	}
	want := make([]string, len(inputs))
	for i, s := range inputs {
		want[i] = string(Append(nil, s))
	}

	const goroutines = 64
	const iterations = 200

	var wg sync.WaitGroup
	errs := make(chan string, goroutines)
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(seed int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(int64(seed)))
			for i := 0; i < iterations; i++ {
				idx := rng.Intn(len(inputs))
				got := string(Append(nil, inputs[idx]))
				if got != want[idx] {
					errs <- got
					return
				}
			}
		}(g)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatalf("concurrent Append produced unexpected output: %q", e)
	}
}

// TestConcurrentAppendSharedDstIsolated confirms that concurrent callers
// using their own dst slices never observe cross-goroutine corruption
// (each goroutine's buffer is its own; the package holds no mutable
// shared state that a concurrent caller could clobber).
func TestConcurrentAppendSharedDstIsolated(t *testing.T) {
	const goroutines = 32
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			s := strings.Repeat("北京", id%5+1)
			want := string(Append(nil, s))
			for i := 0; i < 100; i++ {
				got := string(Append(nil, s))
				if got != want {
					t.Errorf("goroutine %d: got %q, want %q", id, got, want)
					return
				}
			}
		}(g)
	}
	wg.Wait()
}

func TestAppendASCII(t *testing.T) {
	cases := []string{
		"",
		"hello world",
		"Hello, World! 123",
		"The quick brown fox jumps over the lazy dog.",
		string([]byte{0, 1, 2, 3, 0x1f, 0x20, 0x7e}), // control chars and printable ASCII up to but not including DEL
	}
	for _, c := range cases {
		got := string(Append(nil, c))
		if got != c {
			t.Errorf("Append(%q) = %q, want %q (ASCII must pass through unchanged)", c, got, c)
		}
	}
}

// TestDELIsDropped documents a subtle but deliberate compatibility
// point: the original library's ASCII fast path uses `r <
// unicode.MaxASCII` (0x7f), so DEL (0x7f) itself is not fast-pathed --
// it goes through the transliteration table, whose entry for DEL is
// empty, so DEL is dropped. This rewrite matches that exactly rather
// than the more "obvious" utf8.RuneSelf (0x80) cutoff.
func TestDELIsDropped(t *testing.T) {
	got := string(Append(nil, "a\x7fb"))
	want := "ab"
	if got != want {
		t.Errorf("Append(%q) = %q, want %q (DEL should be dropped, matching the original library)", "a\x7fb", got, want)
	}
}

func TestAppendKnownTransliterations(t *testing.T) {
	cases := []struct{ in, want string }{
		{"kožušček", "kozuscek"},
		{"café", "cafe"},
		{"naïve", "naive"},
		{"Björk", "Bjork"},
		{"北京", "Bei Jing "},
		{"日本語", "Ri Ben Yu "},
		{"Москва", "Moskva"},
		{"Ελλάδα", "Ellada"},
		{"日本語 hello", "Ri Ben Yu  hello"},
	}
	for _, c := range cases {
		got := string(Append(nil, c.in))
		if got != c.want {
			t.Errorf("Append(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAppendGrowsExistingBuffer(t *testing.T) {
	dst := []byte("prefix:")
	got := string(Append(dst, "abc"))
	want := "prefix:abc"
	if got != want {
		t.Errorf("Append with existing prefix = %q, want %q", got, want)
	}
}

func TestAppendReturnedSliceSharesBackingArrayWhenCapacityAllows(t *testing.T) {
	dst := make([]byte, 0, 64)
	out := Append(dst, "hello")
	// Write a canary into the original backing array via dst's capacity
	// region and confirm out sees it -- proves no reallocation occurred.
	dst = dst[:cap(dst)]
	dst[0] = 'X'
	if out[:1][0] != 'X' {
		t.Errorf("expected Append to reuse dst's backing array when capacity allows")
	}
}

func TestAppendBytesMatchesAppend(t *testing.T) {
	inputs := []string{"hello", "kožušček", "北京市", "", "a\x80b"} // \x80 is invalid UTF-8 continuation byte alone
	for _, in := range inputs {
		want := string(Append(nil, in))
		got := string(AppendBytes(nil, []byte(in)))
		if got != want {
			t.Errorf("AppendBytes(%q) = %q, want %q (mismatch with Append)", in, got, want)
		}
	}
}

func TestInvalidUTF8(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a\x80b", "ab"}, // lone continuation byte dropped
		{"\xff\xfe", ""}, // two invalid bytes dropped
		{"valid\xc3\x28invalid", "valid(invalid"}, // \xc3 alone is invalid; the following ASCII '(' still passes through
	}
	for _, c := range cases {
		got := string(Append(nil, c.in))
		if got != c.want {
			t.Errorf("Append(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestModeKeep(t *testing.T) {
	// A rune with no table entry (e.g. a private-use or unassigned
	// codepoint) should be dropped under ModeSkip and preserved under
	// ModeKeep.
	r := rune(0xF0000) // supplementary private-use area, beyond maxRune
	in := "a" + string(r) + "b"

	skip := string(AppendMode(nil, in, ModeSkip))
	if skip != "ab" {
		t.Errorf("ModeSkip: got %q, want %q", skip, "ab")
	}

	keep := string(AppendMode(nil, in, ModeKeep))
	if keep != in {
		t.Errorf("ModeKeep: got %q, want %q", keep, in)
	}
}

func TestModeKeepInvalidByte(t *testing.T) {
	in := "a\x80b"
	got := string(AppendMode(nil, in, ModeKeep))
	if got != in {
		t.Errorf("ModeKeep with invalid byte: got %q, want %q", got, in)
	}
}

func TestSizeHintExact(t *testing.T) {
	cases := []string{
		"",
		"hello world",
		"kožušček",
		"北京市 nice to meet you 日本語",
		strings.Repeat("北", 500),
		"a\x80b\xffvalid",
	}
	for _, c := range cases {
		want := len(Append(nil, c))
		got := SizeHint(c)
		if got != want {
			t.Errorf("SizeHint(%q) = %d, want %d", c, got, want)
		}
	}
}

func TestSizeHintModeKeep(t *testing.T) {
	cases := []string{
		"",
		"hello",
		string(rune(0xF0000)),
		"a\x80b",
	}
	for _, c := range cases {
		want := len(AppendMode(nil, c, ModeKeep))
		got := SizeHintMode(c, ModeKeep)
		if got != want {
			t.Errorf("SizeHintMode(%q, ModeKeep) = %d, want %d", c, got, want)
		}
	}
}

func TestValidUTF8(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"hello", true},
		{"kožušček", true},
		{"北京", true},
		{"a\x80b", false},
		{"", true},
	}
	for _, c := range cases {
		got := ValidUTF8(c.in)
		if got != c.want {
			t.Errorf("ValidUTF8(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestUnidecodeConvenience(t *testing.T) {
	got := Unidecode("kožušček")
	want := "kozuscek"
	if got != want {
		t.Errorf("Unidecode(%q) = %q, want %q", "kožušček", got, want)
	}
}

func TestRuneAtMaxBoundary(t *testing.T) {
	// maxRune itself and one above it: the +1 must be skipped/kept
	// according to mode, never crash or index out of range.
	within := string(rune(maxRune))
	beyond := string(rune(maxRune + 1))

	_ = Append(nil, within)
	_ = Append(nil, beyond)
	_ = AppendMode(nil, beyond, ModeKeep)
}

func TestAllRunesNoPanic(t *testing.T) {
	// Sweep a broad sample of the codepoint space (every rune would be
	// slow in short mode) to make sure appendRune/runeSize never index
	// out of range regardless of section/position combination.
	for r := rune(0); r < 0x110000; r += 37 { // coprime-ish stride for coverage
		if !utf8.ValidRune(r) {
			continue
		}
		s := string(r)
		_ = Append(nil, s)
		_ = AppendMode(nil, s, ModeKeep)
		_ = SizeHint(s)
		_ = SizeHintMode(s, ModeKeep)
	}
}

func TestEmptyInput(t *testing.T) {
	if got := Append(nil, ""); len(got) != 0 {
		t.Errorf("Append(nil, \"\") = %q, want empty", got)
	}
	if got := SizeHint(""); got != 0 {
		t.Errorf("SizeHint(\"\") = %d, want 0", got)
	}
}

func TestAppendIsDeterministic(t *testing.T) {
	in := "kožušček 北京市 Москва café naïve"
	first := string(Append(nil, in))
	for i := 0; i < 100; i++ {
		got := string(Append(nil, in))
		if got != first {
			t.Fatalf("iteration %d: Append(%q) = %q, want %q (non-deterministic output)", i, in, got, first)
		}
	}
}

// TestZeroAllocations is a permanent regression guard for the library's
// core promise: when dst already has enough spare capacity, Append,
// AppendBytes, SizeHint, and ValidUTF8 must not allocate on the heap at
// all, for any mix of ASCII and non-ASCII input.
func TestZeroAllocations(t *testing.T) {
	inputs := []string{
		strings.Repeat("The quick brown fox jumps over the lazy dog. ", 20),
		strings.Repeat("kožušček Příliš žluťoučký kůň úpěl ďábelské ódy ", 20),
		strings.Repeat("北京市朝阳区建国路甲92号 ", 20),
		strings.Repeat("日本語のテキストです ", 20),
		strings.Repeat("한국어 텍스트입니다 ", 20),
		strings.Repeat("mixed 日本語 and English and Ελληνικά and Русский ", 20),
	}

	for _, in := range inputs {
		in := in
		buf := make([]byte, 0, SizeHint(in)+64) // ample spare capacity
		bytesIn := []byte(in)

		allocs := testing.AllocsPerRun(200, func() {
			buf = Append(buf[:0], in)
		})
		if allocs != 0 {
			t.Errorf("Append(%.20q...) allocated %.1f times per call, want 0", in, allocs)
		}

		allocs = testing.AllocsPerRun(200, func() {
			buf = AppendBytes(buf[:0], bytesIn)
		})
		if allocs != 0 {
			t.Errorf("AppendBytes(%.20q...) allocated %.1f times per call, want 0", in, allocs)
		}

		allocs = testing.AllocsPerRun(200, func() {
			_ = SizeHint(in)
		})
		if allocs != 0 {
			t.Errorf("SizeHint(%.20q...) allocated %.1f times per call, want 0", in, allocs)
		}

		allocs = testing.AllocsPerRun(200, func() {
			_ = ValidUTF8(in)
		})
		if allocs != 0 {
			t.Errorf("ValidUTF8(%.20q...) allocated %.1f times per call, want 0", in, allocs)
		}
	}
}

// TestZeroAllocationsModeKeep confirms the same zero-allocation
// guarantee holds under ModeKeep, including for input containing
// codepoints outside the transliteration table (which must be
// re-encoded via utf8.AppendRune -- itself allocation-free since it
// operates on the caller-provided slice).
func TestZeroAllocationsModeKeep(t *testing.T) {
	in := "abc" + string(rune(0xf0000)) + "def" + string(rune(0x10ffff))
	buf := make([]byte, 0, 256)

	allocs := testing.AllocsPerRun(200, func() {
		buf = AppendMode(buf[:0], in, ModeKeep)
	})
	if allocs != 0 {
		t.Errorf("AppendMode(ModeKeep) allocated %.1f times per call, want 0", allocs)
	}
}
