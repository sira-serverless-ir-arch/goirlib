package filter

import (
	"github.com/sira-serverless-ir-arch/goirlib/filter/stemmer"
	"github.com/sira-serverless-ir-arch/goirlib/language"
	"github.com/sira-serverless-ir-arch/porter2"
	"sync"
)

type Stemmer struct {
	Language  language.Language
	Algorithm stemmer.StemmerAlgorithm
}

func NewStemmer(algorithm stemmer.StemmerAlgorithm) Filter {
	return &Stemmer{
		Algorithm: algorithm,
	}

}

func (s *Stemmer) Process(text []string) []string {

	if s.Algorithm != stemmer.Snowball {
		panic("Invalid Stemmer Algorithm")
	}

	textSize := len(text)
	temp := make([]string, textSize)

	var wg sync.WaitGroup
	wg.Add(textSize)
	for i, word := range text {
		go func(i int, word string) {
			defer wg.Done()
			temp[i] = porter2.Stem(word)
		}(i, word)
	}
	wg.Wait()
	return temp
}
