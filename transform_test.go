package text

import (
	"bytes"
	"reflect"
	"testing"
)

var detectEncodingTests = []struct {
	Input  []byte
	Result Encoding
	Error  string
}{
	{
		Input: []byte{},
		Error: "cannot detect character encoding",
	},
	{
		Input: []byte("ab"),
		Error: "cannot detect character encoding",
	},
	{
		Input: []byte("abc"),
		Error: "cannot detect character encoding",
	},
	{
		Input:  append([]byte(UTF8BOM), []byte("abc")...),
		Result: UTF8M,
	},
}

func TestDetectEncoding(t *testing.T) {
	for _, v := range detectEncodingTests {
		result, err := DetectEncoding(bytes.NewReader(v.Input))
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("result = %#v, want %#v for %q", result, v.Result, v.Input)
		}
	}
}

func TestEncode(t *testing.T) {
	s := "日本語"

	enc := UTF8
	expect := []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e}
	result, _ := Encode(s, enc)
	b := []byte(result)

	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %q in %s", b, expect, s, enc)
	}

	enc = UTF8M
	expect = []byte{0xef, 0xbb, 0xbf, 0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e}
	result, _ = Encode(s, enc)
	b = []byte(result)

	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %q in %s", b, expect, s, enc)
	}

	enc = SJIS
	expect = []byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}
	result, _ = Encode(s, enc)
	b = []byte(result)

	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %q in %s", b, expect, s, enc)
	}
}

func TestDecode(t *testing.T) {
	expect := []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e}

	s := string([]byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e})
	enc := UTF8
	result, _ := Decode(s, enc)
	b := []byte(result)
	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %v in %s", b, expect, s, enc)
	}

	s = string([]byte{0xef, 0xbb, 0xbf, 0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e})
	enc = UTF8M
	result, _ = Decode(s, enc)
	b = []byte(result)
	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %v in %s", b, expect, s, enc)
	}

	s = string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea})
	enc = SJIS
	result, _ = Decode(s, enc)
	b = []byte(result)
	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %v in %s", b, expect, s, enc)
	}
}
