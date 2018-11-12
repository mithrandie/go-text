package csv

import (
	"bytes"
	"strings"

	"github.com/mithrandie/go-text"
)

type Encoder struct {
	Delimiter     rune
	LineBreak     text.LineBreak
	WithoutHeader bool
	Encoding      text.Encoding

	header    []Field
	recordSet [][]Field
	fieldLen  int
	lineBreak string
	buf       bytes.Buffer
}

func NewEncoder(recordCounts int) *Encoder {
	return &Encoder{
		Delimiter:     ',',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Encoding:      text.UTF8,
		fieldLen:      0,
		recordSet:     make([][]Field, 0, recordCounts),
	}
}

func (e *Encoder) SetHeader(header []Field) {
	e.header = header
	if e.fieldLen < len(header) {
		e.fieldLen = len(header)
	}
}

func (e *Encoder) AppendRecord(record []Field) {
	e.recordSet = append(e.recordSet, record)
	if e.fieldLen < len(record) {
		e.fieldLen = len(record)
	}
}

func (e *Encoder) Encode() (string, error) {
	if e.fieldLen < 1 {
		return "", nil
	}

	e.lineBreak = e.LineBreak.Value()

	lines := make([]string, 0, len(e.recordSet)+1)

	if !e.WithoutHeader {
		lines = append(lines, e.formatRecord(e.header))
	}

	for _, record := range e.recordSet {
		lines = append(lines, e.formatRecord(record))
	}

	return text.Encode(strings.Join(lines, e.lineBreak), e.Encoding)
}

func (e *Encoder) formatRecord(record []Field) string {
	e.buf.Reset()

	for i := 0; i < e.fieldLen; i++ {
		if 0 < i {
			e.buf.WriteRune(e.Delimiter)
		}

		if i < len(record) {
			if record[i].Quote || e.includeDelimiter(record[i].Contents) {
				e.buf.WriteRune('"')
				runes := []rune(record[i].Contents)
				pos := 0

				for {
					if len(runes) <= pos {
						break
					}

					r := runes[pos]
					switch r {
					case '"':
						e.buf.WriteRune(r)
						e.buf.WriteRune(r)
					default:
						e.buf.WriteRune(r)
					}

					pos++
				}
				e.buf.WriteRune('"')
			} else {
				e.buf.WriteString(record[i].Contents)
			}
		}
	}
	return e.buf.String()
}

func (e *Encoder) includeDelimiter(s string) bool {
	for _, r := range s {
		if r == e.Delimiter {
			return true
		}
	}
	return false
}
