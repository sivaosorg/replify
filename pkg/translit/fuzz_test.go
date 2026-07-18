package translit

import (
	"testing"
	"unicode/utf8"
)

// FuzzAppend feeds arbitrary byte sequences (valid or invalid UTF-8,
// arbitrary length) through Append/AppendMode and checks invariants
// that must hold no matter what: no panic, output stays within the
// promised byte budget (ASCII in, ASCII/UTF-8 out depending on mode),
// and repeated calls are deterministic.
func FuzzAppend(f *testing.F) {
	seeds := []string{
		"",
		"hello world",
		"kožušček",
		"北京市朝阳区",
		"日本語のテスト",
		"한국어 테스트",
		"Ελληνικά",
		"Москва, Россия",
		"\x00\x01\x02",
		"\xff\xfe\xfd",
		"valid\xc3\x28invalid",
		string(rune(0x10FFFF)),
		string(rune(0xeffff)),
		string(rune(0xf0000)),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		out1 := Append(nil, s)
		out2 := Append(nil, s)
		if string(out1) != string(out2) {
			t.Fatalf("Append is not deterministic for %q", s)
		}

		for _, b := range out1 {
			if b >= utf8.RuneSelf {
				t.Fatalf("Append(%q) produced non-ASCII byte 0x%x in ModeSkip output", s, b)
			}
		}

		keep := AppendMode(nil, s, ModeKeep)
		if utf8.ValidString(s) && !utf8.Valid(keep) {
			t.Fatalf("AppendMode(%q, ModeKeep) produced invalid UTF-8 from valid input: %q", s, keep)
		}

		hint := SizeHint(s)
		if hint != len(out1) {
			t.Fatalf("SizeHint(%q) = %d, but Append produced %d bytes", s, hint, len(out1))
		}

		hintKeep := SizeHintMode(s, ModeKeep)
		if hintKeep != len(keep) {
			t.Fatalf("SizeHintMode(%q, ModeKeep) = %d, but AppendMode produced %d bytes", s, hintKeep, len(keep))
		}

		// AppendBytes must agree with Append for the same input.
		if got := AppendBytes(nil, []byte(s)); string(got) != string(out1) {
			t.Fatalf("AppendBytes(%q) = %q, want %q (disagrees with Append)", s, got, out1)
		}
	})
}

// FuzzValidUTF8 checks ValidUTF8 agrees with the standard library and
// never panics.
func FuzzValidUTF8(f *testing.F) {
	f.Add("hello")
	f.Add("\xff\xfe")
	f.Add("kožušček")
	f.Fuzz(func(t *testing.T, s string) {
		if ValidUTF8(s) != utf8.ValidString(s) {
			t.Fatalf("ValidUTF8(%q) disagrees with utf8.ValidString", s)
		}
	})
}
