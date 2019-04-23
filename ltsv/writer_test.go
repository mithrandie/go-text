package ltsv

import (
	"bytes"
	"testing"

	"github.com/mithrandie/go-text"
)

var writerWriteTests = []struct {
	Name      string
	Header    []string
	Records   [][]string
	LineBreak text.LineBreak
	Encoding  text.Encoding
	Expect    string
	NewError  string
	Error     string
}{
	{
		Name:      "Empty Data",
		Records:   [][]string{},
		LineBreak: text.LF,
		Encoding:  text.UTF8,
		Expect:    "",
	},
	{
		Name:   "LTSV",
		Header: []string{"c1", "c2", "c3"},
		Records: [][]string{
			{
				"-1",
				"a",
				"",
			},
			{
				"2.0123",
				"2016-02-01T16:00:00.123456-07:00",
				"abc,de",
			},
		},
		LineBreak: text.LF,
		Encoding:  text.UTF8,
		Expect: "c1:-1\tc2:a\tc3:\n" +
			"c1:2.0123\tc2:2016-02-01T16:00:00.123456-07:00\tc3:abc,de",
	},
	{
		Name:   "Encode to SJIS",
		Header: []string{"c1", "c2", "c3"},
		Records: [][]string{
			{
				"-1",
				"日本語",
				"",
			},
			{
				"2.0123",
				"2016-02-01T16:00:00.123456-07:00",
				"abc,de",
			},
		},
		LineBreak: text.LF,
		Encoding:  text.SJIS,
		Expect: "c1:-1\tc2:" + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "\tc3:\n" +
			"c1:2.0123\tc2:2016-02-01T16:00:00.123456-07:00\tc3:abc,de",
	},
	{
		Name:   "Encode to UTF8M",
		Header: []string{"c1", "c2", "c3"},
		Records: [][]string{
			{
				"-1",
				"a",
				"",
			},
			{
				"2.0123",
				"2016-02-01T16:00:00.123456-07:00",
				"abc,de",
			},
		},
		LineBreak: text.LF,
		Encoding:  text.UTF8M,
		Expect: string(text.UTF8BOMS()) +
			"c1:-1\tc2:a\tc3:\n" +
			"c1:2.0123\tc2:2016-02-01T16:00:00.123456-07:00\tc3:abc,de",
	},
	{
		Name:   "Field Length Error",
		Header: []string{"c1", "c2", "c3", "c4"},
		Records: [][]string{
			{
				"-1",
				"a",
				"",
			},
			{
				"2.0123",
				"2016-02-01T16:00:00.123456-07:00",
				"abc,de",
			},
		},
		LineBreak: text.LF,
		Encoding:  text.UTF8,
		Error:     "field length does not match",
	},
	{
		Name:   "Unpermitted Character in Label",
		Header: []string{"c1:", "c2", "c3"},
		Records: [][]string{
			{
				"-1",
				"a",
				"",
			},
			{
				"2.0123",
				"2016-02-01T16:00:00.123456-07:00",
				"abc,de",
			},
		},
		LineBreak: text.LF,
		Encoding:  text.UTF8,
		NewError:  "unpermitted character in label: U+003A",
	},
	{
		Name:   "Unpermitted Character in Field Value",
		Header: []string{"c1", "c2", "c3"},
		Records: [][]string{
			{
				"-1",
				"a\t",
				"",
			},
			{
				"2.0123",
				"2016-02-01T16:00:00.123456-07:00",
				"abc,de",
			},
		},
		LineBreak: text.LF,
		Encoding:  text.UTF8,
		Error:     "unpermitted character in field-value: U+0009",
	},
}

func TestWriter_Write(t *testing.T) {
	for _, v := range writerWriteTests {
		b := new(bytes.Buffer)

		w, err := NewWriter(b, v.Header, v.LineBreak, v.Encoding)
		if err != nil {
			if v.NewError == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.NewError != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.NewError)
			}
			continue
		} else if v.NewError != "" {
			t.Errorf("%s: no error, want error %q", v.Name, v.NewError)
			continue
		}

		for _, r := range v.Records {
			err = w.Write(r)
			if err != nil {
				break
			}
		}
		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		} else if v.Error != "" {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		_ = w.Flush()

		result := b.String()

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
