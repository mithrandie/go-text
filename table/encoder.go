package table

import (
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
	LineBreak            string
	EastAsianEncoding    bool
	CountDiacriticalSign bool
	Encoding             text.Encoding

	// GFM or Org Table
	WithoutHeader bool

	// GFM Table only
	alignments []text.FieldAlignment

	header    []*Field
	recordSet [][]*Field
	fieldLen  int
	buf       bytes.Buffer
}

func NewEncoder(format Format, recordCounts int) *Encoder {
	return &Encoder{
		Format:               format,
		LineBreak:            text.LF.Value(),
		EastAsianEncoding:    false,
		CountDiacriticalSign: false,
		Encoding:             text.UTF8,
		WithoutHeader:        false,
		fieldLen:             0,
		recordSet:            make([][]*Field, 0, recordCounts),
	}
}

func (e *Encoder) SetHeader(header []*Field) {
	e.header = e.prepareRecord(header)
	if e.fieldLen < len(header) {
		e.fieldLen = len(header)
	}
}

func (e *Encoder) SetFieldAlignments(alignments []text.FieldAlignment) {
	e.alignments = alignments
}

func (e *Encoder) AppendRecord(record []*Field) {
	e.recordSet = append(e.recordSet, e.prepareRecord(record))
	if e.fieldLen < len(record) {
		e.fieldLen = len(record)
	}
}

func (e *Encoder) prepareRecord(record []*Field) []*Field {
	for _, f := range record {
		e.prepareField(f)
	}
	return record
}

func (e *Encoder) prepareField(field *Field) *Field {
	lines := strings.Split(e.escape(field.Contents), e.LineBreak)

	width := 0
	for _, v := range lines {
		l := text.Width(v, e.EastAsianEncoding, e.CountDiacriticalSign)
		if width < l {
			width = l
		}
	}

	field.Lines = lines
	field.Width = width
	return field
}

func (e *Encoder) Encode() (string, error) {
	if e.fieldLen < 1 {
		return "", nil
	}

	lines := make([]string, 0, len(e.recordSet)+4)

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

		if e.Format == PlainTable {
			lines = append(lines, e.formatTextHR(fieldWidths))
		}
		lines = append(lines, e.formatRecord(e.header, fieldWidths))

		switch e.Format {
		case GFMTable:
			lines = append(lines, e.formatGfmHR(fieldWidths))
		case OrgTable:
			lines = append(lines, e.formatOrgHR(fieldWidths))
		default:
			lines = append(lines, e.formatTextHR(fieldWidths))
		}
	}

	if 0 < len(e.recordSet) {
		for _, record := range e.recordSet {
			lines = append(lines, e.formatRecord(record, fieldWidths))
		}

		if e.Format == PlainTable {
			lines = append(lines, e.formatTextHR(fieldWidths))
		}
	}

	return text.Encode(strings.Join(lines, e.LineBreak), e.Encoding)
}

func (e *Encoder) formatRecord(record []*Field, widths []int) string {
	lineLen := 0
	for _, f := range record {
		n := len(f.Lines)
		if lineLen < n {
			lineLen = n
		}
	}

	lines := make([]string, 0, lineLen)

	for lineIdx := 0; lineIdx < lineLen; lineIdx++ {
		e.buf.Reset()

		for i := 0; i < e.fieldLen; i++ {
			e.buf.WriteRune(VLine)
			e.buf.WriteRune(PadChar)

			if len(record) <= i || len(record[i].Lines) <= lineIdx || len(record[i].Lines[lineIdx]) < 1 {
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), widths[i]+1))
				continue
			}

			padLen := widths[i] - text.Width(record[i].Lines[lineIdx], e.EastAsianEncoding, e.CountDiacriticalSign)
			cellAlign := record[i].Alignment
			if cellAlign == text.LeftAligned && text.IsRightToLeftLetters(record[i].Lines[lineIdx]) {
				cellAlign = text.RightAligned
			}

			switch cellAlign {
			case text.Centering:
				halfPadLen := padLen / 2
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), halfPadLen))
				e.buf.WriteString(record[i].Lines[lineIdx])
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), (padLen-halfPadLen)+1))
			case text.RightAligned:
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen))
				e.buf.WriteString(record[i].Lines[lineIdx])
				e.buf.WriteRune(PadChar)
			default:
				e.buf.WriteString(record[i].Lines[lineIdx])
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen+1))
			}
		}
		e.buf.WriteRune(VLine)
		lines = append(lines, e.buf.String())
	}

	return strings.Join(lines, e.LineBreak)
}

func (e *Encoder) formatTextHR(widths []int) string {
	e.buf.Reset()

	for _, w := range widths {
		e.buf.WriteRune(CrossLine)
		e.buf.Write(bytes.Repeat([]byte(string(HLine)), w+2))
	}
	e.buf.WriteRune(CrossLine)
	return e.buf.String()
}

func (e *Encoder) formatGfmHR(widths []int) string {
	e.buf.Reset()

	for i, w := range widths {
		align := text.NotAligned
		if i < len(e.alignments) {
			align = e.alignments[i]
		}

		e.buf.WriteRune(VLine)
		e.buf.WriteRune(PadChar)
		switch align {
		case text.Centering:
			e.buf.WriteRune(AlignSign)
			e.buf.Write(bytes.Repeat([]byte(string(HLine)), w-2))
			e.buf.WriteRune(AlignSign)
		case text.RightAligned:
			e.buf.Write(bytes.Repeat([]byte(string(HLine)), w-1))
			e.buf.WriteRune(AlignSign)
		case text.LeftAligned:
			e.buf.WriteRune(AlignSign)
			e.buf.Write(bytes.Repeat([]byte(string(HLine)), w-1))
		default:
			e.buf.Write(bytes.Repeat([]byte(string(HLine)), w))
		}
		e.buf.WriteRune(PadChar)
	}
	e.buf.WriteRune(VLine)
	return e.buf.String()
}

func (e *Encoder) formatOrgHR(widths []int) string {
	e.buf.Reset()

	e.buf.WriteRune(VLine)
	for i, w := range widths {
		if 0 < i {
			e.buf.WriteRune(CrossLine)
		}
		e.buf.Write(bytes.Repeat([]byte(string(HLine)), w+2))
	}
	e.buf.WriteRune(VLine)
	return e.buf.String()
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
				e.buf.WriteString(e.LineBreak)
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
