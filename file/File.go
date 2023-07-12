package file

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

func RenameFile(currentPath, newPath string) error {
	err := os.Rename(currentPath, newPath)
	if err != nil {
		return fmt.Errorf("not possible to rename file: %v", err)
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Delete(path string) error {
	err := os.Remove(path)
	return err
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

func SecureSaveFile(path, fileName string, buff []byte) error {

	currentPath := filepath.Join(path, fileName)
	tempPath := filepath.Join(path, fileName+"_temp")
	oldPath := filepath.Join(path, fileName+"_old")

	err := SaveFileOnDisk(tempPath, buff)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	if Exists(currentPath) {
		err = RenameFile(currentPath, oldPath)
		if err != nil {
			return fmt.Errorf("failed to rename current file to old file: %w", err)
		}
	}

	err = RenameFile(tempPath, currentPath)
	if err != nil {
		return fmt.Errorf("failed to rename temp file to current file: %w", err)
	}

	if Exists(oldPath) {
		err = Delete(oldPath)
		if err != nil {
			return fmt.Errorf("failed to delete old file: %w", err)
		}
	}

	return nil
}
