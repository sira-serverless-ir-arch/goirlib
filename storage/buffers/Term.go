// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package buffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Term struct {
	_tab flatbuffers.Table
}

func GetRootAsTerm(buf []byte, offset flatbuffers.UOffsetT) *Term {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Term{}
	x.Init(buf, n+offset)
	return x
}

func FinishTermBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsTerm(buf []byte, offset flatbuffers.UOffsetT) *Term {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Term{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedTermBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Term) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Term) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Term) Key() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Term) Values(obj *Document) *Document {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Document)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func TermStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func TermAddKey(builder *flatbuffers.Builder, key flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(key), 0)
}
func TermAddValues(builder *flatbuffers.Builder, values flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(values), 0)
}
func TermEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
