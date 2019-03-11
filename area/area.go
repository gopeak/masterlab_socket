/**
 *  场景管理
 *  创建多个channel,一个channel对应一个publisher,chanel从Hub订阅消息后分发给客户端
 *
 */

package area

import (
	"fmt"
	"net"
	"time"
	"sync/atomic"
	"encoding/json"
	"math/rand"
	"masterlab_socket/global"
	"masterlab_socket/golog"
	"masterlab_socket/lib/websocket"
	"masterlab_socket/lib/syncmap"
	"masterlab_socket/protocol"
	"masterlab_socket/util"
)

// 所有的场景名称列表
var Areas = make([]string, 0, 1000)

// 场景集合
var AreasMap *syncmap.SyncMap

// 一个全局的场景
var GlobalArea   *AreaType

// 所有的用户连接对象
var AllConns *syncmap.SyncMap
var AllWsConns *syncmap.SyncMap

// 用户加入过的场景列表
var UserJoinedAreas *syncmap.SyncMap

// 用户结构
type Session struct {
	IP string
	User string
	LoggedIn bool
	KickOut  bool
	Sid string
	ConnectTime int64
	PacketTime  int64
}

type AreaType struct {
	// 唯一标识符
	Id string
	// 场景名称
	Name string
	// 当前场景包含的socket连接对象
	Conns *syncmap.SyncMap
	// 当前场景包含的websocket连接对象
	WsConns *syncmap.SyncMap

	CreateTime int64

}


// 预创建多个场景
func InitConfig() {

	AreasMap   = syncmap.New()
	AllConns   = syncmap.New()
	AllWsConns = syncmap.New()

	for _, area_id := range global.Config.Area.Init_area {
		Create(area_id, area_id)
	}
	GlobalArea = new(AreaType)
	GlobalArea.Id = "global"
	GlobalArea.Name = "全局场景"
	GlobalArea.Conns = syncmap.New()
	GlobalArea.WsConns = syncmap.New()
	GlobalArea.CreateTime = time.Now().Unix()
	AreasMap.Set("global",GlobalArea)
}

// 获取场景列表
func Gets(  ) map[string]string{

	var areas_map map[string]string
	areas_map = make(map[string]string)  //字典的创建
	var area *AreaType
	for item := range AreasMap.IterItems(){
		area = item.Value.((*AreaType))
		areas_map[item.Key] = area.Name
	}
	fmt.Println( "area Gets:", areas_map )
	return areas_map
}

// 创建一个场景
func Create(area_id string, name string) {

	Areas = append(Areas, area_id)
	area := new(AreaType)
	area.Id = area_id
	area.Name = name
	area.WsConns = syncmap.New()
	area.Conns = syncmap.New()
	area.CreateTime = time.Now().Unix()
	AreasMap.Set( area_id,area)
}

// 创建一个场景
func Get( area_id string ) *AreaType{

	v,ok := AreasMap.Get(area_id)
	if ok {
		return v.(*AreaType)
	}
	return nil
}
// 删除一个场景
func Remove(id string) {
	// 1.删除名称
	for index, elem := range Areas {
		if elem==id {
			Areas = append(Areas[:index],Areas[index+1:]...)
			return
		}
	}
	// 删除场景对象
	AreasMap.Delete( id )
}

// 检查是否已经创建了场景
func CheckExist(area_id string) bool {
	return AreasMap.Has(area_id)
}

func AddSid(sid string, area_id string) bool {

	area_id = util.TrimStr( area_id )
	sid = util.TrimStr( sid )
	exist := CheckExist(area_id)
	//fmt.Println( area_id," CheckChannelExist:", exist )
	if !exist {
		return false
	}
	user_conn := GetConn( sid )
	user_wsconn :=  GetWsConn( sid )
	fmt.Println( "AreaAddSid user_conn:",sid, user_conn )
	// 会话如果属于socket
	if user_conn != nil {
		Subscribe(area_id, user_conn, sid)
	}
	// 会话如果属于websocket
	if user_wsconn != nil {
		WsSubscribe( area_id, user_wsconn, sid )
	}
	// 该用户加入过的场景列表
	var userJoinedChannels = make([]string, 0, 1000)
	tmp, ok := UserJoinedAreas.Get(sid)
	if ok {
		userJoinedChannels = tmp.([]string)
	}
	userJoinedChannels = append(userJoinedChannels, area_id)
	UserJoinedAreas.Set(sid, userJoinedChannels)
	//}
	return true

}

/**
 *  socket连接 加入到场景中
 */
func Subscribe(area_id string, conn *net.TCPConn, sid string) {

	area := Get( area_id )
	if( area==nil  ) {
		golog.Error( "Area  ",area_id," no exist! "  )
		return
	}else{
		if( area.Conns.Size()<=0 ){
			area.Conns = syncmap.New()
		}
		if  !area.Conns.Has(sid) {
			area.Conns.Set(sid, conn)
		}
		AreasMap.Set( area_id,area )
	}
}

/**
 *  websocket连接 加入到场景中
 */
func WsSubscribe(area_id string, ws *websocket.Conn, sid string) {

	area := Get( area_id )
	if( area==nil ) {
		golog.Error( "Area  ",area_id," no exist! "  )
		return
	}else{
		if( area.WsConns.Size()<=0 ){
			area.WsConns = syncmap.New()
		}
		if  !area.WsConns.Has(sid) {
			area.WsConns.Set(sid, ws)
		}
		AreasMap.Set( area_id,area )
	}
}


func GetSids( area_id string) []string {
	ret := make([]string,0)
	area := Get( area_id )
	if( area!=nil ){
		for tmp := range area.Conns.IterItems(){
			ret=append(ret,tmp.Key)
		}
		for tmp := range area.WsConns.IterItems(){
			ret=append(ret,tmp.Key)
		}
	}
	return ret
}



/**
 *  检查用户是否加入到场景中
 */
func CheckUserJoined(area_id string, sid string) bool {

	area := Get( area_id )
	if( area!=nil ) {
		if  area.Conns.Has(sid) {
			return true
		}
		if  area.WsConns.Has(sid) {
			return true
		}
	}
	return false

}


/**
 *  用户退出某个场景
 */
func UnSubscribe(area_id string, sid string) {

	area := Get( area_id )
	if( area!=nil ) {
		area.Conns.Delete( sid )
		area.WsConns.Delete( sid )
		AreasMap.Set( area_id, area )
	}

}

// 用户退出所有场景
func UserUnSubscribe(user_sid string) {

	for index, _ := range Areas {
		UnSubscribe(Areas[index], user_sid)
	}
	UnSubGlobal( user_sid )
}

/**
 *  在场景中广播消息
 */
func Broatcast( sid string,area_id string, msg []byte ) {

	area := Get( area_id )
	if( area==nil ) {
		golog.Error("AreasMap no found :",area_id)
		return
	}
	var conn *net.TCPConn
	protocolPacket := new(protocol.Pack)
	protocolPacket.Init()
	// socket部分
	for item := range area.Conns.IterItems() {
		conn = item.Value.(*net.TCPConn)
		buf,_ := protocolPacket.WrapBroatcastResp( area_id, sid, msg  )
		//fmt.Println( "Broatcast:",  string(buf) )
		n,err:=conn.Write( buf )
		if err!=nil {
			golog.Error("Broatcast conn.Write err :",err.Error()," expect ", len(buf),", but only write:",n )
		}
	}

	var wsconn *websocket.Conn
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	for item := range area.WsConns.IterItems() {
		wsconn = item.Value.(*websocket.Conn)
		buf, _ := json.Marshal(protocolJson.WrapBroatcastRespObj( area_id, sid, msg) )
		_,err:= wsconn.Write( buf )
		if err!=nil {
			golog.Error("Broatcast wsconn.Write err: ", err.Error() )
		}
	}
}

/**
 *  在场景中广播消息
 */
func BroatcastGlobal( sid string, msg []byte ) {

	var conn *net.TCPConn
	//fmt.Println("场景里有:", GlobalArea.Conns.Size(),"个conn连接")
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	for item := range GlobalArea.Conns.IterItems() {
		conn = item.Value.(*net.TCPConn)
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,_ := protocolPacket.WrapBroatcastResp( "global", sid, msg  )
		conn.Write( buf )
	}
	//fmt.Println("广播里有:", GlobalArea.Conns.Size(),"个ws连接")
	var wsconn *websocket.Conn
	for item := range GlobalArea.WsConns.IterItems() {
		wsconn = item.Value.(*websocket.Conn)
		buf, _ := json.Marshal(protocolJson.WrapBroatcastRespObj( "global", sid, msg) )
		go wsconn.Write( buf )
	}
}

func UnSubGlobal( sid string ) {

	GlobalArea.Conns.Delete( sid )
	GlobalArea.WsConns.Delete( sid )

}

/**
 *  点对点发送消息
 */
func Push(  to_sid string ,from_sid string,to_data []byte ) {

	conn :=  GetConn(to_sid)
	if( conn!=nil ) {
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,err := protocolPacket.WrapPushResp(  to_sid, from_sid,to_data )
		if err!=nil {
			fmt.Println( "protocolPacket.WrapPushResp:",err.Error() )
		}
		_,err =conn.Write( buf )
		if err!=nil {
			fmt.Println( "Push conn.Write err:",err.Error() )
		}
		return
	}

	ws:=GetWsConn(to_sid)
	if( ws!=nil ) {
		protocolJson := new(protocol.Json)
		protocolJson.Init()
		buf  :=  protocolJson.WrapPushResp( to_sid, from_sid, to_data )
		fmt.Println( "push, to_sid:", to_sid , string(buf))
		_,err:=ws.Write( buf )
		if err!=nil {
			fmt.Println( "wsconn.Write err:",err.Error() )
		}
		return
	}
}




func GetConn(sid string) *net.TCPConn {

	conn, ok := AllConns.Get(sid)
	if !ok {
		return nil
	} else {
		return conn.(*net.TCPConn)
	}
}

func DeleteConn(sid string) {

	AllConns.Delete(sid)

}

func GetWsConn(sid string) *websocket.Conn {
	wsconn, ok := AllWsConns.Get(sid)
	if !ok {
		return nil
	} else {
		return wsconn.(*websocket.Conn)
	}
}

func DeleteWsConn(sid string) {

	AllWsConns.Delete(sid)

}

func DeleteUserssion(sid string) {

	global.UserSessions.Delete(sid)

}

func ConnRegister(conn *net.TCPConn, sid string) {

	//SubscribeChannel("area-global", conn, user_sid)

	AllConns.Set( sid, conn )

	_, ok := global.UserSessions.Get(sid)
	if !ok {
		data := &Session{
			conn.RemoteAddr().String(),
			"{}",
			true,  // 登录成功
			false, // 是否被踢出
			sid,
			time.Now().Unix(), //加入时间
			time.Now().Unix(),
		}
		global.UserSessions.Set(sid, data)
	}

}

func WsConnRegister(ws *websocket.Conn, user_sid string) {

	golog.Debug("user_sid: ", user_sid)
	//SubscribeWsChannel("area-global", ws, user_sid)

	AllWsConns.Set( user_sid, ws )

	_, ok := global.UserSessions.Get(user_sid)
	if !ok {
		data := &Session{
			ws.RemoteAddr().String(),
			"{}",
			true,  // 登录成功
			false, // 是否被踢出
			user_sid,
			time.Now().Unix(), //加入时间
			time.Now().Unix(),
		}
		global.UserSessions.Set(user_sid, data)
	}

}


func DeleteSession(sid string) {

	global.UserSessions.Delete(sid)
}

func DeleteUserJoinedAreas(sid string) {

	UserJoinedAreas.Delete(sid)

}

func FreeConn(conn *net.TCPConn, sid string) {

	conn.Close()
	golog.Warn("Sid closing:", sid)
	DeleteConn(sid)
	DeleteSession(sid)
	DeleteUserJoinedAreas(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	UserUnSubscribe(sid)
	AllConns.Delete( sid )
}

func FreeWsConn(ws *websocket.Conn, sid string) {

	//ws.Write([]byte{'E', 'O', 'F'})
	ws.Close()
	golog.Warn("Sid closing:", sid)
	DeleteWsConn(sid)
	DeleteSession(sid)
	DeleteUserJoinedAreas(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	UserUnSubscribe(sid)
}

/**
 * 检查
 */
func CheckSid(sid string) bool {

	return true
	_, exist := global.UserSessions.Get(sid)
	return exist
}

func CreateSid() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sid := fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(999999))
	return sid
}
