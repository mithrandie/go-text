package fixedlen

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/go-text"
)

var readerReadHeaderTests = []struct {
	Name               string
	Input              string
	DelimiterPositions []int
	Encoding           text.Encoding
	WithoutNull        bool
	Output             []string
	Error              string
}{
	{
		Name:               "ReadHeader",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		Encoding:           text.UTF8,
		WithoutNull:        false,
		Output:             []string{"ab", "cde", "fghi"},
	},
	{
		Name:               "ReadHeader with reversed position",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{6, 2, 9},
		Encoding:           text.UTF8,
		WithoutNull:        false,
		Error:              "invalid delimiter position: [6, 2, 9]",
	},
}

func TestFixedLengthReader_ReadHeader(t *testing.T) {
	for _, v := range readerReadHeaderTests {
		r, _ := NewReader(strings.NewReader(v.Input), v.DelimiterPositions, v.Encoding)

		r.WithoutNull = v.WithoutNull

		header, err := r.ReadHeader()

		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
		}
		if !reflect.DeepEqual(header, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, header, v.Output)
		}
	}
}

var readerReadAllTests = []struct {
	Name               string
	Input              string
	DelimiterPositions []int
	WithoutNull        bool
	SingleLine         bool
	Encoding           text.Encoding
	Output             [][]text.RawText
	ExpectLineBreak    text.LineBreak
	Error              string
}{
	{
		Name:               "ReadAll",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll with empty fields",
		Input:              "ab     \nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), nil, nil},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll with empty fields without nulls",
		Input:              "ab\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        true,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText(""), text.RawText("")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll with nil",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: nil,
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{},
			{},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll with empty positions",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{},
			{},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll with trimming value",
		Input:              "abcdefghi\nk   lm  no    pqurst",
		DelimiterPositions: []int{3, 11, 20},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("abc"), text.RawText("defghi"), nil},
			{text.RawText("k"), text.RawText("lm  no"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll with cutting length",
		Input:              "abcdefghi\nklmnopqurst\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqur")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqur")},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll from SJIS Text",
		Input:              "abcde" + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 11, 15},
		WithoutNull:        false,
		Encoding:           text.SJIS,
		Output: [][]text.RawText{
			{text.RawText("abcde"), text.RawText("日本語"), text.RawText("fghi")},
			{text.RawText("klmno"), text.RawText("pqurst"), nil},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "ReadAll LineBreak CR",
		Input:              "abcdefghi\rklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.CR,
	},
	{
		Name:               "ReadAll LineBreak CRLF",
		Input:              "abcdefghi\r\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.CRLF,
	},
	{
		Name:               "ReadAll LineBreak CR with cutting length",
		Input:              "abcdefghi\rklmnopqurst\rklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqur")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqur")},
		},
		ExpectLineBreak: text.CR,
	},
	{
		Name:               "ReadAll LineBreak CRLF with cutting length",
		Input:              "abcdefghi\r\nklmnopqurst\r\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqur")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqur")},
		},
		ExpectLineBreak: text.CRLF,
	},
	{
		Name:               "ReadAll with negative position",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{-2, 5},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Error:              "invalid delimiter position: [-2, 5]",
	},
	{
		Name:               "ReadAll with reversed position",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{6, 2, 9},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Error:              "invalid delimiter position: [6, 2, 9]",
	},
	{
		Name:               "ReadAll with position error",
		Input:              "abcde日本語fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 10, 15},
		WithoutNull:        false,
		Encoding:           text.UTF8,
		Error:              "cannot delimit lines in a byte array of a character",
	},
	{
		Name:               "ReadAll from SJIS Text with position error",
		Input:              "abcde" + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 10, 15},
		WithoutNull:        false,
		Encoding:           text.SJIS,
		Error:              "cannot delimit lines in a byte array of a character",
	},
	{
		Name:               "UTF-8 with BOM",
		Input:              string(text.UTF8BOM()) + "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           text.UTF8M,
		Output: [][]text.RawText{
			{text.RawText("ab"), text.RawText("cde"), text.RawText("fghi")},
			{text.RawText("kl"), text.RawText("mno"), text.RawText("pqurst")},
		},
		ExpectLineBreak: text.LF,
	},
	{
		Name:               "BOM does not exist error",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           text.UTF8M,
		Error:              "byte order mark for UTF-8 does not exist",
	},
	{
		Name:               "ReadAll Without LineBreak",
		Input:              "aaabbbbcccccdddeeeefffff",
		DelimiterPositions: []int{3, 7, 12},
		WithoutNull:        false,
		SingleLine:         true,
		Encoding:           text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("aaa"), text.RawText("bbbb"), text.RawText("ccccc")},
			{text.RawText("ddd"), text.RawText("eeee"), text.RawText("fffff")},
		},
		ExpectLineBreak: "",
	},
}

func TestFixedLengthReader_ReadAll(t *testing.T) {
	for _, v := range readerReadAllTests {
		r, err := NewReader(strings.NewReader(v.Input), v.DelimiterPositions, v.Encoding)
		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}

		r.WithoutNull = v.WithoutNull
		r.SingleLine = v.SingleLine

		records, err := r.ReadAll()

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
		if !reflect.DeepEqual(records, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, records, v.Output)
		}
		if r.DetectedLineBreak != v.ExpectLineBreak {
			t.Errorf("%s: detected line-break = %s, want %s", v.Name, r.DetectedLineBreak, v.ExpectLineBreak)
		}
	}
}
