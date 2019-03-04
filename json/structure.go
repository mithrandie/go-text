package json

import (
	"strconv"
	"strings"
)

const (
	BeginArray     = '['
	BeginObject    = '{'
	EndArray       = ']'
	EndObject      = '}'
	NameSeparator  = ':'
	ValueSeparator = ','
	QuotationMark  = '"'
	EscapeMark     = '\\'
)

var WhiteSpaces = []rune{
	32, //Space
	9,  //Horizontal Tab
	10, //Line Feed
	13, //Carriage Return
}

const (
	FalseValue = "false"
	TrueValue  = "true"
	NullValue  = "null"
)

type Structure interface {
	Encode() string
}

type Object struct {
	Members []ObjectMember
}

func NewObject(capacity int) Object {
	return Object{
		Members: make([]ObjectMember, 0, capacity),
	}
}

func (obj *Object) Len() int {
	return len(obj.Members)
}

func (obj *Object) Add(key string, val Structure) {
	obj.Members = append(obj.Members, ObjectMember{Key: key, Value: val})
}

func (obj *Object) Exists(key string) bool {
	for _, m := range obj.Members {
		if m.Key == key {
			return true
		}
	}
	return false
}

func (obj *Object) Value(key string) Structure {
	for _, m := range obj.Members {
		if m.Key == key {
			return m.Value
		}
	}
	return nil
}

func (obj *Object) Update(key string, val Structure) {
	for i, m := range obj.Members {
		if m.Key == key {
			obj.Members[i].Value = val
			break
		}
	}
}

func (obj Object) Encode() string {
	strs := make([]string, 0, obj.Len())
	for _, m := range obj.Members {
		strs = append(strs, Quote(Escape(m.Key))+string(NameSeparator)+m.Value.Encode())
	}
	return string(BeginObject) + strings.Join(strs[:], string(ValueSeparator)) + string(EndObject)
}

type ObjectMember struct {
	Key   string
	Value Structure
}

type Array []Structure

func (ar Array) Encode() string {
	strs := make([]string, 0, len(ar))
	for _, v := range ar {
		strs = append(strs, v.Encode())
	}
	return string(BeginArray) + strings.Join(strs[:], string(ValueSeparator)) + string(EndArray)
}

type Number float64

func (n Number) Encode() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n Number) Raw() float64 {
	return float64(n)
}

type Float float64

func (n Float) Encode() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n Float) Raw() float64 {
	return float64(n)
}

type Integer int64

func (n Integer) Encode() string {
	return strconv.FormatInt(int64(n), 10)
}

func (n Integer) Raw() int64 {
	return int64(n)
}

type String string

func (s String) Encode() string {
	return Quote(Escape(string(s)))
}

func (s String) Raw() string {
	return string(s)
}

type Boolean bool

func (b Boolean) Encode() string {
	if b {
		return TrueValue
	}
	return FalseValue
}

func (b Boolean) Raw() bool {
	return bool(b)
}

type Null struct{}

func (n Null) Encode() string {
	return NullValue
}

func Quote(s string) string {
	return string(QuotationMark) + s + string(QuotationMark)
}
