package csv

import (
	"bytes"
	"testing"

	"github.com/mithrandie/go-text"
)

var writerWriteTests = []struct {
	Name      string
	Records   [][]Field
	Delimiter rune
	LineBreak text.LineBreak
	Encoding  text.Encoding
	Expect    string
}{
	{
		Name:      "Empty Data",
		Records:   [][]Field{},
		Delimiter: ',',
		LineBreak: text.LF,
		Encoding:  text.UTF8,
		Expect:    "",
	},
	{
		Name: "CSV",
		Records: [][]Field{
			{
				{Contents: "c1", Quote: true},
				{Contents: "c2\nsecond line", Quote: true},
				{Contents: "c3", Quote: true},
			},
			{
				{Contents: "-1", Quote: false},
				{Contents: "", Quote: false},
				{Contents: "true", Quote: false},
			},
			{
				{Contents: "2.0123", Quote: false},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Quote: true},
				{Contents: "abc,de\"f", Quote: false},
			},
		},
		Delimiter: ',',
		LineBreak: text.LF,
		Expect: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abc,de\"\"f\"",
	},
	{
		Name: "TSV",
		Records: [][]Field{
			{
				{Contents: "c1", Quote: true},
				{Contents: "c2\nsecond line", Quote: true},
				{Contents: "c3", Quote: true},
			},
			{
				{Contents: "-1", Quote: false},
				{Contents: "", Quote: false},
				{Contents: "true", Quote: false},
			},
			{
				{Contents: "2.0123", Quote: false},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Quote: true},
				{Contents: "abc,de\"f", Quote: false},
			},
		},
		Delimiter: '\t',
		LineBreak: text.LF,
		Expect: "\"c1\"\t\"c2\nsecond line\"\t\"c3\"\n" +
			"-1\t\ttrue\n" +
			"2.0123\t\"2016-02-01T16:00:00.123456-07:00\"\tabc,de\"f",
	},
	{
		Name: "Encode to SJIS",
		Records: [][]Field{
			{
				{Contents: "c1", Quote: true},
				{Contents: "c2\nsecond line", Quote: true},
				{Contents: "c3", Quote: true},
			},
			{
				{Contents: "-1", Quote: false},
				{Contents: "", Quote: false},
				{Contents: "true", Quote: false},
			},
			{
				{Contents: "2.0123", Quote: false},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Quote: true},
				{Contents: "日本語", Quote: false},
			},
		},
		Delimiter: ',',
		LineBreak: text.LF,
		Encoding:  text.SJIS,
		Expect: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\"," + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}),
	},
	{
		Name: "Encode to UTF8M",
		Records: [][]Field{
			{
				{Contents: "c1", Quote: true},
				{Contents: "c2\nsecond line", Quote: true},
				{Contents: "c3", Quote: true},
			},
			{
				{Contents: "-1", Quote: false},
				{Contents: "", Quote: false},
				{Contents: "true", Quote: false},
			},
			{
				{Contents: "2.0123", Quote: false},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Quote: true},
				{Contents: "abc,de\"f", Quote: false},
			},
		},
		Delimiter: ',',
		LineBreak: text.LF,
		Encoding:  text.UTF8M,
		Expect: string(text.UTF8BOM()) + "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abc,de\"\"f\"",
	},
}

func TestWriter_Write(t *testing.T) {
	for _, v := range writerWriteTests {
		w := new(bytes.Buffer)

		e, _ := NewWriter(w, v.LineBreak, v.Encoding)
		e.Delimiter = v.Delimiter
		for _, r := range v.Records {
			_ = e.Write(r)
		}
		_ = e.Flush()

		result := w.String()

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
