package text

type Encoder interface {
	SetHeader([]*Field)
	AppendRecord([]*Field)
	Encode() (string, error)
}

type Field struct {
	Contents  string
	Alignment FieldAlignment

	Lines []string
	Width int
}

func NewField(contents string, alignment FieldAlignment) *Field {
	return &Field{
		Contents:  contents,
		Alignment: alignment,
	}
}
