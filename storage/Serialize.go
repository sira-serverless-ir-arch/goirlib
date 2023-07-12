package storage

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sira-serverless-ir-arch/goirlib/model"
	"github.com/sira-serverless-ir-arch/goirlib/storage/buffers"
	"log"
)

func SerializeFieldSizeLength(name string, size int32, length int32) []byte {
	b := flatbuffers.NewBuilder(0)

	nameOffset := b.CreateString(name)

	buffers.FieldMetricsStart(b)
	buffers.FieldMetricsAddName(b, nameOffset)
	buffers.FieldMetricsAddSize(b, size)
	buffers.FieldMetricsAddLength(b, length)
	fieldSizeLength := buffers.FieldMetricsEnd(b)

	b.Finish(fieldSizeLength)

	return b.FinishedBytes()
}

func DeserializeFieldSizeLength(buf []byte) (string, int32, int32) {
	fieldSizeLength := buffers.GetRootAsFieldMetrics(buf, 0)

	name := string(fieldSizeLength.Name())
	size := fieldSizeLength.Size()
	length := fieldSizeLength.Length()

	return name, size, length
}

func SerializeNumberFieldTerm(data map[string]int) []byte {
	b := flatbuffers.NewBuilder(0)

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

func DeserializeNumberFieldTerm(buf []byte) map[string]int {
	fieldTerm := buffers.GetRootAsFieldTerm(buf, 0)

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

func SerializeFieldMap(data map[string]model.Field) []byte {
	b := flatbuffers.NewBuilder(0)

	var fieldEntries []flatbuffers.UOffsetT

	for key, field := range data {
		keyOffset := b.CreateString(key)
		nameOffset := b.CreateString(field.Name)

		var tfOffsets []flatbuffers.UOffsetT
		for tfKey, tfValue := range field.TF {
			tfKeyOffset := b.CreateString(tfKey)
			buffers.TermFrequencyStart(b)
			buffers.TermFrequencyAddKey(b, tfKeyOffset)
			buffers.TermFrequencyAddValue(b, int32(tfValue))
			tfOffset := buffers.TermFrequencyEnd(b)
			tfOffsets = append(tfOffsets, tfOffset)
		}

		buffers.FieldStartTfVector(b, len(tfOffsets))
		for i := len(tfOffsets) - 1; i >= 0; i-- {
			b.PrependUOffsetT(tfOffsets[i])
		}
		tfVector := b.EndVector(len(tfOffsets))

		buffers.FieldStart(b)
		buffers.FieldAddName(b, nameOffset)
		buffers.FieldAddLength(b, int32(field.Length))
		buffers.FieldAddTf(b, tfVector)
		fieldOffset := buffers.FieldEnd(b)

		buffers.FieldEntryStart(b)
		buffers.FieldEntryAddKey(b, keyOffset)
		buffers.FieldEntryAddValue(b, fieldOffset)
		fieldEntryOffset := buffers.FieldEntryEnd(b)

		fieldEntries = append(fieldEntries, fieldEntryOffset)
	}

	buffers.RootFieldEntryStartEntriesVector(b, len(fieldEntries))
	for i := len(fieldEntries) - 1; i >= 0; i-- {
		b.PrependUOffsetT(fieldEntries[i])
	}
	entriesVector := b.EndVector(len(fieldEntries))

	buffers.RootFieldEntryStart(b)
	buffers.RootFieldEntryAddEntries(b, entriesVector)
	rootOffset := buffers.RootFieldEntryEnd(b)

	b.Finish(rootOffset)

	return b.FinishedBytes()
}

func DeserializeFieldMap(buf []byte) map[string]model.Field {
	root := buffers.GetRootAsRootFieldEntry(buf, 0)

	var fieldMap buffers.FieldMap
	var field buffers.Field
	var termFrequency buffers.TermFrequency

	data := make(map[string]model.Field, root.EntriesLength())
	for i := 0; i < root.EntriesLength(); i++ {
		if root.Entries(&fieldMap, i) {
			key := string(fieldMap.Key())

			fieldMap.Value(&field)
			fieldModel := model.Field{
				Name:   string(field.Name()),
				Length: int(field.Length()),
				TF:     make(map[string]int, field.TfLength()),
			}
			for j := 0; j < field.TfLength(); j++ {
				if field.Tf(&termFrequency, j) {
					tfKey := string(termFrequency.Key())
					fieldModel.TF[tfKey] = int(termFrequency.Value())
				}
			}
			data[key] = fieldModel
		}
	}

	return data
}

func SerializeIndex(data map[string]map[string]bool) []byte {
	b := flatbuffers.NewBuilder(0)

	var terms []flatbuffers.UOffsetT

	for term, documents := range data {
		keyOffset := b.CreateString(term)

		var valsOffsets []flatbuffers.UOffsetT
		for documentId := range documents {
			valOffset := b.CreateString(documentId)
			valsOffsets = append(valsOffsets, valOffset)
		}

		buffers.DocumentStartValuesVector(b, len(valsOffsets))
		for i := len(valsOffsets) - 1; i >= 0; i-- {
			b.PrependUOffsetT(valsOffsets[i])
		}
		valuesVector := b.EndVector(len(valsOffsets))

		buffers.DocumentStart(b)
		buffers.DocumentAddValues(b, valuesVector)
		document := buffers.DocumentEnd(b)

		buffers.TermStart(b)
		buffers.TermAddKey(b, keyOffset)
		buffers.TermAddValues(b, document)
		term := buffers.TermEnd(b)

		terms = append(terms, term)
	}

	buffers.IndexStartEntriesVector(b, len(terms))
	for i := len(terms) - 1; i >= 0; i-- {
		b.PrependUOffsetT(terms[i])
	}
	entriesVector := b.EndVector(len(terms))

	buffers.IndexStart(b)
	buffers.IndexAddEntries(b, entriesVector)
	index := buffers.IndexEnd(b)

	b.Finish(index)

	return b.FinishedBytes()
}

func DeserializeIndex(buf []byte) map[string]map[string]bool {
	index := buffers.GetRootAsIndex(buf, 0)

	data := make(map[string]map[string]bool)

	termsLen := index.EntriesLength()

	for i := 0; i < termsLen; i++ {
		term := new(buffers.Term)
		if !index.Entries(term, i) {
			log.Fatalf("Failed to get TermSize")
		}

		key := string(term.Key())

		document := new(buffers.Document)
		term.Values(document)

		valuesLen := document.ValuesLength()
		values := make(map[string]bool)

		for j := 0; j < valuesLen; j++ {
			value := string(document.Values(j))
			values[value] = true
		}

		data[key] = values
	}

	return data
}
