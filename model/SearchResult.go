package model

type SearchResult struct {
	AvgDocLength      float64
	FieldDocuments    map[string]Field
	TotalDocuments    float64
	NumFieldsWithTerm map[string]int
}
