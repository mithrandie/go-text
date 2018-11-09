package text

import "testing"

var widthTests = []struct {
	String               string
	EastAsianEncoding    bool
	CountDiacriticalSign bool
	Expect               int
}{
	{
		String:               "日本語\nabc",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		Expect:               9,
	},
	{
		String:               "日本語\033[33mab\033[0mc",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		Expect:               9,
	},
	{
		String:               "日本語abc",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		Expect:               9,
	},
	{
		String:               "العَرَبِيَّة",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		Expect:               7,
	},
	{
		String:               "العَرَبِيَّة",
		EastAsianEncoding:    true,
		CountDiacriticalSign: true,
		Expect:               12,
	},
	{
		String:               "(´・ω・｀)",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		Expect:               12,
	},
	{
		String:               "(´・ω・｀)",
		EastAsianEncoding:    false,
		CountDiacriticalSign: false,
		Expect:               10,
	},
}

func TestWidth(t *testing.T) {
	for _, v := range widthTests {
		result := Width(v.String, v.EastAsianEncoding, v.CountDiacriticalSign)
		if result != v.Expect {
			t.Errorf("width = %d, want %d for %q, %t, %t", result, v.Expect, v.String, v.EastAsianEncoding, v.CountDiacriticalSign)
		}
	}
}

var runeByteSizeTests = []struct {
	Rune     rune
	Encoding Encoding
	Expect   int
}{
	{
		Rune:     '日',
		Encoding: UTF8,
		Expect:   3,
	},
	{
		Rune:     'a',
		Encoding: UTF8,
		Expect:   1,
	},
	{
		Rune:     '\n',
		Encoding: UTF8,
		Expect:   1,
	},
	{
		Rune:     '日',
		Encoding: SJIS,
		Expect:   2,
	},
	{
		Rune:     'a',
		Encoding: SJIS,
		Expect:   1,
	},
	{
		Rune:     '\n',
		Encoding: SJIS,
		Expect:   1,
	},
}

func TestRuneByteSize(t *testing.T) {
	for _, v := range runeByteSizeTests {
		result := RuneByteSize(v.Rune, v.Encoding)
		if result != v.Expect {
			t.Errorf("byte size = %d, want %d for %q, %s", result, v.Expect, v.Rune, v.Encoding)
		}
	}
}

var byteSizeTests = []struct {
	String   string
	Encoding Encoding
	Expect   int
}{
	{
		String:   "日本語abc",
		Encoding: UTF8,
		Expect:   12,
	},
	{
		String:   "日本語abc",
		Encoding: SJIS,
		Expect:   9,
	},
}

func TestByteSize(t *testing.T) {
	for _, v := range byteSizeTests {
		result := ByteSize(v.String, v.Encoding)
		if result != v.Expect {
			t.Errorf("byte size = %d, want %d for %q, %s", result, v.Expect, v.String, v.Encoding)
		}
	}
}

func TestIsRightToLeftLetters(t *testing.T) {
	var expect bool

	s := ""
	expect = false
	result := IsRightToLeftLetters(s)
	if result != expect {
		t.Errorf("right-to-left letters = %t, want %t for %q", result, expect, s)
	}

	s = "العَرَبِيَّة"
	expect = true
	result = IsRightToLeftLetters(s)
	if result != expect {
		t.Errorf("right-to-left letters = %t, want %t for %q", result, expect, s)
	}
}
