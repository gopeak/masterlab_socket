package golang

import (
	"masterlab_socket/protocol"
	"masterlab_socket/lib/websocket"
	"net"
)

type TaskType struct {
	Conn *net.TCPConn

	WsConn *websocket.Conn

	ReqType string

	ReqHeader *protocol.ReqHeader

	Data []byte

}

type ReturnType struct {
	Ret string `json:"ret"`

	Type string `json:"type"`

	Sid string `json:"sid"`

	Msg string `json:"msg"`
}

func (this *TaskType) Init( conn *net.TCPConn, req_obj *protocol.ReqRoot ) *TaskType {

	this.Conn = conn
	this.ReqType = req_obj.Type
	this.ReqHeader = &req_obj.Header
	this.Data   = req_obj.Data
	return this
}


func (this *TaskType) WsInit(wsconn *websocket.Conn, req_obj *protocol.ReqRoot) *TaskType {


	this.ReqType = req_obj.Type
	this.ReqHeader = &req_obj.Header
	this.Data   = req_obj.Data
	this.WsConn = wsconn
	return this
}