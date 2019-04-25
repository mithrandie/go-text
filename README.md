# go-text

Package text and the sub-packages are Go libraries to operate text.

[![Build Status](https://travis-ci.org/mithrandie/go-text.svg?branch=master)](https://travis-ci.org/mithrandie/go-text)
[![License: MIT](https://img.shields.io/badge/License-MIT-lightgrey.svg)](https://opensource.org/licenses/MIT)

## Supported Character Encodings
- UTF-8
- UTF-16
- Shift-JIS

## Sub Packages

### color
Supports ANSI escape sequences.

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/mithrandie/go-text/color"
)

const (
	BlueColor   = "blue"
	YellowColor = "yellow"
)

type Config struct {
	Palette              color.PaletteConfig `json:"palette"`
}

var jsonConfig = `
{
  "palette": {
    "effectors": {
      "color1": {
        "effects": [
          "Bold"
        ],
        "foreground": "Blue",
        "background": null
      },
      "color2": {
        "effects": [],
        "foreground": "Magenta",
        "background": null
      }
    }
  }
}
`

func main() {
	message := "message"
	
	// Use JSON Configuration 
	conf := &Config{} 
	if err := json.Unmarshal([]byte(jsonConfig), conf); err != nil {
		panic(err)
	}
	
	palette, err := color.GeneratePalette(conf.Palette)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(palette.Render("color1", message))

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
	defer func() {
		if err = fp.Close(); err != nil {
			panic(err.Error())
		}
	}()
	
	r, _ := csv.NewReader(fp, text.UTF8)
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
	defer func() {
		if err = wfp.Close(); err != nil {
			panic(err.Error())
		}
	}()
	
	w, err := csv.NewWriter(wfp, lineBreak, text.SJIS)
	if err != nil {
		panic(err.Error())
	}
	w.Delimiter = ','
	
	for _, record := range recordSet {
		r := make([]csv.Field, 0, len(record))
		for _, field := range record {
			r = append(r, csv.NewField(string(field), false))
		}
		if err := w.Write(r); err != nil {
			panic("write error")
		}
	}
	if err = w.Flush(); err != nil {
		panic(err)
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
	defer func() {
		if err = fp.Close(); err != nil {
			panic(err.Error())
		}
	}()
	
	r, _ := fixedlen.NewReader(fp, []int{5, 10, 45, 60}, text.UTF8)
	r.WithoutNull = true
	recordSet, err := r.ReadAll()
	if err != nil {
		panic("fixed-length read error")
	}
	
	lineBreak := r.DetectedLineBreak
	
	wfp, err := os.Create("example_new.txt")
	if err != nil {
		panic("file open error")
	}
	defer func() {
		if err = wfp.Close(); err != nil {
			panic(err.Error())
		}
	}()

	w, err := fixedlen.NewWriter(wfp, []int{5, 10, 45, 60}, lineBreak, text.SJIS)
	if err != nil {
		panic(err.Error)
	}
	
	for _, record := range recordSet {
		r := make([]fixedlen.Field, 0, len(record))
		for _, field := range record {
			r = append(r, fixedlen.NewField(string(field), text.NotAligned))
		}
		if err = w.Write(r); err != nil {
			panic(err.Error())
		}
	}
	if err = w.Flush(); err != nil {
		panic(err.Error())
	}
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

### ltsv
Supports reading and writing LTSV Format.

```go
package main

import (
	"os"
	
	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/ltsv"
)

func main() {
	fp, err := os.Open("example.txt")
	if err != nil {
		panic("file open error")
	}
	defer func() {
		if err = fp.Close(); err != nil {
			panic(err.Error())
		}
	}()
	
	r, _ := ltsv.NewReader(fp, text.UTF8)
	r.WithoutNull = true
	recordSet, err := r.ReadAll()
	if err != nil {
		panic("ltsv read error")
	}
	
	header := r.Header.Fields()
	lineBreak := r.DetectedLineBreak
	
	wfp, err := os.Create("example_new.ltsv")
	if err != nil {
		panic("file open error")
	}
	defer func() {
		if err = wfp.Close(); err != nil {
			panic(err.Error())
		}
	}()

	w, err := ltsv.NewWriter(wfp, header, lineBreak, text.UTF8)
	if err != nil {
		panic("ltsv writer generation error")
	}
	
	for _, record := range recordSet {
		r := make([]string, 0, len(record))
		for _, field := range record {
			r = append(r, string(field))
		}
		if err = w.Write(r); err != nil {
			panic(err.Error())
		}
	}
	if err = w.Flush(); err != nil {
		panic(err.Error())
	}
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