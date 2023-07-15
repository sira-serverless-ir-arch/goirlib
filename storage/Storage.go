package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/model"
)

type Storage interface {
	GetDocuments(fieldName string, term string) (model.Set, bool)
	//Passa os Ids de documentos e retorna os fields
	GetFields(documentId []string, fieldName string) map[string]model.Field
	SaveOrUpdate(documentId string, field model.Field)
	//GetFieldDocumentTest(documentId string) map[string]model.Field
	GetIndex(fieldName string) map[string]*model.Set
	//Soma de todos os TF para um campo e especifico
	GetFieldLength(fieldName string) int
	//Total de documentos indexados para o campo
	GetFieldSize(fieldName string) int
	//Retorna quandos documentos possui um determinando atributo
	GetNumberFieldTerm(fieldName string, terms []string) map[string]int
}
