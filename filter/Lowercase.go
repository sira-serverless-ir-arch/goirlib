package filter

import "strings"

type Lowercase struct {
}

func NewLowercase() Filter {
	return Lowercase{}
}

func (l Lowercase) Process(text []string) []string {
	for i, t := range text {
		text[i] = strings.ToLower(t)
	}
	return text
}
