package main

import (
	"github.com/robfig/cron"
	json_orgin "encoding/json"
	"masterlab_socket/global"
	"os"
	"strings"
	_"masterlab_socket/lib/websocket"
	"masterlab_socket/protocol"
)

type Api struct {

	Init func()

}

// 获取服务器的根路径
func (api *Api)GetBase() string {

	dir, err:= os.Getwd()
	if err != nil {
		LogError("GetBase Error ", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}


func (api *Api)GetEnableStatus() bool {
	if global.AppConfig.Enable <= 0 {
		return false
	} else {
		return true
	}
}

func (api *Api)Enable() bool {

	global.AppConfig.Enable = 1
	return true
}

func (api *Api)Disable() bool {

	global.AppConfig.Enable = 0
	return true
}

func (api *Api)AddCron(expression string, exefnc func()) bool {

	if cron, ok := global.Crons[expression]; ok {
		LogInfo("cron exist :", cron)
		return false
	}
	c := cron.New()
	c.AddFunc(expression, exefnc)
	c.Start()
	global.Crons[expression] = c
	return true
}

func (api *Api)RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true
}

func (api *Api)Get(key string) bool {

	return true
}

func (api *Api)Set(key string, value string) bool {

	return true
}

func (api *Api)GetSession(sid string) string {
	session,exist := global.UserSessions.Get(sid)
	if !exist {
		return "{}"
	}
	str,err := json_orgin.Marshal(session)
	if( err!=nil){
		LogError("Api GetSession json Marshal err:",err.Error())
		return "{}"
	}
	return string(str)
}

func (api *Api)Kick(sid string) bool {

	protocolPacket := new(protocol.Pack)
	protocolPacket.Init()

	user_conn := AreaGetConn(sid)
	if user_conn != nil {
		// 通知消息退出
		buf,_ := protocolPacket.WrapRespErr( "kicked" )
		user_conn.Write( buf )
		AreaFreeConn(user_conn,sid )
	}

	user_wsconn := AreaGetWsConn(sid)
	if user_wsconn != nil {
		// 通知消息退出
		protocolJson:= new(protocol.Json)
		protocolJson.Init()
		go user_wsconn.Write( protocolJson.WrapRespErr("kicked") )
		AreaFreeWsConn( user_wsconn,sid)
	}
	AreaUserUnSubscribe(sid)
	AreaDeleteUserssion(sid)

	return true
}

func (api *Api)CreateArea(id string, name string) bool {

	AreaCreate(id, name)
	return true
}

func (api *Api)RemoveArea(id string) bool {

	AreaRemove(id)
	return true
}

func (api *Api)GetAreas() map[string]string {

	return AreaGets()
}

func (api *Api)GetAreasKey() []string {

	return Areas
}

func (api *Api)GetSidsByArea(channel_id string) string {

	buf,err:= json_orgin.Marshal(AreaGetSids(channel_id))
	if err!=nil {
		return string(buf)
	}else{
		return "[]"
	}
}

func (api *Api)AreaAddSid(sid string, area_id string) bool {

	return  AreaAddSid( sid , area_id )
}

func (api *Api)AreaKickSid(sid string, area_id string) bool {

	AreaUnSubscribe( area_id,sid)
	return true
}

func (api *Api)Push( to_sid, from_sid   string , data_buf []byte ) bool {
	AreaPush( to_sid, from_sid, data_buf )
	return true
}

func (api *Api)PushBySids(from_sid string,to_sids []string, data_buf []byte ) bool {

	for _,to_sid:=   range to_sids {
		AreaPush( to_sid, from_sid, data_buf )
	}
	return true
}

func (this *Api) Broadcast( sid string, area_id string, msg []byte) bool {

	AreaBroatcast( sid, area_id ,msg)
	return true
}

func (this *Api) UpdateSession( sid string, data string ) bool {

	tmp, user_session_exist := global.UserSessions.Get(sid)
	var user_session *Session
	if user_session_exist {
		user_session = tmp.(*Session)
		user_session.User = data
		global.UserSessions.Set(sid, user_session)
	}
	return true
}

func (api *Api)BroadcastAll( msg []byte ) bool {
	AreaBroatcastGlobal("GM",msg )
	return true
}

func (api *Api)GetUserJoinedAreas( sid string ) string {

	buf,err:=json_orgin.Marshal(AreaGetSids(sid))
	if( err!=nil ) {
		return  "[]"
	}
	return  string( buf )
}

func (api *Api)GetAllSession( ) string {

	var UserSessions = map[string]*Session{}
	for item := range global.UserSessions.IterItems() {
		UserSessions[item.Key] = item.Value.(*Session)
	}
	ret, _ := json_orgin.Marshal(UserSessions)
	return  string(ret)
}

