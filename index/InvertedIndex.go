package index

import (
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"runtime"
	"sync"
)

type Index interface {
	Process(documentId string, field model.Field)
	Search(terms []string, fieldName string) model.SearchResult
	GetIndex(fieldName string) map[string]*storage.Set
	GetFieldDocument(documentId string) map[string]model.Field
}

type InvertedIndex struct {
	Storage storage.Storage
}

func NewIndex(storage storage.Storage) Index {
	return &InvertedIndex{
		Storage: storage,
	}
}

func (i *InvertedIndex) GetFieldDocument(documentId string) map[string]model.Field {
	return i.Storage.GetFieldDocumentTest(documentId)
}

func (i *InvertedIndex) GetIndex(fieldName string) map[string]*storage.Set {
	return i.Storage.GetIndex(fieldName)
}

func (i *InvertedIndex) Search(terms []string, fieldName string) model.SearchResult {

	maxRoutines := runtime.NumCPU() * 2

	sumDl := 0.0
	count := 0.0
	result := model.SearchResult{
		NumFieldsWithTerm: make(map[string]int),
	}

	documentIdCh := make(chan string, len(terms))
	semaphoreCh := make(chan struct{}, maxRoutines)

	var wg1 sync.WaitGroup
	for _, term := range terms {
		wg1.Add(1)
		go func(term, fieldName string) {
			defer wg1.Done()

			semaphoreCh <- struct{}{} // acquire semaphore
			if documents, ok := i.Storage.GetDocuments(fieldName, term); ok {
				for id := range documents.GetData() {
					documentIdCh <- id
				}
			}
			<-semaphoreCh
		}(term, fieldName)
	}

	go func() {
		wg1.Wait()
		close(documentIdCh)
	}()

	numBatches := 30
	batches := make([][]string, numBatches)
	hasDocument := make(map[string]bool)
	ixx := 0
	for id := range documentIdCh {
		if _, ok := hasDocument[id]; ok {
			continue
		}
		hasDocument[id] = true
		batchIndex := ixx % numBatches
		batches[batchIndex] = append(batches[batchIndex], id)
		ixx++
		count++
	}

	fieldCh := make(chan map[string]model.Field, len(terms))

	var wg2 sync.WaitGroup
	for _, documentsId := range batches {
		wg2.Add(1)
		go func(documentsId []string, fieldName string) {
			defer wg2.Done()

			semaphoreCh <- struct{}{} // acquire semaphore
			fields := i.Storage.GetFields(documentsId, fieldName)
			for documentId, field := range fields {
				fieldCh <- map[string]model.Field{
					documentId: field,
				}
			}
			<-semaphoreCh // release semaphore
		}(documentsId, fieldName)
	}

	go func() {
		wg2.Wait()
		close(fieldCh)
	}()

	fieldDocuments := make(map[string]model.Field)
	for fieldDocument := range fieldCh {
		for documentId, field := range fieldDocument {
			fieldDocuments[documentId] = field
			sumDl += float64(field.Length)
		}
	}

	result.TotalDocuments = count
	result.FieldDocuments = fieldDocuments
	result.AvgDocLength = sumDl / count

	return result
}

//func (i *InvertedIndex) Search(terms []string, fieldName string) model.SearchResult {
//
//	sumDl := 0.0
//	count := 0.0
//	result := model.SearchResult{
//		NumFieldsWithTerm: make(map[string]int),
//	}
//
//	documentIdCh := make(chan string, len(terms))
//
//	var wg1 sync.WaitGroup
//	wg1.Add(len(terms))
//	for _, term := range terms {
//
//		go func(term, fieldName string) {
//			defer wg1.Done()
//			if documents, ok := i.Storage.GetDocuments(fieldName, term); ok {
//				for id := range documents.GetData() {
//					documentIdCh <- id
//				}
//			}
//		}(term, fieldName)
//
//	}
//
//	go func() {
//		wg1.Wait()
//		close(documentIdCh)
//	}()
//
//	numBatches := 30
//	batches := make([][]string, numBatches)
//	hasDocument := make(map[string]bool)
//	ixx := 0
//	for id := range documentIdCh {
//		if _, ok := hasDocument[id]; ok {
//			continue
//		}
//		hasDocument[id] = true
//		batchIndex := ixx % numBatches
//		batches[batchIndex] = append(batches[batchIndex], id)
//		ixx++
//		count++
//	}
//
//	var wg2 sync.WaitGroup
//	wg2.Add(len(batches))
//	fieldCh := make(chan map[string]model.Field, len(terms))
//	for _, documentsId := range batches {
//		go func(documentsId []string, fieldName string) {
//			defer wg2.Done()
//			fields := i.Storage.GetFields(documentsId, fieldName)
//			for documentId, field := range fields {
//				fieldCh <- map[string]model.Field{
//					documentId: field,
//				}
//			}
//
//		}(documentsId, fieldName)
//	}
//
//	go func() {
//		wg2.Wait()
//		close(fieldCh)
//	}()
//
//	fieldDocuments := make(map[string]model.Field)
//	for fieldDocument := range fieldCh {
//		for documentId, field := range fieldDocument {
//			fieldDocuments[documentId] = field
//			sumDl += float64(field.Length)
//		}
//	}
//
//	result.TotalDocuments = count
//	result.FieldDocuments = fieldDocuments
//	result.AvgDocLength = sumDl / count
//
//	return result
//}

func (i *InvertedIndex) Process(documentId string, field model.Field) {
	i.Storage.SaveOrUpdate(documentId, field)
}
