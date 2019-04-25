package text

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func utf8bom(s string) []byte {
	ret, _ := Encode([]byte(s), UTF8M)
	return []byte(ret)
}

func utf16be(s string) []byte {
	ret, _ := Encode([]byte(s), UTF16BE)
	return []byte(ret)
}

func utf16le(s string) []byte {
	ret, _ := Encode([]byte(s), UTF16LE)
	return []byte(ret)
}

func utf16bebom(s string) []byte {
	ret, _ := Encode([]byte(s), UTF16BEM)
	return []byte(ret)
}

func utf16lebom(s string) []byte {
	ret, _ := Encode([]byte(s), UTF16LEM)
	return []byte(ret)
}

func sjis(s string) []byte {
	ret, _ := Encode([]byte(s), SJIS)
	return []byte(ret)
}

var detectEncodingTests = []struct {
	Input  []byte
	Result Encoding
	Error  string
}{
	{
		Input:  []byte{},
		Result: UTF8,
	},
	{
		Input:  []byte("ab"),
		Result: UTF8,
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
		if result != v.Result {
			t.Errorf("result = %s, want %s for %q", result, v.Result, v.Input)
		}
	}
}

var detectInSpecifiedEncodingTests = []struct {
	Input    []byte
	Encoding Encoding
	Result   Encoding
	Error    string
}{
	{
		Input:    []byte{},
		Encoding: AUTO,
		Result:   UTF8,
	},
	{
		Input:    []byte("ab"),
		Encoding: AUTO,
		Result:   UTF8,
	},
	{
		Input:    utf8bom("abc"),
		Encoding: AUTO,
		Result:   UTF8M,
	},
	{
		Input:    []byte("abc"),
		Encoding: AUTO,
		Result:   UTF8,
	},
	{
		Input:    utf16bebom("abc"),
		Encoding: AUTO,
		Result:   UTF16BEM,
	},
	{
		Input:    utf16lebom("abc"),
		Encoding: AUTO,
		Result:   UTF16LEM,
	},
	{
		Input:    utf16lebom("abc"),
		Encoding: UTF16LEM,
		Result:   UTF16LEM,
	},
	{
		Input:    []byte("abc"),
		Encoding: UTF8,
		Result:   UTF8,
	},
	{
		Input:    utf16be("abc"),
		Encoding: UTF16,
		Result:   UTF16BE,
	},
}

func TestDetectInSpecifiedEncoding(t *testing.T) {
	for _, v := range detectInSpecifiedEncodingTests {
		result, err := DetectInSpecifiedEncoding(bytes.NewReader(v.Input), v.Encoding)
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
		if result != v.Result {
			t.Errorf("result = %s, want %s for %q", result, v.Result, v.Input)
		}
	}
}

var inferEncodingTests = []struct {
	Input  []byte
	Expect Encoding
	Error  string
}{
	{
		Input:  []byte{},
		Expect: UTF8,
	},
	{
		Input:  []byte("abc"),
		Expect: UTF8,
	},
	{
		Input:  []byte("日本語"),
		Expect: UTF8,
	},
	{
		Input:  []byte("🍺"),
		Expect: UTF8,
	},
	{
		Input:  utf16be("abc"),
		Expect: UTF16BE,
	},
	{
		Input:  utf16le("abc"),
		Expect: UTF16LE,
	},
	{
		Input:  utf16le("abcdefg"),
		Expect: UTF16LE,
	},
	{
		Input:  utf16be("🙈🙉🙊"),
		Expect: UTF16BE,
	},
	{
		Input:  utf16le("🙈🙉🙊"),
		Expect: UTF16LE,
	},
	{
		Input:  sjis("日本語ｱｲｳ"),
		Expect: SJIS,
	},
	{
		Input: []byte{0xd8, 0x00, 0xd8, 0x00},
		Error: "cannot detect character encoding",
	},
}

func TestInferEncoding(t *testing.T) {
	for _, v := range inferEncodingTests {
		result, err := InferEncoding(v.Input, true)
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
		if result != v.Expect {
			t.Errorf("result = %s, want %s for %q", result, v.Expect, v.Input)
		}
	}
}

var getTransformWriterTests = []struct {
	Encoding Encoding
	Expect   []byte
	Error    string
}{
	{
		Encoding: UTF8,
		Expect:   []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e},
	},
	{
		Encoding: UTF8M,
		Expect:   []byte{0xef, 0xbb, 0xbf, 0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e},
	},
	{
		Encoding: UTF16,
		Expect:   []byte{0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16BE,
		Expect:   []byte{0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16LE,
		Expect:   []byte{0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: UTF16BEM,
		Expect:   []byte{0xfe, 0xff, 0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16LEM,
		Expect:   []byte{0xff, 0xfe, 0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: SJIS,
		Expect:   []byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea},
	},
	{
		Encoding: AUTO,
		Error:    "invalid character encoding",
	},
}

func TestGetTransformWriter(t *testing.T) {
	src := []byte("日本語")

	buf := new(bytes.Buffer)

	for _, v := range getTransformWriterTests {
		w, err := GetTransformWriter(buf, v.Encoding)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Encoding)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Encoding)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Encoding)
			continue
		}

		_, _ = w.Write(src)
		result, err := ioutil.ReadAll(buf)
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %s, want %s for %q", result, v.Expect, v.Encoding)
		}
	}
}

var encodeTests = []struct {
	Encoding Encoding
	Expect   []byte
	Error    string
}{
	{
		Encoding: UTF8,
		Expect:   []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e},
	},
	{
		Encoding: UTF8M,
		Expect:   []byte{0xef, 0xbb, 0xbf, 0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e},
	},
	{
		Encoding: UTF16,
		Expect:   []byte{0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16BE,
		Expect:   []byte{0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16LE,
		Expect:   []byte{0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: UTF16BEM,
		Expect:   []byte{0xfe, 0xff, 0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16LEM,
		Expect:   []byte{0xff, 0xfe, 0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: SJIS,
		Expect:   []byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea},
	},
	{
		Encoding: AUTO,
		Error:    "invalid character encoding",
	},
}

func TestEncode(t *testing.T) {
	src := []byte("日本語")

	for _, v := range encodeTests {
		result, err := Encode(src, v.Encoding)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Encoding)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Encoding)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Encoding)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %s, want %s for %q", result, v.Expect, v.Encoding)
		}
	}
}

var decodeTests = []struct {
	Encoding Encoding
	Source   []byte
	Error    string
}{
	{
		Encoding: UTF8,
		Source:   []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e},
	},
	{
		Encoding: UTF8M,
		Source:   []byte{0xef, 0xbb, 0xbf, 0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e},
	},
	{
		Encoding: UTF16,
		Source:   []byte{0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16,
		Source:   []byte{0xfe, 0xff, 0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16,
		Source:   []byte{0xff, 0xfe, 0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: UTF16BE,
		Source:   []byte{0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16LE,
		Source:   []byte{0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: UTF16BEM,
		Source:   []byte{0xfe, 0xff, 0x65, 0xe5, 0x67, 0x2c, 0x8a, 0x9e},
	},
	{
		Encoding: UTF16LEM,
		Source:   []byte{0xff, 0xfe, 0xe5, 0x65, 0x2c, 0x67, 0x9e, 0x8a},
	},
	{
		Encoding: SJIS,
		Source:   []byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea},
	},
	{
		Encoding: AUTO,
		Error:    "invalid character encoding",
	},
}

func TestDecode(t *testing.T) {
	expect := []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e}

	for _, v := range decodeTests {
		result, err := Decode(v.Source, v.Encoding)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Encoding)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Encoding)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Encoding)
			continue
		}
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result = %s, want %s for %q", result, expect, v.Encoding)
		}
	}
}
