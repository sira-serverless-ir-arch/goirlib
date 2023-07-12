package storage

import (
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

func (d *DiskIO) LoadIndexOnHD() (map[string]map[string]*model.Set, error) {

	fields := file.ListDirectories(d.RootFolder)
	index := make(map[string]map[string]*model.Set)

	for _, field := range fields {
		path := filepath.Join(d.RootFolder, field, file.IndexFile)
		buf, err := file.ReadFileOnDisk(path)
		if err != nil {
			return nil, err
		}

		termDocuments := make(map[string]*model.Set)
		count := 0
		for term, documents := range DeserializeIndex(file.DecompressData(buf)) {
			set := model.NewSet()
			for documentId := range documents {
				set.Add(documentId)
			}
			termDocuments[term] = set
			count++
		}
		index[field] = termDocuments
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
	TermSize  map[string]int
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

	go func() {
		for fieldDocument := range data {
			localBuffer.Lock()
			localBuffer.data[fieldDocument.DocumentId] = fieldDocument.Field
			localBuffer.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			localBuffer.Lock()
			for documentId, fieldMap := range localBuffer.data {

				path := filepath.Join(d.RootFolder, file.Documents)
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
	}()
}

func (d *DiskIO) SaveNumberFieldTermOnDisk(data chan NumberFieldTermTransferData) {

	buffer := bufferNumberFieldTerm{
		data: make(map[string]map[string]int),
	}

	go func() {
		for fieldTerm := range data {
			buffer.Lock()
			buffer.data[fieldTerm.FieldName] = fieldTerm.TermSize
			buffer.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
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
	}()
}

func (d *DiskIO) SaveFieldSizeLengthOnDisc(data chan FieldSizeLengthTransferData) {

	buffer := bufferFieldLengthSize{
		Size:   make(map[string]int),
		Length: make(map[string]int),
	}

	go func() {
		for field := range data {
			buffer.Lock()
			buffer.Length[field.FieldName] = field.Length
			buffer.Size[field.FieldName] = field.Size
			buffer.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
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
	}()
}

type IndexTransferData struct {
	FieldName  string
	IndexField map[string]*model.Set
}

type BufferIndex struct {
	sync.Mutex
	buffer map[string]map[string]*model.Set
}

func (d *DiskIO) SaveIndexOnDisk(indexCh chan IndexTransferData) {

	bufferData := BufferIndex{
		buffer: make(map[string]map[string]*model.Set),
	}

	go func() {
		for index := range indexCh {
			bufferData.Lock()
			bufferData.buffer[index.FieldName] = index.IndexField
			bufferData.Unlock()
		}
	}()

	semaphore := make(chan struct{}, 5) // Limit to 5 goroutines
	go func() {
		for {
			time.Sleep(5 * time.Second)
			bufferData.Lock()
			for fieldName, data := range bufferData.buffer {
				semaphore <- struct{}{}
				go func(fieldName string, terms map[string]*model.Set) {
					defer func() { <-semaphore }()

					path := filepath.Join(d.RootFolder, fieldName)
					file.CreteDirIfNotExist(path)

					tempData := make(map[string]map[string]bool)

					for term, documents := range terms {
						tempData[term] = documents.GetData()
					}

					buff := SerializeIndex(tempData)
					err := file.SecureSaveFile(path, file.IndexFile, file.CompressData(buff))
					if err != nil {
						log.Fatalf(err.Error())
					}

				}(fieldName, data)
				delete(bufferData.buffer, fieldName)
			}
			bufferData.Unlock()
		}
	}()
}
