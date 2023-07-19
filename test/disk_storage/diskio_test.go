package disk_storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"log"
	"os"
	"testing"
)

func TestLoadTermDocumentsFromIndex(t *testing.T) {

	if file.Exists("data7/") {
		err := os.RemoveAll("data7/")
		if err != nil {
			log.Fatal(err)
		}
	}

	store, err := storage.NewDiskStore("data7/", 5)
	indexDocs(documents, store)

	//reload from index file on storage_disk
	store, _ = storage.NewDiskStore("data7/", 10)
	index, err := storage.LoadTermDocumentsFromIndex("summary", "data7/")
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

	index, err = storage.LoadTermDocumentsFromIndex("Id", "data7/")
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

	index, err = storage.LoadTermDocumentsFromIndex("name", "data7/")
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

	if file.Exists("data8/") {
		err := os.RemoveAll("data8/")
		if err != nil {
			log.Fatal(err)
		}
	}

	store, err := storage.NewDiskStore("data8/", 5)
	if err != nil {
		panic(err)
	}
	indexDocs(documents, store)

	nDocuments := []string{
		"{\"Id\":\"4\",\"name\":\"Maria Joana\",\"summary\":\"house of Jack\"}",
		"{\"Id\":\"5\",\"name\":\"Isabella Taliba\",\"summary\":\"Iceland land of fire\"}",
	}

	store, err = storage.NewDiskStore("data8/", 5)
	if err != nil {
		panic(err)
	}
	indexDocs(nDocuments, store)

	//reload index file from storage_disk
	store, _ = storage.NewDiskStore("data8/", 10)
	index, _ := storage.LoadTermDocumentsFromIndex("name", "data8/")

	if index["isabella"].Size() != 2 {
		t.Errorf("Expected 2, got %d", index["isabella"].Size())
	}

	if index["taliba"].Size() != 2 {
		t.Errorf("Expected 2, got %d", index["taliba"].Size())
	}

	index, _ = storage.LoadTermDocumentsFromIndex("summary", "data8/")
	if index["iceland"].Size() != 1 {
		t.Errorf("Expected 1, got %d", index["iceland"].Size())
	}

	if index["hous"].Size() != 3 {
		t.Errorf("Expected 3, got %d", index["hous"].Size())
	}
}
