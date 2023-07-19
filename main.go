package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sira-serverless-ir-arch/goirlib/cache"
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

func maincc() {
	lru := cache.NewAsyncMap[float64]()

	go func() {
		for i := 0; i < 1000; i++ {
			lru.Put(fmt.Sprintf("teste %v", i), float64(i))
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			lru.Put(fmt.Sprintf("teste %v", i), float64(i))
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			lru.Put(fmt.Sprintf("teste %v", i), float64(i))
		}
	}()

	go func() {
		for {
			for i := 0; i < 1000; i++ {
				lru.Get(fmt.Sprintf("teste %v", i))
			}
		}
	}()

	go func() {
		for {
			for i := 0; i < 1000; i++ {
				lru.Get(fmt.Sprintf("teste %v", i))
			}
		}
	}()

	go func() {
		for {
			for i := 0; i < 1000; i++ {
				lru.Get(fmt.Sprintf("teste %v", i))
			}
		}
	}()

	go func() {
		for {
			for i := 0; i < 1000; i++ {
				lru.Get(fmt.Sprintf("teste %v", i))
			}
		}
	}()

	time.Sleep(20 * time.Second)

	fmt.Println()
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	store, err := storage.NewDiskStore("disk_storage/", 300000)
	if err != nil {
		panic(err)
	}

	inv := index.NewIndex(store)
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
		pId := Preprocessing(id)
		pTitle := Preprocessing(document.Title)
		pText := Preprocessing(document.Text)

		inv.IndexDocument(id, model.NormalizedDocument{
			Fields: []model.Field{
				{
					Name:   "id",
					Length: len(pId),
					TF:     metric.TermFrequency(pId),
				},
				{
					Name:   "title",
					Length: len(pTitle),
					TF:     metric.TermFrequency(pTitle),
				},
				{
					Name:   "text",
					Length: len(pText),
					TF:     metric.TermFrequency(pText),
				},
			},
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
	err = r.Run()
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

func mainxcrt() {

	store, err := storage.NewDiskStore("disk_storage/", 1000)
	if err != nil {
		panic(err)
	}
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

		var fields []model.Field
		for k, v := range flatted {
			r := Preprocessing(fmt.Sprintf("%s", v))
			fields = append(fields, model.Field{
				Name:   k,
				Length: len(r),
				TF:     metric.TermFrequency(r),
			})
		}
		invertedIndex.IndexDocument(id, model.NormalizedDocument{
			Fields: fields,
		})

		if doci%100 == 0 {
			fmt.Println(doci)
		}
	}

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
