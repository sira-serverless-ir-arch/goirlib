package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sira-serverless-ir-arch/goirlib/field"
	"github.com/sira-serverless-ir-arch/goirlib/filter"
	"github.com/sira-serverless-ir-arch/goirlib/filter/stemmer"
	"github.com/sira-serverless-ir-arch/goirlib/index"
	"github.com/sira-serverless-ir-arch/goirlib/language"
	"github.com/sira-serverless-ir-arch/goirlib/metric"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/search"
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"github.com/sira-serverless-ir-arch/goirlib/testecolection"
	"github.com/sira-serverless-ir-arch/goirlib/tokenizer"
	"time"
)

func Preprocessing(text string) []string {
	r := tokenizer.NewStandard().Tokenize(text)
	r = filter.NewLowercase().Process(r)
	r = filter.NewStopWords(language.GetWords(language.English)).Process(r)
	r = filter.NewStemmer(stemmer.Snowball).Process(r)
	return filter.NewASCII().Process(r)
}

func main() {

	invertedIndex := index.NewIndex(storage.NewMemory())

	for _, document := range testecolection.GetTextDocuments() {
		obs := field.StringToObject(document)
		flatted := field.Flatten(obs)

		//Cria um ID se não existir 1
		id := field.GetID(flatted)
		if id == "" {
			id = uuid.New().String()
		}

		var fields []model.Field

		//processa cada field idividualmente
		for k, v := range flatted {
			r := Preprocessing(fmt.Sprintf("%s", v))
			f := model.Field{
				Name:   k,
				Length: len(r),
				TF:     metric.TermFrequency(r),
			}
			fields = append(fields, f)
			invertedIndex.Process(id, f)
		}
	}

	start := time.Now()
	s := search.NewStandard(invertedIndex) //search.NewAsyncStandard(invertedIndex, 10)
	results := s.Search([]model.Query{
		//{
		//	FieldName: "Id",
		//	Terms:     Preprocessing("6fc61b08-8cf8-4583-b3ac-64d132990637"),
		//},
		//{
		//	FieldName: "Title",
		//	Terms:     Preprocessing("galaxies from Swift UV"),
		//},
		{
			FieldName: "Summary",
			Terms:     Preprocessing("galaxies from Swift UV"),
			Boost:     1.2,
		},
		//{
		//	FieldName: "Summary",
		//	Terms:     Preprocessing("Deep Belief Nets for Topic Modeling"),
		//},
		//{
		//	FieldName: "Summary",
		//	Terms:     Preprocessing("We give a brief expositio"),
		//},
	})

	end := time.Since(start)
	fmt.Println("time", end)

	fmt.Println(results[0:5])

}
