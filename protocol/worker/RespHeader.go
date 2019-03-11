// automatically generated, do not modify

package worker

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type RespHeader struct {
	_tab flatbuffers.Table
}

func GetRootAsRespHeader(buf []byte, offset flatbuffers.UOffsetT) *RespHeader {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RespHeader{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *RespHeader) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RespHeader) Cmd() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RespHeader) Seq() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RespHeader) Sid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RespHeader) Gzip() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RespHeader) Status() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func RespHeaderStart(builder *flatbuffers.Builder) { builder.StartObject(5) }
func RespHeaderAddCmd(builder *flatbuffers.Builder, cmd flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(cmd), 0) }
func RespHeaderAddSeq(builder *flatbuffers.Builder, seq int32) { builder.PrependInt32Slot(1, seq, 0) }
func RespHeaderAddSid(builder *flatbuffers.Builder, sid flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(sid), 0) }
func RespHeaderAddGzip(builder *flatbuffers.Builder, gzip byte) { builder.PrependByteSlot(3, gzip, 0) }
func RespHeaderAddStatus(builder *flatbuffers.Builder, status int32) { builder.PrependInt32Slot(4, status, 0) }
func RespHeaderEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
