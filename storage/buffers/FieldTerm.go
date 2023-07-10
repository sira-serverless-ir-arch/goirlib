// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package buffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type FieldTerm struct {
	_tab flatbuffers.Table
}

func GetRootAsFieldTerm(buf []byte, offset flatbuffers.UOffsetT) *FieldTerm {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &FieldTerm{}
	x.Init(buf, n+offset)
	return x
}

func FinishFieldTermBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsFieldTerm(buf []byte, offset flatbuffers.UOffsetT) *FieldTerm {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &FieldTerm{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedFieldTermBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *FieldTerm) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FieldTerm) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *FieldTerm) Entries(obj *TermSize, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *FieldTerm) EntriesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func FieldTermStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func FieldTermAddEntries(builder *flatbuffers.Builder, entries flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(entries), 0)
}
func FieldTermStartEntriesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func FieldTermEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
