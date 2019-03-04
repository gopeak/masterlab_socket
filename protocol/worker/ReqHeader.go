// automatically generated, do not modify

package worker

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ReqHeader struct {
	_tab flatbuffers.Table
}

func GetRootAsReqHeader(buf []byte, offset flatbuffers.UOffsetT) *ReqHeader {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ReqHeader{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *ReqHeader) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ReqHeader) Seq() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ReqHeader) Sid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReqHeader) NoResp() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ReqHeader) Token() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReqHeader) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReqHeader) Gzip() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func ReqHeaderStart(builder *flatbuffers.Builder) { builder.StartObject(6) }
func ReqHeaderAddSeq(builder *flatbuffers.Builder, seq int32) { builder.PrependInt32Slot(0, seq, 0) }
func ReqHeaderAddSid(builder *flatbuffers.Builder, sid flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(sid), 0) }
func ReqHeaderAddNoResp(builder *flatbuffers.Builder, noResp byte) { builder.PrependByteSlot(2, noResp, 0) }
func ReqHeaderAddToken(builder *flatbuffers.Builder, token flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(token), 0) }
func ReqHeaderAddVersion(builder *flatbuffers.Builder, version flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(version), 0) }
func ReqHeaderAddGzip(builder *flatbuffers.Builder, gzip byte) { builder.PrependByteSlot(5, gzip, 0) }
func ReqHeaderEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
