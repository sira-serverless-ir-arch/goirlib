package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/cache"
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage/disk"
)

type Disk struct {
	indexCh            chan disk.IndexTransferData
	fieldCh            chan disk.FieldSizeLengthTransferData
	numberFieldTermCh  chan disk.NumberFieldTermTransferData
	documentFieldCh    chan disk.DocumentFieldTransferData
	rootFolder         string
	Index              map[string]map[string]*model.Set
	FieldDocument      map[string]map[string]model.Field
	FieldLength        map[string]int
	FieldSize          map[string]int
	NumberFieldTerm    map[string]map[string]int
	FieldDocumentCache *cache.LRUCache[map[string]model.Field]
}

func NewDisk(rootFolder string, cacheSize int) Storage {

	indexCh := make(chan disk.IndexTransferData, 100)
	fieldCh := make(chan disk.FieldSizeLengthTransferData)
	numberFieldTermCh := make(chan disk.NumberFieldTermTransferData)
	documentFieldCh := make(chan disk.DocumentFieldTransferData)

	d := &Disk{
		indexCh:            indexCh,
		fieldCh:            fieldCh,
		numberFieldTermCh:  numberFieldTermCh,
		documentFieldCh:    documentFieldCh,
		rootFolder:         rootFolder,
		Index:              make(map[string]map[string]*model.Set),
		FieldDocument:      make(map[string]map[string]model.Field),
		FieldLength:        make(map[string]int),
		FieldSize:          make(map[string]int),
		NumberFieldTerm:    make(map[string]map[string]int),
		FieldDocumentCache: cache.NewLRUCache[map[string]model.Field](cacheSize),
	}

	file.CreteDirIfNotExist(rootFolder)

	d.LoadIndexOnHD(rootFolder)
	d.LoadFieldSizeLengthOnHD(rootFolder)
	d.LoadNumberFieldTermOnHD(rootFolder)

	go disk.SaveFieldSizeLengthOnDisc(rootFolder, fieldCh)
	go disk.SaveIndexOnDisk(rootFolder, indexCh)
	go disk.SaveNumberFieldTermOnDisk(rootFolder, numberFieldTermCh)
	go disk.SaveDocumentFieldOnDisk(rootFolder, documentFieldCh)

	return d
}

func (d *Disk) GetNumberFieldTerm(fieldName string, terms []string) map[string]int {
	temp := make(map[string]int)
	for _, term := range terms {
		temp[term] = d.NumberFieldTerm[fieldName][term]
	}
	return temp
}

func (d *Disk) GetFieldLength(fieldName string) int {
	return d.FieldLength[fieldName]
}

func (d *Disk) GetFieldSize(fieldName string) int {
	return d.FieldSize[fieldName]
}

func (d *Disk) GetFieldDocumentTest(documentId string) map[string]model.Field {
	return d.FieldDocument[documentId]
}

func (d *Disk) GetIndex(fieldName string) map[string]*model.Set {
	return d.Index[fieldName]
}

func (d *Disk) GetDocuments(fieldName string, term string) (model.Set, bool) {
	if indexField, ok := d.Index[fieldName]; ok {
		if set, ok := indexField[term]; ok {
			return *set, true
		}
	}
	return model.Set{}, false
}

func (d *Disk) GetFields(documentId []string, fieldName string) map[string]model.Field {

	fields := make(map[string]model.Field)

	for _, id := range documentId {
		if fieldDocumentPtr, ok := d.FieldDocumentCache.Get(id); ok {
			fieldDocument := *fieldDocumentPtr
			if field, ok := fieldDocument[fieldName]; ok {
				fields[id] = field
			}
		} else {
			if fieldDocument, ok := d.LoadFieldDocumentOnHD(d.rootFolder, id); ok {
				d.FieldDocumentCache.Put(id, fieldDocument)
				if field, ok := fieldDocument[fieldName]; ok {
					fields[id] = field
				}
			}
		}
	}

	return fields

}

func (d *Disk) UpdateFieldSizeLength(field model.Field) {
	d.FieldSize[field.Name] += 1
	d.FieldLength[field.Name] += field.Length

	d.fieldCh <- disk.FieldSizeLengthTransferData{
		FieldName: field.Name,
		Length:    d.FieldLength[field.Name],
		Size:      d.FieldSize[field.Name],
	}
}

func (d *Disk) UpdateNumberFieldTerm(field model.Field) {
	fieldTerm := d.NumberFieldTerm[field.Name]
	if fieldTerm == nil {
		fieldTerm = make(map[string]int)
	}

	for term := range field.TF {
		fieldTerm[term] += 1
	}

	d.NumberFieldTerm[field.Name] = fieldTerm
	d.numberFieldTermCh <- disk.NumberFieldTermTransferData{
		FieldName: field.Name,
		TermSize:  fieldTerm,
	}
}

func (d *Disk) UpdateFieldDocument(documentId string, field model.Field) {
	if _, ok := d.FieldDocument[documentId]; !ok {
		d.FieldDocument[documentId] = make(map[string]model.Field)
	}

	d.FieldDocument[documentId][field.Name] = field
	d.FieldDocumentCache.Put(documentId, d.FieldDocument[documentId])

	d.documentFieldCh <- disk.DocumentFieldTransferData{
		DocumentId: documentId,
		Field:      d.FieldDocument[documentId],
	}
}

func (d *Disk) UpdateIndex(documentId string, field model.Field) {
	indexField := d.Index[field.Name]
	if indexField == nil {
		indexField = make(map[string]*model.Set)
		d.Index[field.Name] = indexField
	}

	for term := range field.TF {
		set := indexField[term]
		if set == nil {
			set = model.NewSet()
			set.Add(documentId)
			indexField[term] = set
		} else {
			set.Add(documentId)
		}
	}

	d.indexCh <- disk.IndexTransferData{
		FieldName:  field.Name,
		IndexField: indexField,
	}
}

func (d *Disk) SaveOrUpdate(documentId string, field model.Field) {

	d.UpdateFieldSizeLength(field)
	d.UpdateIndex(documentId, field)
	d.UpdateNumberFieldTerm(field)
	d.UpdateFieldDocument(documentId, field)
}
