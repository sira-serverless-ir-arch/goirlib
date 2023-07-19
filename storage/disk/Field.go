package disk

//type FieldSizeLengthTransferData struct {
//	FieldName string
//	Size      int
//	Length    int
//}
//
//type NumberFieldTermTransferData struct {
//	FieldName string
//	TermSize  map[string]int
//}
//
//type DocumentFieldTransferData struct {
//	DocumentId string
//	Field      map[string]model.Field
//}
//
//type BufferFieldDocument struct {
//	sync.Mutex
//	disk_storage map[string]map[string]model.Field
//}
//
//type bufferNumberFieldTerm struct {
//	sync.Mutex
//	disk_storage map[string]map[string]int
//}
//
//type bufferFieldLengthSize struct {
//	sync.Mutex
//	Length map[string]int
//	Size   map[string]int
//}
//
//func SaveDocumentFieldOnDisk(rootFolder string, disk_storage chan DocumentFieldTransferData) {
//
//	localBuffer := BufferFieldDocument{
//		disk_storage: make(map[string]map[string]model.Field),
//	}
//
//	go func() {
//		for fieldDocument := range disk_storage {
//			localBuffer.Lock()
//			localBuffer.disk_storage[fieldDocument.DocumentId] = fieldDocument.Field
//			localBuffer.Unlock()
//		}
//	}()
//
//	go func() {
//		for {
//			time.Sleep(5 * time.Second)
//			localBuffer.Lock()
//			for documentId, fieldMap := range localBuffer.disk_storage {
//
//				path := filepath.Join(rootFolder, file.Documents, file.DocumentsMetrics)
//				file.CreteDirIfNotExist(path)
//
//				buf := SerializeFieldMap(fieldMap)
//				path = filepath.Join(path, documentId)
//
//				err := file.SaveFileOnDisk(path, file.CompressData(buf))
//				if err != nil {
//					log.Fatalf(err.Error())
//				}
//				delete(localBuffer.disk_storage, documentId)
//			}
//			localBuffer.Unlock()
//		}
//	}()
//}
//
//func SaveNumberFieldTermOnDisk(rootFolder string, disk_storage chan NumberFieldTermTransferData) {
//
//	buffer := bufferNumberFieldTerm{
//		disk_storage: make(map[string]map[string]int),
//	}
//
//	go func() {
//		for fieldTerm := range disk_storage {
//			buffer.Lock()
//			buffer.disk_storage[fieldTerm.FieldName] = fieldTerm.TermSize
//			buffer.Unlock()
//		}
//	}()
//
//	go func() {
//		for {
//			time.Sleep(5 * time.Second)
//			buffer.Lock()
//			for fieldName, termSize := range buffer.disk_storage {
//
//				path := filepath.Join(rootFolder, fieldName)
//				file.CreteDirIfNotExist(path)
//
//				buf := SerializeNumberFieldTerm(termSize)
//				path = filepath.Join(path, file.NumberFieldTerm)
//				err := file.SaveFileOnDisk(path, file.CompressData(buf))
//				if err != nil {
//					log.Fatalf(err.Error())
//				}
//				delete(buffer.disk_storage, fieldName)
//			}
//			buffer.Unlock()
//		}
//	}()
//}
//
//func SaveFieldSizeLengthOnDisc(rootFolder string, disk_storage chan FieldSizeLengthTransferData) {
//
//	buffer := bufferFieldLengthSize{
//		Size:   make(map[string]int),
//		Length: make(map[string]int),
//	}
//
//	go func() {
//		for field := range disk_storage {
//			buffer.Lock()
//			buffer.Length[field.FieldName] = field.Length
//			buffer.Size[field.FieldName] = field.Size
//			buffer.Unlock()
//		}
//	}()
//
//	go func() {
//		for {
//			time.Sleep(5 * time.Second)
//			buffer.Lock()
//			for fieldName := range buffer.Length {
//				length := buffer.Length[fieldName]
//				size := buffer.Size[fieldName]
//
//				path := filepath.Join(rootFolder, fieldName)
//				file.CreteDirIfNotExist(path)
//
//				buf := SerializeFieldSizeLength(fieldName, int32(size), int32(length))
//
//				path = filepath.Join(path, file.MetricsFile)
//				err := file.SaveFileOnDisk(path, file.CompressData(buf))
//				if err != nil {
//					log.Fatalf(err.Error())
//				}
//
//				delete(buffer.Size, fieldName)
//				delete(buffer.Length, fieldName)
//			}
//			buffer.Unlock()
//		}
//	}()
//
//}
//
//func SerializeFieldSizeLength(name string, size int32, length int32) []byte {
//	b := flatbuffers.NewBuilder(0)
//
//	nameOffset := b.CreateString(name)
//
//	buffers.FieldMetricsStart(b)
//	buffers.FieldMetricsAddName(b, nameOffset)
//	buffers.FieldMetricsAddSize(b, size)
//	buffers.FieldMetricsAddLength(b, length)
//	fieldSizeLength := buffers.FieldMetricsEnd(b)
//
//	b.Finish(fieldSizeLength)
//
//	return b.FinishedBytes()
//}
//
//func DeserializeFieldSizeLength(buf []byte) (string, int32, int32) {
//	fieldSizeLength := buffers.GetRootAsFieldMetrics(buf, 0)
//
//	name := string(fieldSizeLength.Name())
//	size := fieldSizeLength.Size()
//	length := fieldSizeLength.Length()
//
//	return name, size, length
//}
//
//func SerializeNumberFieldTerm(disk_storage map[string]int) []byte {
//	b := flatbuffers.NewBuilder(0)
//
//	var termSizes []flatbuffers.UOffsetT
//	for term, size := range disk_storage {
//		termKey := b.CreateString(term)
//
//		buffers.TermSizeStart(b)
//		buffers.TermSizeAddKey(b, termKey)
//		buffers.TermSizeAddValue(b, int32(size))
//		termSize := buffers.TermSizeEnd(b)
//
//		termSizes = append(termSizes, termSize)
//	}
//
//	buffers.FieldTermStartEntriesVector(b, len(termSizes))
//	for i := len(termSizes) - 1; i >= 0; i-- {
//		b.PrependUOffsetT(termSizes[i])
//	}
//	entriesVector := b.EndVector(len(termSizes))
//
//	buffers.FieldTermStart(b)
//	buffers.FieldTermAddEntries(b, entriesVector)
//	fieldTerm := buffers.FieldTermEnd(b)
//
//	b.Finish(fieldTerm)
//
//	return b.FinishedBytes()
//}
//
//func DeserializeNumberFieldTerm(buf []byte) map[string]int {
//	fieldTerm := buffers.GetRootAsFieldTerm(buf, 0)
//
//	entriesLength := fieldTerm.EntriesLength()
//	disk_storage := make(map[string]int, entriesLength)
//
//	var termSize buffers.TermSize
//	for i := 0; i < entriesLength; i++ {
//		if fieldTerm.Entries(&termSize, i) {
//			disk_storage[string(termSize.Key())] = int(termSize.Value())
//		}
//	}
//
//	return disk_storage
//}
//
//func SerializeFieldMap(disk_storage map[string]model.Field) []byte {
//	b := flatbuffers.NewBuilder(0)
//
//	var fieldEntries []flatbuffers.UOffsetT
//
//	for key, field := range disk_storage {
//		keyOffset := b.CreateString(key)
//		nameOffset := b.CreateString(field.Name)
//
//		var tfOffsets []flatbuffers.UOffsetT
//		for tfKey, tfValue := range field.TF {
//			tfKeyOffset := b.CreateString(tfKey)
//			buffers.TermFrequencyStart(b)
//			buffers.TermFrequencyAddKey(b, tfKeyOffset)
//			buffers.TermFrequencyAddValue(b, int32(tfValue))
//			tfOffset := buffers.TermFrequencyEnd(b)
//			tfOffsets = append(tfOffsets, tfOffset)
//		}
//
//		buffers.FieldStartTfVector(b, len(tfOffsets))
//		for i := len(tfOffsets) - 1; i >= 0; i-- {
//			b.PrependUOffsetT(tfOffsets[i])
//		}
//		tfVector := b.EndVector(len(tfOffsets))
//
//		buffers.FieldStart(b)
//		buffers.FieldAddName(b, nameOffset)
//		buffers.FieldAddLength(b, int32(field.Length))
//		buffers.FieldAddTf(b, tfVector)
//		fieldOffset := buffers.FieldEnd(b)
//
//		buffers.FieldEntryStart(b)
//		buffers.FieldEntryAddKey(b, keyOffset)
//		buffers.FieldEntryAddValue(b, fieldOffset)
//		fieldEntryOffset := buffers.FieldEntryEnd(b)
//
//		fieldEntries = append(fieldEntries, fieldEntryOffset)
//	}
//
//	buffers.RootFieldEntryStartEntriesVector(b, len(fieldEntries))
//	for i := len(fieldEntries) - 1; i >= 0; i-- {
//		b.PrependUOffsetT(fieldEntries[i])
//	}
//	entriesVector := b.EndVector(len(fieldEntries))
//
//	buffers.RootFieldEntryStart(b)
//	buffers.RootFieldEntryAddEntries(b, entriesVector)
//	rootOffset := buffers.RootFieldEntryEnd(b)
//
//	b.Finish(rootOffset)
//
//	return b.FinishedBytes()
//}
//
//func DeserializeFieldMap(buf []byte) map[string]model.Field {
//	root := buffers.GetRootAsRootFieldEntry(buf, 0)
//
//	var fieldMap buffers.FieldMap
//	var field buffers.Field
//	var termFrequency buffers.TermFrequency
//
//	disk_storage := make(map[string]model.Field, root.EntriesLength())
//	for i := 0; i < root.EntriesLength(); i++ {
//		if root.Entries(&fieldMap, i) {
//			key := string(fieldMap.Key())
//
//			fieldMap.Value(&field)
//			fieldModel := model.Field{
//				Name:   string(field.Name()),
//				Length: int(field.Length()),
//				TF:     make(map[string]int, field.TfLength()),
//			}
//			for j := 0; j < field.TfLength(); j++ {
//				if field.Tf(&termFrequency, j) {
//					tfKey := string(termFrequency.Key())
//					fieldModel.TF[tfKey] = int(termFrequency.Value())
//				}
//			}
//			disk_storage[key] = fieldModel
//		}
//	}
//
//	return disk_storage
//}
