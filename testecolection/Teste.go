package testecolection

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title   string `xml:"title"`
	Summary string `xml:"summary"`
	Author  Author `xml:"author"`
	Link    []Link `xml:"link"`
}

type Author struct {
	Name string `xml:"name"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
}

func GetArvixData() string {
	url := "https://export.arxiv.org/api/query?search_query=all:deep&start=0&max_results=10000"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação HTTP:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

type Document struct {
	Id      string
	Title   string
	Summary string
	Author  string
}

func CriarArquivo() {

	xmlData := GetArvixData()

	var feed Feed
	err := xml.Unmarshal([]byte(xmlData), &feed)
	if err != nil {
		fmt.Println("Erro ao analisar o XML:", err)
		return
	}

	fmt.Println("Título do Feed:", feed.Title)
	fmt.Println()

	arquivo, err := os.Create("testecolection/arxiv.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer func(arquivo *os.File) {
		err := arquivo.Close()
		if err != nil {

		}
	}(arquivo)

	for _, entry := range feed.Entries {
		GravarLinhasNoArquivo(arquivo, Document{
			Id:      uuid.New().String(),
			Title:   entry.Title,
			Summary: entry.Summary,
			Author:  entry.Author.Name,
		})
	}
}

func GravarLinhasNoArquivo(arquivo *os.File, document Document) {

	escritor := bufio.NewWriter(arquivo)

	marshal, _ := json.Marshal(document)

	_, err := escritor.WriteString(string(marshal) + "\n")
	if err != nil {
		fmt.Println(err)
	}

	err = escritor.Flush()
	if err != nil {
		fmt.Println(err)
	}

}

func GetTextDocuments() []string {
	// Abrir o arquivo para leitura
	arquivo, err := os.Open("testecolection/arxiv.txt")
	if err != nil {
		panic(err)
	}
	defer func(arquivo *os.File) {
		err := arquivo.Close()
		if err != nil {

		}
	}(arquivo)

	// Criar um leitor bufio para ler o arquivo linha a linha
	leitor := bufio.NewReader(arquivo)

	var documents []string
	for {
		linha, err := leitor.ReadString('\n')
		if err != nil {
			// Verificar se é o fim do arquivo
			if err.Error() == "EOF" {
				break
			}
			panic(err)
		}

		documents = append(documents, linha)
	}

	return documents
}

func GetDocuments() []Document {
	// Abrir o arquivo para leitura
	arquivo, err := os.Open("testecolection/arxiv.txt")
	if err != nil {
		panic(err)
	}
	defer func(arquivo *os.File) {
		err := arquivo.Close()
		if err != nil {

		}
	}(arquivo)

	// Criar um leitor bufio para ler o arquivo linha a linha
	leitor := bufio.NewReader(arquivo)

	var documents []Document
	for {
		linha, err := leitor.ReadString('\n')
		if err != nil {
			// Verificar se é o fim do arquivo
			if err.Error() == "EOF" {
				break
			}
			panic(err)
		}
		document := Document{}
		err = json.Unmarshal([]byte(linha), &document)
		if err != nil {
			panic(err)
		}

		documents = append(documents, document)
	}

	return documents
}
