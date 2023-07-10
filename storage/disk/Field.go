package disk

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sira-serverless-ir-arch/goirlib/file"
	"github.com/sira-serverless-ir-arch/goirlib/storage/buffers"
	"log"
	"path/filepath"
	"sync"
	"time"
)

var (
	mFz = sync.Mutex{}
)

type FieldSizeLengthTransferData struct {
	FieldName string
	Size      int
	Length    int
}

func SaveFieldSizeLengthOnDisc(rootFolder string, data chan FieldSizeLengthTransferData) {

	bufferFieldLength := make(map[string]int)
	bufferFieldSize := make(map[string]int)

	go func() {
		for field := range data {
			mFz.Lock()
			bufferFieldLength[field.FieldName] = field.Length
			bufferFieldSize[field.FieldName] = field.Size
			mFz.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			m.Lock()
			for fieldName := range bufferFieldLength {
				length := bufferFieldLength[fieldName]
				size := bufferFieldSize[fieldName]

				path := filepath.Join(rootFolder, fieldName)
				file.CreteDirIfNotExist(path)

				buf := SerializeFieldSizeLength(fieldName, int32(size), int32(length))

				path = filepath.Join(path, file.MetricsFile)
				err := file.SaveFileOnDisk(path, file.CompressData(buf))
				if err != nil {
					log.Fatalf(err.Error())
				}

				delete(bufferFieldSize, fieldName)
				delete(bufferFieldLength, fieldName)
			}
			m.Unlock()
		}
	}()

}

func SerializeFieldSizeLength(name string, size int32, length int32) []byte {
	b := flatbuffers.NewBuilder(0)

	nameOffset := b.CreateString(name)

	buffers.FieldMetrixStart(b)
	buffers.FieldMetrixAddName(b, nameOffset)
	buffers.FieldMetrixAddSize(b, size)
	buffers.FieldMetrixAddLength(b, length)
	fieldSizeLength := buffers.FieldMetrixEnd(b)

	// Finalizando o buffer.
	b.Finish(fieldSizeLength)

	return b.FinishedBytes()
}

func DeserializeFieldSizeLength(buf []byte) (string, int32, int32) {
	fieldSizeLength := buffers.GetRootAsFieldMetrix(buf, 0)

	name := string(fieldSizeLength.Name())
	size := fieldSizeLength.Size()
	length := fieldSizeLength.Length()

	return name, size, length
}

func SerializeFieldTerm(data map[string]int) []byte {
	b := flatbuffers.NewBuilder(1024)

	var termSizes []flatbuffers.UOffsetT
	for term, size := range data {
		termKey := b.CreateString(term)

		buffers.TermSizeStart(b)
		buffers.TermSizeAddKey(b, termKey)
		buffers.TermSizeAddValue(b, int32(size))
		termSize := buffers.TermSizeEnd(b)

		termSizes = append(termSizes, termSize)
	}

	buffers.FieldTermStartEntriesVector(b, len(termSizes))
	for i := len(termSizes) - 1; i >= 0; i-- {
		b.PrependUOffsetT(termSizes[i])
	}
	entriesVector := b.EndVector(len(termSizes))

	buffers.FieldTermStart(b)
	buffers.FieldTermAddEntries(b, entriesVector)
	fieldTerm := buffers.FieldTermEnd(b)

	b.Finish(fieldTerm)

	return b.FinishedBytes()
}

func DeserializeFieldTerm(buf []byte) map[string]int {
	fieldTerm := buffers.GetRootAsFieldTerm(buf, 1024)

	entriesLength := fieldTerm.EntriesLength()
	data := make(map[string]int, entriesLength)

	var termSize buffers.TermSize
	for i := 0; i < entriesLength; i++ {
		if fieldTerm.Entries(&termSize, i) {
			data[string(termSize.Key())] = int(termSize.Value())
		}
	}

	return data
}
