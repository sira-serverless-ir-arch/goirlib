package metric

func TermFrequency(terms []string) map[string]int {
	frequencyMap := make(map[string]int)
	for _, term := range terms {
		frequencyMap[term]++
	}
	return frequencyMap
}
