// automatically generated, do not modify

package hub

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type HubResp struct {
	_tab flatbuffers.Table
}

func GetRootAsHubResp(buf []byte, offset flatbuffers.UOffsetT) *HubResp {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &HubResp{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *HubResp) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *HubResp) Cmd() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *HubResp) Err() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *HubResp) ReqId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *HubResp) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func HubRespStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func HubRespAddCmd(builder *flatbuffers.Builder, cmd flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(cmd), 0) }
func HubRespAddErr(builder *flatbuffers.Builder, err flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(err), 0) }
func HubRespAddReqId(builder *flatbuffers.Builder, reqId flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(reqId), 0) }
func HubRespAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(data), 0) }
func HubRespEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
