package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/model"
)

type Memory struct {
	Index         map[string]map[string]*Set
	FieldDocument map[string]map[string]model.Field
}

func (m *Memory) GetFieldDocumentTest(documentId string) map[string]model.Field {
	return m.FieldDocument[documentId]
}

func (m *Memory) GetIndex(fieldName string) map[string]*Set {
	return m.Index[fieldName]
}

func NewMemory() Storage {
	return &Memory{
		Index:         make(map[string]map[string]*Set),
		FieldDocument: make(map[string]map[string]model.Field),
	}
}

func (m *Memory) GetDocuments(fieldName string, term string) (Set, bool) {
	if indexField, ok := m.Index[fieldName]; ok {
		if set, ok := indexField[term]; ok {
			return *set, true
		}
	}
	return Set{}, false
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

	indexField := m.Index[field.Name]

	if indexField == nil {
		indexField = make(map[string]*Set)
		m.Index[field.Name] = indexField
	}

	for key := range field.TF {
		set := indexField[key]
		if set == nil {
			set = NewSet()
			set.Add(documentId)
			indexField[key] = set
		} else {
			set.Add(documentId)
		}
	}
}

func (m *Memory) createFieldDocument(documentId string, field model.Field) {
	if fieldIndex, ok := m.FieldDocument[documentId]; ok {
		fieldIndex[field.Name] = field
	} else {
		m.FieldDocument[documentId] = map[string]model.Field{field.Name: field}
	}
}

//func (m *Memory) createFieldDocument(documentId string, field model.Field) {
//	fieldIndex := m.FieldDocument[documentId]
//	if fieldIndex == nil {
//		m.FieldDocument[documentId] = map[string]model.Field{
//			field.Name: field,
//		}
//	} else {
//		fieldIndex[field.Name] = field
//		m.FieldDocument[documentId] = fieldIndex
//	}
//
//}
