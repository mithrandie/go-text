package json

import (
	"strings"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/color"
)

func NewJsonPalette() *color.Palette {
	objKey := color.NewEffector()
	objKey.SetEffect(color.Bold)
	objKey.SetFGColor(color.Blue)

	str := color.NewEffector()
	str.SetFGColor(color.Green)

	num := color.NewEffector()
	num.SetFGColor(color.Magenta)

	b := color.NewEffector()
	b.SetFGColor(color.Yellow)
	b.SetEffect(color.Bold)

	n := color.NewEffector()
	n.SetFGColor(color.BrightBlack)

	p := color.NewPalette()
	p.SetEffector(ObjectKeyEffect, objKey)
	p.SetEffector(StringEffect, str)
	p.SetEffector(NumberEffect, num)
	p.SetEffector(BooleanEffect, b)
	p.SetEffector(NullEffect, n)

	return p
}

type Encoder struct {
	EscapeType   EscapeType
	PrettyPrint  bool
	LineBreak    text.LineBreak
	IndentSpaces int
	Palette      *color.Palette

	nameSeparator string
	lineBreak     string

	decoder *Decoder
}

func NewEncoder() *Encoder {
	return &Encoder{
		EscapeType:    Backslash,
		PrettyPrint:   false,
		LineBreak:     text.LF,
		IndentSpaces:  2,
		Palette:       nil,
		nameSeparator: string(NameSeparator),
		decoder:       NewDecoder(),
	}
}

func (e *Encoder) Encode(structure Structure) string {
	if e.PrettyPrint {
		e.lineBreak = e.LineBreak.Value()
		e.nameSeparator = string(NameSeparator) + " "
		if e.Palette != nil {
			e.Palette.Enable()
		}
	} else {
		e.lineBreak = ""
		e.nameSeparator = string(NameSeparator)
		if e.Palette != nil {
			e.Palette.Disable()
		}
	}

	return e.encodeStructure(structure, 0)
}

func (e *Encoder) encodeStructure(structure Structure, depth int) string {
	var indent string
	var elementIndent string
	if e.PrettyPrint {
		indent = strings.Repeat(" ", e.IndentSpaces*depth)
		elementIndent = strings.Repeat(" ", e.IndentSpaces*(depth+1))
	}

	var encoded string

	switch structure.(type) {
	case Object:
		obj := structure.(Object)
		strs := make([]string, 0, obj.Len())
		for _, member := range obj.Members {
			strs = append(
				strs,
				elementIndent+
					e.effect(ObjectKeyEffect, e.formatString(member.Key))+
					e.nameSeparator+
					e.encodeStructure(member.Value, depth+1),
			)
		}
		encoded = string(BeginObject) +
			e.lineBreak +
			strings.Join(strs[:], string(ValueSeparator)+e.lineBreak) +
			e.lineBreak +
			indent + string(EndObject)
	case Array:
		array := structure.(Array)
		strs := make([]string, 0, len(array))
		for _, v := range array {
			strs = append(strs, elementIndent+e.encodeStructure(v, depth+1))
		}
		if len(strs) < 1 {
			encoded = string(BeginArray) + string(EndArray)
		} else {
			encoded = string(BeginArray) +
				e.lineBreak +
				strings.Join(strs[:], string(ValueSeparator)+e.lineBreak) +
				e.lineBreak +
				indent + string(EndArray)
		}
	case Number, Float, Integer:
		encoded = e.effect(NumberEffect, structure.Encode())
	case String:
		str := structure.(String).Raw()
		if 0 < len(str) {
			if decoded, _, err := e.decoder.Decode(str); err == nil && isComplexType(decoded) {
				encoded = e.encodeStructure(decoded, depth)
			} else {
				encoded = e.effect(StringEffect, e.formatString(str))
			}
		} else {
			encoded = e.effect(StringEffect, e.formatString(str))
		}
	case Boolean:
		encoded = e.effect(BooleanEffect, structure.Encode())
	case Null:
		encoded = e.effect(NullEffect, structure.Encode())
	}

	return encoded
}

func (e *Encoder) formatString(s string) string {
	var escaped string

	switch e.EscapeType {
	case AllWithHexDigits:
		escaped = EscapeAll(s)
	case HexDigits:
		escaped = EscapeWithHexDigits(s)
	default:
		escaped = Escape(s)
	}

	return string(QuotationMark) + escaped + string(QuotationMark)
}

func (e *Encoder) effect(key string, s string) string {
	if e.Palette == nil {
		return s
	}
	return e.Palette.Render(key, s)
}

func isComplexType(s Structure) bool {
	return isObject(s) || isArray(s)
}

func isObject(s Structure) bool {
	_, ok := s.(Object)
	return ok
}

func isArray(s Structure) bool {
	_, ok := s.(Array)
	return ok
}
