package fixedlen

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/mithrandie/go-text"
)

const PadChar = ' '

type Encoder struct {
	DelimiterPositions DelimiterPositions
	LineBreak          string
	WithoutHeader      bool
	Encoding           text.Encoding

	header    []*Field
	recordSet [][]*Field
	fieldLen  int
	buf       bytes.Buffer
}

func NewEncoder(recordCounts int) *Encoder {
	return &Encoder{
		DelimiterPositions: nil,
		LineBreak:          text.LF.Value(),
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		fieldLen:           0,
		recordSet:          make([][]*Field, 0, recordCounts),
	}
}

func (e *Encoder) SetHeader(header []*Field) {
	e.header = header
	if e.fieldLen < len(header) {
		e.fieldLen = len(header)
	}
}

func (e *Encoder) AppendRecord(record []*Field) {
	e.recordSet = append(e.recordSet, record)
	if e.fieldLen < len(record) {
		e.fieldLen = len(record)
	}
}

func (e *Encoder) Encode() (string, error) {
	prevPos := 0
	for _, endPos := range e.DelimiterPositions {
		if endPos < 0 || endPos <= prevPos {
			return "", errors.New(fmt.Sprintf("invalid delimiter position: %s", e.DelimiterPositions))
		}
		prevPos = endPos
	}

	insertSpace := false
	if e.DelimiterPositions == nil {
		lenPerField := e.measureLengthPerField()
		e.DelimiterPositions = make(DelimiterPositions, 0, len(lenPerField))
		pos := 0
		for _, l := range lenPerField {
			pos = pos + l
			e.DelimiterPositions = append(e.DelimiterPositions, pos)
		}

		insertSpace = true
	}

	if len(e.DelimiterPositions) < 1 {
		return "", nil
	}

	var err error
	e.buf.Reset()

	if !e.WithoutHeader {
		if err = e.formatRecord(e.header, insertSpace); err != nil {
			return e.buf.String(), err
		}
	}

	for _, record := range e.recordSet {
		if 0 < e.buf.Len() {
			e.buf.WriteString(e.LineBreak)
		}

		if err = e.formatRecord(record, insertSpace); err != nil {
			return e.buf.String(), err
		}
	}

	return text.Encode(e.buf.String(), e.Encoding)
}

func (e *Encoder) formatRecord(record []*Field, insertSpace bool) error {
	start := 0
	for i, end := range e.DelimiterPositions {
		if insertSpace && 0 < i {
			e.buf.WriteRune(PadChar)
		}

		size := end - start
		if i < len(record) {
			if err := e.addField(record[i], size); err != nil {
				return err
			}
		} else {
			e.buf.Write(bytes.Repeat([]byte(string(PadChar)), size))
		}
		start = end
	}
	return nil
}

func (e *Encoder) addField(field *Field, fieldSize int) error {
	size := text.ByteSize(field.Contents, e.Encoding)
	if fieldSize < size {
		return errors.New(fmt.Sprintf("value is too long: %q for %d byte(s) length field", field.Contents, fieldSize))
	}

	padLen := fieldSize - size

	switch field.Alignment {
	case text.Centering:
		halfPadLen := padLen / 2
		e.buf.Write(bytes.Repeat([]byte(string(PadChar)), halfPadLen))
		e.buf.WriteString(field.Contents)
		e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen-halfPadLen))
	case text.RightAligned:
		e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen))
		e.buf.WriteString(field.Contents)
	default:
		e.buf.WriteString(field.Contents)
		e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen))
	}

	return nil
}

func (e *Encoder) measureLengthPerField() []int {
	length := make([]int, e.fieldLen)

	if !e.WithoutHeader {
		for i, v := range e.header {
			length[i] = text.ByteSize(v.Contents, e.Encoding)
		}
	}

	for _, record := range e.recordSet {
		for i, v := range record {
			l := text.ByteSize(v.Contents, e.Encoding)
			if length[i] < l {
				length[i] = l
			}
		}
	}
	return length
}
