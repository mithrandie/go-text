package csv

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/go-text"
)

var readAllTests = []struct {
	Name        string
	Encoding    text.Encoding
	Delimiter   rune
	WithoutNull bool
	Input       string
	Output      [][]text.RawText
	LineBreak   text.LineBreak
	EnclosedAll bool
	Error       string
}{
	{
		Name:     "NewLineLF",
		Input:    "a,b,c\nd,e,f",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("e"), text.RawText("f")},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "NewLineCR",
		Input:    "a,b,c\rd,e,f",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("e"), text.RawText("f")},
		},
		LineBreak: text.CR,
	},
	{
		Name:     "NewLineCRLF",
		Input:    "a,b,c\r\nd,e,f",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("e"), text.RawText("f")},
		},
		LineBreak: text.CRLF,
	},
	{
		Name:      "TabDelimiter",
		Delimiter: '\t',
		Input:     "a\tb\tc\nd\te\tf",
		Encoding:  text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("e"), text.RawText("f")},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "QuotedString",
		Input:    "a,\"b\",\"ccc\ncc\"\nd,e,",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("ccc\ncc")},
			{text.RawText("d"), text.RawText("e"), nil},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "EscapeDoubleQuote",
		Input:    "a,\"b\",\"ccc\"\"cc\"\nd,e,\"\"",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("ccc\"cc")},
			{text.RawText("d"), text.RawText("e"), text.RawText("")},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "DoubleQuoteInNoQuoteField",
		Input:    "a,b,ccc\"cc\nd,e,",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("ccc\"cc")},
			{text.RawText("d"), text.RawText("e"), nil},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "SingleValue",
		Input:    "a",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a")},
		},
		LineBreak: "",
	},
	{
		Name:     "Trailing empty lines",
		Input:    "a,b,c\nd,e,f\n\n",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("e"), text.RawText("f")},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "Different Line Breaks",
		Input:    "a,b,\"c\r\nd\"\ne,f,g",
		Encoding: text.UTF8,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c\r\nd")},
			{text.RawText("e"), text.RawText("f"), text.RawText("g")},
		},
		LineBreak: text.LF,
	},
	{
		Name:     "Decode Character Code",
		Encoding: text.SJIS,
		Input:    "a,b,c\nd," + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + ",f",
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("日本語"), text.RawText("f")},
		},
		LineBreak: text.LF,
	},
	{
		Name:        "Without Null",
		Input:       "\"a\",\"b\",\"1\"\n\"d\",,2",
		Encoding:    text.UTF8,
		WithoutNull: true,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("1")},
			{text.RawText("d"), text.RawText(""), text.RawText("2")},
		},
		LineBreak:   text.LF,
		EnclosedAll: true,
	},
	{
		Name:     "ExtraneousQuote",
		Input:    "a,\"b\",\"ccc\ncc\nd,e,",
		Encoding: text.UTF8,
		Error:    "line 3, column 5: extraneous \" in field",
	},
	{
		Name:     "UnexpectedQuote",
		Input:    "a,\"b\",\"ccc\"cc\nd,e,",
		Encoding: text.UTF8,
		Error:    "line 1, column 11: unexpected \" in field",
	},
	{
		Name:     "NumberOfFieldsIsLess",
		Input:    "a,b,c\nd,e\nf,g,h",
		Encoding: text.UTF8,
		Error:    "line 2, column 0: wrong number of fields in line",
	},
	{
		Name:     "NumberOfFieldsIsGreater",
		Input:    "a,b,c\nd,e,f,g\nh,i,j",
		Encoding: text.UTF8,
		Error:    "line 2, column 6: wrong number of fields in line",
	},
	{
		Name:     "UTF8 with BOM",
		Input:    string([]byte(text.UTF8BOM)) + "a,b,c\nd,e,f",
		Encoding: text.UTF8M,
		Output: [][]text.RawText{
			{text.RawText("a"), text.RawText("b"), text.RawText("c")},
			{text.RawText("d"), text.RawText("e"), text.RawText("f")},
		},
		LineBreak: text.LF,
	},
}

func TestReader_ReadAll(t *testing.T) {
	for _, v := range readAllTests {
		r, err := NewReader(strings.NewReader(v.Input), v.Encoding)
		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}

		if v.Delimiter != 0 {
			r.Delimiter = v.Delimiter
		}
		r.WithoutNull = v.WithoutNull

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
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if !reflect.DeepEqual(records, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, records, v.Output)
			t.Errorf("%s: records = %#v, want %#v", v.Name, records, v.Output)
		}

		if r.DetectedLineBreak != v.LineBreak {
			t.Errorf("%s: line break = %q, want %q", v.Name, r.DetectedLineBreak, v.LineBreak)
		}

		if r.EnclosedAll != v.EnclosedAll {
			t.Errorf("%s: enclosed all = %t, want %t", v.Name, r.EnclosedAll, v.EnclosedAll)
		}
	}
}

func TestReader_ReadHeader(t *testing.T) {
	input := "h1,h2 ,h3\na,b,c\nd,e,f"
	outHeader := []string{"h1", "h2 ", "h3"}
	output := [][]text.RawText{
		{text.RawText("a"), text.RawText("b"), text.RawText("c")},
		{text.RawText("d"), text.RawText("e"), text.RawText("f")},
	}

	r, _ := NewReader(strings.NewReader(input), text.UTF8)

	header, err := r.ReadHeader()
	if err != nil {
		t.Errorf("unexpected error %q", err.Error())
	}
	if !reflect.DeepEqual(header, outHeader) {
		t.Errorf("header = %q, want %q", header, outHeader)
	}

	records, err := r.ReadAll()
	if err != nil {
		t.Errorf("unexpected error %q", err.Error())
	}
	if !reflect.DeepEqual(records, output) {
		t.Errorf("records = %q, want %q", records, output)
	}

	input = "h1,\"h2 ,h3\na,b,c\nd,e,f"
	expectErr := "line 3, column 6: extraneous \" in field"

	r, _ = NewReader(strings.NewReader(input), text.UTF8)
	_, err = r.ReadHeader()
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q", err.Error(), expectErr)
	}
}

var readerReadAllBenchmarkText = strings.Repeat("aaaaaa,\"bbbbbb\",cccccc\n", 10000)

func BenchmarkReader_ReadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(readerReadAllBenchmarkText)
		reader, _ := NewReader(r, text.UTF8)
		reader.Delimiter = ','
		reader.WithoutNull = false
		_, _ = reader.ReadAll()
	}
}
