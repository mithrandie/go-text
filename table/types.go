package table

type TableFormat int

const (
	PlainTable TableFormat = iota
	GFMTable
	OrgTable
)
