// automatically generated, do not modify

package worker

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type req_header struct {
	_tab flatbuffers.Table
}

func GetRootAsreq_header(buf []byte, offset flatbuffers.UOffsetT) *req_header {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &req_header{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *req_header) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *req_header) Seq() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *req_header) Sid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *req_header) NoResp() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *req_header) Token() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *req_header) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *req_header) Gzip() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func req_headerStart(builder *flatbuffers.Builder) { builder.StartObject(6) }
func req_headerAddSeq(builder *flatbuffers.Builder, seq int32) { builder.PrependInt32Slot(0, seq, 0) }
func req_headerAddSid(builder *flatbuffers.Builder, sid flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(sid), 0) }
func req_headerAddNoResp(builder *flatbuffers.Builder, noResp byte) { builder.PrependByteSlot(2, noResp, 0) }
func req_headerAddToken(builder *flatbuffers.Builder, token flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(token), 0) }
func req_headerAddVersion(builder *flatbuffers.Builder, version flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(version), 0) }
func req_headerAddGzip(builder *flatbuffers.Builder, gzip byte) { builder.PrependByteSlot(5, gzip, 0) }
func req_headerEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
