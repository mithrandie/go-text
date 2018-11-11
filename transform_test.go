package text

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	s := "日本語"

	enc := UTF8
	expect := []byte{0xe6, 0x97, 0xa5, 0xe6, 0x9c, 0xac, 0xe8, 0xaa, 0x9e}
	result, _ := Encode(s, enc)
	b := []byte(result)

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

	s = string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea})
	enc = SJIS
	result, _ = Decode(s, enc)
	b = []byte(result)
	if !reflect.DeepEqual(b, expect) {
		t.Errorf("result = %v, want %v for %v in %s", b, expect, s, enc)
	}
}
