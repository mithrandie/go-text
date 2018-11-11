package color

import (
	"strconv"
	"strings"
)

const (
	color256FGPrefix = "38;5;"
	color256BGPrefix = "48;5;"
	colorRgbFGPrefix = "38;2;"
	colorRgbBGPrefix = "48;2;"
)

const resetSequence = "\033[0m"

type color interface {
	FG() string
	BG() string
}

type color8 struct {
	color Code
}

func (c color8) FG() string {
	return strconv.Itoa(int(c.color))
}

func (c color8) BG() string {
	return strconv.Itoa(int(c.color) + 10)
}

type color256 struct {
	color int
}

func (c color256) FG() string {
	return color256FGPrefix + strconv.Itoa(c.color)
}

func (c color256) BG() string {
	return color256BGPrefix + strconv.Itoa(c.color)
}

type colorRGB struct {
	red   int
	green int
	blue  int
}

func (c colorRGB) FG() string {
	return colorRgbFGPrefix + strconv.Itoa(c.red) + ";" + strconv.Itoa(c.green) + ";" + strconv.Itoa(c.blue)
}

func (c colorRGB) BG() string {
	return colorRgbBGPrefix + strconv.Itoa(c.red) + ";" + strconv.Itoa(c.green) + ";" + strconv.Itoa(c.blue)
}

type Effector struct {
	Enclose bool

	effects []EffectCode
	fgColor color
	bgColor color

	sequence string
}

func NewEffector() *Effector {
	return &Effector{
		Enclose: true,
	}
}

func (e *Effector) SetEffect(code ...EffectCode) {
	e.effects = code
	e.GenerateSequence()
}

func (e *Effector) SetFGColor(color Code) {
	e.fgColor = color8{
		color: color,
	}
	e.GenerateSequence()
}

func (e *Effector) SetBGColor(color Code) {
	e.bgColor = color8{
		color: color,
	}
	e.GenerateSequence()
}

func (e *Effector) SetFG256Color(color int) {
	e.fgColor = color256{
		color: color,
	}
	e.GenerateSequence()
}

func (e *Effector) SetBG256Color(color int) {
	e.bgColor = color256{
		color: color,
	}
	e.GenerateSequence()
}

func (e *Effector) SetFGRGBColor(red int, green int, blue int) {
	e.fgColor = colorRGB{
		red:   red,
		green: green,
		blue:  blue,
	}
	e.GenerateSequence()
}

func (e *Effector) SetBGRGBColor(red int, green int, blue int) {
	e.bgColor = colorRGB{
		red:   red,
		green: green,
		blue:  blue,
	}
	e.GenerateSequence()
}

func (e *Effector) GenerateSequence() {
	params := make([]string, 0, len(e.effects))

	if e.bgColor != nil {
		params = append(params, e.bgColor.BG())
	}

	if e.fgColor != nil {
		params = append(params, e.fgColor.FG())
	}

	if len(params) < 1 && len(e.effects) < 1 {
		e.effects = []EffectCode{Reset}
	}

	for _, v := range e.effects {
		params = append(params, strconv.Itoa(int(v)))
	}

	e.sequence = "\033[" + strings.Join(params, ";") + "m"
}

func (e *Effector) Render(s string) string {
	if !UseEffect || len(e.sequence) < 1 {
		return s
	}

	if e.Enclose {
		return e.sequence + s + resetSequence
	}
	return e.sequence + s
}
