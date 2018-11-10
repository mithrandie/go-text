package text

type TableFormat int

const (
	PlainTable TableFormat = iota
	GFMTable
	OrgTable
)
