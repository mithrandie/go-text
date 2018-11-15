package color

import (
	"testing"
)

var effectCodeStringTests = []struct {
	Code   EffectCode
	Expect string
}{
	{
		Code:   Reset,
		Expect: "Reset",
	},
	{
		Code:   Bold,
		Expect: "Bold",
	},
	{
		Code:   Faint,
		Expect: "Faint",
	},
	{
		Code:   Italic,
		Expect: "Italic",
	},
	{
		Code:   Underline,
		Expect: "Underline",
	},
	{
		Code:   SlowBlink,
		Expect: "SlowBlink",
	},
	{
		Code:   RapidBlink,
		Expect: "RapidBlink",
	},
	{
		Code:   ReverseVideo,
		Expect: "ReverseVideo",
	},
	{
		Code:   Conceal,
		Expect: "Conceal",
	},
	{
		Code:   CrossedOut,
		Expect: "CrossedOut",
	},
	{
		Code:   EffectCode(-1),
		Expect: "",
	},
}

func TestEffectCode_String(t *testing.T) {
	for _, v := range effectCodeStringTests {
		result := v.Code.String()
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %#v", result, v.Expect, v.Code)
		}
	}
}

var codeStringTests = []struct {
	Code   Code
	Expect string
}{
	{
		Code:   Black,
		Expect: "Black",
	},
	{
		Code:   Red,
		Expect: "Red",
	},
	{
		Code:   Green,
		Expect: "Green",
	},
	{
		Code:   Yellow,
		Expect: "Yellow",
	},
	{
		Code:   Blue,
		Expect: "Blue",
	},
	{
		Code:   Magenta,
		Expect: "Magenta",
	},
	{
		Code:   Cyan,
		Expect: "Cyan",
	},
	{
		Code:   White,
		Expect: "White",
	},
	{
		Code:   BrightBlack,
		Expect: "BrightBlack",
	},
	{
		Code:   BrightRed,
		Expect: "BrightRed",
	},
	{
		Code:   BrightGreen,
		Expect: "BrightGreen",
	},
	{
		Code:   BrightYellow,
		Expect: "BrightYellow",
	},
	{
		Code:   BrightBlue,
		Expect: "BrightBlue",
	},
	{
		Code:   BrightMagenta,
		Expect: "BrightMagenta",
	},
	{
		Code:   BrightCyan,
		Expect: "BrightCyan",
	},
	{
		Code:   BrightWhite,
		Expect: "BrightWhite",
	},
	{
		Code:   DefaultColor,
		Expect: "DefaultColor",
	},
	{
		Code:   Code(-1),
		Expect: "",
	},
}

func TestCode_String(t *testing.T) {
	for _, v := range codeStringTests {
		result := v.Code.String()
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %#v", result, v.Expect, v.Code)
		}
	}
}

var parseEffectCodeTests = []struct {
	Input  string
	Expect EffectCode
	Error  string
}{
	{
		Input:  "Reset",
		Expect: Reset,
	},
	{
		Input:  "Bold",
		Expect: Bold,
	},
	{
		Input:  "Faint",
		Expect: Faint,
	},
	{
		Input:  "Italic",
		Expect: Italic,
	},
	{
		Input:  "Underline",
		Expect: Underline,
	},
	{
		Input:  "SlowBlink",
		Expect: SlowBlink,
	},
	{
		Input:  "RapidBlink",
		Expect: RapidBlink,
	},
	{
		Input:  "ReverseVideo",
		Expect: ReverseVideo,
	},
	{
		Input:  "Conceal",
		Expect: Conceal,
	},
	{
		Input:  "CrossedOut",
		Expect: CrossedOut,
	},
	{
		Input: "invalid",
		Error: "\"invalid\" cannot convert to EffectCode",
	},
}

func TestParseEffectCode(t *testing.T) {
	for _, v := range parseEffectCodeTests {
		result, err := ParseEffectCode(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %q", result, v.Expect, v.Input)
		}
	}
}

var parseColorCodeTests = []struct {
	Input  string
	Expect Code
	Error  string
}{
	{
		Input:  "Black",
		Expect: Black,
	},
	{
		Input:  "Red",
		Expect: Red,
	},
	{
		Input:  "Green",
		Expect: Green,
	},
	{
		Input:  "Yellow",
		Expect: Yellow,
	},
	{
		Input:  "Blue",
		Expect: Blue,
	},
	{
		Input:  "Magenta",
		Expect: Magenta,
	},
	{
		Input:  "Cyan",
		Expect: Cyan,
	},
	{
		Input:  "White",
		Expect: White,
	},
	{
		Input:  "BrightBlack",
		Expect: BrightBlack,
	},
	{
		Input:  "BrightRed",
		Expect: BrightRed,
	},
	{
		Input:  "BrightGreen",
		Expect: BrightGreen,
	},
	{
		Input:  "BrightYellow",
		Expect: BrightYellow,
	},
	{
		Input:  "BrightBlue",
		Expect: BrightBlue,
	},
	{
		Input:  "BrightMagenta",
		Expect: BrightMagenta,
	},
	{
		Input:  "BrightCyan",
		Expect: BrightCyan,
	},
	{
		Input:  "BrightWhite",
		Expect: BrightWhite,
	},
	{
		Input:  "DefaultColor",
		Expect: DefaultColor,
	},
	{
		Input: "invalid",
		Error: "\"invalid\" cannot convert to Color Code",
	},
}

func TestParseColorCode(t *testing.T) {
	for _, v := range parseColorCodeTests {
		result, err := ParseColorCode(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %q", result, v.Expect, v.Input)
		}
	}
}
