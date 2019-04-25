package text

import (
	"testing"
)

var widthTests = []struct {
	String               string
	EastAsianEncoding    bool
	CountDiacriticalSign bool
	CountFormatCode      bool
	Expect               int
}{
	{
		String:               "日本語\nabc",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               9,
	},
	{
		String:               "日本語\033[33mab\033[0mc",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               9,
	},
	{
		String:               "日本語abc",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               9,
	},
	{
		String:               "العَرَبِيَّة",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               7,
	},
	{
		String:               "العَرَبِيَّة",
		EastAsianEncoding:    true,
		CountDiacriticalSign: true,
		CountFormatCode:      false,
		Expect:               12,
	},
	{
		String:               "(´・ω・｀)",
		EastAsianEncoding:    true,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               12,
	},
	{
		String:               "(´・ω・｀)",
		EastAsianEncoding:    false,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               10,
	},
	{
		String:               "abc" + string(0x200b) + "def",
		EastAsianEncoding:    false,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Expect:               6,
	},
	{
		String:               "abc" + string(0x200b) + "def",
		EastAsianEncoding:    false,
		CountDiacriticalSign: false,
		CountFormatCode:      true,
		Expect:               7,
	},
}

func TestWidth(t *testing.T) {
	for _, v := range widthTests {
		result := Width(v.String, v.EastAsianEncoding, v.CountDiacriticalSign, v.CountFormatCode)
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
	{
		String:   "日本語🍺",
		Encoding: UTF16,
		Expect:   10,
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

	s = "\033[33m" + "العَرَبِيَّة" + "\033[0m"
	expect = true
	result = IsRightToLeftLetters(s)
	if result != expect {
		t.Errorf("right-to-left letters = %t, want %t for %q", result, expect, s)
	}

	s = "\033[33m1 " + "العَرَبِيَّة" + "\033[0m"
	expect = true
	result = IsRightToLeftLetters(s)
	if result != expect {
		t.Errorf("right-to-left letters = %t, want %t for %q", result, expect, s)
	}
}
