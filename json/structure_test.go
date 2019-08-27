package json

import (
	"reflect"
	"testing"
)

func TestObject(t *testing.T) {
	obj := NewObject(0)

	exists := obj.Exists("str")
	if exists {
		t.Errorf("exists = %t, want %t", exists, false)
	}

	obj.Add("str", String("value"))

	exists = obj.Exists("str")
	if !exists {
		t.Errorf("exists = %t, want %t", exists, true)
	}

	value := obj.Value("str")
	var expectValue Structure = String("value")
	if !reflect.DeepEqual(value, expectValue) {
		t.Errorf("value = %#v, want %#v", value, expectValue)
	}

	value = obj.Value("notexist")
	expectValue = nil
	if !reflect.DeepEqual(value, expectValue) {
		t.Errorf("value = %#v, want %#v", value, expectValue)
	}

	obj.Update("str", String("updated"))
	value = obj.Value("str")
	expectValue = String("updated")
	if !reflect.DeepEqual(value, expectValue) {
		t.Errorf("value = %#v, want %#v", value, expectValue)
	}

	ar := Array{
		Number(1),
		Number(2),
		Number(3),
		Float(4.56),
		Integer(789),
	}
	obj.Add("ar", ar)
	obj.Add("null", Null{})
	obj.Add("bool", Boolean(false))

	encoded := obj.Encode()
	expectJson := "{" +
		"\"str\":\"updated\"," +
		"\"ar\":[1,2,3,4.56,789]," +
		"\"null\":null," +
		"\"bool\":false" +
		"}"
	if encoded != expectJson {
		t.Errorf("json = %s, want %s", encoded, expectJson)
	}

	keys := obj.Keys()
	expectKeys := []string{"str", "ar", "null", "bool"}
	if !reflect.DeepEqual(keys, expectKeys) {
		t.Errorf("keys = %s, want %s", keys, expectKeys)
	}

	values := obj.Values()
	expectValues := []Structure{
		String("updated"),
		Array{
			Number(1),
			Number(2),
			Number(3),
			Float(4.56),
			Integer(789),
		},
		Null{},
		Boolean(false),
	}
	if !reflect.DeepEqual(values, expectValues) {
		t.Errorf("values = %s, want %s", values, expectValues)
	}

	appendedKeys := make([]string, 0, obj.Len())
	obj.Range(func(key string, value Structure) bool {
		if key == "null" {
			return false
		}
		appendedKeys = append(appendedKeys, key)
		return true
	})
	expectAppendedKeys := []string{"str", "ar"}
	if !reflect.DeepEqual(appendedKeys, expectAppendedKeys) {
		t.Errorf("appended keys = %s, want %s", appendedKeys, expectAppendedKeys)
	}
}
