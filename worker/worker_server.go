package worker

import (
	"bufio"
	"fmt"
	"flag"
	"net"
	"reflect"
	"strconv"
	"time"
	"encoding/json"
	"masterlab_socket/util"
	"masterlab_socket/area"
	"masterlab_socket/global"
	"masterlab_socket/golog"
	"masterlab_socket/protocol"
	"masterlab_socket/worker/golang"
	"github.com/BurntSushi/toml"
)



type WorkerConfigType struct {

	Loglevel     string		`toml:"loglevel"`
	SingleMode   bool	  	`toml:"single_mode"`
	Servers [][]string       	`toml:"servers"`
	ToHub []string  		`toml:"connect_to_hub"`

}


var WorkerConfig   WorkerConfigType

// 初始化worker服务
func InitWorkerServer() {

	var err error
	var filepath string
	flag.StringVar(&filepath,"worker", "worker.toml", "worker.toml's file path")
	_, err = toml.DecodeFile( filepath, &WorkerConfig )

	if  err != nil {
		fmt.Println("worker.toml.DecodeFile error:", err.Error())
		return
	}
	for _, data := range WorkerConfig.Servers {

		if len( data )<=2 {
			fmt.Println("worker.toml servers length err:" ,data )
			continue
		}
		host := data[0]
		port_str := data[1]
		worker_language  := data[2]
		port, _ := strconv.Atoi(port_str)
		if worker_language == "go" {
			go WorkerServer(host, port)
		}
	}
	time.Sleep( 1*time.Second)
	golang.InitReqHubPool( WorkerConfig.ToHub )
}


/**
 * 监听客户端连接
 */
func WorkerServer(host string, port int) {

	fmt.Println("WorkerServer :", host, port)
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(host), (port), ""})
	if err != nil {
		golog.Error("ListenTCP Exception:", err.Error())
		return
	}

	// 处理客户端连接
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			golog.Error("AcceptTCP Exception::", err.Error(), time.Now().UnixNano())
			break
		}
		// 校验ip地址
		conn.SetKeepAlive(true)
		//conn.SetDeadline(30*time.Second)
		defer conn.Close()
		//conn.SetNoDelay(false)
		golog.Info("RemoteAddr:", conn.RemoteAddr().String())

		go handleWorker(conn)

	} //end for {
}

func handleWorker(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	defer func() {
		err := recover()
		if err != nil {
			conn.Close()
			fmt.Println( "handleWorker err :", err)
		}
	}()
	for {
		_type,header_buf,data_buf,_, err :=protocol.DecodePacket( reader )
		if err != nil {
			golog.Error("handleWorker protocol.DecodePacket err:",err.Error() )
			conn.Close()
			break
		}

		if util.Int2String(int(_type)) == protocol.TypePing{
			protocolPack := new(protocol.Pack)
			protocolPack.Init()
			protocolPack.WrapResp( protocol.TypePing,"",0,0,[]byte("pong") )
			conn.Close()
			break
		}
		//fmt.Println( "HandleWorkerStr str: ",string(header_buf), string(data_buf) )
		go func(header_buf []byte, data_buf []byte,conn *net.TCPConn) {

			protocolPack := new(protocol.Pack)
			protocolPack.Init()
			req_obj, _ := protocolPack.GetReqObj( _type,header_buf, data_buf)
			Invoker(conn, req_obj)

		}( header_buf, data_buf, conn )
	}
}

func Invoker(conn *net.TCPConn, req_obj *protocol.ReqRoot) interface{} {

	task_obj := new(golang.TaskType).Init(conn, req_obj)

	invoker_ret := InvokeObjectMethod(task_obj, req_obj.Header.Cmd)
	//fmt.Println("invoker_ret", invoker_ret)
	// 判断是否需要响应数据
	if req_obj.Type == protocol.TypeReq && !req_obj.Header.NoResp {
		protocolPack := new(protocol.Pack)
		protocolPack.Init()
		invoker_ret_buf := util.Convert2Byte( invoker_ret )
		switch invoker_ret.(type) {
		case golang.ReturnType:
			tmp :=  invoker_ret.(golang.ReturnType)
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
	case golang.ReturnType:
		return ret.Interface().(golang.ReturnType)
	default:
		fmt.Println("vtype:", vtype)
		golog.Error( "返回的类型无法处理:",vtype)
	}
	return ""

}
