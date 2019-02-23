package text

import (
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var utf8bom = []byte{0xef, 0xbb, 0xbf}

func UTF8BOM() []byte {
	return utf8bom
}

func DetectEncoding(r io.ReadSeeker) (Encoding, error) {
	r.Seek(0, io.SeekStart)

	lead := make([]byte, 3)
	n, err := r.Read(lead)
	r.Seek(0, io.SeekStart)

	if err == nil && n == 3 && reflect.DeepEqual(UTF8BOM(), lead) {
		return UTF8M, nil
	}
	return "", errors.New("cannot detect character encoding")
}

func SkipBOM(r io.Reader, enc Encoding) (io.Reader, error) {
	if enc == UTF8M {
		lead := make([]byte, 3)
		n, err := r.Read(lead)
		if err != nil || n != 3 || !reflect.DeepEqual(UTF8BOM(), lead) {
			return r, errors.New("byte order mark for UTF-8 does not exist")
		}
	}
	return r, nil
}

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

// Get a writer to transform character encoding from UTF-8 to another encoding.
func GetTransformWriter(w io.Writer, enc Encoding) io.Writer {
	switch enc {
	case SJIS:
		return transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
	default:
		return w
	}
}

// Encode a string from UTF-8 to another encoding.
func Encode(str string, enc Encoding) (string, error) {
	if enc == UTF8 || enc == UTF8M {
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
	if enc == UTF8 || enc == UTF8M {
		return str, nil
	}

	r := GetTransformDecoder(strings.NewReader(str), enc)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
