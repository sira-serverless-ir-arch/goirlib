package storage

//
//import (
//	"fmt"
//	"github.com/sira-serverless-ir-arch/goirlib/cache"
//	"github.com/sira-serverless-ir-arch/goirlib/model"
//	"log"
//	"time"
//)
//
//type IndexDisk struct {
//	diskIO  DiskIO
//	indexCh chan IndexTransferData
//	index   *cache.Shard[*model.Set]
//}
//
//func NewIndexDisk(diskIO DiskIO, indexCh chan IndexTransferData) (*IndexDisk, error) {
//
//	startTime := time.Now()
//	index, err := diskIO.LoadIndexOnHD()
//	if err != nil {
//		return nil, err
//	}
//	fmt.Printf("The index was loaded to memory in %v\n", time.Now().Sub(startTime))
//
//	i := &IndexDisk{
//		diskIO:  diskIO,
//		indexCh: indexCh,
//		index:   nil,
//	}
//
//	go diskIO.SaveIndexOnDisk(indexCh)
//
//	go func() {
//		//merge index fragments to index
//		for {
//			time.Sleep(15 * time.Second)
//			_, err := diskIO.LoadIndexOnHD()
//			if err != nil {
//				log.Fatalf("Error merge index: %v", err)
//			}
//			fmt.Println("Merge documents")
//		}
//	}()
//
//	return i, nil
//}
//
//func (i *IndexDisk) UpdateIndex(documentId string, field model.Field) {
//
//	indexField, _ := i.index.Get(field.Name)
//
//	for term := range field.TF {
//		setPtr, ok := indexField.Get(term) //indexField[term]
//		if !ok {
//			set := model.NewSet()
//			set.Add(documentId)
//			indexField.Set(term, set) //indexField[term] = set
//			i.indexCh <- IndexTransferData{
//				FieldName:  field.Name,
//				Term:       term,
//				DocumentId: documentId,
//			}
//		} else {
//			set := *setPtr
//			set.Add(documentId)
//			i.indexCh <- IndexTransferData{
//				FieldName:  field.Name,
//				Term:       term,
//				DocumentId: documentId,
//			}
//		}
//	}
//}
//
//func (i *IndexDisk) GetIndex(fieldName string) map[string]*model.Set {
//	temp := make(map[string]*model.Set)
//
//	if indexField, ok := i.index.Get(fieldName); ok {
//
//		for key := range indexField.GetData() {
//			if set, ok := indexField.Get(key); ok {
//				temp[key] = *set
//			}
//		}
//	}
//
//	return temp
//}
//
//func (i *IndexDisk) GetDocuments(fieldName string, term string) (model.Set, bool) {
//	if indexField, ok := i.index.Get(fieldName); ok {
//		if set, ok := indexField.Get(term); ok {
//			return **set, true
//		}
//	}
//	return model.Set{}, false
//}
