package codegen

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

// ─────────────────────────────────────────────
// New()
// ─────────────────────────────────────────────

func TestNew_Defaults(t *testing.T) {
	g, err := New()
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	opts := g.Options()
	if opts.Length != 8 {
		t.Errorf("default length: want 8, got %d", opts.Length)
	}
	if opts.Charset != CharsetAlphanumeric {
		t.Errorf("default charset: want %q, got %q", CharsetAlphanumeric, opts.Charset)
	}
	if opts.Prefix != "" {
		t.Errorf("default prefix: want \"\", got %q", opts.Prefix)
	}
	if opts.Suffix != "" {
		t.Errorf("default suffix: want \"\", got %q", opts.Suffix)
	}
}

func TestNew_WithOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		wantOpts Options
	}{
		{
			name:     "with length",
			opts:     []Option{WithLength(12)},
			wantOpts: Options{Length: 12, Charset: CharsetAlphanumeric},
		},
		{
			name:     "with charset numeric",
			opts:     []Option{WithCharset(CharsetNumeric)},
			wantOpts: Options{Length: 8, Charset: CharsetNumeric},
		},
		{
			name:     "with charset alpha upper",
			opts:     []Option{WithCharset(CharsetAlphaUpper)},
			wantOpts: Options{Length: 8, Charset: CharsetAlphaUpper},
		},
		{
			name:     "with prefix",
			opts:     []Option{WithPrefix("ORD-")},
			wantOpts: Options{Length: 8, Charset: CharsetAlphanumeric, Prefix: "ORD-"},
		},
		{
			name:     "with suffix",
			opts:     []Option{WithSuffix("-VN")},
			wantOpts: Options{Length: 8, Charset: CharsetAlphanumeric, Suffix: "-VN"},
		},
		{
			name: "with all options combined",
			opts: []Option{
				WithLength(10),
				WithCharset(CharsetAlphaUpper),
				WithPrefix("ORD-"),
				WithSuffix("-2024"),
			},
			wantOpts: Options{
				Length:  10,
				Charset: CharsetAlphaUpper,
				Prefix:  "ORD-",
				Suffix:  "-2024",
			},
		},
		{
			name:     "with custom charset",
			opts:     []Option{WithCustomCharset("ABC123")},
			wantOpts: Options{Length: 8, Charset: "ABC123"},
		},
		{
			name:     "with custom charset deduplication",
			opts:     []Option{WithCustomCharset("AABBCC123")},
			wantOpts: Options{Length: 8, Charset: "ABC123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := New(tt.opts...)
			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}
			opts := g.Options()
			if opts.Length != tt.wantOpts.Length {
				t.Errorf("length: want %d, got %d", tt.wantOpts.Length, opts.Length)
			}
			if opts.Charset != tt.wantOpts.Charset {
				t.Errorf("charset: want %q, got %q", tt.wantOpts.Charset, opts.Charset)
			}
			if opts.Prefix != tt.wantOpts.Prefix {
				t.Errorf("prefix: want %q, got %q", tt.wantOpts.Prefix, opts.Prefix)
			}
			if opts.Suffix != tt.wantOpts.Suffix {
				t.Errorf("suffix: want %q, got %q", tt.wantOpts.Suffix, opts.Suffix)
			}
		})
	}
}

func TestNew_InvalidOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr error
	}{
		{
			name:    "zero length",
			opts:    []Option{WithLength(0)},
			wantErr: ErrInvalidLength,
		},
		{
			name:    "negative length",
			opts:    []Option{WithLength(-5)},
			wantErr: ErrInvalidLength,
		},
		{
			name:    "empty charset via WithCharset",
			opts:    []Option{WithCharset("")},
			wantErr: ErrEmptyCharset,
		},
		{
			name:    "empty charset via WithCustomCharset",
			opts:    []Option{WithCustomCharset("")},
			wantErr: ErrEmptyCharset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.opts...)
			if err == nil {
				t.Fatal("New() expected error, got nil")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error: want %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// ─────────────────────────────────────────────
// MustNew()
// ─────────────────────────────────────────────

func TestMustNew_Success(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustNew() unexpected panic: %v", r)
		}
	}()
	g := MustNew(WithLength(10), WithCharset(CharsetAlphanumericUpper))
	if g == nil {
		t.Fatal("MustNew() returned nil")
	}
}

func TestMustNew_PanicOnInvalidLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustNew() expected panic, got none")
		}
	}()
	MustNew(WithLength(0))
}

func TestMustNew_PanicOnEmptyCharset(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustNew() expected panic, got none")
		}
	}()
	MustNew(WithCharset(""))
}

// ─────────────────────────────────────────────
// Generate()
// ─────────────────────────────────────────────

func TestGenerate_OutputLength(t *testing.T) {
	tests := []struct {
		name         string
		opts         []Option
		wantTotalLen int
	}{
		{"default (8 chars)", nil, 8},
		{"length 1", []Option{WithLength(1)}, 1},
		{"length 16", []Option{WithLength(16)}, 16},
		{"length 32", []Option{WithLength(32)}, 32},
		{"with 4-char prefix", []Option{WithLength(8), WithPrefix("ORD-")}, 12},
		{"with 3-char suffix", []Option{WithLength(8), WithSuffix("-VN")}, 11},
		{"with prefix+suffix", []Option{WithLength(8), WithPrefix("P-"), WithSuffix("-S")}, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := New(tt.opts...)
			if err != nil {
				t.Fatalf("New() error: %v", err)
			}
			code, err := g.Generate()
			if err != nil {
				t.Fatalf("Generate() error: %v", err)
			}
			got := len([]rune(code))
			if got != tt.wantTotalLen {
				t.Errorf("total length: want %d, got %d (code=%q)", tt.wantTotalLen, got, code)
			}
		})
	}
}

func TestGenerate_CharsetAdherence(t *testing.T) {
	charsets := []struct {
		name    string
		charset Charset
	}{
		{"numeric", CharsetNumeric},
		{"alpha lower", CharsetAlphaLower},
		{"alpha upper", CharsetAlphaUpper},
		{"alpha mixed", CharsetAlpha},
		{"alphanumeric", CharsetAlphanumeric},
		{"alphanumeric upper", CharsetAlphanumericUpper},
		{"alphanumeric lower", CharsetAlphanumericLower},
		{"custom", Charset("XYZ!@#789")},
	}

	for _, tc := range charsets {
		t.Run(tc.name, func(t *testing.T) {
			g, err := New(WithLength(200), WithCharset(tc.charset))
			if err != nil {
				t.Fatalf("New() error: %v", err)
			}
			code, err := g.Generate()
			if err != nil {
				t.Fatalf("Generate() error: %v", err)
			}
			for _, ch := range code {
				if !strings.ContainsRune(string(tc.charset), ch) {
					t.Errorf("code contains %q which is not in charset %q", ch, tc.charset)
				}
			}
		})
	}
}

func TestGenerate_HasPrefix(t *testing.T) {
	prefix := "ORD-"
	g := MustNew(WithLength(8), WithPrefix(prefix))
	for i := 0; i < 20; i++ {
		code, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}
		if !strings.HasPrefix(code, prefix) {
			t.Errorf("code %q does not start with prefix %q", code, prefix)
		}
	}
}

func TestGenerate_HasSuffix(t *testing.T) {
	suffix := "-VN"
	g := MustNew(WithLength(8), WithSuffix(suffix))
	for i := 0; i < 20; i++ {
		code, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}
		if !strings.HasSuffix(code, suffix) {
			t.Errorf("code %q does not end with suffix %q", code, suffix)
		}
	}
}

func TestGenerate_RandomPart_Length(t *testing.T) {
	prefix := "ORD-"
	suffix := "-VN"
	length := 10
	g := MustNew(WithLength(length), WithPrefix(prefix), WithSuffix(suffix))

	code, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	inner := strings.TrimPrefix(strings.TrimSuffix(code, suffix), prefix)
	if len([]rune(inner)) != length {
		t.Errorf("random part length: want %d, got %d (inner=%q)", length, len([]rune(inner)), inner)
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	g := MustNew(WithLength(12))
	n := 1000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		code, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate() error at iteration %d: %v", i, err)
		}
		if _, dup := seen[code]; dup {
			t.Errorf("duplicate code detected: %q", code)
		}
		seen[code] = struct{}{}
	}
}

func TestGenerate_SingleCharCharset(t *testing.T) {
	g := MustNew(WithLength(5), WithCharset("A"))
	code, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if code != "AAAAA" {
		t.Errorf("want %q, got %q", "AAAAA", code)
	}
}

// ─────────────────────────────────────────────
// GenerateN()
// ─────────────────────────────────────────────

func TestGenerateN_ValidCounts(t *testing.T) {
	g := MustNew(WithLength(8))
	for _, n := range []int{1, 5, 50, 500} {
		t.Run("", func(t *testing.T) {
			codes, err := g.GenerateN(n)
			if err != nil {
				t.Fatalf("GenerateN(%d) error: %v", n, err)
			}
			if len(codes) != n {
				t.Errorf("GenerateN(%d): want %d codes, got %d", n, n, len(codes))
			}
			for i, code := range codes {
				if len([]rune(code)) != 8 {
					t.Errorf("codes[%d] length: want 8, got %d (code=%q)", i, len([]rune(code)), code)
				}
			}
		})
	}
}

func TestGenerateN_InvalidCount(t *testing.T) {
	g := MustNew(WithLength(8))
	for _, n := range []int{0, -1, -100} {
		t.Run("", func(t *testing.T) {
			codes, err := g.GenerateN(n)
			if err == nil {
				t.Fatalf("GenerateN(%d) expected error, got nil", n)
			}
			if !errors.Is(err, ErrInvalidCount) {
				t.Errorf("error: want ErrInvalidCount, got %v", err)
			}
			if codes != nil {
				t.Errorf("GenerateN(%d) codes: want nil, got %v", n, codes)
			}
		})
	}
}

func TestGenerateN_AllCodesAdherToCharset(t *testing.T) {
	charset := CharsetAlphanumericUpper
	g := MustNew(WithLength(8), WithCharset(charset))
	codes, err := g.GenerateN(100)
	if err != nil {
		t.Fatalf("GenerateN() error: %v", err)
	}
	for i, code := range codes {
		for _, ch := range code {
			if !strings.ContainsRune(string(charset), ch) {
				t.Errorf("codes[%d]=%q: contains %q not in charset", i, code, ch)
			}
		}
	}
}

// ─────────────────────────────────────────────
// SetOptions()
// ─────────────────────────────────────────────

func TestSetOptions_UpdatesConfig(t *testing.T) {
	g := MustNew(WithLength(8), WithCharset(CharsetAlpha))

	err := g.SetOptions(WithLength(16), WithCharset(CharsetNumeric))
	if err != nil {
		t.Fatalf("SetOptions() error: %v", err)
	}

	opts := g.Options()
	if opts.Length != 16 {
		t.Errorf("length after SetOptions: want 16, got %d", opts.Length)
	}
	if opts.Charset != CharsetNumeric {
		t.Errorf("charset after SetOptions: want %q, got %q", CharsetNumeric, opts.Charset)
	}

	code, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() after SetOptions error: %v", err)
	}
	if len([]rune(code)) != 16 {
		t.Errorf("code length after SetOptions: want 16, got %d", len([]rune(code)))
	}
	for _, ch := range code {
		if !strings.ContainsRune(string(CharsetNumeric), ch) {
			t.Errorf("code %q contains %q not in numeric charset", code, ch)
		}
	}
}

func TestSetOptions_InvalidDoesNotMutate(t *testing.T) {
	original := Options{Length: 8, Charset: CharsetAlphanumeric}
	g := MustNew(WithLength(original.Length), WithCharset(original.Charset))

	err := g.SetOptions(WithLength(-1))
	if err == nil {
		t.Fatal("SetOptions() expected error for invalid length, got nil")
	}
	if !errors.Is(err, ErrInvalidLength) {
		t.Errorf("error: want ErrInvalidLength, got %v", err)
	}

	opts := g.Options()
	if opts.Length != original.Length {
		t.Errorf("length after failed SetOptions: want %d, got %d", original.Length, opts.Length)
	}
	if opts.Charset != original.Charset {
		t.Errorf("charset after failed SetOptions: want %q, got %q", original.Charset, opts.Charset)
	}
}

func TestSetOptions_EmptyCharsetDoesNotMutate(t *testing.T) {
	g := MustNew(WithLength(8))
	original := g.Options()

	err := g.SetOptions(WithCharset(""))
	if err == nil {
		t.Fatal("SetOptions() expected error for empty charset, got nil")
	}
	if !errors.Is(err, ErrEmptyCharset) {
		t.Errorf("error: want ErrEmptyCharset, got %v", err)
	}

	opts := g.Options()
	if opts.Charset != original.Charset {
		t.Errorf("charset mutated after failed SetOptions")
	}
}

// ─────────────────────────────────────────────
// Options() — snapshot immutability
// ─────────────────────────────────────────────

func TestOptions_SnapshotIsImmutable(t *testing.T) {
	g := MustNew(WithLength(8))

	snapshot := g.Options()
	snapshot.Length = 999
	snapshot.Charset = "XYZ"

	current := g.Options()
	if current.Length != 8 {
		t.Errorf("Generator length mutated via snapshot: want 8, got %d", current.Length)
	}
	if current.Charset != CharsetAlphanumeric {
		t.Errorf("Generator charset mutated via snapshot")
	}
}

// ─────────────────────────────────────────────
// Package-level Generate()
// ─────────────────────────────────────────────

func TestPackageGenerate_Defaults(t *testing.T) {
	code, err := Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if len([]rune(code)) != 8 {
		t.Errorf("length: want 8, got %d", len([]rune(code)))
	}
}

func TestPackageGenerate_WithOptions(t *testing.T) {
	code, err := Generate(WithLength(12), WithPrefix("ORD-"), WithCharset(CharsetAlphanumericUpper))
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if !strings.HasPrefix(code, "ORD-") {
		t.Errorf("code %q missing prefix 'ORD-'", code)
	}
	if len([]rune(code)) != 16 { // 4 + 12
		t.Errorf("total length: want 16, got %d", len([]rune(code)))
	}
}

func TestPackageGenerate_InvalidOptions(t *testing.T) {
	_, err := Generate(WithLength(0))
	if err == nil {
		t.Fatal("Generate() expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidLength) {
		t.Errorf("error: want ErrInvalidLength, got %v", err)
	}
}

// ─────────────────────────────────────────────
// Charset helpers
// ─────────────────────────────────────────────

func TestCharset_String(t *testing.T) {
	c := CharsetNumeric
	if c.String() != "0123456789" {
		t.Errorf("String(): want %q, got %q", "0123456789", c.String())
	}
}

func TestCharset_Len(t *testing.T) {
	tests := []struct {
		charset Charset
		want    int
	}{
		{CharsetNumeric, 10},
		{CharsetAlphaLower, 26},
		{CharsetAlphaUpper, 26},
		{CharsetAlpha, 52},
		{CharsetAlphanumeric, 62},
		{CharsetAlphanumericUpper, 36},
		{CharsetAlphanumericLower, 36},
	}
	for _, tt := range tests {
		if got := tt.charset.Len(); got != tt.want {
			t.Errorf("charset %q Len(): want %d, got %d", tt.charset, tt.want, got)
		}
	}
}

// ─────────────────────────────────────────────
// deduplicateRunes (internal)
// ─────────────────────────────────────────────

func TestDeduplicateRunes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"A", "A"},
		{"AABB", "AB"},
		{"ABCABC", "ABC"},
		{"123123", "123"},
		{"AABBCC123", "ABC123"},
	}
	for _, tt := range tests {
		got := deduplicateRunes(tt.input)
		if got != tt.want {
			t.Errorf("deduplicateRunes(%q): want %q, got %q", tt.input, tt.want, got)
		}
	}
}

// ─────────────────────────────────────────────
// validateOptions (internal)
// ─────────────────────────────────────────────

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr error
	}{
		{"valid", Options{Length: 8, Charset: CharsetAlphanumeric}, nil},
		{"zero length", Options{Length: 0, Charset: CharsetAlphanumeric}, ErrInvalidLength},
		{"negative length", Options{Length: -1, Charset: CharsetAlphanumeric}, ErrInvalidLength},
		{"empty charset", Options{Length: 8, Charset: ""}, ErrEmptyCharset},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptions(tt.opts)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validateOptions(): want %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// ─────────────────────────────────────────────
// Thread safety tests
// ─────────────────────────────────────────────

func TestConcurrent_Generate(t *testing.T) {
	g := MustNew(WithLength(12), WithCharset(CharsetAlphanumericUpper))

	const goroutines = 200
	const perGoroutine = 50

	results := make([]string, goroutines*perGoroutine)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				code, err := g.Generate()
				if err != nil {
					t.Errorf("goroutine %d Generate() error: %v", base, err)
					return
				}
				results[base*perGoroutine+j] = code
			}
		}(i)
	}

	wg.Wait()

	for i, code := range results {
		if len([]rune(code)) != 12 {
			t.Errorf("results[%d] length: want 12, got %d", i, len([]rune(code)))
		}
	}
}

func TestConcurrent_GenerateN(t *testing.T) {
	g := MustNew(WithLength(8))

	const goroutines = 50
	const n = 20

	var mu sync.Mutex
	var allCodes []string
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes, err := g.GenerateN(n)
			if err != nil {
				t.Errorf("concurrent GenerateN() error: %v", err)
				return
			}
			mu.Lock()
			allCodes = append(allCodes, codes...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	want := goroutines * n
	if len(allCodes) != want {
		t.Errorf("total codes: want %d, got %d", want, len(allCodes))
	}
}

func TestConcurrent_SetOptionsAndGenerate(t *testing.T) {
	g := MustNew(WithLength(8))

	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		lengths := []int{8, 12, 16, 10}
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				_ = g.SetOptions(WithLength(lengths[i%len(lengths)]))
				i++
			}
		}
	}()

	errCh := make(chan error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				if _, err := g.Generate(); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	genDone := make(chan struct{})
	go func() {
		close(genDone)
	}()

	<-genDone
	close(stop)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("Generate() during concurrent SetOptions: %v", err)
	}
}

func TestConcurrent_Options(t *testing.T) {
	g := MustNew(WithLength(8))
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = g.Options()
		}()
	}
	wg.Wait()
}

// ─────────────────────────────────────────────
// Benchmark
// ─────────────────────────────────────────────

func BenchmarkGenerate_Default(b *testing.B) {
	g := MustNew()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Generate()
	}
}

func BenchmarkGenerate_Length32(b *testing.B) {
	g := MustNew(WithLength(32))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Generate()
	}
}

func BenchmarkGenerate_Parallel(b *testing.B) {
	g := MustNew(WithLength(12))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = g.Generate()
		}
	})
}

func BenchmarkGenerateN_100(b *testing.B) {
	g := MustNew(WithLength(12))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.GenerateN(100)
	}
}
