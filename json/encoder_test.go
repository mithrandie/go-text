package json

import (
	"math"
	"testing"

	"github.com/mithrandie/go-text"
)

var encoderEncodeTests = []struct {
	Input          Structure
	Escape         EscapeType
	PrettyPrint    bool
	NanInfHandling NanInfHandling
	FloatFormat    FloatFormat
	LineBreak      text.LineBreak
	UsePalette     bool
	Expect         string
	Error          string
}{
	{
		Input:       String("abc"),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "\"abc\"",
	},
	{
		Input:       String("1"),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "\"1\"",
	},
	{
		Input:       String("01"),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "\"01\"",
	},
	{
		Input:       Number(-1.234),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "-1.234",
	},
	{
		Input:          Number(0.0000000000123),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToNull,
		FloatFormat:    ENotationForLargeExponents,
		LineBreak:      text.LF,
		Expect:         "1.23e-11",
	},
	{
		Input:          Number(0.0000000000123),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: CreateError,
		FloatFormat:    ENotationForLargeExponents,
		LineBreak:      text.LF,
		Expect:         "1.23e-11",
	},
	{
		Input:          Number(0.0000000000123),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToStringNotation,
		FloatFormat:    ENotationForLargeExponents,
		LineBreak:      text.LF,
		Expect:         "1.23e-11",
	},
	{
		Input:          Number(-1.234),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: CreateError,
		LineBreak:      text.LF,
		Expect:         "-1.234",
	},
	{
		Input:          Number(-1.234),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToStringNotation,
		LineBreak:      text.LF,
		Expect:         "-1.234",
	},
	{
		Input:          Number(math.Inf(1)),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToNull,
		LineBreak:      text.LF,
		Expect:         "null",
	},
	{
		Input:          Number(math.Inf(1)),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: CreateError,
		LineBreak:      text.LF,
		Error:          "+Inf is not supported",
	},
	{
		Input:          Number(math.Inf(1)),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToStringNotation,
		LineBreak:      text.LF,
		Expect:         "+Inf",
	},
	{
		Input:          Number(math.NaN()),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToNull,
		LineBreak:      text.LF,
		Expect:         "null",
	},
	{
		Input:          Number(math.NaN()),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: CreateError,
		LineBreak:      text.LF,
		Error:          "NaN is not supported",
	},
	{
		Input:          Number(math.NaN()),
		Escape:         Backslash,
		PrettyPrint:    false,
		NanInfHandling: ConvertToStringNotation,
		LineBreak:      text.LF,
		Expect:         "NaN",
	},
	{
		Input:       Float(-1.234),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "-1.234",
	},
	{
		Input:       Integer(1234),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "1234",
	},
	{
		Input:       Boolean(true),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "true",
	},
	{
		Input:       Boolean(false),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "false",
	},
	{
		Input:       Null{},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "null",
	},
	{
		Input: Array{
			String("value1"),
			String("value2"),
			String("value3"),
		},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "[\"value1\",\"value2\",\"value3\"]",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "{\"key1\":\"value1\",\"key2\":\"value2\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key\"1",
					Value: String("value\"1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "{\"key\\\"1\":\"value\\\"1\",\"key2\":\"value2\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key\"1",
					Value: String("value\"1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Escape:      HexDigits,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "{\"key\\u00221\":\"value\\u00221\",\"key2\":\"value2\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key\"1",
					Value: String("value\"1"),
				},
			},
		},
		Escape:      AllWithHexDigits,
		PrettyPrint: false,
		LineBreak:   text.LF,
		Expect:      "{\"\\u006b\\u0065\\u0079\\u0022\\u0031\":\"\\u0076\\u0061\\u006c\\u0075\\u0065\\u0022\\u0031\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
				{
					Key: "key2",
					Value: Array{
						Object{
							Members: []ObjectMember{
								{
									Key:   "akey1",
									Value: Boolean(true),
								},
								{
									Key:   "akey2",
									Value: Null{},
								},
								{
									Key:   "akey3",
									Value: Array{},
								},
							},
						},
						Object{
							Members: []ObjectMember{
								{
									Key:   "akey1",
									Value: Number(-2.3e-6),
								},
								{
									Key: "akey2",
									Value: Array{
										String("A"),
										String("B"),
										String("C"),
									},
								},
								{
									Key:   "akey3",
									Value: String(""),
								},
							},
						},
					},
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		Expect: "{\n" +
			"  \"key1\": \"value1\",\n" +
			"  \"key2\": [\n" +
			"    {\n" +
			"      \"akey1\": true,\n" +
			"      \"akey2\": null,\n" +
			"      \"akey3\": []\n" +
			"    },\n" +
			"    {\n" +
			"      \"akey1\": -0.0000023,\n" +
			"      \"akey2\": [\n" +
			"        \"A\",\n" +
			"        \"B\",\n" +
			"        \"C\"\n" +
			"      ],\n" +
			"      \"akey3\": \"\"\n" +
			"    }\n" +
			"  ]\n" +
			"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
				{
					Key:   "key2",
					Value: String("[1, 2, 3]"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		Expect: "{\n" +
			"  \"key1\": \"value1\",\n" +
			"  \"key2\": [\n" +
			"    1,\n" +
			"    2,\n" +
			"    3\n" +
			"  ]\n" +
			"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.CRLF,
		Expect: "{\r\n" +
			"  \"key1\": \"value1\",\r\n" +
			"  \"key2\": \"value2\"\r\n" +
			"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		UsePalette:  true,
		Expect: "{\n" +
			"  \x1b[34;1m\"key1\"\x1b[0m: \x1b[32m\"value1\"\x1b[0m\n" +
			"}",
	},
	{
		Input:       Number(1),
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		UsePalette:  true,
		Expect:      "\x1b[35m1\x1b[0m",
	},
	{
		Input:       String("abc"),
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		UsePalette:  true,
		Expect:      "\x1b[32m\"abc\"\x1b[0m",
	},
	{
		Input:       Boolean(true),
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		UsePalette:  true,
		Expect:      "\x1b[33;1mtrue\x1b[0m",
	},
	{
		Input:       Null{},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   text.LF,
		UsePalette:  true,
		Expect:      "\x1b[90mnull\x1b[0m",
	},
	{
		Input:       Number(1),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   text.LF,
		UsePalette:  true,
		Expect:      "1",
	},
}

func TestEncoder_Encode(t *testing.T) {
	palette := NewJsonPalette()

	for _, v := range encoderEncodeTests {
		e := NewEncoder()

		e.EscapeType = v.Escape
		e.PrettyPrint = v.PrettyPrint
		e.NanInfHandling = v.NanInfHandling
		e.FloatFormat = v.FloatFormat
		e.LineBreak = v.LineBreak
		if v.UsePalette {
			e.Palette = palette
		} else {
			e.Palette = nil
		}

		result, err := e.Encode(v.Input)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err.Error(), v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if result != v.Expect {
			t.Errorf("result = %q, want %q for EscapeType:%d PrettyPrint:%t Input:%#v", result, v.Expect, v.Escape, v.PrettyPrint, v.Input)
		}
	}
}
