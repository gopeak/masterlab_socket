package worker

import (
	"masterlab_socket/golog"
	"masterlab_socket/global"
	"masterlab_socket/area"
	"masterlab_socket/lib/syncmap"
	"masterlab_socket/util"
	"masterlab_socket/protocol"
	"masterlab_socket/hub"
	"fmt"
	"time"
	"net"
	"bufio"
	"strconv"
	"github.com/robfig/cron"
	"encoding/json"
)


type Sdk struct {

	Connected bool

	HubConn *net.TCPConn

	ReqType string

	ReqHeader *protocol.ReqHeader

	Data []byte

}

type PushReqHub struct {
	Sid bool
	Msg string
	Info map[string]string
}

type AfterWorkCallback func(   resp_buf string ) (string)

var ReqSeqCallbacks *syncmap.SyncMap

var ReqHubConns  =  make( []*net.TCPConn, 0 )

var InitialCap  int

var ToHub []string


func (sdk *Sdk) Init( _type string, req_header *protocol.ReqHeader, data []byte ) *Sdk{

	sdk.ReqHeader = req_header
	sdk.ReqType   = _type
	sdk.Data = data
	sdk.Connected = false
	return sdk
}

func (sdk *Sdk) InitCmd(cmd string,sid string,reqid int,data []byte) *Sdk{

	req_header_obj := protocol.ReqHeader{}
	req_header_obj.Cmd = cmd
	req_header_obj.SeqId = reqid
	req_header_obj.Sid = sid

	req_obj := &protocol.ReqRoot{}
	req_obj.Header = req_header_obj
	req_obj.Type = protocol.TypeReq
	req_obj.Data = data
	sdk.Data = data
	sdk.Connected = false
	return sdk
}

// 数据连接
func (sdk *Sdk) connect( ) bool{

	if sdk.HubConn!=nil {
		return true
	}
	hub_host := ToHub[0]
	hub_port_str := ToHub[1]
	ip_port := hub_host + ":" + hub_port_str

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ip_port)
	hubconn, err_req := net.DialTCP("tcp", nil, tcpAddr)
	if( err_req!=nil ){
		sdk.HubConn=nil
		return false
	}
	sdk.HubConn = hubconn
	return true

}


func   InitReqHubPool( to_hub []string  ) {

	// create a factory() to be used with channel based pool
	ReqSeqCallbacks = syncmap.New()

	InitialCap  = 10
	fmt.Println( "global.Config.ToWorker",global.Config.ToWorker.Servers)
	ToHub = to_hub
	factory    := func() (*net.TCPConn, error) {

		ip_port := to_hub[0] + ":" + to_hub[1]

		tcpAddr, _ := net.ResolveTCPAddr("tcp4", ip_port)
		hubconn, err_req := net.DialTCP("tcp", nil, tcpAddr)
		//fmt.Println( "InitConnectionHubPool hubconn ", hubconn )

		return hubconn,err_req
	}
	for i := 0; i < InitialCap; i++ {

		var err_req error
		conn, err_req:= factory()
		if( err_req!=nil ) {
			golog.Error( "InitConnectionHubPool hubconn  err:", err_req.Error() )
			continue
		}
		ReqHubConns = append( ReqHubConns, conn )
		go handleReqHubResponse( conn )
	}
}


// 侦听Hub server返回的数据，然后回调worker的函数
func  handleReqHubResponse(conn *net.TCPConn) {
	time.Sleep( 2*time.Second)
	reader := bufio.NewReader(conn)
	defer func() {
		err := recover()
		if err != nil {
			conn.Close()
			fmt.Println( "ReadHubResp err :", err)
		}
	}()
	for {
		cmd_buf,sid_buf,seq_buf,data_buf,err := protocol.HubUnPack( reader)
		if err != nil {
			fmt.Println( "handleReqHubResponse protocol.Unpack error: ",string(sid_buf), err.Error())
			conn.Close()
			break
		}

		callback_key:=string(cmd_buf) + string(seq_buf)
		if string(seq_buf)=="0"{
			ReqSeqCallbacks.Delete( callback_key )
			continue

		}
		/*
		fmt.Println( "callback_sid:", string(sid_buf) )
		fmt.Println( "callback_key:", callback_key )
		fmt.Println( "callback data:", string(data_buf)  )
		*/
		_item,ok := ReqSeqCallbacks.Get( callback_key )
		if( ok ) {
			callback := _item.( AfterWorkCallback )
			fmt.Println( "callback func :", callback  )
			callback( string(data_buf) )
			ReqSeqCallbacks.Delete( callback_key )
		}
	}
}


// 向Hub请求数据并监听返回,该请求将会阻塞除非等待返回超时
func (sdk *Sdk) ReqHubAsync( req_cmd string , data []byte ,handler AfterWorkCallback  ) (string,bool) {

	seq_id := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)
	req_buf,err:= protocol.HubPack( req_cmd,"",seq_id, data )
	if err != nil {
		golog.Error( "ReqHubAsync protocol.HubPack err:" , err.Error() )
		return err.Error(),false
	}
	index := util.RandInt64(0, int64(len(ReqHubConns)))
	req_hub_conn  := ReqHubConns[index]

	if( req_hub_conn==nil  ){
		golog.Error( "req_hub_conn is nil "  )
		return "", false
	}
	callback_key:=req_cmd + seq_id
	ReqSeqCallbacks.Set( callback_key, handler )
	_,err = req_hub_conn.Write( req_buf )
	if err!=nil {
		golog.Error( "ReqHubAsync req_hub_conn.Write err:" , err.Error() )
		return err.Error() ,false
	}
	return "ok",true
}


// 向Hub请求数据并监听返回,该请求将会阻塞除非等待返回超时
func (sdk *Sdk) ReqHub( req_cmd string , data []byte ) (string,bool) {

	seq_id := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)
	req_buf,err:= protocol.HubPack( req_cmd, sdk.ReqHeader.Sid, seq_id, data )
	if err != nil {
		golog.Error( "ReqHub protocol.HubPack err:" , err.Error() )
		return err.Error(),false
	}

	sdk.connect()
	_,err=sdk.HubConn.Write( req_buf )
	if( err!=nil ) {
		return "sdk.HubConn.Write err",false
	}
	reader := bufio.NewReader(sdk.HubConn)
	defer func() {
		err := recover()
		if err != nil {
			sdk.HubConn.Close()
			fmt.Println( "ReqHub err :", err)
		}
	}()
	for {

		cmd_buf,sid_buf,seq_buf,data_buf,err := protocol.HubUnPack( reader)
		if err != nil {
			fmt.Println( "ReqHub protocol.Unpack error: ", string(sid_buf), err.Error() )
			sdk.HubConn.Close()
			break
		}

		select {

		case <- time.After(5 * time.Second):
			return "timeout 5 second ",false

		default:
			if string(cmd_buf)+string(seq_buf) == req_cmd+seq_id{
				//fmt.Println( "ReqHub:",string(data_buf))
				// 如果服务返回错误
				sdk.HubConn.Close()
				return string(data_buf),true
			}
		}
	}

	return "",false
}

func (sdk *Sdk) PushHub( req_cmd string , data []byte ) bool {

	req_buf,err:= protocol.HubPack( req_cmd, sdk.ReqHeader.Sid, "", data )
	if err != nil {
		golog.Error( "ReqHub protocol.HubPack err:" , err.Error() )
		return false
	}
	sdk.connect()
	_,err=sdk.HubConn.Write( req_buf )
	if( err!=nil ) {
		return false
	}
	return true
}



// 获取服务器的根路径
func (sdk *Sdk)  GetBase() string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetBase()
	}

	ret,ok :=sdk.ReqHub( "GetBase",[]byte("")  )
	if ok {
		return ret
	}
	return ""

}

// 获取服务启用状态
func (sdk *Sdk) GetEnableStatus() bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetEnableStatus()
	}
	ret,ok:= sdk.ReqHub( "GetEnableStatus",[]byte("")  )
	if( !ok ){
		return false
	}
	if( ret=="1" ){
		return true
	}else{
		return false
	}
}

func (sdk *Sdk) Enable() bool {
	if( global.SingleMode ) {
		global.AppConfig.Enable = 1
		return true
	}
	return sdk.PushHub( "Enable", []byte("") )
}

func (sdk *Sdk) Disable() bool {

	if( global.SingleMode ) {
		global.AppConfig.Enable = 0
		return true
	}
	return sdk.PushHub( "Disable",[]byte("") )
}

func (sdk *Sdk) AddCron(expression string, exefnc func()) bool {

	if cron, ok := global.Crons[expression]; ok {
		golog.Info("cron exist :", cron)
		return false
	}
	c := cron.New()
	c.AddFunc(expression, exefnc)
	c.Start()
	global.Crons[expression] = c
	return true

}

func (sdk *Sdk) RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}
	return true
}

func (sdk *Sdk) Get(key string) string {
	if( global.SingleMode ) {
		str,err:=hub.Get(key)
		if err!=nil {
			golog.Error("Redis Get err:",err.Error())
			return ""
		}
		return str
	}
	ret,ok := sdk.ReqHub( "Get",[]byte(key) )
	if( !ok ) {
		return ""
	}
	return ret
}

func (sdk *Sdk) Set(key string, value string,expire int) bool {

	if( global.SingleMode ) {
		ret,err:=hub.Set(key,value,expire)
		if err!=nil {
			golog.Error("Redis Set err:",err.Error())
			return false
		}
		return ret
	}
	json:=fmt.Sprintf(`{"key":"%s","value":"%s","expire":%d}`,key,value,expire)
	ret:= sdk.PushHub( "Set",[]byte(json) )
	return ret
}

// 该方法仅在单机模式下调用
func (sdk *Sdk) GetSessionType(sid string) *area.Session  {

	session,exist := global.UserSessions.Get(sid)
	if !exist {
		return nil
	}
	return session.(*area.Session)
}

func (sdk *Sdk) GetSession(sid string)  string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetSession( sid )
	}
	ret,ok := sdk.ReqHub( "GetSession",[]byte(sid)   )
	if !ok{
		return ""
	}
	return ret

}

func (sdk *Sdk) Kick(sid string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Kick( sid )
	}
	return sdk.PushHub( "Kick",[]byte(sid) )
}

func (sdk *Sdk) CreateArea(id string, name string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.CreateArea( id,name )
	}
	json:=fmt.Sprintf(`{"id":"%s","name":"%s","expire":%d}`,id,name)
	return sdk.PushHub( "CreateArea",[]byte(json)  )

}

func (sdk *Sdk) RemoveArea(id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.RemoveArea( id )
	}
	return sdk.PushHub( "RemoveArea",[]byte(id) )
}

func (sdk *Sdk) GetAreas() map[string]string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetAreas(  )
	}
	var areas map[string]string
	ret,ok:= sdk.ReqHub( "GetAreas",[]byte("")  )
	//fmt.Println( "sdk ReqHub:",ret )
	if !ok {
		return areas
	}
	err:=json.Unmarshal( []byte(ret), &areas )
	if err!=nil {
		fmt.Println( "sdk GetAreas err :",err.Error() )
	}
	return areas
}

func (sdk *Sdk) GetAreasStr() string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		buf,err := json.Marshal(  api.GetAreas(  ) )
		if err!=nil {
			return "{}"
		}
		return string(buf)
	}
	ret,ok:= sdk.ReqHub( "GetAreas",[]byte("")  )
	if !ok {
		return "{}"
	}
	return ret
}

func (sdk *Sdk) GetAreasKey()  []string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetAreasKey(  )
	}
	var areas = make([]string,0 )
	ret,_ := sdk.ReqHub( "GetAreasKey",[]byte("")  )
	json.Unmarshal( []byte(ret), areas )

	return areas
}



func (sdk *Sdk) GetSidsByArea(channel_id string) string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetSidsByArea( channel_id )
	}
	ret,ok :=  sdk.ReqHub( "GetSidsByArea",[]byte(channel_id)  )
	if( !ok ) {
		return "{}"
	}
	return ret

}

func (sdk *Sdk) AreaAddSid(sid string, area_id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.AreaAddSid( sid, area_id  )
	}
	json:=fmt.Sprintf(`{"sid":"%s","area_id":"%s"}`,sid, area_id )
	return sdk.PushHub( "AreaAddSid",[]byte(json) )

}

func (sdk *Sdk) AreaKickSid( sid string, area_id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.AreaKickSid( sid, area_id  )
	}
	json:=fmt.Sprintf(`{"sid":"%s","area_id":"%s"}`,sid, area_id )
	return sdk.PushHub( "AreaKickSid",[]byte(json))

}

func (sdk *Sdk) Push( from_sid string ,to_sid string , to_data  []byte ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Push ( from_sid,to_sid, to_data  )
	}

	return sdk.PushHub( "Push",to_data )

}

func (sdk *Sdk) PushBySids(from_sid string,to_sids []string, data []byte) bool {

	for _,to_sid:=   range to_sids {
		sdk.Push(from_sid, to_sid, data )
	}
	return true

}

func (sdk *Sdk) Broatcast(sid string ,area_id string,  data []byte ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Broadcast( sid,area_id, data  )
	}
	return sdk.PushHub( "Broatcast",data )

}


func (sdk *Sdk) BroadcastAll( msg []byte ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.BroadcastAll( msg )
	}
	return sdk.PushHub( "BroadcastAll", msg )

}


func (sdk *Sdk) UpdateSession( sid string, data string ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.UpdateSession( sid, data )
	}
	json:=fmt.Sprintf(`{"sid":"%s","data":"%s"}`,sid, data )
	return sdk.PushHub( "UpdateSession",[]byte(json) )

}

func (sdk *Sdk)GetUserJoinedAreas(sid string) string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetUserJoinedAreas(sid)
	}

	ret,ok :=sdk.ReqHub( "GetUserJoinedAreas",[]byte(sid) )
	if ok {
		return ret
	}
	return ""

}

func (sdk *Sdk)GetAllSession( ) string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetAllSession()
	}

	ret,ok :=sdk.ReqHub( "GetAllSession",[]byte(""))
	if ok {
		return ret
	}
	return ""

}
