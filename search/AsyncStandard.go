package search

//
//import (
//	"goir/index"
//	"goir/model"
//	"goir/similarity"
//	"sort"
//	"sync"
//)
//
//type AsyncStandard struct {
//	Index      index.Index
//	numWorkers int
//}
//
//func NewAsyncStandard(index index.Index, numWorkers int) Search {
//	return &AsyncStandard{
//		Index:      index,
//		numWorkers: numWorkers,
//	}
//}
//
//type Job struct {
//	Query           []string
//	K1              float64
//	B               float64
//	AvgDocLength    float64
//	TotalDocuments  float64
//	NumFieldsWithTerm *map[string]int
//	Document        model.NormalizedDocument
//}
//
//func (s *AsyncStandard) Search(query []string) []model.DocumentScore {
//
//	r := s.Index.Search(query)
//
//	jobs := make(chan Job, len(r.Documents))
//	results := make(chan model.DocumentScore, len(r.Documents))
//
//	wg := &sync.WaitGroup{}
//
//	// Start workers
//	for i := 0; i < s.numWorkers; i++ {
//		wg.Add(1)
//		go worker(jobs, results, wg)
//	}
//
//	// Send jobs to workers
//	for _, document := range r.Documents {
//		jobs <- Job{
//			Document:        document,
//			Query:           query,
//			K1:              1.2,
//			B:               0.75,
//			AvgDocLength:    r.AvgDocLength,
//			TotalDocuments:  r.TotalDocuments,
//			NumFieldsWithTerm: &r.NumFieldsWithTerm,
//		}
//	}
//	close(jobs)
//
//	wg.Wait()
//	close(results)
//
//	documentsScores := make([]model.DocumentScore, 0, len(r.Documents))
//	for result := range results {
//		documentsScores = append(documentsScores, result)
//	}
//
//	sort.Slice(documentsScores, func(i, j int) bool {
//		return documentsScores[i].Score > documentsScores[j].Score
//	})
//
//	return documentsScores
//}
//
//func worker(jobs <-chan Job, results chan<- model.DocumentScore, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	for job := range jobs {
//		score := similarity.BM25(job.Query, job.K1, job.B, job.AvgDocLength, job.TotalDocuments, *job.NumFieldsWithTerm, job.Document)
//		result := model.DocumentScore{
//			Id:    job.Document.Id,
//			Score: score,
//		}
//		results <- result
//	}
//}
