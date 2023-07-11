package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sira-serverless-ir-arch/goirlib/field"
	"github.com/sira-serverless-ir-arch/goirlib/filter"
	"github.com/sira-serverless-ir-arch/goirlib/filter/stemmer"
	"github.com/sira-serverless-ir-arch/goirlib/index"
	"github.com/sira-serverless-ir-arch/goirlib/language"
	"github.com/sira-serverless-ir-arch/goirlib/metric"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/search"
	"github.com/sira-serverless-ir-arch/goirlib/storage"
	"github.com/sira-serverless-ir-arch/goirlib/testecolection"
	"github.com/sira-serverless-ir-arch/goirlib/tokenizer"
	"net/http"
	"time"
)

type Document struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Result struct {
	Total          int           `json:"total"`
	Duration       string        `json:"duration"`
	Algorithm      string        `json:"algorithm"`
	SemanticSearch bool          `json:"semanticSearch"`
	QueryResults   []QueryResult `json:"queryResults"`
}

type QueryResult struct {
	Similarity float64  `json:"similarity"`
	Document   Document `json:"document"`
}

func isNotDocumentComplete(doc Document) bool {
	if doc.Id == "" {
		return true
	}
	return false
}

var id = 0

func mainxx() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	inv := index.NewIndex(storage.NewMemory())
	sr := search.NewStandard(inv)

	r.POST("/nir", func(c *gin.Context) {

		var document Document
		if err := c.ShouldBindJSON(&document); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if isNotDocumentComplete(document) {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid document": 400})
			return
		}

		id := document.Id
		r := Preprocessing(id)
		inv.IndexDocument(id, model.Field{
			Name:   "id",
			Length: len(r),
			TF:     metric.TermFrequency(r),
		})

		r = Preprocessing(document.Title)
		inv.IndexDocument(id, model.Field{
			Name:   "title",
			Length: len(r),
			TF:     metric.TermFrequency(r),
		})

		r = Preprocessing(document.Text)
		inv.IndexDocument(id, model.Field{
			Name:   "text",
			Length: len(r),
			TF:     metric.TermFrequency(r),
		})

		c.JSON(http.StatusCreated, gin.H{"success": 201})

	})
	r.GET("/nir", func(c *gin.Context) {

		paramPairs := c.Request.URL.Query()

		start := time.Now()
		terms := Preprocessing(paramPairs.Get("query"))
		rs := sr.Search([]model.Query{
			{
				FieldName: "title",
				Terms:     terms,
				Boost:     1,
			},
			{
				FieldName: "text",
				Terms:     terms,
				Boost:     1,
			},
		})

		duration := time.Since(start)
		id += 1
		//fmt.Println("Qid:", id, len(rs))
		if len(rs) == 0 {
			body := Result{
				Duration:       duration.String(),
				Total:          len(rs),
				Algorithm:      "B25F",
				SemanticSearch: false,
			}
			c.JSON(http.StatusOK, body)
			return
		}

		queryResults := make([]QueryResult, len(rs))
		for i, score := range rs {
			queryResults[i] = QueryResult{
				Similarity: score.Score,
				Document: Document{
					Id: score.Id,
				},
			}
		}

		body := Result{
			Duration:       duration.String(),
			Total:          10,
			Algorithm:      "BM25",
			SemanticSearch: false,
			QueryResults:   queryResults[0:10],
		}
		c.JSON(http.StatusOK, body)
	})
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func Preprocessing(text string) []string {
	r := tokenizer.NewStandard().Tokenize(text)
	r = filter.NewLowercase().Process(r)
	r = filter.NewStopWords(language.GetWords(language.English)).Process(r)
	r = filter.NewStemmer(stemmer.Snowball).Process(r)
	return filter.NewASCII().Process(r)
}

func main() {

	store := storage.NewDisk("data/", 1000)
	invertedIndex := index.NewIndex(store)

	doci := 0
	for _, document := range testecolection.GetTextDocuments() {
		doci += 1
		obs := field.StringToObject(document)
		flatted := field.Flatten(obs)

		//Cria um ID se não existir 1
		id := field.GetID(flatted)
		if id == "" {
			id = uuid.New().String()
		}

		//processa cada field idividualmente
		for k, v := range flatted {
			r := Preprocessing(fmt.Sprintf("%s", v))
			f := model.Field{
				Name:   k,
				Length: len(r),
				TF:     metric.TermFrequency(r),
			}
			invertedIndex.IndexDocument(id, f)
		}
	}

	fmt.Println("Chegou aqui?")

	for {
		time.Sleep(5 * time.Second)
	}

	//start := time.Now()
	//s := search.NewStandard(invertedIndex) //search.NewAsyncStandard(invertedIndex, 10)
	//results := s.Search([]model.Query{
	//	{
	//		FieldName: "Title",
	//		Terms:     Preprocessing("galaxies from Swift UV"),
	//		Boost:     1,
	//	},
	//	{
	//		FieldName: "Summary",
	//		Terms:     Preprocessing("galaxies from Swift UV"),
	//		Boost:     1,
	//	},
	//})
	//
	//end := time.Since(start)
	//fmt.Println("time", end)
	//
	//fmt.Println(results[0:5])

}
