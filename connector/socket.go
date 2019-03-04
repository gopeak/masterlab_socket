package connector

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"
	"masterlab_socket/area"
	"masterlab_socket/global"
	"masterlab_socket/golog"
	"masterlab_socket/protocol"
	"masterlab_socket/worker"
	"masterlab_socket/util"
	"masterlab_socket/worker/golang"
	"encoding/json"
)


/**
 * 监听客户端连接
 */
func SocketConnector(ip string, port int) {

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		golog.Error("ListenTCP Exception:", err.Error())
		return
	}
	// 初始化
	golog.Debug("Game Connetor Server :", ip, port)
	//go statTick()
	listenAcceptTCP(listen)
}

/**
 *  处理客户端连接
 */
func listenAcceptTCP(listen *net.TCPListener) {

	for {
		conn, err := listen.AcceptTCP()
		defer conn.Close()
		if err != nil {
			golog.Error("AcceptTCP Exception::", err.Error())
			continue
		}
		atomic.AddInt32(&global.SumConnections, 1)
		conn.SetNoDelay(false)

		// 校验ip地址
		conn.SetKeepAlive(true)

		go handleClient(conn, area.CreateSid())
		//go handleClientMsgSingle( conn ,CreateSid() )

	} //end for {
}


func responseProcess( conn *net.TCPConn,  headerr_buf, data_buf []byte)  {

	protocolPack := new(protocol.Pack)
	protocolPack.Init()
	resp_header, err := protocolPack.GetRespHeaderObj(  headerr_buf )
	if err!=nil{
		golog.Error( "responseProcess protocolPack.GetRespHeaderObj err: ", err.Error(),string(data_buf) )
		return
	}
	//fmt.Println("responseProcess resp_obj.Data: ", string(data_buf) )

	if global.IsAuthCmd(resp_header.Cmd) {

		var ret golang.ReturnType
		//data_buf = util.TrimX001( data_buf )
		err := json.Unmarshal( data_buf ,&ret)
		if err!=nil{
			//fmt.Println("AuthCmd return json err: ", err.Error(),string(data_buf)  )
			golog.Error( "AuthCmd return json err: ", err.Error(),string(data_buf)  )
			return
		}
		//fmt.Println("AuthCmd: ", ret.Ret,string(data_buf) )
		if ret.Ret == "ok" {
			if conn != nil {
				area.ConnRegister( conn, ret.Sid )
			}
		}
	}
}

// 性能测试的呵呵检测单机效能
func handleClientMsgSingle(conn *net.TCPConn, sid string) {

	//声明一个管道用于接收解包的数据
	qps := 0 // make(chan int64, 0)

	reader := bufio.NewReader(conn)
	protocolPacket := new(protocol.Pack)
	protocolPacket.Init()
	//defer conn.Close()
	for {
		if !global.Config.Enable {
			buf,_ := protocolPacket.WrapResp( "Info", "", 0 , 200, []byte(global.DISBALE_RESPONSE) )
			conn.Write( buf )
			conn.Close()
			break
		}
		_,header, data, _,err := protocol.DecodePacket( reader )
		if err != nil {
			conn.Close()
			break
		}
		qps++
		if qps%100 == 0 {
			fmt.Println("qps: ", qps)
		}
		atomic.AddInt64(&global.Qps, 1)


		req_obj,err := protocolPacket.GetReqHeaderObj( header )
		buf,_ := protocolPacket.WrapResp( "GetUserSession", req_obj.Sid, req_obj.SeqId , 200, data )
		conn.Write( buf )

	}
}

func handleClient(conn *net.TCPConn, sid string) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	last_sid := ""
	defer area.FreeConn(conn, last_sid)
	protocolPacket := new(protocol.Pack)
	protocolPacket.Init()
	for {
		if !global.Config.Enable {
			buf,_ := protocolPacket.WrapResp( "Info", last_sid, 0 , 200, []byte(global.DISBALE_RESPONSE) )
			conn.Write( buf )
			area.FreeConn(conn, last_sid)
			break
		}
		_type,header,data,all_buf,err := protocol.DecodePacket( reader )
		if err!=nil {
			golog.Error("SocketHandle protocol.DecodePacket err : "  + err.Error())
			buf,_ := protocolPacket.WrapResp( "Error", last_sid, 0 , 500, []byte(global.ERROR_RESPONSE) )
			conn.Write( buf )
			area.FreeConn(conn, last_sid)
			return
		}
		req_obj ,err := protocolPacket.GetReqObj( _type,header,data )
		if err != nil {
			golog.Error("protocolPacket.GetReqObj err : "  + err.Error())
			area.FreeConn(conn, last_sid)
			break
		}
		last_sid = req_obj.Header.Sid
		ret, ret_err := dispatchMsg( req_obj, conn ,all_buf)
		if ret_err != nil {
			if ret < 0 {
				fmt.Println(ret_err.Error())
				continue
			}
			if ret == 0 {
				fmt.Println(ret_err.Error())
				break
			}
		}

	}
}


func directInvoker( conn *net.TCPConn, req_obj *protocol.ReqRoot ) interface{} {

	task_obj := new(golang.TaskType).Init(conn, req_obj)
	invoker_ret := worker.InvokeObjectMethod(task_obj, req_obj.Header.Cmd)
	//fmt.Println("invoker_ret", invoker_ret)
	// 判断是否需要响应数据
	if req_obj.Type == protocol.TypeReq && !req_obj.Header.NoResp {
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()

		data_buf := util.Convert2Byte( invoker_ret )
		buf,_ := protocolPacket.WrapResp( req_obj.Header.Cmd, req_obj.Header.Sid, req_obj.Header.SeqId , 200, data_buf )
		conn.Write( buf )

		if global.IsAuthCmd(req_obj.Header.Cmd) {
			var return_obj golang.ReturnType
			return_obj = invoker_ret.(golang.ReturnType)
			if return_obj.Ret == "ok" {
				if conn != nil {
					area.ConnRegister(conn, return_obj.Sid)
				}
				//fmt.Println("handleWorkerResponse AuthCmd sid: ", req_obj.Header.Cmd, return_obj.Sid )
			}
		}
	}
	return invoker_ret
}


/**
 * 根据消息类型分发处理
 */
func dispatchMsg(req_obj *protocol.ReqRoot, conn *net.TCPConn,all_buf []byte) (int, error) {

	var err error
	//  认证检查,
	if !global.IsAuthCmd(req_obj.Header.Cmd) && !area.CheckSid(req_obj.Header.Sid) {
		area.FreeConn(conn, req_obj.Header.Sid)
		err = errors.New("Auth failed!")
		return 0, err
	}


	directInvoker( conn ,req_obj )

	return 1, nil
}



func checkError(err error) {
	if err != nil {
		golog.Error(os.Stderr, "Connector error: %s", err.Error())
	}
}

func statTick() {

	timer := time.Tick(1000 * time.Millisecond)
	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		fmt.Println(time.Now().Unix(), " Connections: ", global.SumConnections, "  Qps: ", global.Qps)
	}
}

func userTick(conn *net.TCPConn) {

	timer := time.Tick(5000 * time.Millisecond)
	protocolPacket := new(protocol.Pack)
	protocolPacket.Init()
	for _ = range timer {
		buf,_ := protocolPacket.WrapResp( "ping", "", 0 , 200, util.Int64ToBytes(time.Now().Unix()) )
		_,err := conn.Write( buf )
		if err!=nil{
			golog.Error( "Socket user_tick err:",err.Error() )
			break
		}
	}
}
