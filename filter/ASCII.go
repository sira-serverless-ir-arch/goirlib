package filter

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

type ASCII struct {
}

func NewASCII() Filter {
	return &ASCII{}
}

func (A *ASCII) Process(text []string) []string {
	temp := make([]string, len(text))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for i, word := range text {
		result, _, _ := transform.String(t, word)
		temp[i] = result
	}

	return temp
}
