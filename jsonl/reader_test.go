package jsonl

import (
	"github.com/mithrandie/go-text/json"
	"reflect"
	"strings"
	"testing"
)

var readerReadAllTests = []struct {
	Label      string
	Input      string
	UseInteger bool
	Expect     []json.Structure
	EscapeType json.EscapeType
	Error      string
}{
	{
		Label:  "Empty String",
		Input:  "",
		Expect: []json.Structure{},
	},
	{
		Label: "Read Object List",
		Input: "{\"key\":\"value\", \"key2\":\"value2\"}\n" +
			"{\"key\":\"value3\", \"key2\":\"value4\"}\n",
		Expect: []json.Structure{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value3"),
					},
					{
						Key:   "key2",
						Value: json.String("value4"),
					},
				},
			},
		},
	},
	{
		Label: "Read Array List",
		Input: "[\"value1\", \"value2\"]\n" +
			"[\"value3\", \"value4\"]",
		Expect: []json.Structure{
			json.Array{
				json.String("value1"),
				json.String("value2"),
			},
			json.Array{
				json.String("value3"),
				json.String("value4"),
			},
		},
	},
	{
		Label: "Read Mixed List",
		Input: "{\"key\":\"value\", \"key2\":\"value2\"}\n" +
			"[\"value3\", \"value4\"]\n" +
			"1",
		Expect: []json.Structure{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
					},
				},
			},
			json.Array{
				json.String("value3"),
				json.String("value4"),
			},
			json.Number(1),
		},
	},
	{
		Label: "Ignore empty lines",
		Input: "{\"key\":\"value\", \"key2\":\"value2\"}\n" +
			"\n" +
			"{\"key\":\"value3\", \"key2\":\"value4\"}\n" +
			"\n",
		Expect: []json.Structure{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value3"),
					},
					{
						Key:   "key2",
						Value: json.String("value4"),
					},
				},
			},
		},
	},
	{
		Label: "Read Object List using Hex Digit Encoding",
		Input: "{\"key\":\"value\", \"key2\":\"value2\"}\n" +
			"{\"key\":\"value\\u0033\", \"key2\":\"value4\"}",
		Expect: []json.Structure{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value3"),
					},
					{
						Key:   "key2",
						Value: json.String("value4"),
					},
				},
			},
		},
		EscapeType: json.AllWithHexDigits,
	},
	{
		Label: "Invalid JSON Line",
		Input: "{\"key\":\"value\", \"key2\":\"value2\"}\n" +
			"A{\"key\":\"value\\u0033\", \"key2\":\"value4\"}",
		Error: "line 2, column 1: unexpected token \"A\"",
	},
}

func TestReader_ReadAll(t *testing.T) {
	for _, v := range readerReadAllTests {
		r := NewReader(strings.NewReader(v.Input))
		r.SetUseInteger(v.UseInteger)

		st, et, err := r.ReadAll()

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Label)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Label)
			}
			continue
		}

		if !reflect.DeepEqual(st, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", st, v.Expect, v.Label)
		}

		if et != v.EscapeType {
			t.Errorf("escape type = %d, want %d for %q", et, v.EscapeType, v.Label)
		}
	}
}
