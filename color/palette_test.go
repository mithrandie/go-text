package color

import (
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
