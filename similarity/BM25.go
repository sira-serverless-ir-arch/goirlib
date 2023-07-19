package similarity

import (
	"github.com/sira-serverless-ir-arch/goirlib/metric"
	"github.com/sira-serverless-ir-arch/goirlib/model"
)

func BM25(query []string, k1, b, boost, avgDocLength, totalDocs float64, numDocsWithTerm map[string]int, field model.Field) float64 {
	score := 0.0
	if boost == 0 {
		boost = 1
	}
	for _, term := range query {
		idf := metric.IdfBM25(totalDocs, float64(numDocsWithTerm[term]))
		frequency := float64(field.TF[term])
		numerator := frequency * (k1 + 1)
		denominator := frequency + k1*(1-b+b*(float64(field.Length)/avgDocLength))
		score += boost * idf * (numerator / denominator)

		//if doc_id == "485" {
		//	//fmt.Println("*******************************")
		//	//fmt.Println("doc_id", doc_id)
		//	fmt.Println("term", term)
		//	//fmt.Println("totalDocs", totalDocs)
		//	//fmt.Println("numDocsWithTerm", numDocsWithTerm[term])
		//	//fmt.Println("idf", idf)
		//	//fmt.Println("frequency", field.TF[term])
		//	fmt.Println("field.Length", field.Length)
		//	//fmt.Println("avgDocLength", avgDocLength)
		//	fmt.Println("bm25", idf*(numerator/denominator))
		//}

	}

	return score
}

//var idfCache = cache.NewShardMap[float64](10)
//var similarityCache = cache.NewShardMap[float64](10)
//
//func CalcSimilarity(k1, b, boost, avgDocLength, idf float64, term string, field model.Field) float64 {
//
//	iMap, ok := similarityCache.Get(field.Name)
//
//	if ok {
//		if similarityPtr, ok := iMap.Get(term); ok {
//			return *similarityPtr
//		} else {
//			frequency := float64(field.TF[term])
//			numerator := frequency * (k1 + 1)
//			denominator := frequency + k1*(1-b+b*(float64(field.Length)/avgDocLength))
//			similarityPtr := boost * idf * (numerator / denominator)
//			iMap.Put(term, similarityPtr)
//			return similarityPtr
//		}
//	}
//
//	panic("Panico na similaridade")
//
//}
//
//func CalcIDF(totalDocs, numDocsWithTerm float64, term, fieldName string) float64 {
//
//	iMap, ok := idfCache.Get(fieldName)
//
//	if ok {
//		if idfPtr, ok := iMap.Get(term); ok {
//			return *idfPtr
//		} else {
//			idf := metric.IdfBM25(totalDocs, numDocsWithTerm)
//			iMap.Put(term, idf)
//			return idf
//		}
//	}
//
//	panic("Panico no IDF")
//
//	//else {
//	//	idf := metric.IdfBM25(totalDocs, numDocsWithTerm)
//	//	iMap := cache.NewCacheRCU[float64]()
//	//	iMap.Put(term, idf)
//	//	idfCache.Put()
//	//	return idf
//	//}
//
//}
//
//func BM25(query []string, k1, b, boost, avgDocLength, totalDocs float64, numDocsWithTerm map[string]int, field model.Field) float64 {
//	score := 0.0
//	if boost == 0 {
//		boost = 1
//	}
//	for _, term := range query {
//		idf := CalcIDF(totalDocs, float64(numDocsWithTerm[term]), term, field.Name)
//		//metric.IdfBM25(totalDocs, float64(numDocsWithTerm[term]))
//
//		score += CalcSimilarity(k1, b, boost, avgDocLength, idf, term, field)
//
//		//if doc_id == "485" {
//		//	//fmt.Println("*******************************")
//		//	//fmt.Println("doc_id", doc_id)
//		//	fmt.Println("term", term)
//		//	//fmt.Println("totalDocs", totalDocs)
//		//	//fmt.Println("numDocsWithTerm", numDocsWithTerm[term])
//		//	//fmt.Println("idf", idf)
//		//	//fmt.Println("frequency", field.TF[term])
//		//	fmt.Println("field.Length", field.Length)
//		//	//fmt.Println("avgDocLength", avgDocLength)
//		//	fmt.Println("bm25", idf*(numerator/denominator))
//		//}
//
//	}
//
//	return score
//}
