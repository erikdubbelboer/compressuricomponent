package compressuricomponent

import (
	"net/url"
	"strconv"
	"strings"
)

type Encoder struct {
	*strings.Replacer
}

func NewEncoder(dict []string) Encoder {
	a := make([]string, len(dict)*2)

	for i, _ := range dict {
		a = append(a, dict[i])
		a = append(a, "~"+strconv.Itoa(i)+"~")
	}

	return Encoder{strings.NewReplacer(a...)}
}

func (e Encoder) Encode(s string) string {
	s = url.QueryEscape(s)
	s = e.Replace(s)
	return s
}

type Decoder struct {
	*strings.Replacer
}

func NewDecoder(dict []string) Decoder {
	a := make([]string, len(dict)*2)

	for i, _ := range dict {
		a = append(a, "~"+strconv.Itoa(i)+"~")
		a = append(a, dict[i])
	}

	return Decoder{strings.NewReplacer(a...)}
}

type appendSliceWriter []byte

// Write writes to the buffer to satisfy io.Writer.
func (w *appendSliceWriter) Write(p []byte) (int, error) {
	*w = append(*w, p...)
	return len(p), nil
}

// WriteString writes to the buffer without string->[]byte->string allocations.
func (w *appendSliceWriter) WriteString(s string) (int, error) {
	*w = append(*w, s...)
	return len(s), nil
}

func (e Decoder) Decode(s string) (string, error) {
	w := make(appendSliceWriter, 0, len(s)*10)
	e.WriteString(&w, s)
	return url.QueryUnescape(string(w))
}
