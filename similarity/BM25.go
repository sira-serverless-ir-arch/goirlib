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

//D é um documento.
//Q é a consulta de pesquisa que é um conjunto de termos {q1, q2,..., qn}.
//f(qi, D) é a frequência do termo qi no documento D (isto é, o TF).
//|D| é o comprimento do documento D em palavras.
//avgdl é o comprimento médio do documento na coleção de documentos.
//k1 e b são parâmetros livres, geralmente escolhidos, no contexto da Recuperação de Informações na Web, como k1 = 1.2 e b = 0.75.
//IDF(qi) é o IDF (Inverse Document Frequency) do termo qi. No caso do BM25, o
//IDF é calculado como log((Total Number Of Documents - Number Of Documents with term t in it + 0.5) / (Number Of Documents with term t in it + 0.5)).
