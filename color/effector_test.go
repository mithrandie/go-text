package color

import (
	"testing"
)

var effectorRenderTests = []struct {
	Name       string
	Text       string
	Enclose    bool
	EffectCode []EffectCode
	FGColor    Code
	BGColor    Code
	FG256Color int
	BG256Color int
	FGRGBColor []int
	BGRGBColor []int
	Expect     string
}{
	{
		Name:    "No Effect",
		Text:    "text",
		Enclose: true,
		Expect:  "text",
	},
	{
		Name:       "Bold and Italic",
		Text:       "text",
		Enclose:    true,
		EffectCode: []EffectCode{Bold, Italic},
		Expect:     "\033[1;3mtext\033[0m",
	},
	{
		Name:       "Empty EffectCode",
		Text:       "text",
		Enclose:    true,
		EffectCode: []EffectCode{},
		Expect:     "\033[0mtext\033[0m",
	},
	{
		Name:       "Unenclosed",
		Text:       "text",
		Enclose:    false,
		EffectCode: []EffectCode{Bold},
		Expect:     "\033[1mtext",
	},
	{
		Name:       "Color",
		Text:       "text",
		Enclose:    true,
		EffectCode: []EffectCode{Bold},
		FGColor:    Cyan,
		BGColor:    Red,
		Expect:     "\033[41;36;1mtext\033[0m",
	},
	{
		Name:       "256-Color",
		Text:       "text",
		Enclose:    true,
		EffectCode: []EffectCode{Bold},
		FG256Color: 29,
		BG256Color: 100,
		Expect:     "\033[48;5;100;38;5;29;1mtext\033[0m",
	},
	{
		Name:       "RGB Color",
		Text:       "text",
		Enclose:    true,
		EffectCode: []EffectCode{Bold},
		FGRGBColor: []int{160, 180, 200},
		BGRGBColor: []int{150, 130, 100},
		Expect:     "\033[48;2;150;130;100;38;2;160;180;200;1mtext\033[0m",
	},
}

func TestEffector_Render(t *testing.T) {
	for _, v := range effectorRenderTests {
		e := NewEffector()
		e.Enclose = v.Enclose
		if v.EffectCode != nil {
			e.SetEffect(v.EffectCode...)
		}
		if 0 < v.FGColor {
			e.SetFGColor(v.FGColor)
		}
		if 0 < v.BGColor {
			e.SetBGColor(v.BGColor)
		}
		if 0 < v.FG256Color {
			e.SetFG256Color(v.FG256Color)
		}
		if 0 < v.BG256Color {
			e.SetBG256Color(v.BG256Color)
		}
		if v.FGRGBColor != nil {
			e.SetFGRGBColor(v.FGRGBColor[0], v.FGRGBColor[1], v.FGRGBColor[2])
		}
		if v.BGRGBColor != nil {
			e.SetBGRGBColor(v.BGRGBColor[0], v.BGRGBColor[1], v.BGRGBColor[2])
		}

		result := e.Render(v.Text)
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %s", result, v.Expect, v.Name)
		}
	}
}
