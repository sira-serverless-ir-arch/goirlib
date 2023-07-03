package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/model"
)

type Storage interface {
	GetDocuments(fieldName string, term string) (Set, bool)
	GetFields(documentId []string, fieldName string) map[string]model.Field
	SaveOrUpdate(documentId string, field model.Field)
	GetFieldDocumentTest(documentId string) map[string]model.Field
	GetIndex(fieldName string) map[string]*Set
	//Soma de todos os TF para um campo e especifico
	GetFieldLength(fieldName string) int
	//Total de documentos indexados para o campo
	GetFieldSize(fieldName string) int
}
