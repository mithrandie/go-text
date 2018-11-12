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
	}
	obj.Add("ar", ar)
	obj.Add("null", Null{})
	obj.Add("bool", Boolean(false))

	encoded := obj.Encode()
	expectJson := "{" +
		"\"str\":\"updated\"," +
		"\"ar\":[1,2,3]," +
		"\"null\":null," +
		"\"bool\":false" +
		"}"
	if encoded != expectJson {
		t.Errorf("json = %s, want %s", encoded, expectJson)
	}
}
