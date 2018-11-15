package color

import (
	"reflect"
	"testing"
)

var paletteRenderTests = []struct {
	Name   string
	Text   string
	Enable bool
	Effect string
	Expect string
}{
	{
		Name:   "Render Error",
		Text:   "text",
		Enable: true,
		Effect: "error",
		Expect: "\033[31;1mtext\033[0m",
	},
	{
		Name:   "Effect Not Exist",
		Text:   "text",
		Enable: true,
		Effect: "notexist",
		Expect: "text",
	},
	{
		Name:   "Disabled",
		Text:   "text",
		Enable: false,
		Effect: "error",
		Expect: "text",
	},
}

func TestPalette_Render(t *testing.T) {
	errKey := "error"
	infoKey := "info"

	errorEffector := NewEffector()
	errorEffector.SetFGColor(Red)
	errorEffector.SetEffect(Bold)

	infoEffector := NewEffector()
	infoEffector.SetFGColor(Green)

	palette := NewPalette()
	palette.SetEffector(errKey, errorEffector)
	palette.SetEffector(infoKey, infoEffector)

	for _, v := range paletteRenderTests {
		if v.Enable {
			palette.Enable()
		} else {
			palette.Disable()
		}

		result := palette.Render(v.Effect, v.Text)
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %s", result, v.Expect, v.Name)
		}
	}
}

var paletteExportConfig = []struct {
	Palette Palette
	Expect  PaletteConfig
}{
	{
		Palette: Palette{
			effects: map[string]*Effector{
				"color1": {
					Enclose: true,
					effects: []EffectCode{Bold, Italic},
					fgColor: color8{color: Cyan},
					bgColor: nil,
				},
				"color2": {
					Enclose: true,
					effects: []EffectCode{},
					fgColor: nil,
					bgColor: color8{color: Blue},
				},
			},
		},
		Expect: PaletteConfig{
			Effectors: map[string]EffectorConfig{
				"color1": {
					Effects:    []string{"Bold", "Italic"},
					Foreground: "Cyan",
					Background: nil,
				},
				"color2": {
					Effects:    []string{},
					Foreground: nil,
					Background: "Blue",
				},
			},
		},
	},
}

func TestPalette_ExportConfig(t *testing.T) {
	for _, v := range paletteExportConfig {
		result := v.Palette.ExportConfig()
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %#v", result, v.Expect, v.Palette)
		}
	}
}

var generatePaletteTests = []struct {
	Config PaletteConfig
	Expect *Palette
	Error  string
}{
	{
		Config: PaletteConfig{
			Effectors: map[string]EffectorConfig{
				"color1": {
					Effects:    []string{"Bold", "Italic"},
					Foreground: "Cyan",
					Background: nil,
				},
				"color2": {
					Effects:    []string{},
					Foreground: nil,
					Background: "Blue",
				},
			},
		},
		Expect: &Palette{
			effects: map[string]*Effector{
				"color1": {
					Enclose:  true,
					effects:  []EffectCode{Bold, Italic},
					fgColor:  color8{color: Cyan},
					bgColor:  nil,
					sequence: "\033[36;1;3m",
				},
				"color2": {
					Enclose:  true,
					effects:  []EffectCode{},
					fgColor:  nil,
					bgColor:  color8{color: Blue},
					sequence: "\033[44m",
				},
			},
			useEffects: true,
		},
	},
	{
		Config: PaletteConfig{
			Effectors: map[string]EffectorConfig{
				"color1": {
					Effects:    []string{"Bold", "Italic"},
					Foreground: "Cyan",
					Background: nil,
				},
				"color2": {
					Effects:    []string{},
					Foreground: nil,
					Background: "invalid",
				},
			},
		},
		Error: "\"invalid\" cannot convert to Color Code",
	},
}

func TestGeneratePalette(t *testing.T) {
	for _, v := range generatePaletteTests {
		result, err := GeneratePalette(v.Config)
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
