// automatically generated, do not modify

package hub

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type HubReq struct {
	_tab flatbuffers.Table
}

func GetRootAsHubReq(buf []byte, offset flatbuffers.UOffsetT) *HubReq {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &HubReq{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *HubReq) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *HubReq) Cmd() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *HubReq) Sid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *HubReq) ReqId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *HubReq) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func HubReqStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func HubReqAddCmd(builder *flatbuffers.Builder, cmd flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(cmd), 0) }
func HubReqAddSid(builder *flatbuffers.Builder, sid flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(sid), 0) }
func HubReqAddReqId(builder *flatbuffers.Builder, reqId flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(reqId), 0) }
func HubReqAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(data), 0) }
func HubReqEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
