package disk

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sira-serverless-ir-arch/goirlib/storage/buffers"
	"log"
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

				folder := rootFolder + fieldName
				CreteDirIfNotExist(folder)

				buf := SerializeFieldSizeLength(fieldName, int32(size), int32(length))
				err := SaveFileOnDisk(folder, "metrics", buf)
				if err != nil {
					log.Fatalf(err.Error())
				}

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
