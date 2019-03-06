package masterlab_socket

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"encoding/json"
	"masterlab_socket/area"
	"masterlab_socket/global"
	"masterlab_socket/golog"
	"masterlab_socket/lib/websocket"
	"masterlab_socket/protocol"
	"masterlab_socket/worker"
	"masterlab_socket/util"
)

func WebsocketConnector(ip string, port int) {

	golog.Info("Websocket Connetor bind :", ip, port)

	var addr = flag.String("addr", fmt.Sprintf(":%d", port), "http service address")

	http.Handle("/ws", websocket.Handler(WebsocketHandleClient))

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/web/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))
	// 初始化群组
	worker.InitGlobalGroup()
	// http请求处理
	worker.InitHandler()

	log.Fatal(http.ListenAndServe(*addr, nil))

}

/**
 *  处理客户端连接
 */
func WebsocketHandleClient(wsconn *websocket.Conn) {

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
	configAddr := global.GetRandWorkerAddr()
	fmt.Println("ip_port:", configAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", configAddr)
	checkError(err)
	req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//defer req_conn.Close()
	checkError(err)
	go wsHandleWorkerResponse(wsconn, req_conn)
	last_sid := ""
	// 监听客户端发送的数据
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	defer wsconn.Close()
	for {
		var buf []byte
		if err = websocket.Message.Receive(wsconn, &buf); err != nil {
			fmt.Println(" websocket.Message.Receive error:", last_sid, "  -->", err.Error())
			area.FreeWsConn(wsconn, last_sid)
			break
		}
		req_obj, err := protocolJson.GetReqObj(buf)
		if err != nil {
			golog.Error("1.WebsocketHandle protocolJson.GetReqObj err : " + err.Error())

			continue
		}
		last_sid = req_obj.Header.Sid
		fmt.Println("req_obj.Header.Cmd: " +  req_obj.Header.Cmd)

		fmt.Println( "WebsocketHandleClient Receive: ", string(buf)  )
		//fmt.Println( "WebsocketHandleClient req_obj header: ", req_obj.Header  )
		go func(req_obj *protocol.ReqRoot, wsconn *websocket.Conn, req_conn *net.TCPConn) {

			ret, ret_err := wsDspatchMsg(req_obj, wsconn, req_conn)
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

		}(req_obj, wsconn, req_conn)

	}
}

func wsHandleWorkerResponse(wsconn *websocket.Conn, req_conn *net.TCPConn) {

	reader := bufio.NewReader(req_conn)
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	for {
		_type,header_buf,data_buf,_, err := protocol.DecodePacket( reader )
		if err != nil {
			golog.Error( "wsHandleWorkerResponse protocol.DecodePacket err: ", err.Error() )
			req_conn.Close()
			break
		}
		fmt.Println("wsHandleWorkerResponse  data :", _type, string(header_buf), string(data_buf) )
		wsResponseProcess( wsconn,header_buf, data_buf  )
		buf :=protocolJson.WrapResp( header_buf,data_buf,200,"" )
		fmt.Println( "protocolJson.WrapResp:",string(buf) )
		go wsconn.Write( buf )
	}
}

func wsResponseProcess(wsconn *websocket.Conn, header_buf []byte, data_buf []byte) {

	protocolPack := new(protocol.Pack)
	protocolPack.Init()
	resp_header, err := protocolPack.GetRespHeaderObj(  header_buf )
	if err!=nil{
		golog.Error( "wsResponseProcess protocolPack.GetRespHeaderObj err: ", err.Error() )
		return
	}
	fmt.Println("handleWorkerResponse resp_obj.Data: ", resp_header.Cmd )

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

func wsDirectInvoker( wsconn *websocket.Conn, req_obj *protocol.ReqRoot) interface{} {

	task_obj := new(worker.TaskType).WsInit(wsconn, req_obj)
	invoker_ret := worker.InvokeObjectMethod(task_obj, req_obj.Header.Cmd)
	//fmt.Println("invoker_ret", invoker_ret)
	// 判断是否需要响应数据
	if req_obj.Type == protocol.TypeReq && !req_obj.Header.NoResp {
		protocolJson := new(protocol.Json)
		protocolJson.Init()
		data_buf := util.Convert2Byte( invoker_ret )
		resp_obj:= protocolJson.WrapRespObj( req_obj ,data_buf, 200 )
		buf,_ := json.Marshal(resp_obj)
		wsconn.Write( buf )

		if global.IsAuthCmd(req_obj.Header.Cmd) {
			var return_obj worker.ReturnType
			return_obj = invoker_ret.(worker.ReturnType)
			if return_obj.Ret == "ok" {
				if wsconn != nil {
					area.WsConnRegister(wsconn, return_obj.Sid)
				}
				fmt.Println("wsHandleWorkerResponse AuthCmd sid: ", req_obj.Header.Cmd, return_obj.Sid )
			}
		}
	}
	return invoker_ret
}


/**
 * 根据消息类型分发处理
 */
func wsDspatchMsg(req_obj *protocol.ReqRoot, wsconn *websocket.Conn, req_conn *net.TCPConn) (int, error) {

	var err error
	// 认证检查,
	if !global.IsAuthCmd(req_obj.Header.Cmd) && !area.CheckSid(req_obj.Header.Sid) {
		area.FreeWsConn(wsconn, req_obj.Header.Sid)
		err = errors.New("认证失败")
		return 0, err
	}
	// 判断单机模式下不需要请求worker
	if global.SingleMode {
		wsDirectInvoker( wsconn ,req_obj )
		return  1, nil
	}

	protocolPack := new(protocol.Pack)
	protocolPack.Init()
	buf,_ := protocolPack.WrapReqWithHeader( &req_obj.Header , req_obj.Data)
	// 提交给worker
	if req_conn != nil {
		go req_conn.Write(buf)
	}

	return 1, nil
}
