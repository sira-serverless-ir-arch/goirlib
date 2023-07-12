package test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sira-serverless-ir-arch/goirlib/field"
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/filter"
	"github.com/sira-serverless-ir-arch/goirlib/filter/stemmer"
	"github.com/sira-serverless-ir-arch/goirlib/language"
	"github.com/sira-serverless-ir-arch/goirlib/metric"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"github.com/sira-serverless-ir-arch/goirlib/tokenizer"
	"log"
	"os"
	"testing"
	"time"
)

func Preprocessing(text string) []string {
	r := tokenizer.NewStandard().Tokenize(text)
	r = filter.NewLowercase().Process(r)
	r = filter.NewStopWords(language.GetWords(language.English)).Process(r)
	r = filter.NewStemmer(stemmer.Snowball).Process(r)
	return filter.NewASCII().Process(r)
}

func indexDocs(store storage.Storage) {
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

		//Armazena os documentos
		for k, v := range flatted {
			r := Preprocessing(fmt.Sprintf("%s", v))
			f := model.Field{
				Name:   k,
				Length: len(r),
				TF:     metric.TermFrequency(r),
			}
			fields = append(fields, f)
			store.SaveOrUpdate(id, f)
		}
	}

	time.Sleep(7 * time.Second)
}

var store storage.Storage

func getTestStore() storage.Storage {
	if store != nil {
		return store
	}

	if file.Exists("data/") {
		err := os.RemoveAll("data/")
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error
	store, err = storage.NewDiskStore("data/", 20)
	if err != nil {
		panic(err)
	}
	indexDocs(store)
	return store
}

func TestStorage(t *testing.T) {

	indx := getTestStore().GetIndex("summary")

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

}

func TestNumberFieldsTerm(t *testing.T) {

	r := getTestStore().GetNumberFieldTerm("Id", []string{"1"})

	if r["1"] != 1 {
		t.Errorf("Expected 1")
	}
	if r["200"] != 0 {
		t.Errorf("Expected 0")
	}
	if r["200"] != 0 {
		t.Errorf("Expected 0")
	}

	r = store.GetNumberFieldTerm("summary", Preprocessing("dragon"))
	if r["dragon"] != 2 {
		t.Errorf("Expected 2")
	}

	r = store.GetNumberFieldTerm("summary", Preprocessing("house"))
	if r["hous"] != 2 {
		t.Errorf("Expected 2")
	}

	r = store.GetNumberFieldTerm("summary", Preprocessing("gorgonia"))
	if r["gorgonia"] != 1 {
		t.Errorf("Expected 1")
	}

	time.Sleep(10 * time.Second)

}

func TestGetFieldSize(t *testing.T) {

	r := getTestStore().GetFieldSize("Id")
	if r != 3 {
		t.Errorf("Expected 3")
	}
}

func TestGetDocuments(t *testing.T) {

	set, _ := getTestStore().GetDocuments("Id", "1")
	if !set.GetData()["1"] {
		t.Errorf("Expected true")
	}

	term := Preprocessing("house")[0]
	set, _ = getTestStore().GetDocuments("summary", term)

	if !set.GetData()["1"] {
		t.Errorf("Expected true")
	}

	if !set.GetData()["2"] {
		t.Errorf("Expected true")
	}

	term = Preprocessing("dragon")[0]
	set, _ = store.GetDocuments("summary", term)

	if !set.GetData()["2"] {
		t.Errorf("Expected true")
	}

	if !set.GetData()["3"] {
		t.Errorf("Expected true")
	}

	term = Preprocessing("gorgonia")[0]
	set, _ = store.GetDocuments("summary", term)

	if !set.GetData()["3"] {
		t.Errorf("Expected true")
	}

}

func TestGetFieldLength(t *testing.T) {

	r := getTestStore().GetFieldLength("summary")
	if r != 6 {
		t.Errorf("Expected 6")
	}

	r = getTestStore().GetFieldLength("Id")
	if r != 3 {
		t.Errorf("Expected 3")
	}
}

func TestGetFields(t *testing.T) {

	r := getTestStore().GetFields([]string{"1", "3"}, "name")
	for s, m := range r {
		fmt.Println(s, m)
	}

	if r["1"].Length != 4 {
		t.Errorf("Expected 4")
	}

	if r["3"].Length != 1 {
		t.Errorf("Expected 1")
	}

	r = store.GetFields([]string{"1"}, "summary")
	if r["1"].Length != 2 {
		t.Errorf("Expected 2")
	}

}
