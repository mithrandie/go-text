package text

import (
	"reflect"
	"testing"
)

var parseEncodingTests = []struct {
	Input  string
	Expect Encoding
	Error  string
}{
	{
		Input:  "utf8",
		Expect: UTF8,
	},
	{
		Input:  "utf8m",
		Expect: UTF8M,
	},
	{
		Input:  "utf16",
		Expect: UTF16,
	},
	{
		Input:  "utf16be",
		Expect: UTF16BE,
	},
	{
		Input:  "utf16le",
		Expect: UTF16LE,
	},
	{
		Input:  "utf16bem",
		Expect: UTF16BEM,
	},
	{
		Input:  "utf16lem",
		Expect: UTF16LEM,
	},
	{
		Input:  "sjis",
		Expect: SJIS,
	},
	{
		Input:  "auto",
		Expect: AUTO,
	},
	{
		Input: "error",
		Error: "\"error\" cannot convert to Encoding",
	},
}

func TestParseEncoding(t *testing.T) {
	for _, v := range parseEncodingTests {
		result, err := ParseEncoding(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", result, v.Expect, v.Input)
		}
	}
}

var parseLineBreakTests = []struct {
	Input  string
	Expect LineBreak
	Error  string
}{
	{
		Input:  "crlf",
		Expect: CRLF,
	},
	{
		Input:  "cr",
		Expect: CR,
	},
	{
		Input:  "lf",
		Expect: LF,
	},
	{
		Input: "error",
		Error: "\"error\" cannot convert to LineBreak",
	},
}

func TestParseLineBreak(t *testing.T) {
	for _, v := range parseLineBreakTests {
		result, err := ParseLineBreak(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", result, v.Expect, v.Input)
		}
	}
}
