package color

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	color256FGPrefix = "38;5;"
	color256BGPrefix = "48;5;"
	colorRgbFGPrefix = "38;2;"
	colorRgbBGPrefix = "48;2;"
)

type EffectorConfig struct {
	Effects    []string    `json:"effects"`
	Foreground interface{} `json:"foreground"`
	Background interface{} `json:"background"`
}

const resetSequence = "\033[0m"

type color interface {
	FG() string
	BG() string
	ConfigValue() interface{}
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

func (c color8) ConfigValue() interface{} {
	return c.color.String()
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

func (c color256) ConfigValue() interface{} {
	return c.color
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

func (c colorRGB) ConfigValue() interface{} {
	return []int{c.red, c.green, c.blue}
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
	e.generateSequence()
}

func (e *Effector) SetFGColor(color Code) {
	e.fgColor = color8{
		color: color,
	}
	e.generateSequence()
}

func (e *Effector) SetBGColor(color Code) {
	e.bgColor = color8{
		color: color,
	}
	e.generateSequence()
}

func (e *Effector) SetFG256Color(color int) {
	e.fgColor = color256{
		color: color,
	}
	e.generateSequence()
}

func (e *Effector) SetBG256Color(color int) {
	e.bgColor = color256{
		color: color,
	}
	e.generateSequence()
}

func (e *Effector) SetFGRGBColor(red int, green int, blue int) {
	e.fgColor = colorRGB{
		red:   red,
		green: green,
		blue:  blue,
	}
	e.generateSequence()
}

func (e *Effector) SetBGRGBColor(red int, green int, blue int) {
	e.bgColor = colorRGB{
		red:   red,
		green: green,
		blue:  blue,
	}
	e.generateSequence()
}

func (e *Effector) generateSequence() {
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

func (e *Effector) ExportConfig() EffectorConfig {
	effects := make([]string, 0, len(e.effects))
	for _, v := range e.effects {
		effects = append(effects, v.String())
	}

	var fg interface{} = nil
	if e.fgColor != nil {
		fg = e.fgColor.ConfigValue()
	}

	var bg interface{} = nil
	if e.bgColor != nil {
		bg = e.bgColor.ConfigValue()
	}

	return EffectorConfig{
		Effects:    effects,
		Foreground: fg,
		Background: bg,
	}
}

func GenerateEffector(config EffectorConfig) (*Effector, error) {
	e := NewEffector()

	e.effects = make([]EffectCode, 0, len(config.Effects))
	for _, v := range config.Effects {
		c, err := ParseEffectCode(v)
		if err != nil {
			return nil, err
		}
		e.effects = append(e.effects, c)
	}

	fg, err := parseColor(config.Foreground)
	if err != nil {
		return nil, err
	}
	e.fgColor = fg

	bg, err := parseColor(config.Background)
	if err != nil {
		return nil, err
	}
	e.bgColor = bg

	e.generateSequence()
	return e, nil
}

func parseColor(i interface{}) (color, error) {
	if i == nil {
		return nil, nil
	}

	switch i.(type) {
	case string:
		c, err := ParseColorCode(i.(string))
		if err != nil {
			return nil, err
		}
		return color8{color: c}, nil
	case int:
		return color256{color: i.(int)}, nil
	case []int:
		rgb := i.([]int)
		if len(rgb) != 3 {
			return nil, errors.New(fmt.Sprintf("%v cannot convert to color", i))
		}
		return colorRGB{red: rgb[0], green: rgb[1], blue: rgb[2]}, nil
	}
	return nil, errors.New(fmt.Sprintf("%v cannot convert to color", i))
}
