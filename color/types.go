package color

type EffectCode int

const (
	Reset EffectCode = iota
	Bold
	Faint
	Italic
	Underline
	SlowBlink
	RapidBlink
	ReverseVideo
	Conceal
	CrossedOut
)

type Code int

const (
	Black Code = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const DefaultColor Code = 39

const (
	BrightBlack Code = iota + 90
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)
