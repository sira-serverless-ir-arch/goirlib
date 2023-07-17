package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"log"
	"path/filepath"
	"sync"
	"time"
)

type DiskIO struct {
	RootFolder string
}

func NewDiskIO(rootFolder string) DiskIO {
	file.CreteDirIfNotExist(rootFolder)
	return DiskIO{rootFolder}
}

func (d *DiskIO) LoadNumberFieldTermOnHD() (map[string]map[string]int, error) {
	fields := file.ListDirectories(d.RootFolder)
	numberFieldTerm := make(map[string]map[string]int)
	for _, field := range fields {
		path := filepath.Join(d.RootFolder, field, file.NumberFieldTerm)
		buf, err := file.ReadFileOnDisk(path)
		if err != nil {
			return nil, err
		}
		numberFieldTerm[field] = DeserializeNumberFieldTerm(file.DecompressData(buf))
	}
	return numberFieldTerm, nil
}

func (d *DiskIO) LoadFieldSizeLengthOnHD() (map[string]int, map[string]int, error) {
	fields := file.ListDirectories(d.RootFolder)
	fieldSize := make(map[string]int)
	fieldLength := make(map[string]int)
	for _, field := range fields {
		path := filepath.Join(d.RootFolder, field, file.MetricsFile)
		buf, err := file.ReadFileOnDisk(path)
		if err != nil {
			return nil, nil, err
		}

		name, size, length := DeserializeFieldSizeLength(file.DecompressData(buf))
		fieldSize[name] = int(size)
		fieldLength[name] = int(length)
	}
	return fieldSize, fieldLength, nil
}

func (d *DiskIO) LoadFieldDocumentOnHD(documentId string) (map[string]model.Field, bool) {
	path := filepath.Join(d.RootFolder, file.Documents, file.DocumentsMetrics, documentId)
	buf, err := file.ReadFileOnDisk(path)
	if err != nil {
		return nil, false
	}
	return DeserializeFieldMap(file.DecompressData(buf)), true
}

func LoadTermDocumentsFromIndex(field string, rootFolder string) (map[string]*model.Set, error) {
	termDocuments := make(map[string]*model.Set)
	path := filepath.Join(rootFolder, field, file.IndexFile)
	if !file.Exists(path) {
		return termDocuments, nil
	}

	buf, err := file.ReadFileOnDisk(path)
	if err != nil {
		return nil, err
	}
	for term, documents := range DeserializeIndex(file.DecompressData(buf)) {
		set := model.NewSet()
		for documentId := range documents {
			set.Add(documentId)
		}
		termDocuments[term] = set
	}
	return termDocuments, nil
}

func (d *DiskIO) LoadIndexOnHD() (map[string]map[string]*model.Set, error) {
	fields := file.ListDirectories(d.RootFolder)
	index := make(map[string]map[string]*model.Set)

	for _, field := range fields {
		path := filepath.Join(d.RootFolder, field, file.Temp)
		file.CreteDirIfNotExist(path)
		files := file.ListFiles(path)

		termDocuments, err := LoadTermDocumentsFromIndex(field, d.RootFolder)
		if err != nil {
			return nil, err
		}

		for _, fileName := range files {
			filePath := filepath.Join(path, fileName)
			buf, err := file.ReadFileOnDisk(filePath)
			if err != nil {
				return nil, err
			}

			for term, documents := range DeserializeIndex(file.DecompressData(buf)) {
				set, ok := termDocuments[term]
				if !ok {
					set = model.NewSet()
				}

				for documentId := range documents {
					set.Add(documentId)
				}
				termDocuments[term] = set
			}
		}
		index[field] = termDocuments

		tempTermDocuments := make(map[string]map[string]bool)
		for term, documents := range termDocuments {
			tempTermDocuments[term] = documents.GetData()
		}

		buffer := SerializeIndex(tempTermDocuments)
		indexPath := filepath.Join(d.RootFolder, field)
		file.CreteDirIfNotExist(indexPath)

		err = file.SecureSaveFile(indexPath, file.IndexFile, file.CompressData(buffer))
		if err != nil {
			return nil, err
		}

		//Delete temp files
		for _, fileName := range files {
			filePath := filepath.Join(path, fileName)
			err = file.Delete(filePath)
			if err != nil {
				return nil, err
			}
		}
	}
	return index, nil
}

//****** Fields *******/

type FieldSizeLengthTransferData struct {
	FieldName string
	Size      int
	Length    int
}

type NumberFieldTermTransferData struct {
	FieldName string
	Term      string
	Size      int
	//TermSize  map[string]int
}

type DocumentFieldTransferData struct {
	DocumentId string
	Field      map[string]model.Field
}

type BufferFieldDocument struct {
	sync.Mutex
	data map[string]map[string]model.Field
}

type bufferNumberFieldTerm struct {
	sync.Mutex
	data map[string]map[string]int
}

type bufferFieldLengthSize struct {
	sync.Mutex
	Length map[string]int
	Size   map[string]int
}

func (d *DiskIO) SaveDocumentFieldOnDisk(data chan DocumentFieldTransferData) {

	localBuffer := BufferFieldDocument{
		data: make(map[string]map[string]model.Field),
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case fieldDocument := <-data:
			localBuffer.Lock()
			localBuffer.data[fieldDocument.DocumentId] = fieldDocument.Field
			localBuffer.Unlock()
		case <-ticker.C:
			localBuffer.Lock()
			for documentId, fieldMap := range localBuffer.data {
				path := filepath.Join(d.RootFolder, file.Documents, file.DocumentsMetrics)
				file.CreteDirIfNotExist(path)

				buf := SerializeFieldMap(fieldMap)
				err := file.SecureSaveFile(path, documentId, file.CompressData(buf))
				if err != nil {
					log.Fatalf(err.Error())
				}
				delete(localBuffer.data, documentId)
			}
			localBuffer.Unlock()
		}
	}
}

func (d *DiskIO) SaveNumberFieldTermOnDisk(data chan NumberFieldTermTransferData) {

	buffer := bufferNumberFieldTerm{
		data: make(map[string]map[string]int),
	}

	ticker := time.NewTicker(7 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case fieldTerm := <-data:
			buffer.Lock()

			if termSize, ok := buffer.data[fieldTerm.FieldName]; ok {
				termSize[fieldTerm.Term] = fieldTerm.Size
			} else {
				buffer.data[fieldTerm.FieldName] = map[string]int{
					fieldTerm.Term: fieldTerm.Size,
				}
			}

			buffer.Unlock()
		case <-ticker.C:
			buffer.Lock()
			for fieldName, termSize := range buffer.data {

				path := filepath.Join(d.RootFolder, fieldName)
				file.CreteDirIfNotExist(path)

				buf := SerializeNumberFieldTerm(termSize)
				err := file.SecureSaveFile(path, file.NumberFieldTerm, file.CompressData(buf))
				if err != nil {
					log.Fatalf(err.Error())
				}
				delete(buffer.data, fieldName)
			}
			buffer.Unlock()
		}
	}

}

func (d *DiskIO) SaveFieldSizeLengthOnDisc(data chan FieldSizeLengthTransferData) {

	buffer := bufferFieldLengthSize{
		Size:   make(map[string]int),
		Length: make(map[string]int),
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case field := <-data:
			buffer.Lock()
			buffer.Length[field.FieldName] = field.Length
			buffer.Size[field.FieldName] = field.Size
			buffer.Unlock()
		case <-ticker.C:
			buffer.Lock()
			for fieldName := range buffer.Length {
				length := buffer.Length[fieldName]
				size := buffer.Size[fieldName]

				path := filepath.Join(d.RootFolder, fieldName)
				file.CreteDirIfNotExist(path)

				buf := SerializeFieldSizeLength(fieldName, int32(size), int32(length))

				//path = filepath.Join(path, file.MetricsFile)
				err := file.SecureSaveFile(path, file.MetricsFile, file.CompressData(buf))
				if err != nil {
					log.Fatalf(err.Error())
				}

				delete(buffer.Size, fieldName)
				delete(buffer.Length, fieldName)
			}
			buffer.Unlock()
		}
	}

}

// fieldName, term, documentId
type IndexTransferData struct {
	FieldName  string
	Term       string
	DocumentId string
	//IndexField map[string]*model.Set
}

type BufferIndex struct {
	sync.Mutex
	buffer map[string]map[string]map[string]bool
}

func (d *DiskIO) SaveIndexOnDisk(indexCh chan IndexTransferData) {

	bufferData := BufferIndex{
		buffer: make(map[string]map[string]map[string]bool),
	}

	ticker := time.NewTicker(13 * time.Second)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case index := <-indexCh:

			bufferData.Lock()
			indexField, ok := bufferData.buffer[index.FieldName]
			if !ok {
				indexField = make(map[string]map[string]bool)
				bufferData.buffer[index.FieldName] = indexField
			}

			termSet, ok := indexField[index.Term]
			if !ok {
				termSet = make(map[string]bool)
				indexField[index.Term] = termSet
			}

			count += 1
			termSet[index.DocumentId] = true
			bufferData.Unlock()

		case <-ticker.C:
			bufferData.Lock()
			for fieldName, terms := range bufferData.buffer {
				buffer := SerializeIndex(terms)
				path := filepath.Join(d.RootFolder, fieldName, file.Temp)
				file.CreteDirIfNotExist(path)

				tempFile := uuid.New().String()
				path = filepath.Join(path, tempFile)
				err := file.SaveFileOnDisk(path, file.CompressData(buffer))
				//err := file.SecureSaveFile(path, tempFile, file.CompressData(buffer))
				if err != nil {
					log.Fatalf(err.Error())
				}
			}
			fmt.Println("Gravou", count)
			count = 0
			bufferData.buffer = make(map[string]map[string]map[string]bool)
			bufferData.Unlock()
		}
	}

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//
	//}()
	//wg.Wait()
}
