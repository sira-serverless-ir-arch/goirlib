package test

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
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"github.com/sira-serverless-ir-arch/goirlib/tokenizer"
	"testing"
)

func Preprocessing(text string) []string {
	r := tokenizer.NewStandard().Tokenize(text)
	r = filter.NewLowercase().Process(r)
	r = filter.NewStopWords(language.GetWords(language.English)).Process(r)
	r = filter.NewStemmer(stemmer.Snowball).Process(r)
	return filter.NewASCII().Process(r)
}

func TestIndex(t *testing.T) {
	invertedIndex := index.NewIndex(storage.NewMemory())

	documents := []string{
		"{\"Id\":\"1\",\"name\":\"taliba jose da silva\",\"summary\":\"house of house\"}",
		"{\"Id\":\"2\",\"name\":\"Tatiane Rodrigues\",\"summary\":\"house of dragon\"}",
		"{\"Id\":\"3\",\"name\":\"Isabella\",\"summary\":\"dragon of gorgonia\"}",
	}
	//testecolection.GetTextDocuments()
	for _, document := range documents {
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
			if k == field.ID {
				continue
			}
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

	indx := invertedIndex.GetIndex("summary")
	if indx["hous"] == nil {
		t.Errorf("Expected word hous")
	}

	if indx["hous"].GetData()["1"] == false {
		t.Errorf("Expected true for key %s", "1")
	}

	if indx["hous"].GetData()["2"] == false {
		t.Errorf("Expected true for key %s", "2")
	}

	if indx["hous"].GetData()["3"] == true {
		t.Errorf("Expected false for key %s", "3")
	}

	if indx["dragon"].GetData()["3"] == false {
		t.Errorf("Expected true for key %s", "3")
	}

	if indx["gorgonia"].GetData()["3"] == false {
		t.Errorf("Expected true for key %s", "3")
	}

	fieldx := invertedIndex.GetFieldDocument("1")

	if fieldx["name"].Length < 4 {
		t.Error("Expected Length 4")
	}

	if fieldx["summary"].Length < 2 {
		t.Error("Expected Length 2")
	}

	if fieldx["summary"].TF["hous"] != 2 {
		t.Error("Expected TF 2 for hous")
	}

	fieldx = invertedIndex.GetFieldDocument("2")

	if fieldx["name"].Length < 2 {
		t.Error("Expected Length 2")
	}

	if fieldx["summary"].Length < 2 {
		t.Error("Expected Length 2")
	}

	if fieldx["summary"].TF["dragon"] != 1 {
		t.Error("Expected TF 1 for dragon")
	}

	if fieldx["summary"].TF["hous"] != 1 {
		t.Error("Expected TF 1 for house")
	}

	if fieldx["name"].TF["tatian"] != 1 {
		t.Error("Expected TF 1 for tatian")
	}

	if fieldx["name"].TF["rodrigu"] != 1 {
		t.Error("Expected TF 1 for rodrigu")
	}
}
