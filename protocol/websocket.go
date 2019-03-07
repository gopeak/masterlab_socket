package protocol

import (
	"encoding/json"
	"fmt"
	"masterlab_socket/util"
)

type Json struct {
	ProtocolObj ProtocolType
	Data        []byte
}

func (this *Json) Init() *Json {

	this.ProtocolObj = ProtocolType{}
	this.ProtocolObj.ReqObj  = ReqRoot{}
	this.ProtocolObj.RespObj = ResponseRoot{}
	this.ProtocolObj.BroatcastObj = BroatcastRoot{}
	this.ProtocolObj.PushObj = PushRoot{}
	return this
}



func (this *Json) GetReqObj(data []byte) (*ReqRoot, error) {

	this.Data = data
	stb := &ReqRoot{}
	err := json.Unmarshal( data, stb )
 	stb.Data = main.Convert2Byte( stb.WsData )
	fmt.Println( "main.Convert2Byte( stb.WsData ):",string(stb.Data) )
	return stb, err


}

func (this *Json) GetRespObj(data []byte) (*ResponseRoot, error) {
	this.Data = data
	stb := &ResponseRoot{}
	err := json.Unmarshal(data,  stb)

	//this.ProtocolObj.RespObj = stb
	return stb, err
}

func (this *Json) GetBroatcastObj(data []byte) (*BroatcastRoot, error) {
	this.Data = data
	stb := &BroatcastRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.BroatcastObj = stb
	return stb, err
}

func (this *Json) GetPushObj(data []byte) (*PushRoot, error) {
	this.Data = data
	stb := &PushRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.PushObj = stb
	return stb, err
}

func (this *Json) WrapRespObj( req_obj *ReqRoot, invoker_ret []byte, status int ) ResponseRoot {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = req_obj.Header.Cmd
	resp_header_obj.SeqId = req_obj.Header.SeqId
	resp_header_obj.Gzip = req_obj.Header.Gzip
	resp_header_obj.Sid = req_obj.Header.Sid
	resp_header_obj.Status = status
	this.ProtocolObj.RespObj.Header =resp_header_obj
	this.ProtocolObj.RespObj.Data = invoker_ret
	this.ProtocolObj.RespObj.Type = TypeResp

	return this.ProtocolObj.RespObj
}

func (this *Json) WrapResp(   header []byte, data []byte, status int, msg string )  []byte  {

	header = main.TrimX001( header )
	data_str := string(data)
	if( main.TrimStr(data_str)==""){
		data_str = `""`
	}
	header_str := string(header)
	if( main.TrimStr(header_str)==""){
		header_str = "{}"
	}
	return []byte(fmt.Sprintf(`{"type":"%s","status":%d,"msg":"%s","header":%s,"data":%s}`,
			TypeResp, status, msg, header_str ,data_str ))


}

func (this *Json) WrapPushRespObj(to_sid string, from_sid string , data []byte ) PushRoot {

	push_header_obj := PushHeader{}
	push_header_obj.Sid = from_sid
	push_obj := PushRoot{}
	push_obj.Data  = string(data)
	push_obj.Header =push_header_obj
	push_obj.Type  = TypePush
	/*
	var map_interface map[string]interface{}
	var map_str map[string]string
	err:=json.Unmarshal( data, map_interface )
	if err==nil {
		push_obj.Data = map_interface
	}else{
		err=json.Unmarshal( data, map_str )
		if err==nil {
			push_obj.Data = map_str
		}
	}
	fmt.Println( "ws WrapPushRespObj to_data_buf:", string(data))
	*/
	return push_obj
}

func (this *Json) WrapPushResp(to_sid string, from_sid string , data []byte ) []byte {

	data_str := string(data)
	if( main.TrimStr(data_str)==""){
		data_str = `""`
	}
	return []byte(fmt.Sprintf(`{"header":{"sid":"%s"},"type":"%s","data":%s  }`,
		 from_sid ,TypePush, data_str ))

}


func (this *Json) WrapBroatcastRespObj(area_id string, from_sid string , data []byte) BroatcastRoot {

	broatcast_header_obj := BroatcastHeader{}
	broatcast_header_obj.Sid = from_sid
	broatcast_header_obj.AreaId = area_id

	broatcast_obj := BroatcastRoot{}
	broatcast_obj.Header =broatcast_header_obj
	broatcast_obj.Data  = string(data)
	broatcast_obj.Type  = TypeBroatcast

	return broatcast_obj
}

/**
 * 封包返回客户端错误的消息
 */
func (this *Json) WrapRespErr(err string) []byte {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = "WrapRespErr"
	resp_header_obj.SeqId = 0
	resp_header_obj.Sid = ""
	resp_header_obj.Status = 500
	this.ProtocolObj.RespObj.Header =resp_header_obj
	this.ProtocolObj.RespObj.Data = []byte(err)
	this.ProtocolObj.RespObj.Type = TypeError

	buf,_ := json.Marshal( this.ProtocolObj.RespObj )

	return buf
}


