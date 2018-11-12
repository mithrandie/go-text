# go-text

Package text and the sub-packages are Go libraries to operate text.

[![Build Status](https://travis-ci.org/mithrandie/go-text.svg?branch=master)](https://travis-ci.org/mithrandie/go-text)
[![License: MIT](https://img.shields.io/badge/License-MIT-lightgrey.svg)](https://opensource.org/licenses/MIT)

## Supported Character Encodings
- UTF-8
- Shift-JIS

## Sub Packages

### color
Supports ANSI escape sequences.

```go
package main

import (
	"fmt"
	
	"github.com/mithrandie/go-text/color"
)

const (
	BlueColor   = "blue"
	YellowColor = "yellow"
)

func main() {
	message := "message"
	
	// Use Effector
	e := color.NewEffector()
	e.SetFGColor(color.Red)
	e.SetEffect(color.Bold, color.Italic)
	
	fmt.Println(e.Render(message))
	
	// Use Palette that bundles multiple effectors.
	blue := color.NewEffector()
	blue.SetFGColor(color.Blue)
	yellow := color.NewEffector()
	yellow.SetFGColor(color.Blue)
	
	palette := color.NewPalette()
	palette.SetEffector(BlueColor, blue)
	palette.SetEffector(YellowColor, yellow)
	
	fmt.Println(palette.Render(BlueColor, message))
}
```

### csv
Supports reading and writing CSV Format.


```go
package main

import (
	"os"
	
	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/csv"
)

func main() {
	fp, err := os.Open("example.csv")
	if err != nil {
		panic("file open error")
	}
	defer fp.Close()
	
	r := csv.NewReader(fp, text.UTF8)
	r.Delimiter = ','
	r.WithoutNull = true
	recordSet, err := r.ReadAll()
	if err != nil {
		panic("csv read error")
	}
	
	lineBreak := r.DetectedLineBreak
	
	wfp, err := os.Create("example_new.csv")
	if err != nil {
		panic("file open error")
	}
	defer wfp.Close()
	e := csv.NewWriter(wfp, lineBreak, text.SJIS)
	e.Delimiter = ','
	
	for _, record := range recordSet {
		r := make([]csv.Field, 0, len(record))
		for _, field := range record {
			r = append(r, csv.NewField(string(field), false))
		}
		if err := e.Write(r); err != nil {
			panic("write error")
		}
	}
}
```

### fixedlen
Supports reading and writing Fixed-Length Format.

```go
package main

import (
	"os"
	
	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/fixedlen"
)

func main() {
	fp, err := os.Open("example.txt")
	if err != nil {
		panic("file open error")
	}
	defer fp.Close()
	
	r := fixedlen.NewReader(fp, []int{5, 10, 45, 60}, text.UTF8)
	r.WithoutNull = true
	recordSet, err := r.ReadAll()
	if err != nil {
		panic("fixed-length read error")
	}
	
	lineBreak := r.DetectedLineBreak
	
	e := fixedlen.NewEncoder(len(recordSet))
	e.DelimiterPositions = []int{5, 10, 45, 60}
	e.LineBreak = lineBreak
	e.WithoutHeader = true
	e.Encoding = text.SJIS
	
	for _, record := range recordSet {
		r := make([]fixedlen.Field, 0, len(record))
		for _, field := range record {
			r = append(r, fixedlen.NewField(string(field), text.NotAligned))
		}
		e.AppendRecord(r)
	}
	csv, _ := e.Encode()
	
	wfp, err := os.Create("example_new.txt")
	if err != nil {
		panic("file open error")
	}
	defer wfp.Close()
	
	wfp.WriteString(csv)	
}
```

### json
Supports reading and writing JSON Format.

```go
package main

import (
	"fmt"
	"io/ioutil"
	
	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"
)

func main() {
	data, err := ioutil.ReadFile("example.json")
	if err != nil {
		panic("file open error")
	}
	
	d := json.NewDecoder()
	structure, escapeType, err := d.Decode(string(data))
	if err != nil {
		panic("json decode error")
	}
	
	e := json.NewEncoder()
	e.EscapeType = escapeType
	e.LineBreak = text.LF
	e.PrettyPrint = true
	e.Palette = json.NewJsonPalette()
	
	encoded := e.Encode(structure)
	fmt.Println(encoded)
}
```

### table
Supports writing text tables.

```go
package main

import (
	"fmt"
	
	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/table"
)

func main() {
	header := []table.Field{
		table.NewField("c1", text.Centering),
		table.NewField("c2", text.Centering),
		table.NewField("c3", text.Centering),
	}
	
	recordSet := [][]table.Field{
		{
			table.NewField("1", text.RightAligned),
			table.NewField("abc", text.LeftAligned),
			table.NewField("true", text.NotAligned),
		},
		{
			table.NewField("2", text.RightAligned),
			table.NewField("def", text.LeftAligned),
			table.NewField("true", text.NotAligned),
		},
		{
			table.NewField("3", text.RightAligned),
			table.NewField("ghi", text.LeftAligned),
			table.NewField("true", text.NotAligned),
		},
	}
	
	alignments := []text.FieldAlignment{
		text.RightAligned,
		text.LeftAligned,
		text.NotAligned,
	}
	
	e := table.NewEncoder(table.GFMTable, len(recordSet))
	e.LineBreak = text.LF
	e.EastAsianEncoding = true
	e.CountDiacriticalSign = false
	e.WithoutHeader = false
    
	
	e.SetHeader(header)
	for _, record := range recordSet {
		e.AppendRecord(record)
	}
	e.SetFieldAlignments(alignments)
	
	encoded, _ := e.Encode()
	fmt.Println(encoded)
}
```