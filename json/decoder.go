package json

import (
	"errors"
	"fmt"
)

type Decoder struct {
	// Returns numeric values as Integer or Float instead of Number if true.
	UseInteger bool
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d Decoder) Decode(src string) (Structure, EscapeType, error) {
	st, et, err := ParseJson(src, d.UseInteger)
	if err != nil {
		se := err.(*SyntaxError)
		return st, et, errors.New(fmt.Sprintf("line %d, column %d: %s", se.Line, se.Column, se.Error()))
	}
	return st, et, nil
}
