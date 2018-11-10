package table

import (
	"testing"

	"github.com/mithrandie/go-text"
)

var encoderEncodeTests = []struct {
	Name                 string
	Format               TableFormat
	Header               []*text.Field
	Records              [][]*text.Field
	Alignments           []text.FieldAlignment
	LineBreak            string
	EastAsiaEncoding     bool
	CountDiacriticalSign bool
	WithoutHeader        bool
	Expect               string
}{
	{
		Name:                 "Empty Fields",
		Format:               PlainTable,
		Header:               []*text.Field{},
		Records:              [][]*text.Field{},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     false,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect:               "",
	},
	{
		Name:   "Empty RecordSet",
		Format: PlainTable,
		Header: []*text.Field{
			{Contents: "c1", Alignment: text.Centering},
			{Contents: "c2", Alignment: text.Centering},
		},
		Records:              [][]*text.Field{},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     false,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect: "" +
			"+----+----+\n" +
			"| c1 | c2 |\n" +
			"+----+----+",
	},
	{
		Name:   "Text Table",
		Format: PlainTable,
		Header: []*text.Field{
			{Contents: "c1", Alignment: text.Centering},
			{Contents: "c2\nsecond line", Alignment: text.Centering},
			{Contents: "c3", Alignment: text.Centering},
		},
		Records: [][]*text.Field{
			{
				{Contents: "-1", Alignment: text.RightAligned},
				{Contents: "UNKNOWN", Alignment: text.Centering},
				{Contents: "false", Alignment: text.Centering},
			},
			{
				{Contents: "2.0123", Alignment: text.RightAligned},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Alignment: text.LeftAligned},
				{Contents: "abcdef", Alignment: text.LeftAligned},
			},
			{
				{Contents: "34567890", Alignment: text.RightAligned},
				{Contents: " ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n", Alignment: text.LeftAligned},
				{Contents: "NULL", Alignment: text.Centering},
			},
		},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     true,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect: "" +
			"+----------+-------------------------------------+--------+\n" +
			"|    c1    |                 c2                  |   c3   |\n" +
			"|          |             second line             |        |\n" +
			"+----------+-------------------------------------+--------+\n" +
			"|       -1 |               UNKNOWN               | false  |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00    | abcdef |\n" +
			"| 34567890 |  ab|cdefghijklmnopqrstuvwxyzabcdefg |  NULL  |\n" +
			"|          | hi\"jk日本語あアｱＡ（                |        |\n" +
			"|          |                                     |        |\n" +
			"+----------+-------------------------------------+--------+",
	},
	{
		Name:   "GFM Table",
		Format: GFMTable,
		Header: []*text.Field{
			{Contents: "c1", Alignment: text.Centering},
			{Contents: "c2\nsecond line", Alignment: text.Centering},
			{Contents: "c3", Alignment: text.Centering},
			{Contents: "c4", Alignment: text.Centering},
		},
		Records: [][]*text.Field{
			{
				{Contents: "-1", Alignment: text.RightAligned},
				{Contents: "", Alignment: text.Centering},
				{Contents: "false", Alignment: text.Centering},
			},
			{
				{Contents: "2.0123", Alignment: text.RightAligned},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Alignment: text.LeftAligned},
				{Contents: "abcdef", Alignment: text.LeftAligned},
			},
			{
				{Contents: "34567890", Alignment: text.RightAligned},
				{Contents: " ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n", Alignment: text.LeftAligned},
				{Contents: "", Alignment: text.Centering},
			},
		},
		Alignments: []text.FieldAlignment{
			text.RightAligned,
			text.Centering,
			text.LeftAligned,
		},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     true,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |  c4  |\n" +
			"| -------: | :-------------------------------------------------------------------: | :----- | ---- |\n" +
			"|       -1 |                                                                       | false  |      |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |      |\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |        |      |",
	},
	{
		Name:   "Org Table",
		Format: OrgTable,
		Header: []*text.Field{
			{Contents: "c1", Alignment: text.Centering},
			{Contents: "c2\nsecond line", Alignment: text.Centering},
			{Contents: "c3", Alignment: text.Centering},
		},
		Records: [][]*text.Field{
			{
				{Contents: "-1", Alignment: text.RightAligned},
				{Contents: "", Alignment: text.Centering},
				{Contents: "false", Alignment: text.Centering},
			},
			{
				{Contents: "2.0123", Alignment: text.RightAligned},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Alignment: text.LeftAligned},
				{Contents: "abcdef", Alignment: text.LeftAligned},
			},
			{
				{Contents: "34567890", Alignment: text.RightAligned},
				{Contents: " ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n", Alignment: text.LeftAligned},
				{Contents: "", Alignment: text.Centering},
			},
		},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     true,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |\n" +
			"|----------+-----------------------------------------------------------------------+--------|\n" +
			"|       -1 |                                                                       | false  |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |        |",
	},
	{
		Name:   "Right To Left Letters",
		Format: PlainTable,
		Header: []*text.Field{
			{Contents: "c1", Alignment: text.Centering},
			{Contents: "c2", Alignment: text.Centering},
		},
		Records: [][]*text.Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "العَرَبِيَّة", Alignment: text.LeftAligned},
			},
			{
				{Contents: "2.012", Alignment: text.RightAligned},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Alignment: text.LeftAligned},
			},
		},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     true,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect: "" +
			"+--------+----------------------------------+\n" +
			"|   c1   |                c2                |\n" +
			"+--------+----------------------------------+\n" +
			"| abc    |                          العَرَبِيَّة |\n" +
			"|  2.012 | 2016-02-01T16:00:00.123456-07:00 |\n" +
			"+--------+----------------------------------+",
	},
	{
		Name:   "Defferent Length Records",
		Format: PlainTable,
		Header: []*text.Field{
			{Contents: "c1", Alignment: text.Centering},
			{Contents: "c2", Alignment: text.Centering},
			{Contents: "c3", Alignment: text.Centering},
		},
		Records: [][]*text.Field{
			{
				{Contents: "-1", Alignment: text.RightAligned},
				{Contents: "UNKNOWN", Alignment: text.Centering},
				{Contents: "false", Alignment: text.Centering},
			},
			{
				{Contents: "2.0123", Alignment: text.RightAligned},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Alignment: text.LeftAligned},
			},
		},
		LineBreak:            text.LF.Value(),
		EastAsiaEncoding:     false,
		CountDiacriticalSign: false,
		WithoutHeader:        false,
		Expect: "" +
			"+--------+----------------------------------+--------+\n" +
			"|   c1   |                c2                |   c3   |\n" +
			"+--------+----------------------------------+--------+\n" +
			"|     -1 |             UNKNOWN              | false  |\n" +
			"| 2.0123 | 2016-02-01T16:00:00.123456-07:00 |        |\n" +
			"+--------+----------------------------------+--------+",
	},
}

func TestEncoder_Encode(t *testing.T) {
	for _, v := range encoderEncodeTests {
		var e *Encoder
		e = NewEncoder(v.Format, len(v.Records))
		e.LineBreak = v.LineBreak
		e.EastAsianEncoding = v.EastAsiaEncoding
		e.CountDiacriticalSign = v.CountDiacriticalSign
		e.WithoutHeader = v.WithoutHeader

		e.SetHeader(v.Header)
		for _, r := range v.Records {
			e.AppendRecord(r)
		}
		if v.Alignments != nil {
			e.SetFieldAlignments(v.Alignments)
		}

		result, _ := e.Encode()

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
