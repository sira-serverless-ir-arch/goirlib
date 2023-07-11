package file

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
)

const (
	NumberFieldTerm  = "nfterm"
	MetricsFile      = "metrics"
	IndexFile        = "index"
	Documents        = "docs"
	DocumentsMetrics = "metrics"
	DocumentsRaw     = "raw"
)

func CompressData(buf []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(buf); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	return b.Bytes()
}

func DecompressData(buf []byte) []byte {
	r := bytes.NewReader(buf)
	gr, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}
	defer func(gr *gzip.Reader) {
		err := gr.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(gr)

	data, err := io.ReadAll(gr)
	if err != nil {
		panic(err)
	}

	return data
}

func SaveFileOnDisk(path string, buf []byte) error {
	err := os.WriteFile(path, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadFileOnDisk(path string) ([]byte, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func CreteDirIfNotExist(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
}

func ListDirectories(path string) []string {
	var directories []string

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != Documents {
			directories = append(directories, file.Name())
		}
	}

	return directories
}
