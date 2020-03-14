package connector

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"masterlab_socket/area"
	"masterlab_socket/global"
	"masterlab_socket/golog"
	"masterlab_socket/lib/websocket"
	"masterlab_socket/protocol"
	"masterlab_socket/util"
	"masterlab_socket/worker"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

func (this *Connector) Websocket(ip string, port int) {

	fmt.Println("Websocket  Server bind :", ip, port)
	var addr = flag.String("addr", fmt.Sprintf(":%d", port), "http service address")

	http.Handle("/ws", websocket.Handler(this.WebsocketHandleClient))

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/web", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))
	// 初始化群组
	// worker.InitGlobalGroup()
	// http请求处理
	// worker.InitHandler()

	log.Fatal(http.ListenAndServe(*addr, nil))
}

/**
 *  处理客户端连接
 */
func (this *Connector)WebsocketHandleClient(wsconn *websocket.Conn) {

	var max_conns int32
	fmt.Println(" websocke client connect:", wsconn.RemoteAddr())
	//remoteAddr :=conn.RemoteAddr()
	atomic.AddInt32(&global.SumConnections, 1)

	// 检查是否超过最大连接数
	max_conns = int32(global.Config.Connector.MaxConections)
	if max_conns > 0 && global.SumConnections > max_conns {
		protocolJson := new(protocol.Json)
		protocolJson.Init()
		protocolJson.WrapRespErr( global.ERROR_MAX_CONNECTIONS )
		return
	}

	last_sid := ""
	// 监听客户端发送的数据
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	defer wsconn.Close()
	for {
		var buf []byte
		if err := websocket.Message.Receive(wsconn, &buf); err != nil {
			fmt.Println(" websocket.Message.Receive error:", last_sid, "  -->", err.Error())
			area.FreeWsConn(wsconn, last_sid)
			break
		}
		fmt.Println( "WebsocketHandleClient Receive: ", string(buf)  )
		req_obj, err := protocolJson.GetReqObj(buf)
		if err != nil {
			golog.Error("1.WebsocketHandle protocolJson.GetReqObj err : " + err.Error())
			continue
		}
		last_sid = req_obj.Header.Sid
		fmt.Println("req_obj.Header.Cmd: " +  req_obj.Header.Cmd)

		//fmt.Println( "WebsocketHandleClient Receive: ", string(buf)  )
		//fmt.Println( "WebsocketHandleClient req_obj header: ", req_obj.Header  )
		go func(req_obj *protocol.ReqRoot, wsconn *websocket.Conn) {

			ret, ret_err := this.wsDspatchMsg(req_obj, wsconn)
			if ret_err != nil {
				if ret < 0 {
					fmt.Println(ret_err.Error())
					return
				}
				if ret == 0 {
					fmt.Println(ret_err.Error())
					return
				}
			}

		}(req_obj, wsconn)

	}
}


func (this *Connector)wsResponseProcess(wsconn *websocket.Conn, header_buf []byte, data_buf []byte) {

	protocolPack := new(protocol.Pack)
	protocolPack.Init()
	resp_header, err := protocolPack.GetRespHeaderObj(  header_buf )
	if err!=nil{
		golog.Error( "wsResponseProcess protocolPack.GetRespHeaderObj err: ", err.Error() )
		return
	}
	fmt.Println("wsResponseProcess resp_obj.Data: ", resp_header.Cmd )

	if global.IsAuthCmd(resp_header.Cmd) {
		var ret worker.ReturnType
		data_buf = util.TrimX001( data_buf )
		err := json.Unmarshal( data_buf ,&ret)
		if err!=nil{
			fmt.Println("AuthCmd return json err: ", err.Error(), string(data_buf)  )
		}
		fmt.Println("AuthCmd: ", ret.Ret, string(data_buf) )
		if ret.Ret == "ok" {
			if wsconn != nil {
				area.WsConnRegister( wsconn, resp_header.Sid )
			}
			fmt.Println("wsResponseProcess AuthCmd sid: ", resp_header.Cmd, ret.Sid )
		}
	}
}

func Convert2Byte( invoker_ret interface{}) []byte{

	var data_buf []byte

	switch invoker_ret.(type) {      //多选语句switch
	case string:
		data_buf = []byte( invoker_ret.(string) )
	case int:
		data_buf = []byte(strconv.Itoa( invoker_ret.(int) ))
	case worker.ReturnType:
		data_buf,_ = json.Marshal( invoker_ret.(worker.ReturnType) )
	}
	return data_buf
}

func (this *Connector)wsDirectInvoker( wsconn *websocket.Conn, req_obj *protocol.ReqRoot) interface{} {

	task_obj := new(worker.TaskType).WsInit(wsconn, req_obj)
	invoker_ret := worker.InvokeObjectMethod(task_obj, req_obj.Header.Cmd)
	//fmt.Println("invoker_ret", invoker_ret)
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	data_buf := []byte("");
	switch invoker_ret.(type) {
		 case  worker.ReturnType:
			 data_buf,_ = json.Marshal( invoker_ret.(worker.ReturnType) )
		 default:
			 data_buf = util.Convert2Byte( invoker_ret )
	}
	res_buf:= protocolJson.WrapResp( req_obj ,data_buf, 200, "" )
	wsconn.Write(res_buf)
	// 判断是否需要响应数据
	if req_obj.Type == protocol.TypeReq && !req_obj.Header.NoResp {

		if global.IsAuthCmd(req_obj.Header.Cmd) {
			var return_obj worker.ReturnType
			return_obj = invoker_ret.(worker.ReturnType)
			if return_obj.Ret == "ok" {
				if wsconn != nil {
					area.WsConnRegister(wsconn, return_obj.Sid)
				}
				fmt.Println("wsDirectInvoker AuthCmd sid: ", req_obj.Header.Cmd, return_obj.Sid )
			}

		}
	}
	return invoker_ret
}


/**
 * 根据消息类型分发处理
 */
func (this *Connector)wsDspatchMsg(req_obj *protocol.ReqRoot, wsconn *websocket.Conn) (int, error) {

	var err error
	// 认证检查,
	if !global.IsAuthCmd(req_obj.Header.Cmd) && !area.CheckSid(req_obj.Header.Sid) {
		area.FreeWsConn(wsconn, req_obj.Header.Sid)
		err = errors.New("认证失败")
		return 0, err
	}
	// 判断单机模式下不需要请求worker
	this.wsDirectInvoker( wsconn ,req_obj )
	return  1, nil
}
