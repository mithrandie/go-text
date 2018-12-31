package ltsv

import (
	"unicode"
)

var LabelTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x002D, 0x002D, 1}, // Hyphen-minux
		{0x002E, 0x002E, 1}, // Full Stop
		{0x0030, 0x0039, 1}, // ASCII Digits
		{0x0041, 0x005A, 1}, // Latin Alphabet Upper Case
		{0x005F, 0x005F, 1}, // Low Line
		{0x0061, 0x007A, 1}, // Latin Alphabet Lower Case
	},
}

var FieldValueTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0001, 0x0008, 1},
		{0x000B, 0x000B, 1}, // Vertical Tab
		{0x000C, 0x000C, 1}, // Form Feed
		{0x000E, 0xFFFF, 1},
	},
	R32: []unicode.Range32{
		{0x10000, 0xfffff, 1},
	},
}
