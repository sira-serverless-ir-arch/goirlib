package test

import (
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"os"
	"testing"
)

func TestLoadTermDocumentsFromIndex(t *testing.T) {

	err := os.RemoveAll("data/")
	if err != nil {
		panic(err)
	}

	store, _ := storage.NewDiskStore("data/", 10)
	indexDocs(documents, store)
	//time.Sleep(50 * )

	index, err := storage.LoadTermDocumentsFromIndex("summary", "data/")
	if err != nil {
		panic(err)
	}

	if index["hous"].Size() != 2 {
		t.Errorf("Expected 2, got %d", index["hous"].Size())
	}

	documents := index["hous"].GetData()
	if !documents["1"] {
		t.Errorf("Expected true")
	}
	if !documents["2"] {
		t.Errorf("Expected true")
	}

	if index["dragon"].Size() != 2 {
		t.Errorf("Expected 2, got %d", index["dragon"].Size())
	}

	if index["gorgonia"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["gorgonia"].Size())
	}

	documents = index["gorgonia"].GetData()
	if !documents["3"] {
		t.Errorf("Expected true")
	}

	index, err = storage.LoadTermDocumentsFromIndex("Id", "data/")
	if err != nil {
		panic(err)
	}

	if index["1"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["1"].Size())
	}

	if index["2"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["1"].Size())
	}

	if index["2"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["1"].Size())
	}

	index, err = storage.LoadTermDocumentsFromIndex("name", "data/")
	if err != nil {
		panic(err)
	}

	if index["isabella"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["isabella"].Size())
	}

	documents = index["isabella"].GetData()
	if !documents["3"] {
		t.Errorf("Expected true")
	}

}

func TestLoadIndexOnHD(t *testing.T) {

	err := os.RemoveAll("data/")
	if err != nil {
		panic(err)
	}

	store, _ = storage.NewDiskStore("data/", 10)
	indexDocs(documents, store)

	nDocuments := []string{
		"{\"Id\":\"4\",\"name\":\"Maria Joana\",\"summary\":\"house of Jack\"}",
		"{\"Id\":\"5\",\"name\":\"Isabella Taliba\",\"summary\":\"Iceland land of fire\"}",
	}

	store, _ = storage.NewDiskStore("data/", 10)
	indexDocs(nDocuments, store)
	index, _ := storage.LoadTermDocumentsFromIndex("name", "data/")

	if index["isabella"].Size() != 2 {
		t.Errorf("Expected 2, got %d", index["isabella"].Size())
	}

	if index["taliba"].Size() != 2 {
		t.Errorf("Expected 2, got %d", index["taliba"].Size())
	}

	index, _ = storage.LoadTermDocumentsFromIndex("summary", "data/")
	if index["iceland"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["iceland"].Size())
	}

	if index["hous"].Size() != 3 {
		t.Errorf("Expected 3, got %d", index["hous"].Size())
	}
}
