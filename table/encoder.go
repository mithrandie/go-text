package table

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/mithrandie/go-text"
)

const (
	VLine             = '|'
	HLine             = '-'
	CrossLine         = '+'
	PadChar           = ' '
	AlignSign         = ':'
	MarkdownLineBreak = "<br />"
	EscapeChar        = '\\'
)

type Encoder struct {
	Format               Format
	LineBreak            text.LineBreak
	EastAsianEncoding    bool
	CountDiacriticalSign bool
	CountFormatCode      bool
	Encoding             text.Encoding

	// GFM or Org Table
	WithoutHeader bool

	// GFM Table only
	alignments []text.FieldAlignment

	header    []Field
	recordSet [][]Field
	fieldLen  int
	lineBreak string
	lineWidth int
	writer    *bufio.Writer
	buf       bytes.Buffer
}

func NewEncoder(format Format, recordCounts int) *Encoder {
	return &Encoder{
		Format:               format,
		LineBreak:            text.LF,
		EastAsianEncoding:    false,
		CountDiacriticalSign: false,
		CountFormatCode:      false,
		Encoding:             text.UTF8,
		WithoutHeader:        false,
		fieldLen:             0,
		recordSet:            make([][]Field, 0, recordCounts),
	}
}

func (e *Encoder) SetHeader(header []Field) {
	e.header = e.prepareRecord(header)
	if e.fieldLen < len(header) {
		e.fieldLen = len(header)
	}
}

func (e *Encoder) SetFieldAlignments(alignments []text.FieldAlignment) {
	e.alignments = alignments
}

func (e *Encoder) AppendRecord(record []Field) {
	e.recordSet = append(e.recordSet, e.prepareRecord(record))
	if e.fieldLen < len(record) {
		e.fieldLen = len(record)
	}
}

func (e *Encoder) prepareRecord(record []Field) []Field {
	for i := range record {
		e.prepareField(&record[i])
	}
	return record
}

func (e *Encoder) prepareField(field *Field) {
	lines := strings.Split(e.escape(field.Contents), "\n")

	width := 0
	for _, v := range lines {
		l := text.Width(v, e.EastAsianEncoding, e.CountDiacriticalSign, e.CountFormatCode)
		if width < l {
			width = l
		}
	}

	field.Lines = lines
	field.Width = width
}

func (e *Encoder) Encode() (string, error) {
	if e.fieldLen < 1 {
		return "", nil
	}

	e.lineBreak = e.LineBreak.Value()

	var err error
	buf := new(bytes.Buffer)

	writer, err := text.GetTransformWriter(buf, e.Encoding)
	if err != nil {
		return "", err
	}
	e.writer = bufio.NewWriter(writer)

	fieldWidths := make([]int, e.fieldLen)

	for _, record := range e.recordSet {
		for i, f := range record {
			fw := f.Width
			if fieldWidths[i] < fw {
				fieldWidths[i] = fw
			}
		}
	}

	if e.Format == PlainTable || !e.WithoutHeader {
		for i, f := range e.header {
			fw := f.Width
			if fieldWidths[i] < fw {
				fieldWidths[i] = fw
			}
			if e.Format == GFMTable {
				if fieldWidths[i] < 3 {
					fieldWidths[i] = 3
				}
			}
			if ((fieldWidths[i] - f.Width) % 2) == 1 {
				fieldWidths[i] = fieldWidths[i] + 1
			}
		}
	}

	e.calculateLineWidth(fieldWidths)

	appended := false
	if e.Format == PlainTable || !e.WithoutHeader {
		if e.Format == PlainTable {
			if err = e.formatTextHR(fieldWidths); err != nil {
				return "", err
			}
			if _, err = e.writer.WriteString(e.lineBreak); err != nil {
				return "", err
			}
		}
		if err = e.formatRecord(e.header, fieldWidths); err != nil {
			return "", err
		}
		if _, err = e.writer.WriteString(e.lineBreak); err != nil {
			return "", err
		}

		switch e.Format {
		case GFMTable:
			err = e.formatGfmHR(fieldWidths)
		case OrgTable:
			err = e.formatOrgHR(fieldWidths)
		default:
			err = e.formatTextHR(fieldWidths)
		}
		if err != nil {
			return "", err
		}
		appended = true
	}

	if 0 < len(e.recordSet) {
		for _, record := range e.recordSet {
			if appended {
				if _, err = e.writer.WriteString(e.lineBreak); err != nil {
					return "", err
				}
			} else {
				appended = true
			}

			if err = e.formatRecord(record, fieldWidths); err != nil {
				return "", err
			}
		}

		if e.Format == PlainTable {
			if _, err = e.writer.WriteString(e.lineBreak); err != nil {
				return "", err
			}
			if err = e.formatTextHR(fieldWidths); err != nil {
				return "", err
			}
		}
	}

	if err = e.writer.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (e *Encoder) formatRecord(record []Field, widths []int) error {
	lineLen := 0
	for _, f := range record {
		n := len(f.Lines)
		if lineLen < n {
			lineLen = n
		}
	}

	for lineIdx := 0; lineIdx < lineLen; lineIdx++ {
		if 0 < lineIdx {
			if _, err := e.writer.WriteString(e.lineBreak); err != nil {
				return err
			}
		}

		line := make([]byte, 0, e.lineWidth*2)
		for i := 0; i < e.fieldLen; i++ {
			line = append(line, VLine)
			line = append(line, PadChar)

			if len(record) <= i || len(record[i].Lines) <= lineIdx || len(record[i].Lines[lineIdx]) < 1 {
				line = append(line, bytes.Repeat([]byte{PadChar}, widths[i]+1)...)
				continue
			}

			padLen := widths[i] - text.Width(record[i].Lines[lineIdx], e.EastAsianEncoding, e.CountDiacriticalSign, e.CountFormatCode)
			cellAlign := record[i].Alignment
			if (cellAlign == text.LeftAligned || cellAlign == text.NotAligned) && text.IsRightToLeftLetters(record[i].Lines[lineIdx]) {
				cellAlign = text.RightAligned
			}

			switch cellAlign {
			case text.Centering:
				halfPadLen := padLen / 2
				line = append(line, bytes.Repeat([]byte(string(PadChar)), halfPadLen)...)
				line = append(line, record[i].Lines[lineIdx]...)
				line = append(line, bytes.Repeat([]byte(string(PadChar)), (padLen-halfPadLen)+1)...)
			case text.RightAligned:
				line = append(line, bytes.Repeat([]byte(string(PadChar)), padLen)...)
				line = append(line, record[i].Lines[lineIdx]...)
				line = append(line, PadChar)
			default:
				line = append(line, record[i].Lines[lineIdx]...)
				line = append(line, bytes.Repeat([]byte(string(PadChar)), padLen+1)...)
			}
		}
		line = append(line, VLine)
		if _, err := e.writer.Write(line); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) formatTextHR(widths []int) error {
	line := make([]byte, 0, e.lineWidth)

	for _, w := range widths {
		line = append(line, CrossLine)
		line = append(line, bytes.Repeat([]byte{HLine}, w+2)...)
	}
	line = append(line, CrossLine)

	_, err := e.writer.Write(line)
	return err
}

func (e *Encoder) formatGfmHR(widths []int) error {
	line := make([]byte, 0, e.lineWidth)

	for i, w := range widths {
		align := text.NotAligned
		if i < len(e.alignments) {
			align = e.alignments[i]
		}

		line = append(line, VLine)
		line = append(line, PadChar)
		switch align {
		case text.Centering:
			line = append(line, AlignSign)
			line = append(line, bytes.Repeat([]byte{HLine}, w-2)...)
			line = append(line, AlignSign)
		case text.RightAligned:
			line = append(line, bytes.Repeat([]byte{HLine}, w-1)...)
			line = append(line, AlignSign)
		case text.LeftAligned:
			line = append(line, AlignSign)
			line = append(line, bytes.Repeat([]byte{HLine}, w-1)...)
		default:
			line = append(line, bytes.Repeat([]byte{HLine}, w)...)
		}
		line = append(line, PadChar)
	}
	line = append(line, VLine)

	_, err := e.writer.Write(line)
	return err
}

func (e *Encoder) formatOrgHR(widths []int) error {
	line := make([]byte, 0, e.lineWidth)

	line = append(line, VLine)
	for i, w := range widths {
		if 0 < i {
			line = append(line, CrossLine)
		}
		line = append(line, bytes.Repeat([]byte{HLine}, w+2)...)
	}
	line = append(line, VLine)

	_, err := e.writer.Write(line)
	return err
}

func (e *Encoder) calculateLineWidth(widths []int) {
	e.lineWidth = 1
	for _, w := range widths {
		e.lineWidth = e.lineWidth + w + 3
	}
}

func (e *Encoder) escape(s string) string {
	e.buf.Reset()

	runes := []rune(s)
	pos := 0

	for {
		if len(runes) <= pos {
			break
		}

		r := runes[pos]
		switch r {
		case '\r':
			if (pos+1) < len(runes) && runes[pos+1] == '\n' {
				pos++
			}
			fallthrough
		case '\n':
			switch e.Format {
			case GFMTable, OrgTable:
				e.buf.WriteString(MarkdownLineBreak)
			default:
				e.buf.WriteRune('\n')
			}
		case VLine:
			switch e.Format {
			case GFMTable, OrgTable:
				e.buf.WriteRune(EscapeChar)
			}
			fallthrough
		default:
			e.buf.WriteRune(r)
		}

		pos++
	}
	return e.buf.String()
}
