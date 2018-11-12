package csv

import (
	"testing"

	"github.com/mithrandie/go-text"
)

var encoderEncodeTests = []struct {
	Name          string
	Header        []Field
	Records       [][]Field
	Delimiter     rune
	LineBreak     text.LineBreak
	WithoutHeader bool
	Encoding      text.Encoding
	Expect        string
}{
	{
		Name:          "Empty Positions",
		Header:        []Field{},
		Records:       [][]Field{},
		Delimiter:     ',',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Encoding:      text.UTF8,
		Expect:        "",
	},
	{
		Name: "Empty RecordSet",
		Header: []Field{
			{Contents: "c1"},
			{Contents: "c2"},
		},
		Records:       [][]Field{},
		Delimiter:     ',',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Encoding:      text.UTF8,
		Expect: "" +
			"c1,c2",
	},
	{
		Name: "CSV",
		Header: []Field{
			{Contents: "c1", Quote: true},
			{Contents: "c2\nsecond line", Quote: true},
			{Contents: "c3", Quote: true},
		},
		Records: [][]Field{
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
		Delimiter:     ',',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Expect: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abc,de\"\"f\"",
	},
	{
		Name: "TSV",
		Header: []Field{
			{Contents: "c1", Quote: true},
			{Contents: "c2\nsecond line", Quote: true},
			{Contents: "c3", Quote: true},
		},
		Records: [][]Field{
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
		Delimiter:     '\t',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Expect: "\"c1\"\t\"c2\nsecond line\"\t\"c3\"\n" +
			"-1\t\ttrue\n" +
			"2.0123\t\"2016-02-01T16:00:00.123456-07:00\"\tabc,de\"f",
	},
	{
		Name: "Uneven Fields",
		Header: []Field{
			{Contents: "c1", Quote: true},
			{Contents: "c2\nsecond line", Quote: true},
		},
		Records: [][]Field{
			{
				{Contents: "-1", Quote: false},
				{Contents: "", Quote: false},
			},
			{
				{Contents: "2.0123", Quote: false},
				{Contents: "2016-02-01T16:00:00.123456-07:00", Quote: true},
				{Contents: "abc,de\"f", Quote: false},
			},
		},
		Delimiter:     ',',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Expect: "\"c1\",\"c2\nsecond line\",\n" +
			"-1,,\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abc,de\"\"f\"",
	},
	{
		Name: "CSV Encode to SJIS",
		Header: []Field{
			{Contents: "c1", Quote: true},
			{Contents: "c2\nsecond line", Quote: true},
			{Contents: "c3", Quote: true},
		},
		Records: [][]Field{
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
		Delimiter:     ',',
		LineBreak:     text.LF,
		WithoutHeader: false,
		Encoding:      text.SJIS,
		Expect: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\"," + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}),
	},
}

func TestEncoder_Encode(t *testing.T) {
	for _, v := range encoderEncodeTests {
		var e *Encoder
		e = NewEncoder(len(v.Records))
		e.Delimiter = v.Delimiter
		e.LineBreak = v.LineBreak.Value()
		e.WithoutHeader = v.WithoutHeader
		e.Encoding = v.Encoding

		e.SetHeader(v.Header)
		for _, r := range v.Records {
			e.AppendRecord(r)
		}

		result, _ := e.Encode()

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
