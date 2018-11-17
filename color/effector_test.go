package color

import (
	"reflect"
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

var effectorExportConfigTests = []struct {
	Effector Effector
	Expect   EffectorConfig
}{
	{
		Effector: Effector{
			Enclose: true,
			effects: []EffectCode{Bold, Italic},
			fgColor: color8{color: Cyan},
			bgColor: nil,
		},
		Expect: EffectorConfig{
			Effects:    []string{"Bold", "Italic"},
			Foreground: "Cyan",
			Background: nil,
		},
	},
	{
		Effector: Effector{
			Enclose: false,
			effects: nil,
			fgColor: nil,
			bgColor: color256{color: 45},
		},
		Expect: EffectorConfig{
			Effects:    []string{},
			Foreground: nil,
			Background: 45,
		},
	},
	{
		Effector: Effector{
			Enclose: true,
			effects: []EffectCode{},
			fgColor: colorRGB{red: 24, green: 165, blue: 45},
			bgColor: nil,
		},
		Expect: EffectorConfig{
			Effects:    []string{},
			Foreground: []int{24, 165, 45},
			Background: nil,
		},
	},
}

func TestEffector_ExportConfig(t *testing.T) {
	for _, v := range effectorExportConfigTests {
		result := v.Effector.ExportConfig()
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %#v", result, v.Expect, v.Effector)
		}
	}
}

var generateEffectorTests = []struct {
	Config EffectorConfig
	Expect *Effector
	Error  string
}{
	{
		Config: EffectorConfig{
			Effects:    []string{"Bold", "Italic"},
			Foreground: "Cyan",
			Background: nil,
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{Bold, Italic},
			fgColor:  color8{color: Cyan},
			bgColor:  nil,
			sequence: "\033[36;1;3m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: nil,
			Background: 45,
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  nil,
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    nil,
			Foreground: []int{24, 165, 45},
			Background: nil,
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  colorRGB{red: 24, green: 165, blue: 45},
			bgColor:  nil,
			sequence: "\033[38;2;24;165;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    nil,
			Foreground: []interface{}{24, 165, 45},
			Background: nil,
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  colorRGB{red: 24, green: 165, blue: 45},
			bgColor:  nil,
			sequence: "\033[38;2;24;165;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: float32(45),
			Background: float64(45),
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  color256{color: 45},
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45;38;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: int(45),
			Background: int8(45),
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  color256{color: 45},
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45;38;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: int16(45),
			Background: int32(45),
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  color256{color: 45},
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45;38;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: int64(45),
			Background: uint(45),
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  color256{color: 45},
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45;38;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: uint8(45),
			Background: uint16(45),
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  color256{color: 45},
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45;38;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{},
			Foreground: uint32(45),
			Background: uint64(45),
		},
		Expect: &Effector{
			Enclose:  true,
			effects:  []EffectCode{},
			fgColor:  color256{color: 45},
			bgColor:  color256{color: 45},
			sequence: "\033[48;5;45;38;5;45m",
		},
	},
	{
		Config: EffectorConfig{
			Effects:    []string{"Bold", "invalid"},
			Foreground: "Cyan",
			Background: nil,
		},
		Error: "\"invalid\" cannot convert to EffectCode",
	},
	{
		Config: EffectorConfig{
			Effects:    []string{"Bold"},
			Foreground: "invalid",
			Background: nil,
		},
		Error: "\"invalid\" cannot convert to Color Code",
	},
	{
		Config: EffectorConfig{
			Effects:    []string{"Bold"},
			Foreground: nil,
			Background: []int{1, 0},
		},
		Error: "[1 0] cannot convert to color",
	},
	{
		Config: EffectorConfig{
			Effects:    []string{"Bold"},
			Foreground: nil,
			Background: true,
		},
		Error: "true cannot convert to color",
	},
}

func TestGenerateEffector(t *testing.T) {
	for _, v := range generateEffectorTests {
		result, err := GenerateEffector(v.Config)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %v", err.Error(), v.Config)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %v", err, v.Error, v.Config)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %v", v.Error, v.Config)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %v", result, v.Expect, v.Config)
		}
	}
}
