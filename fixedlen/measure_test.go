package fixedlen

import (
	"reflect"
	"testing"

	"github.com/mithrandie/go-text"
)

var meagureMeasureTests = []struct {
	Name     string
	Records  [][]Field
	Encoding text.Encoding
	Expect   DelimiterPositions
}{
	{
		Name: "Measure",
		Records: [][]Field{
			{
				{Contents: "abcabc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		Encoding: text.UTF8,
		Expect:   DelimiterPositions{6, 9, 12},
	},
	{
		Name: "Uneven Records",
		Records: [][]Field{
			{
				{Contents: "abcabc", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
				{Contents: "pqr", Alignment: text.LeftAligned},
			},
		},
		Encoding: text.UTF8,
		Expect:   DelimiterPositions{6, 9, 12, 15},
	},
	{
		Name: "Multibyte Characters",
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "日本語", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		Encoding: text.UTF8,
		Expect:   DelimiterPositions{3, 12, 15},
	},
	{
		Name: "Multibyte Characters in Shift-JIS",
		Records: [][]Field{
			{
				{Contents: "abc", Alignment: text.LeftAligned},
				{Contents: "日本語", Alignment: text.LeftAligned},
				{Contents: "def", Alignment: text.LeftAligned},
			},
			{
				{Contents: "ghi", Alignment: text.LeftAligned},
				{Contents: "jkl", Alignment: text.LeftAligned},
				{Contents: "mno", Alignment: text.LeftAligned},
			},
		},
		Encoding: text.SJIS,
		Expect:   DelimiterPositions{3, 9, 12},
	},
}

func TestMeasure_Measure(t *testing.T) {
	for _, v := range meagureMeasureTests {
		m := NewMeasure()
		m.Encoding = v.Encoding

		for _, record := range v.Records {
			m.Measure(record)
		}

		positions := m.GeneratePositions()

		if !reflect.DeepEqual(positions, v.Expect) {
			t.Errorf("%s: positions = %v, want %v", v.Name, positions, v.Expect)
		}
	}
}
