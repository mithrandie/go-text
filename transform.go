package text

import (
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetTransformEncoder(r io.Reader, enc Encoding) io.Reader {
	switch enc {
	case SJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	default:
		return r
	}
}

func GetTransformDecoder(r io.Reader, enc Encoding) io.Reader {
	switch enc {
	case SJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	default:
		return r
	}
}

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
