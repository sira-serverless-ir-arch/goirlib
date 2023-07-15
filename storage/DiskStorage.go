package storage

import (
	"fmt"
	"github.com/sira-serverless-ir-arch/goirlib/cache"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"sync"
	"time"
)

type Disk struct {
	indexCh            chan IndexTransferData
	fieldCh            chan FieldSizeLengthTransferData
	numberFieldTermCh  chan NumberFieldTermTransferData
	documentFieldCh    chan DocumentFieldTransferData
	index              map[string]map[string]*model.Set
	fieldLength        map[string]int
	fieldSize          map[string]int
	numberFieldTerm    map[string]map[string]int
	fieldDocumentCache *cache.LRUCache[map[string]model.Field]
	diskIO             DiskIO
	mu                 sync.RWMutex
}

func NewDiskStore(rootFolder string, cacheSize int) (Storage, error) {

	indexCh := make(chan IndexTransferData, 100)
	fieldCh := make(chan FieldSizeLengthTransferData)
	numberFieldTermCh := make(chan NumberFieldTermTransferData)
	documentFieldCh := make(chan DocumentFieldTransferData)

	diskIO := NewDiskIO(rootFolder)

	startTime := time.Now()
	index, err := diskIO.LoadIndexOnHD()
	if err != nil {
		return nil, err
	}
	fmt.Printf("The index was loaded to memory in %v\n", time.Now().Sub(startTime))

	startTime = time.Now()
	fieldSize, fieldLength, err := diskIO.LoadFieldSizeLengthOnHD()
	if err != nil {
		return nil, err
	}
	fmt.Printf("The metrics was loaded to memory in %v\n", time.Now().Sub(startTime))

	startTime = time.Now()
	numberFieldTerm, err := diskIO.LoadNumberFieldTermOnHD()
	if err != nil {
		return nil, err
	}
	fmt.Printf("The nfterm was loaded to memory in %v\n", time.Now().Sub(startTime))

	d := &Disk{
		diskIO:             diskIO,
		indexCh:            indexCh,
		fieldCh:            fieldCh,
		numberFieldTermCh:  numberFieldTermCh,
		documentFieldCh:    documentFieldCh,
		index:              index,
		fieldLength:        fieldLength,
		fieldSize:          fieldSize,
		numberFieldTerm:    numberFieldTerm,
		fieldDocumentCache: cache.NewLRUCache[map[string]model.Field](cacheSize),
		mu:                 sync.RWMutex{},
	}

	go diskIO.SaveFieldSizeLengthOnDisc(fieldCh)
	go diskIO.SaveIndexOnDisk(indexCh)
	go diskIO.SaveNumberFieldTermOnDisk(numberFieldTermCh)
	go diskIO.SaveDocumentFieldOnDisk(documentFieldCh)

	go func() {
		//merge index fragments to index
		for {
			time.Sleep(15 * time.Second)
			_, err := diskIO.LoadIndexOnHD()
			if err != nil {
				return
			}
			fmt.Println("Merge documents")
		}
	}()

	return d, nil
}

func (d *Disk) GetNumberFieldTerm(fieldName string, terms []string) map[string]int {
	temp := make(map[string]int)
	for _, term := range terms {
		temp[term] = d.numberFieldTerm[fieldName][term]
	}
	return temp
}

func (d *Disk) GetFieldLength(fieldName string) int {
	return d.fieldLength[fieldName]
}

func (d *Disk) GetFieldSize(fieldName string) int {
	return d.fieldSize[fieldName]
}

func (d *Disk) GetIndex(fieldName string) map[string]*model.Set {
	return d.index[fieldName]
}

func (d *Disk) GetDocuments(fieldName string, term string) (model.Set, bool) {
	if indexField, ok := d.index[fieldName]; ok {
		if set, ok := indexField[term]; ok {
			return *set, true
		}
	}
	return model.Set{}, false
}

func (d *Disk) GetFields(documentId []string, fieldName string) map[string]model.Field {

	fields := make(map[string]model.Field)
	for _, id := range documentId {
		if fieldDocumentPtr, ok := d.fieldDocumentCache.Get(id); ok {
			fieldDocument := *fieldDocumentPtr
			if field, ok := fieldDocument[fieldName]; ok {
				fields[id] = field
			}
		} else {
			if fieldDocument, ok := d.diskIO.LoadFieldDocumentOnHD(id); ok {
				if field, ok := fieldDocument[fieldName]; ok {
					fields[id] = field
				}
				d.fieldDocumentCache.Put(id, fieldDocument)
			}
		}
	}

	return fields

}

func (d *Disk) UpdateFieldSizeLength(field model.Field) {
	d.fieldSize[field.Name] += 1
	d.fieldLength[field.Name] += field.Length

	d.fieldCh <- FieldSizeLengthTransferData{
		FieldName: field.Name,
		Length:    d.fieldLength[field.Name],
		Size:      d.fieldSize[field.Name],
	}
}

func (d *Disk) UpdateNumberFieldTerm(field model.Field) {
	fieldTerm := d.numberFieldTerm[field.Name]
	if fieldTerm == nil {
		fieldTerm = make(map[string]int)
	}

	for term := range field.TF {
		fieldTerm[term] += 1
		d.numberFieldTermCh <- NumberFieldTermTransferData{
			FieldName: field.Name,
			Term:      term,
			Size:      fieldTerm[term],
		}
	}

	d.numberFieldTerm[field.Name] = fieldTerm

}

func (d *Disk) UpdateFieldDocument(documentId string, field model.Field) {

	var fieldIndex map[string]model.Field

	if fieldIndexPtr, ok := d.fieldDocumentCache.Get(documentId); ok {
		fieldIndex = *fieldIndexPtr
	} else {
		fieldIndex = make(map[string]model.Field)
	}

	fieldIndex[field.Name] = field
	d.fieldDocumentCache.Put(documentId, fieldIndex)

	d.documentFieldCh <- DocumentFieldTransferData{
		DocumentId: documentId,
		Field:      fieldIndex,
	}

}

func (d *Disk) UpdateIndex(documentId string, field model.Field) {

	indexField := d.index[field.Name]
	if indexField == nil {
		indexField = make(map[string]*model.Set)
		d.index[field.Name] = indexField
	}

	for term := range field.TF {
		set := indexField[term]
		if set == nil {
			set = model.NewSet()
			set.Add(documentId)
			indexField[term] = set
			d.indexCh <- IndexTransferData{
				FieldName:  field.Name,
				Term:       term,
				DocumentId: documentId,
			}
		} else {
			set.Add(documentId)
			d.indexCh <- IndexTransferData{
				FieldName:  field.Name,
				Term:       term,
				DocumentId: documentId,
			}
		}
	}

}

func (d *Disk) SaveOrUpdate(documentId string, field model.Field) {

	d.UpdateFieldSizeLength(field)
	d.UpdateIndex(documentId, field)
	d.UpdateNumberFieldTerm(field)
	d.UpdateFieldDocument(documentId, field)
}
