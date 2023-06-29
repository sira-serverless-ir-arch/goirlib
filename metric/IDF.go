package metric

import (
	"math"
)

func IdfBM25(totalDocs float64, numDocsWithTerm float64) float64 {
	numerator := totalDocs - numDocsWithTerm + 0.5
	denominator := numDocsWithTerm + 0.5
	return math.Log(1 + (numerator / denominator))
}
