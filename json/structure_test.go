package json

import (
	"math"
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
		Integer(math.MaxInt64),
	}
	obj.Add("ar", ar)
	obj.Add("null", Null{})
	obj.Add("bool", Boolean(false))

	encoded := obj.Encode()
	expectJson := "{" +
		"\"str\":\"updated\"," +
		"\"ar\":[1,2,3,4.56,9223372036854776000]," +
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
			Integer(math.MaxInt64),
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

var numberEncodeTests = []struct {
	value  Float
	expect string
}{
	{
		value:  1.234,
		expect: "1.234",
	},
	{
		value:  Number(math.Inf(1)),
		expect: "null",
	},
	{
		value:  Number(math.Inf(-1)),
		expect: "null",
	},
	{
		value:  Number(math.NaN()),
		expect: "null",
	},
}

func TestNumber_Encode(t *testing.T) {
	for _, v := range numberEncodeTests {
		r := v.value.Encode()
		if r != v.expect {
			t.Errorf("encoded string = %s, want %s", r, v.expect)
		}
	}
}

func TestNumber_IsPositiveInfinity(t *testing.T) {
	f := Float(math.Inf(1))
	ret := f.IsPositiveInfinity()
	if ret != true {
		t.Errorf("IsPositiveInfinity = %t, want %t", ret, true)
	}

	f = Float(math.Inf(-1))
	ret = f.IsPositiveInfinity()
	if ret != false {
		t.Errorf("IsPositiveInfinity = %t, want %t", ret, false)
	}
}

func TestNumber_IsNegativeInfinity(t *testing.T) {
	f := Float(math.Inf(1))
	ret := f.IsNegativeInfinity()
	if ret != false {
		t.Errorf("IsNegativeInfinity = %t, want %t", ret, false)
	}

	f = Float(math.Inf(-1))
	ret = f.IsNegativeInfinity()
	if ret != true {
		t.Errorf("IsNegativeInfinity = %t, want %t", ret, true)
	}
}
