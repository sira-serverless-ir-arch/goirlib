package search

import "github.com/sira-serverless-ir-arch/goirlib/model"

type Search interface {
	Search(query []model.Query) []model.DocumentScore
}
