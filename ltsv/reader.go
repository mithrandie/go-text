package ltsv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mithrandie/go-text"
)

type Header struct {
	list []string
	keys map[string]bool
}

func NewHeader() *Header {
	return &Header{
		list: make([]string, 0, 16),
		keys: make(map[string]bool, 16),
	}
}

func (h *Header) Exists(key string) bool {
	_, ok := h.keys[key]
	return ok
}

func (h *Header) Add(key string) {
	if _, ok := h.keys[key]; !ok {
		h.keys[key] = true
		h.list = append(h.list, key)
	}
}

func (h *Header) Len() int {
	return len(h.list)
}

func (h *Header) Fields() []string {
	return h.list
}

type Record map[string]*bytes.Buffer

func (r Record) Write(key string, value []byte) {
	if _, ok := r[key]; !ok {
		r[key] = new(bytes.Buffer)
	}
	r[key].Reset()
	r[key].Write(value)
}

func (r Record) Clear() {
	for k := range r {
		r[k].Reset()
	}
}

type Reader struct {
	WithoutNull bool

	reader *bufio.Reader
	line   int
	column int

	keyBuf   *bytes.Buffer
	valueBuf *bytes.Buffer
	record   Record

	Header            *Header
	DetectedLineBreak text.LineBreak
}

func NewReader(r io.Reader, enc text.Encoding) *Reader {
	return &Reader{
		WithoutNull: false,
		reader:      bufio.NewReader(text.GetTransformDecoder(r, enc)),
		line:        1,
		column:      0,
		keyBuf:      new(bytes.Buffer),
		valueBuf:    new(bytes.Buffer),
		record:      make(Record),
		Header:      NewHeader(),
	}
}

func (r *Reader) newError(s string) error {
	return errors.New(fmt.Sprintf("line %d, column %d: %s", r.line, r.column, s))
}

func (r *Reader) Read() ([]text.RawText, error) {
	r.record.Clear()

	fieldNum := 0
	for {
		eol, err := r.parseField()
		if err != nil {
			if err == io.EOF {
				if fieldNum < 1 {
					return nil, io.EOF
				}
			} else {
				return nil, err
			}
		}

		if eol && fieldNum < 1 {
			continue
		}

		key := r.keyBuf.String()
		if !r.Header.Exists(key) {
			r.Header.Add(key)
		}
		r.record.Write(key, r.valueBuf.Bytes())

		fieldNum++

		if eol {
			break
		}
	}

	values := make([]text.RawText, 0, r.Header.Len())
	for _, key := range r.Header.Fields() {
		b, ok := r.record[key]
		if !ok || b.Len() < 1 {
			if r.WithoutNull {
				values = append(values, text.RawText{})
			} else {
				values = append(values, nil)
			}
		} else {
			v := make([]byte, b.Len())
			copy(v, b.Bytes())
			values = append(values, v)
		}
	}

	return values, nil
}

func (r *Reader) ReadAll() ([][]text.RawText, error) {
	records := make([][]text.RawText, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	for i := range records {
		for j := len(records[i]); j < r.Header.Len(); j++ {
			if r.WithoutNull {
				records[i] = append(records[i], text.RawText{})
			} else {
				records[i] = append(records[i], nil)
			}
		}
	}

	return records, nil
}

func (r *Reader) parseField() (eol bool, err error) {
	r.keyBuf.Reset()
	r.valueBuf.Reset()

	var lineBreak text.LineBreak
	readingKey := true

ParseFieldLoop:
	for {
		lineBreak = ""

		ch, _, e := r.reader.ReadRune()
		r.column++

		if e != nil {
			if e == io.EOF {
				eol = true
			}
			err = e
			break ParseFieldLoop
		}

		switch ch {
		case '\r':
			nextCh, _, _ := r.reader.ReadRune()
			if nextCh == '\n' {
				lineBreak = text.CRLF
			} else {
				r.reader.UnreadRune()
				lineBreak = text.CR
			}
			ch = '\n'
		case '\n':
			lineBreak = text.LF
		}

		switch ch {
		case '\n':
			if r.DetectedLineBreak == "" {
				r.DetectedLineBreak = lineBreak
			}
			r.line++
			r.column = 0
			eol = true
			fallthrough
		case '\t':
			break ParseFieldLoop
		case ':':
			readingKey = false
		default:
			if readingKey {
				r.keyBuf.WriteRune(ch)
			} else {
				r.valueBuf.WriteRune(ch)
			}
		}
	}

	if 0 < r.keyBuf.Len() && readingKey {
		err = r.newError("missing field separator")
	}

	return
}
