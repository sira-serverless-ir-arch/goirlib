package disk

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage/buffers"
	"log"
	"sync"
	"time"
)

type IndexTransferData struct {
	FieldName  string
	IndexField map[string]*model.Set
}

var (
	m = sync.Mutex{}
)

func SaveIndexOnDisk(rootFolder string, indexCh chan IndexTransferData) {

	bufferData := make(map[string]map[string]*model.Set)

	go func() {
		for index := range indexCh {
			m.Lock()
			bufferData[index.FieldName] = index.IndexField
			m.Unlock()
		}
	}()

	semaphore := make(chan struct{}, 5) // Limit to 5 goroutines
	go func() {
		for {
			time.Sleep(5 * time.Second)
			m.Lock()
			for fieldName, data := range bufferData {
				semaphore <- struct{}{}
				go func(fieldName string, terms map[string]*model.Set) {
					defer func() { <-semaphore }()

					dir := rootFolder + fieldName
					CreteDirIfNotExist(dir)

					tempData := make(map[string]map[string]bool)

					for term, documents := range terms {
						tempData[term] = documents.GetData()
					}

					buff := SerializeIndex(tempData)
					err := SaveFileOnDisk(dir, "index", buff)
					if err != nil {
						log.Fatalf(err.Error())
					}
				}(fieldName, data)
				delete(bufferData, fieldName)
			}
			m.Unlock()
		}
	}()
}

func SerializeIndex(data map[string]map[string]bool) []byte {
	b := flatbuffers.NewBuilder(0)

	var terms []flatbuffers.UOffsetT

	for term, documents := range data {
		keyOffset := b.CreateString(term)

		var valsOffsets []flatbuffers.UOffsetT
		for documentId := range documents {
			valOffset := b.CreateString(documentId)
			valsOffsets = append(valsOffsets, valOffset)
		}

		buffers.DocumentStartValuesVector(b, len(valsOffsets))
		for i := len(valsOffsets) - 1; i >= 0; i-- {
			b.PrependUOffsetT(valsOffsets[i])
		}
		valuesVector := b.EndVector(len(valsOffsets))

		buffers.DocumentStart(b)
		buffers.DocumentAddValues(b, valuesVector)
		document := buffers.DocumentEnd(b)

		buffers.TermStart(b)
		buffers.TermAddKey(b, keyOffset)
		buffers.TermAddValues(b, document)
		term := buffers.TermEnd(b)

		terms = append(terms, term)
	}

	buffers.IndexStartEntriesVector(b, len(terms))
	for i := len(terms) - 1; i >= 0; i-- {
		b.PrependUOffsetT(terms[i])
	}
	entriesVector := b.EndVector(len(terms))

	buffers.IndexStart(b)
	buffers.IndexAddEntries(b, entriesVector)
	index := buffers.IndexEnd(b)

	b.Finish(index)

	return b.FinishedBytes()
}

func DeserializeIndex(buf []byte) map[string]map[string]bool {
	index := buffers.GetRootAsIndex(buf, 0)

	data := make(map[string]map[string]bool)

	termsLen := index.EntriesLength()

	for i := 0; i < termsLen; i++ {
		term := new(buffers.Term)
		if !index.Entries(term, i) {
			log.Fatalf("Failed to get Term")
		}

		key := string(term.Key())

		document := new(buffers.Document)
		term.Values(document)

		valuesLen := document.ValuesLength()
		values := make(map[string]bool)

		for j := 0; j < valuesLen; j++ {
			value := string(document.Values(j))
			values[value] = true
		}

		data[key] = values
	}

	return data
}
