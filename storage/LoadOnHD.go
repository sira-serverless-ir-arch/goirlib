package storage

import (
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage/disk"
	"log"
)

func (d *Disk) LoadFieldSizeLengthOnHD(rootFolder string) {
	fields := disk.ListDirectories(rootFolder)
	for _, field := range fields {
		buf, err := disk.ReadFileOnDisk(rootFolder+field, "metrics")
		if err != nil {
			log.Fatalf(err.Error())
		}

		name, size, length := disk.DeserializeFieldSizeLength(buf)
		d.FieldSize[name] = int(size)
		d.FieldLength[name] = int(length)
	}

}

func (d *Disk) LoadIndexOnHD(rootFolder string) {

	fields := disk.ListDirectories(rootFolder)

	for _, field := range fields {
		buf, err := disk.ReadFileOnDisk(rootFolder+field, "index")
		if err != nil {
			log.Fatalf(err.Error())
		}

		termDocuments := make(map[string]*model.Set)
		count := 0
		for term, documents := range disk.DeserializeIndex(buf) {
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
