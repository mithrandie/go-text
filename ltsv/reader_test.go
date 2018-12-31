package ltsv

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/go-text"
)

func TestRecord_Clear(t *testing.T) {
	r := make(Record)
	r["key1"] = new(bytes.Buffer)
	r["key1"].WriteString("a")
	r["key2"] = new(bytes.Buffer)
	r["key2"].WriteString("b")

	r.Clear()
	for k := range r {
		if r[k].Len() != 0 {
			t.Errorf("field length = %d, want %d", len(r), 0)
		}
	}
}

var readAllTests = []struct {
	Name        string
	Encoding    text.Encoding
	WithoutNull bool
	Input       string
	Output      [][]text.RawText
	Fields      []string
	LineBreak   text.LineBreak
	Error       string
}{
	{
		Name:        "LineBreak LF",
		Encoding:    text.UTF8,
		WithoutNull: false,
		Input:       "f1:v1\tf2:v2\tf3:v3\nf1:v4\tf2:v5\tf3:v6\n\n",
		Output: [][]text.RawText{
			{text.RawText("v1"), text.RawText("v2"), text.RawText("v3")},
			{text.RawText("v4"), text.RawText("v5"), text.RawText("v6")},
		},
		Fields:    []string{"f1", "f2", "f3"},
		LineBreak: text.LF,
	},
	{
		Name:        "LineBreak CRLF",
		Encoding:    text.UTF8,
		WithoutNull: false,
		Input:       "f1:v1\tf2:v2\tf3:v3\r\nf1:v4\tf2:v5\tf3:v6",
		Output: [][]text.RawText{
			{text.RawText("v1"), text.RawText("v2"), text.RawText("v3")},
			{text.RawText("v4"), text.RawText("v5"), text.RawText("v6")},
		},
		Fields:    []string{"f1", "f2", "f3"},
		LineBreak: text.CRLF,
	},
	{
		Name:        "LineBreak CR",
		Encoding:    text.UTF8,
		WithoutNull: false,
		Input:       "f1:v1\tf2:v2\tf3:v3\rf1:v4\tf2:v5\tf3:v6",
		Output: [][]text.RawText{
			{text.RawText("v1"), text.RawText("v2"), text.RawText("v3")},
			{text.RawText("v4"), text.RawText("v5"), text.RawText("v6")},
		},
		Fields:    []string{"f1", "f2", "f3"},
		LineBreak: text.CR,
	},
	{
		Name:        "Difference Keys",
		Encoding:    text.UTF8,
		WithoutNull: false,
		Input:       "f1:v1\tf2:v2\tf3:v3\nf3:v6\tf1:v4\tf4:v7",
		Output: [][]text.RawText{
			{text.RawText("v1"), text.RawText("v2"), text.RawText("v3"), nil},
			{text.RawText("v4"), nil, text.RawText("v6"), text.RawText("v7")},
		},
		Fields:    []string{"f1", "f2", "f3", "f4"},
		LineBreak: text.LF,
	},
	{
		Name:        "Without Null",
		Encoding:    text.UTF8,
		WithoutNull: true,
		Input:       "f1:v1\tf2:v2\tf3:v3\nf3:v6\tf1:v4\tf4:v7",
		Output: [][]text.RawText{
			{text.RawText("v1"), text.RawText("v2"), text.RawText("v3"), text.RawText("")},
			{text.RawText("v4"), text.RawText(""), text.RawText("v6"), text.RawText("v7")},
		},
		Fields:    []string{"f1", "f2", "f3", "f4"},
		LineBreak: text.LF,
	},
	{
		Name:        "SJIS",
		Encoding:    text.SJIS,
		WithoutNull: false,
		Input:       "f1:v1\tf2:v2\tf3:v3\nf1:v4\tf2:" + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "\tf3:v6",
		Output: [][]text.RawText{
			{text.RawText("v1"), text.RawText("v2"), text.RawText("v3")},
			{text.RawText("v4"), text.RawText("日本語"), text.RawText("v6")},
		},
		Fields:    []string{"f1", "f2", "f3"},
		LineBreak: text.LF,
	},
	{
		Name:        "Invalid Format",
		Encoding:    text.UTF8,
		WithoutNull: false,
		Input:       "f1:v1\tf2\tf3:v3\nf1:v4\tf2:v5\tf3:v6",
		Error:       "line 1, column 9: missing field separator",
	},
}

func TestReader_ReadAll(t *testing.T) {
	for _, v := range readAllTests {
		r := NewReader(strings.NewReader(v.Input), v.Encoding)
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

		if !reflect.DeepEqual(records, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, records, v.Output)
			t.Errorf("%s: records = %#v, want %#v", v.Name, records, v.Output)
		}

		if r.DetectedLineBreak != v.LineBreak {
			t.Errorf("%s: line break = %q, want %q", v.Name, r.DetectedLineBreak, v.LineBreak)
		}

		if !reflect.DeepEqual(r.Header.Fields(), v.Fields) {
			t.Errorf("%s: fields = %v, want %v", v.Name, r.Header.Fields(), v.Fields)
		}
	}
}
