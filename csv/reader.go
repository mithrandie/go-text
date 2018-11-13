package csv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/mithrandie/go-text"
)

type Reader struct {
	Delimiter   rune
	WithoutNull bool
	Encoding    text.Encoding

	reader *bufio.Reader
	line   int
	column int

	recordBuf     bytes.Buffer
	fieldStartPos []int
	fieldQuoted   []bool

	FieldsPerRecord int

	DetectedLineBreak text.LineBreak
	EnclosedAll       bool
}

func NewReader(r io.Reader, enc text.Encoding) *Reader {
	return &Reader{
		Delimiter:       ',',
		WithoutNull:     false,
		Encoding:        enc,
		reader:          bufio.NewReader(text.GetTransformDecoder(r, enc)),
		line:            1,
		column:          0,
		FieldsPerRecord: 0,
		EnclosedAll:     true,
	}
}

func (r *Reader) newError(s string) error {
	return errors.New(fmt.Sprintf("line %d, column %d: %s", r.line, r.column, s))
}

func (r *Reader) ReadHeader() ([]string, error) {
	record, err := r.parseRecord(true)
	if err != nil {
		return nil, err
	}

	header := make([]string, len(record))
	for i, v := range record {
		header[i] = string(v)
	}
	return header, nil
}

func (r *Reader) Read() ([]text.RawText, error) {
	return r.parseRecord(r.WithoutNull)
}

func (r *Reader) ReadAll() ([][]text.RawText, error) {
	records := make([][]text.RawText, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *Reader) parseRecord(withoutNull bool) ([]text.RawText, error) {
	r.recordBuf.Reset()
	r.fieldStartPos = r.fieldStartPos[:0]
	r.fieldQuoted = r.fieldQuoted[:0]

	fieldIndex := 0
	fieldPosition := 0
	for {
		if 0 < r.FieldsPerRecord && r.FieldsPerRecord <= fieldIndex {
			return nil, r.newError("wrong number of fields in line")
		}

		fieldPosition = r.recordBuf.Len()
		quoted, eol, err := r.parseField()

		if err != nil {
			if err == io.EOF {
				if fieldIndex < 1 && r.recordBuf.Len() < 1 {
					return nil, io.EOF
				}
			} else {
				return nil, err
			}
		}

		if eol && fieldIndex < 1 && r.recordBuf.Len() < 1 {
			continue
		}

		r.fieldStartPos = append(r.fieldStartPos, fieldPosition)
		r.fieldQuoted = append(r.fieldQuoted, quoted)
		fieldIndex++

		if eol {
			break
		}
	}

	if r.FieldsPerRecord < 1 {
		r.FieldsPerRecord = fieldIndex
	} else if fieldIndex < r.FieldsPerRecord {
		r.line--
		return nil, r.newError("wrong number of fields in line")
	}

	record := make([]text.RawText, 0, r.FieldsPerRecord)
	recordStr := make([]byte, r.recordBuf.Len())
	copy(recordStr, r.recordBuf.Bytes())
	for i, pos := range r.fieldStartPos {
		var endPos int
		if i == len(r.fieldStartPos)-1 {
			endPos = r.recordBuf.Len()
		} else {
			endPos = r.fieldStartPos[i+1]
		}

		if !withoutNull && pos == endPos && !r.fieldQuoted[i] {
			record = append(record, nil)
		} else {
			record = append(record, recordStr[pos:endPos])
		}
	}

	return record, nil
}

func (r *Reader) parseField() (bool, bool, error) {
	var eof error
	eol := false
	startPos := r.recordBuf.Len()

	quoted := false
	escaped := false

	var lineBreak text.LineBreak

Read:
	for {
		lineBreak = ""

		ch, _, err := r.reader.ReadRune()
		r.column++

		if err != nil {
			if err == io.EOF {
				if !escaped && quoted {
					return quoted, eol, r.newError("extraneous \" in field")
				}
				eol = true
			}
			return quoted, eol, err
		}

		switch ch {
		case '\r':
			nxtCh, _, _ := r.reader.ReadRune()
			if nxtCh == '\n' {
				lineBreak = text.CRLF
			} else {
				r.reader.UnreadRune()
				lineBreak = text.CR
			}
			ch = '\n'
		case '\n':
			lineBreak = text.LF
		}
		if ch == '\n' {
			r.line++
			r.column = 0
		}

		if quoted {
			if escaped {
				switch ch {
				case '"':
					escaped = false
					r.recordBuf.WriteRune(ch)
					continue
				case r.Delimiter:
					break Read
				case '\n':
					if r.DetectedLineBreak == "" {
						r.DetectedLineBreak = lineBreak
					}
					eol = true
					break Read
				default:
					r.column--
					return quoted, eol, r.newError("unexpected \" in field")
				}
			}

			switch ch {
			case '"':
				escaped = true
			case '\n':
				r.recordBuf.WriteString(lineBreak.Value())
			default:
				r.recordBuf.WriteRune(ch)
			}
			continue
		}

		switch ch {
		case '\n':
			if r.DetectedLineBreak == "" {
				r.DetectedLineBreak = lineBreak
			}
			eol = true
			break Read
		case r.Delimiter:
			break Read
		case '"':
			if startPos == r.recordBuf.Len() {
				quoted = true
			} else {
				r.recordBuf.WriteRune(ch)
			}
		default:
			if r.EnclosedAll && unicode.IsLetter(ch) {
				r.EnclosedAll = false
			}
			r.recordBuf.WriteRune(ch)
		}
	}

	return quoted, eol, eof
}
