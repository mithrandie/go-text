package fixedlen

import (
	"bytes"
	"testing"

	"github.com/mithrandie/go-text"
)

var writerWriteTests = []struct {
	Name               string
	Records            [][]Field
	DelimiterPositions []int
	LineBreak          text.LineBreak
	Encoding           text.Encoding
	InsertSpace        bool
	SingleLine         bool
	Expect             string
	Error              string
}{
	{
		Name: "Empty Positions",
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
		DelimiterPositions: nil,
		LineBreak:          text.LF,
		Encoding:           text.UTF8,
		Expect:             "\n\n",
	},
	{
		Name:               "Empty RecordSet",
		Records:            [][]Field{},
		DelimiterPositions: []int{10, 42},
		LineBreak:          text.LF,
		Encoding:           text.UTF8,
		Expect:             "",
	},
	{
		Name: "Fixed-Length Encode",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		Encoding:           text.UTF8,
		Expect: "" +
			"c1        c2                              c3      \n" +
			"        -1                                false   \n" +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  \n" +
			"true       abc                                    ",
	},
	{
		Name: "Fixed-Length Encode to SJIS",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		Encoding:           text.SJIS,
		Expect: "" +
			"c1   c2        c3   \n" +
			"abc  " + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "    def  \n" +
			"ghi  jkl       mno  ",
	},
	{
		Name: "Encode to UTF8M",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		Encoding:           text.UTF8M,
		Expect: text.UTF8BOM +
			"c1        c2                              c3      \n" +
			"        -1                                false   \n" +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  \n" +
			"true       abc                                    ",
	},
	{
		Name: "Fixed-Length Encode with Empty Field",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		Encoding:           text.UTF8,
		Expect: "" +
			"c1   c2        c3        \n" +
			"abc  def       def       \n" +
			"ghi  jkl       mno       ",
	},
	{
		Name: "Fixed-Length Encode Invalid Positions",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		Encoding:           text.UTF8,
		Error:              "invalid delimiter position: [5, 15, 10]",
	},
	{
		Name: "Fixed-Length Encode Field size too long",
		Records: [][]Field{
			{
				{Contents: "cccccc1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		Encoding:           text.UTF8,
		Error:              "value is too long: \"cccccc1\" for 5 byte(s) length field",
	},
	{
		Name: "Fixed-Length Encode Insert Space",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		DelimiterPositions: []int{6, 9, 12},
		LineBreak:          text.LF,
		Encoding:           text.UTF8,
		InsertSpace:        true,
		Expect: "" +
			"c1     c2  c3 \n" +
			"abcabc def def\n" +
			"ghi    jkl mno",
	},
	{
		Name: "Fixed-Length Encode SingleLine",
		Records: [][]Field{
			{
				{Contents: "c1"},
				{Contents: "c2"},
				{Contents: "c3"},
			},
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
		DelimiterPositions: []int{6, 9, 12},
		LineBreak:          text.LF,
		Encoding:           text.UTF8,
		InsertSpace:        false,
		SingleLine:         true,
		Expect: "" +
			"c1    c2 c3 " +
			"abcabcdefdef" +
			"ghi   jklmno",
	},
}

func TestWriter_Write(t *testing.T) {
	for _, v := range writerWriteTests {
		errOccurred := false

		w := new(bytes.Buffer)

		e, _ := NewWriter(w, v.DelimiterPositions, v.LineBreak, v.Encoding)
		e.InsertSpace = v.InsertSpace
		e.SingleLine = v.SingleLine

		for _, r := range v.Records {
			err := e.Write(r)
			if err != nil {
				errOccurred = true
				if len(v.Error) < 1 {
					t.Errorf("%s: unexpected error %q", v.Name, err.Error())
					break
				} else if err.Error() != v.Error {
					t.Errorf("%s: error %q, want error %q", v.Name, err, v.Error)
					break
				}
			}
		}
		_ = e.Flush()

		if !errOccurred && 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if errOccurred {
			continue
		}

		result := w.String()

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
