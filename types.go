package text

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Encoding string

const (
	UTF8  Encoding = "UTF8"
	UTF8M Encoding = "UTF8M"
	SJIS  Encoding = "SJIS"
)

var EncodingLiteral = map[Encoding]string{
	UTF8:  "UTF8",
	UTF8M: "UTF8M",
	SJIS:  "SJIS",
}

func (e Encoding) String() string {
	return EncodingLiteral[e]
}

type LineBreak string

const (
	CR   LineBreak = "\r"
	LF   LineBreak = "\n"
	CRLF LineBreak = "\r\n"
)

var LineBreakLiteral = map[LineBreak]string{
	CR:   "CR",
	LF:   "LF",
	CRLF: "CRLF",
}

func (lb LineBreak) Value() string {
	return reflect.ValueOf(lb).String()
}

func (lb LineBreak) String() string {
	return LineBreakLiteral[lb]
}

type FieldAlignment int

const (
	NotAligned FieldAlignment = iota
	Centering
	RightAligned
	LeftAligned
)

type RawText []byte

func ParseEncoding(s string) (Encoding, error) {
	var encoding Encoding
	switch strings.ToUpper(s) {
	case "UTF8":
		encoding = UTF8
	case "UTF8M":
		encoding = UTF8M
	case "SJIS":
		encoding = SJIS
	default:
		return encoding, errors.New(fmt.Sprintf("%q cannot convert to Encoding", s))
	}
	return encoding, nil
}

func ParseLineBreak(s string) (LineBreak, error) {
	var lb LineBreak
	switch strings.ToUpper(s) {
	case "CRLF":
		lb = CRLF
	case "CR":
		lb = CR
	case "LF":
		lb = LF
	default:
		return lb, errors.New(fmt.Sprintf("%q cannot convert to LineBreak", s))
	}
	return lb, nil
}
