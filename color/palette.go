package color

type Palette struct {
	effects    map[string]*Effector
	useEffects bool
}

func NewPalette() *Palette {
	return &Palette{
		effects:    make(map[string]*Effector),
		useEffects: true,
	}
}

func (p *Palette) Enable() {
	p.useEffects = true
}

func (p *Palette) Disable() {
	p.useEffects = false
}

func (p *Palette) SetEffector(key string, effector *Effector) {
	p.effects[key] = effector
}

func (p *Palette) Render(key string, text string) string {
	if e, ok := p.effects[key]; ok && p.useEffects {
		return e.Render(text)
	}
	return text
}
