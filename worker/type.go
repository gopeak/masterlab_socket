package worker

import (
	"encoding/json"
	"fmt"
	"masterlab_socket/protocol"
	"masterlab_socket/lib/websocket"
	"masterlab_socket/golog"
	"masterlab_socket/area"
	"masterlab_socket/util"
	"masterlab_socket/global"
	"net"
	"reflect"
)


type WorkerConfigType struct {

	Loglevel     string		`toml:"loglevel"`
	SingleMode   bool	  	`toml:"single_mode"`
	Servers [][]string       	`toml:"servers"`
	ToHub []string  		`toml:"connect_to_hub"`

}


var WorkerConfig   WorkerConfigType


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


func Invoker(conn *net.TCPConn, req_obj *protocol.ReqRoot) interface{} {

	task_obj := new(TaskType).Init(conn, req_obj)

	invoker_ret := InvokeObjectMethod(task_obj, req_obj.Header.Cmd)
	//fmt.Println("invoker_ret", invoker_ret)
	// 判断是否需要响应数据
	if req_obj.Type == protocol.TypeReq && !req_obj.Header.NoResp {
		protocolPack := new(protocol.Pack)
		protocolPack.Init()
		invoker_ret_buf := util.Convert2Byte( invoker_ret )
		switch invoker_ret.(type) {
		case ReturnType:
			tmp :=  invoker_ret.(ReturnType)
			invoker_ret_buf,_ = json.Marshal( tmp )
		}
		buf, _ := protocolPack.WrapResp( req_obj.Header.Cmd, req_obj.Header.Sid ,req_obj.Header.SeqId, 200, invoker_ret_buf )
		conn.Write(buf)
	}
	if global.SingleMode {
		if global.IsAuthCmd(req_obj.Header.Cmd) {
			area.ConnRegister(conn, req_obj.Header.Sid)
		}
	}
	return invoker_ret
}

func InvokeObjectMethod(object interface{}, methodName string ) interface{} {

	inputs := make([]reflect.Value, 0)

	fnc := reflect.ValueOf(object).MethodByName(methodName)
	empty:=reflect.Value{}
	if fnc==empty{
		fmt.Println( " reflect MethodByName " ,methodName ," no found!" )
		//golog.Error( " reflect MethodByName " ,methodName ," no found!" )
		return ""
	}
	ret := fnc.Call( inputs )[0]

	switch vtype := ret.Interface().(type) {

	case nil:
		return nil
	case bool:
		return ret.Interface().(bool)

	case float32:
		return ret.Interface().(float32)
	case float64:
		return ret.Interface().(float32)
	case int:
		return ret.Interface().(int)
	case uint8:
		return ret.Interface().(uint8)
	case uint16:
		return ret.Interface().(uint16)
	case uint32:
		return ret.Interface().(uint32)
	case uint64:
		return ret.Interface().(uint64)
	case int8:
		return ret.Interface().(int8)
	case int16:
		return ret.Interface().(int16)
	case int32:
		return ret.Interface().(int32)
	case int64:
		return ret.Interface().(int64)
	case []byte:
		return  ret.Interface().([]byte)
	case string:
		return  ret.Interface().(string)
	case []string:
		return ret.Interface().([]string)
	case map[string]string:
		return ret.Interface().(map[string]string)
	case map[string]interface{}:
		return ret.Interface().(map[string]interface{})
	case ReturnType:
		return ret.Interface().(ReturnType)
	default:
		fmt.Println("vtype:", vtype)
		golog.Error( "返回的类型无法处理:",vtype)
	}
	return ""

}
