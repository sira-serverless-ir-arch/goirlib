package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage/disk"
	"log"
	"path/filepath"
)

func (d *Disk) LoadNumberFieldTermOnHD(rootFolder string) {
	fields := file.ListDirectories(rootFolder)
	for _, field := range fields {
		path := filepath.Join(rootFolder, field, file.NumberFieldTerm)
		buf, err := file.ReadFileOnDisk(path)
		if err != nil {
			log.Fatalf(err.Error())
		}
		d.NumberFieldTerm[field] = disk.DeserializeNumberFieldTerm(file.DecompressData(buf))
	}
}

func (d *Disk) LoadFieldSizeLengthOnHD(rootFolder string) {
	fields := file.ListDirectories(rootFolder)
	for _, field := range fields {
		path := filepath.Join(rootFolder, field, file.MetricsFile)
		buf, err := file.ReadFileOnDisk(path)
		if err != nil {
			log.Fatalf(err.Error())
		}

		name, size, length := disk.DeserializeFieldSizeLength(file.DecompressData(buf))
		d.FieldSize[name] = int(size)
		d.FieldLength[name] = int(length)
	}
}

func (d *Disk) LoadFieldDocumentOnHD(rootFolder string, documentId string) (map[string]model.Field, bool) {
	path := filepath.Join(rootFolder, file.Documents, file.DocumentsMetrics, documentId)
	buf, err := file.ReadFileOnDisk(path)
	if err != nil {
		return nil, false
	}
	return disk.DeserializeFieldMap(file.DecompressData(buf)), true
}

func (d *Disk) LoadIndexOnHD(rootFolder string) {

	fields := file.ListDirectories(rootFolder)

	for _, field := range fields {
		path := filepath.Join(rootFolder, field, file.IndexFile)
		buf, err := file.ReadFileOnDisk(path)
		if err != nil {
			log.Fatalf(err.Error())
		}

		termDocuments := make(map[string]*model.Set)
		count := 0
		for term, documents := range disk.DeserializeIndex(file.DecompressData(buf)) {
			set := model.NewSet()
			for documentId := range documents {
				set.Add(documentId)
			}
			termDocuments[term] = set
			count++
		}
		d.Index[field] = termDocuments
	}
}
