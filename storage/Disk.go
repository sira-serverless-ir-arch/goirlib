package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage/disk"
)

type Disk struct {
	indexCh         chan disk.IndexTransferData
	fieldCh         chan disk.FieldSizeLengthTransferData
	Index           map[string]map[string]*model.Set
	FieldDocument   map[string]map[string]model.Field
	FieldLength     map[string]int
	FieldSize       map[string]int
	NumberFieldTerm map[string]map[string]int
}

func NewDisk(rootFolder string) Storage {

	indexCh := make(chan disk.IndexTransferData, 100)
	fieldCh := make(chan disk.FieldSizeLengthTransferData)

	d := &Disk{
		indexCh:         indexCh,
		fieldCh:         fieldCh,
		Index:           make(map[string]map[string]*model.Set),
		FieldDocument:   make(map[string]map[string]model.Field),
		FieldLength:     make(map[string]int),
		FieldSize:       make(map[string]int),
		NumberFieldTerm: make(map[string]map[string]int),
	}

	file.CreteDirIfNotExist(rootFolder)

	d.LoadIndexOnHD(rootFolder)
	d.LoadFieldSizeLengthOnHD(rootFolder)

	go disk.SaveFieldSizeLengthOnDisc(rootFolder, fieldCh)
	go disk.SaveIndexOnDisk(rootFolder, indexCh)

	return d
}

func (d Disk) GetDocuments(fieldName string, term string) (model.Set, bool) {
	//TODO implement me
	panic("implement me")
}

func (d *Disk) GetFields(documentId []string, fieldName string) map[string]model.Field {
	//TODO implement me
	panic("implement me")
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

	//numérico de campos que possui o termo

}

func (d Disk) GetFieldDocumentTest(documentId string) map[string]model.Field {
	//TODO implement me
	panic("implement me")
}

func (d Disk) GetIndex(fieldName string) map[string]*model.Set {
	//TODO implement me
	panic("implement me")
}

func (d Disk) GetFieldLength(fieldName string) int {
	//TODO implement me
	panic("implement me")
}

func (d Disk) GetFieldSize(fieldName string) int {
	//TODO implement me
	panic("implement me")
}

func (d Disk) GetNumberFieldTerm(fieldName string, terms []string) map[string]int {
	//TODO implement me
	panic("implement me")
}
