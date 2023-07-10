package search

import (
	"github.com/sira-serverless-ir-arch/goirlib/index"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/similarity"
	"sort"
	"sync"
)

type Standard struct {
	Index index.Index
}

func NewStandard(index index.Index) Search {
	return &Standard{
		Index: index,
	}
}

func (s *Standard) Search(query []model.Query) []model.DocumentScore {

	scoresCh := make(chan []model.DocumentScore, len(query))
	var wg sync.WaitGroup
	wg.Add(len(query))

	for _, q := range query {
		go func(terms []string, fieldName string, boost float64) {
			defer wg.Done()
			r := s.Index.Search(terms, fieldName)
			documentsScores := make([]model.DocumentScore, len(r.FieldDocuments))
			i := 0
			for documentId, field := range r.FieldDocuments {
				documentsScores[i] = model.DocumentScore{
					Id:    documentId,
					Score: similarity.BM25(terms, 1.2, 0.75, boost, r.AvgDocLength, r.TotalDocuments, r.NumFieldsWithTerm, field),
				}

				i += 1
			}
			scoresCh <- documentsScores
		}(q.Terms, q.FieldName, q.Boost)
	}

	go func() {
		wg.Wait()
		close(scoresCh)
	}()

	tempScores := make(map[string]float64)
	for documentScores := range scoresCh {
		for _, document := range documentScores {
			tempScores[document.Id] += document.Score
		}
	}

	documentsScores := make([]model.DocumentScore, len(tempScores))
	i := -1
	for documentId, score := range tempScores {
		i++
		documentsScores[i] = model.DocumentScore{
			Id:    documentId,
			Score: score,
		}
	}

	sort.Slice(documentsScores, func(i, j int) bool {
		return documentsScores[i].Score > documentsScores[j].Score
	})

	return documentsScores
}

//func (s *Standard) Search(query []model.Query) []model.DocumentScore {
//
//	scores := make(map[string][]model.DocumentScore)
//
//	wg := sync.WaitGroup{}
//	wg.Add(len(query))
//	for _, q := range query {
//
//		go func(terms []string, fieldName string) {
//			defer wg.Done()
//			r := s.Index.Search(terms, fieldName)
//			documentsScores := make([]model.DocumentScore, len(r.FieldDocuments))
//			i := 0
//			for id, field := range r.FieldDocuments {
//				documentsScores[i] = model.DocumentScore{
//					Id:    id,
//					Score: similarity.BM25(terms, 1.2, 0.75, r.AvgDocLength, r.TotalDocuments, r.NumFieldsWithTerm, field),
//				}
//				i += 1
//			}
//			scores[fieldName] = documentsScores
//		}(q.Terms, q.FieldName)
//	}
//	wg.Wait()
//
//	tempScores := make(map[string]float64)
//	for _, documentScores := range scores {
//		for _, document := range documentScores {
//			tempScores[document.Id] += document.Score
//		}
//	}
//
//	documentsScores := make([]model.DocumentScore, len(tempScores))
//	i := -1
//	for documentId, score := range tempScores {
//		i++
//		documentsScores[i] = model.DocumentScore{
//			Id:    documentId,
//			Score: score,
//		}
//	}
//	sort.Slice(documentsScores, func(i, j int) bool {
//		return documentsScores[i].Score > documentsScores[j].Score
//	})
//	return documentsScores
//}
