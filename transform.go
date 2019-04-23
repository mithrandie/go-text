package text

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var utf8bom = [3]byte{0xef, 0xbb, 0xbf}

func UTF8BOM() [3]byte {
	return utf8bom
}

func UTF8BOMS() []byte {
	bom := UTF8BOM()
	return bom[0:3]
}

func DetectEncoding(r io.ReadSeeker) (Encoding, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	lead := make([]byte, 3)
	n, err := r.Read(lead)
	if _, e := r.Seek(0, io.SeekStart); e != nil {
		if err != nil {
			e = errors.New(strings.Join([]string{err.Error(), e.Error()}, "\n"))
		}
		return "", e
	}
	if err != nil || n != 3 {
		return "", errors.New("cannot detect character encoding")
	}

	bom := UTF8BOM()
	for i := range bom {
		if bom[i] != lead[i] {
			return "", errors.New("cannot detect character encoding")
		}
	}
	return UTF8M, nil
}

func SkipBOM(r io.Reader, enc Encoding) (io.Reader, error) {
	if enc == UTF8M {
		lead := make([]byte, 3)
		n, err := r.Read(lead)
		if err != nil || n != 3 {
			return r, errors.New("byte order mark for UTF-8 does not exist")
		}

		bom := UTF8BOM()
		for i := range bom {
			if bom[i] != lead[i] {
				return r, errors.New("byte order mark for UTF-8 does not exist")
			}
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
