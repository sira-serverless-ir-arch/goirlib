package filter

type StopWords struct {
	Words map[string]bool
}

func NewStopWords(words map[string]bool) Filter {
	return &StopWords{
		Words: words,
	}
}

func (s *StopWords) Process(text []string) []string {

	var temp []string

	for _, element := range text {
		if s.Words[element] {
			continue
		}
		temp = append(temp, element)
	}

	return temp
}
