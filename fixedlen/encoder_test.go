package fixedlen

import (
	"testing"

	"github.com/mithrandie/go-text"
)

var fixedLengthEncoderEncodeTests = []struct {
	Name               string
	Header             []Field
	Records            [][]Field
	DelimiterPositions []int
	LineBreak          text.LineBreak
	WithoutHeader      bool
	Encoding           text.Encoding
	Expect             string
	Error              string
}{
	{
		Name:               "Empty Positions",
		Header:             []Field{},
		Records:            [][]Field{},
		DelimiterPositions: []int{},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Expect:             "",
	},
	{
		Name: "Empty RecordSet",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
		},
		Records:            [][]Field{},
		DelimiterPositions: []int{10, 42},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Expect: "" +
			"c1        c2                              ",
	},
	{
		Name: "Fixed-Length Encode",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "-1", Alignment: text.RightAligned},
				{Contents: "", Alignment: text.Centering},
				{Contents: "false", Alignment: text.LeftAligned},
			},
			{
				{Contents: "2.0123", Alignment: text.RightAligned},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Alignment: text.LeftAligned},
				{Contents: "abcdef", Alignment: text.LeftAligned},
			},
			{
				{Contents: "true", Alignment: text.LeftAligned},
				{Contents: " abc", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{10, 42, 50},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Expect: "" +
			"c1        c2                              c3      \n" +
			"        -1                                false   \n" +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  \n" +
			"true       abc                                    ",
	},
	{
		Name: "Fixed-Length Encode to SJIS",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "日本語", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.SJIS,
		Expect: "" +
			"c1   c2        c3   \n" +
			"abc  " + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "    def  \n" +
			"ghi  jkl       mno  ",
	},
	{
		Name: "Fixed-Length Encode Without Header",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          text.LF,
		WithoutHeader:      true,
		Encoding:           text.UTF8,
		Expect: "" +
			"abc  def       def  \n" +
			"ghi  jkl       mno  ",
	},
	{
		Name: "Fixed-Length Encode with Empty Field",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{5, 15, 20, 25},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Expect: "" +
			"c1   c2        c3        \n" +
			"abc  def       def       \n" +
			"ghi  jkl       mno       ",
	},
	{
		Name: "Fixed-Length Encode Invalid Positions",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{5, 15, 10},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Error:              "invalid delimiter position: [5, 15, 10]",
	},
	{
		Name: "Fixed-Length Encode Field Length too long in Header",
		Header: []Field{
			{Contents: "cccccc1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Error:              "value is too long: \"cccccc1\" for 5 byte(s) length field",
	},
	{
		Name: "Fixed-Length Encode Field Length too long in Record",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abcabc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Error:              "value is too long: \"abcabc\" for 5 byte(s) length field",
	},
	{
		Name: "Fixed-Length Encode Concatnate Automatically",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
			{Contents: "c3"},
		},
		Records: [][]Field{
			{
				{Contents: "abcabc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: nil,
		LineBreak:          text.LF,
		WithoutHeader:      false,
		Encoding:           text.UTF8,
		Expect: "" +
			"c1     c2  c3 \n" +
			"abcabc def def\n" +
			"ghi    jkl mno",
	},
	{
		Name:   "Fixed-Length Encode Concatnate Automatically Without Header",
		Header: nil,
		Records: [][]Field{
			{
				{Contents: "abcabc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		DelimiterPositions: nil,
		LineBreak:          text.LF,
		WithoutHeader:      true,
		Encoding:           text.UTF8,
		Expect: "" +
			"abcabc def def\n" +
			"ghi    jkl mno",
	},
}

func TestFixedLengthEncoder_Encode(t *testing.T) {
	for _, v := range fixedLengthEncoderEncodeTests {
		e := NewEncoder(len(v.Records))
		e.DelimiterPositions = v.DelimiterPositions
		e.LineBreak = v.LineBreak.Value()
		e.WithoutHeader = v.WithoutHeader
		e.Encoding = v.Encoding

		e.SetHeader(v.Header)
		for _, r := range v.Records {
			e.AppendRecord(r)
		}

		result, err := e.Encode()

		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
		}
		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
