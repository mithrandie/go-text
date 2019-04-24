package text

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const UTF8BOM = "\ufeff"

var ErrUnknownEncoding = errors.New("cannot detect character encoding")

func DetectEncodingByBOM(r io.ReadSeeker) (enc Encoding, err error) {
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return
	}

	lead := make([]byte, 3)
	n, _ := r.Read(lead)
	if n == 3 && lead[0] == UTF8BOM[0] && lead[1] == UTF8BOM[1] && lead[2] == UTF8BOM[2] {
		enc = UTF8M
	} else {
		err = ErrUnknownEncoding
	}

	if _, e := r.Seek(0, io.SeekStart); e != nil {
		err = e
	}
	return
}

func DetectEncoding(r io.ReadSeeker) (enc Encoding, err error) {
	defer func() {
		if _, e := r.Seek(0, io.SeekStart); e != nil {
			err = e
		}
	}()

	if enc, err = DetectEncodingByBOM(r); err == nil || err != ErrUnknownEncoding {
		return
	}

	lead := make([]byte, 1024)
	n, _ := r.Read(lead)
	lead = lead[:n]

	if utf8.Valid(lead) {
		return UTF8, nil
	}

	str := string(lead)
	if decoded, e := Decode(str, SJIS); e == nil {
		if encoded, e := Encode(decoded, SJIS); e == nil {
			if str == encoded {
				return SJIS, nil
			}
		}
	}

	err = ErrUnknownEncoding
	return
}

// Get a reader to transform character encoding from UTF-8 to another encoding.
func GetTransformEncoder(r io.Reader, enc Encoding) io.Reader {
	switch enc {
	case SJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	case UTF8M:
		return transform.NewReader(r, NewUTF8MEncoder())
	default:
		return transform.NewReader(r, unicode.UTF8.NewEncoder())
	}
}

// Get a reader to transform character encoding from any encoding to UTF-8.
func GetTransformDecoder(r io.Reader, enc Encoding) io.Reader {
	switch enc {
	case SJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	case UTF8M:
		return transform.NewReader(r, unicode.BOMOverride(unicode.UTF8.NewDecoder()))
	default:
		return transform.NewReader(r, unicode.UTF8.NewDecoder())
	}
}

// Get a writer to transform character encoding from UTF-8 to another encoding.
func GetTransformWriter(w io.Writer, enc Encoding) io.Writer {
	switch enc {
	case SJIS:
		return transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
	case UTF8M:
		return transform.NewWriter(w, NewUTF8MEncoder())
	default:
		return transform.NewWriter(w, unicode.UTF8.NewEncoder())
	}
}

// Encode a string from UTF-8 to another encoding.
func Encode(str string, enc Encoding) (string, error) {
	r := GetTransformEncoder(strings.NewReader(str), enc)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Decode a string from any encoding to UTF-8.
func Decode(str string, enc Encoding) (string, error) {
	r := GetTransformDecoder(strings.NewReader(str), enc)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type bomPolicy uint8

const (
	writeBOM  bomPolicy = 0x01
	ignoreBOM bomPolicy = 0
)

type UTF8MEncoder struct {
	initialBOMPolicy bomPolicy
	currentBOMPolicy bomPolicy
}

func NewUTF8MEncoder() *encoding.Encoder {
	return &encoding.Encoder{Transformer: &UTF8MEncoder{
		initialBOMPolicy: writeBOM,
		currentBOMPolicy: writeBOM,
	}}
}

func (u *UTF8MEncoder) Reset() {
	u.currentBOMPolicy = u.initialBOMPolicy
}

func (u *UTF8MEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if u.currentBOMPolicy&writeBOM != 0 {
		if len(dst) < 3 {
			return 0, 0, transform.ErrShortDst
		}
		bom := []byte(UTF8BOM)
		dst[0], dst[1], dst[2] = bom[0], bom[1], bom[2]
		u.currentBOMPolicy = ignoreBOM
		nDst = 3
	}

	for i := range src {
		if nDst+1 > len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst] = src[i]
		nDst++
		nSrc++
	}

	return nDst, nSrc, err
}
