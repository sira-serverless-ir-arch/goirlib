package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/model"
)

type Memory struct {
	Index           map[string]map[string]*model.Set
	FieldDocument   map[string]map[string]model.Field
	FieldLength     map[string]int
	FieldSize       map[string]int
	NumberFieldTerm map[string]map[string]int
}

func (m *Memory) UpdateIndex(documentId string, field model.Field) {
	m.FieldSize[field.Name] += 1
	m.FieldLength[field.Name] += field.Length

	indexField := m.Index[field.Name]

	if indexField == nil {
		indexField = make(map[string]*model.Set)
		m.Index[field.Name] = indexField
	}

	for key := range field.TF {
		set := indexField[key]
		if set == nil {
			set = model.NewSet()
			set.Add(documentId)
			indexField[key] = set
		} else {
			set.Add(documentId)
		}
	}
}

func NewMemory() Storage {
	return &Memory{
		Index:           make(map[string]map[string]*model.Set),
		FieldDocument:   make(map[string]map[string]model.Field),
		FieldLength:     make(map[string]int),
		FieldSize:       make(map[string]int),
		NumberFieldTerm: make(map[string]map[string]int),
	}
}

func (m *Memory) GetNumberFieldTerm(fieldName string, terms []string) map[string]int {
	temp := make(map[string]int)
	for _, term := range terms {
		temp[term] = m.NumberFieldTerm[fieldName][term]
	}
	return temp
}

func (m *Memory) GetFieldLength(fieldName string) int {
	return m.FieldLength[fieldName]
}

func (m *Memory) GetFieldSize(fieldName string) int {
	return m.FieldSize[fieldName]
}

func (m *Memory) GetFieldDocumentTest(documentId string) map[string]model.Field {
	return m.FieldDocument[documentId]
}

func (m *Memory) GetIndex(fieldName string) map[string]*model.Set {
	return m.Index[fieldName]
}

func (m *Memory) GetDocuments(fieldName string, term string) (model.Set, bool) {
	if indexField, ok := m.Index[fieldName]; ok {
		if set, ok := indexField[term]; ok {
			return *set, true
		}
	}
	return model.Set{}, false
}

func (m *Memory) GetFields(documentId []string, fieldName string) map[string]model.Field {

	fields := make(map[string]model.Field)
	for _, id := range documentId {
		if fieldDocument, ok := m.FieldDocument[id]; ok {
			if field, ok := fieldDocument[fieldName]; ok {
				fields[id] = field
			}
		}
	}

	return fields

}

func (m *Memory) SaveOrUpdate(documentId string, field model.Field) {
	m.createFieldDocument(documentId, field)

	m.UpdateIndex(documentId, field)

	//numérico de campos que possui o termo
	fieldTerm := m.NumberFieldTerm[field.Name]
	if fieldTerm == nil {
		fieldTerm = make(map[string]int)
	}

	for term := range field.TF {
		fieldTerm[term] += 1
	}
	m.NumberFieldTerm[field.Name] = fieldTerm

}

func (m *Memory) createFieldDocument(documentId string, field model.Field) {
	if fieldIndex, ok := m.FieldDocument[documentId]; ok {
		fieldIndex[field.Name] = field
	} else {
		m.FieldDocument[documentId] = map[string]model.Field{field.Name: field}
	}
}
