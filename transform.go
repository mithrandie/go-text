package text

import (
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Get a reader to transform character encoding from UTF-8 to another encoding.
func GetTransformEncoder(r io.Reader, enc Encoding) io.Reader {
	switch enc {
	case SJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	default:
		return r
	}
}

// Get a reader to transform character encoding from any encoding to UTF-8.
func GetTransformDecoder(r io.Reader, enc Encoding) io.Reader {
	switch enc {
	case SJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	default:
		return r
	}
}

// Encode a string from UTF-8 to another encoding.
func Encode(str string, enc Encoding) (string, error) {
	if enc == UTF8 {
		return str, nil
	}

	r := GetTransformEncoder(strings.NewReader(str), enc)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Decode a string from any encoding to UTF-8.
func Decode(str string, enc Encoding) (string, error) {
	if enc == UTF8 {
		return str, nil
	}

	r := GetTransformDecoder(strings.NewReader(str), enc)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
