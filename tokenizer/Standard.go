package tokenizer

import (
	"strings"
	"unicode"
)

type Standard struct {
}

func NewStandard() Tokenizer {
	return &Standard{}
}

func (s *Standard) Tokenize(text string) []string {
	tokens := strings.FieldsFunc(text, isNotAlnum)
	return tokens
}

func isNotAlnum(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '\''
}
