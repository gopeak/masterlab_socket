//
//  Hub server
//
//

package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
	"encoding/json"
	"masterlab_socket/global"
	"github.com/antonholmquist/jason"
	"masterlab_socket/protocol"
)

type Hub struct {

	Init func()

}

/**
 * 监听客户端连接
 */
func (this *Hub)Server() {

	hub_host := global.Config.Hub.Hub_host
	hub_port, _ := strconv.Atoi(global.Config.Hub.Hub_port)
	fmt.Println("Hub  Server :", hub_host, hub_port)
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(hub_host), hub_port, ""})
	if err != nil {
		LogError("Hub listenTCP Exception:", err.Error())
		return
	}
	this.listen(listen)
}

/**
 *  处理客户端连接
 */
func (this *Hub)listen(listen *net.TCPListener) {

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			LogError("AcceptTCP Exception::", err.Error(), time.Now().UnixNano())
			break
		}
		// 校验ip地址
		conn.SetKeepAlive(true)
		///defer conn.Close()
		conn.SetNoDelay(false)

		//go handleWorkerWithJson( conn  )
		go this.handleHubConn(conn)

	} //end for {

}

func (this *Hub)handleHubConn(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	defer conn.Close()
	for {
		cmd_buf,sid_buf,seq_buf,data_buf,err := protocol.HubUnPack( reader)
		if err != nil {
			//fmt.Println( "handleHubConn protocol.Unpack error: ",string(sid_buf), err.Error())
			conn.Close()
			return
		}
		if  TrimStr(string(cmd_buf))==""{
			//fmt.Println( "handleHubConn cmd empty" )
			LogError( "handleHubConn protocol.HubUnPack err: ","handleHubConn cmd empty"  )
			conn.Close()
			return
		}
		go this.workeDispath( string(cmd_buf) , string(sid_buf), string(seq_buf), data_buf, conn)

	}

}


//  Worker using REQ socket to do load-balancing
//
func (this *Hub)workeDispath(  cmd, sid, seq string,  data_buf []byte, conn *net.TCPConn) {

	//  Process messages as they arrive

	data := string( data_buf )
	api := new(Api)
	//fmt.Println( "hubWorkeDispath cmd:", cmd )

	if cmd == "GetBase" {
		ret_buf := []byte( api.GetBase() )
		write_buf,err:=protocol.HubPack( cmd,sid,seq,ret_buf )
		if err!=nil {
			LogError( "hubWorkeDispath GetBase protocol.HubPack err:", err.Error() )
			return
		}
		_,errw := conn.Write( write_buf )
		if errw!=nil {
			fmt.Println( "hubWorkeDispath GetBase conn.Write err:", errw.Error() )
		}
		return
	}
	if cmd == "GetEnableStatus" {
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( fmt.Sprintf("%d",global.AppConfig.Enable)) )
		conn.Write( write_buf )
		return
	}
	if cmd == "Enable" {
		global.AppConfig.Enable = 1
		//write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( "1" ) )
		//conn.Write( write_buf )
		return
	}
	if cmd == "Disable" {
		global.AppConfig.Enable = 0
		//write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( "1" ) )
		//conn.Write( write_buf )
		return
	}
	if cmd == "Get" {
		str,err:=Get(data)
		if( err!=nil ) {
			write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( "" ) )
			conn.Write( write_buf )
			return
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "Set" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub Set json err:",err_json.Error())
			return
		}
		key,err_key := data_json.GetString("key")
		value,err_v := data_json.GetString("value")
		expire,err_e := data_json.GetInt64("expire")
		if( err_key!=nil || err_v!=nil || err_e!=nil ){
			LogError("Hub Set json err:",err_key.Error()+err_v.Error()+err_e.Error())
			return
		}
		_,err:=Set(key,value,expire)
		if( err!=nil ) {
			LogError("Hub Set err:",err.Error())
			//write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( "0" ) )
			//conn.Write( write_buf )
		}
		//write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( "1" ) )
		//conn.Write( write_buf )
		return

	}

	if cmd == "GetSession" {
		str :=api.GetSession(data)
		fmt.Println( "api.GetSession:",str)
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "Kick" {
		ret :=api.Kick(data)
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "CreateArea" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub Set json err:",err_json.Error())
			return
		}
		id,err1 := data_json.GetString("id")
		name,err2 := data_json.GetString("name")
		if( err1!=nil || err2!=nil )  {
			LogError("Hub Set json err:",err1.Error()+err2.Error() )
			return
		}
		ret:=api.CreateArea( id, name )
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return

	}

	if cmd == "RemoveArea" {
		ret :=api.RemoveArea(data)
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "GetAreas" {
		areas_map :=api.GetAreas()
		areas_buf,_ := json.Marshal( areas_map )
		write_buf,_:=protocol.HubPack( cmd,sid,seq, areas_buf )
		conn.Write( write_buf )
		return
	}
	if cmd == "GetAreasKey" {
		areas_key :=api.GetAreasKey()
		areas_buf,_ := json.Marshal( areas_key )
		write_buf,_:=protocol.HubPack( cmd,sid,seq, areas_buf )
		conn.Write( write_buf )
		return
	}

	if cmd == "GetSidsByArea" {
		str :=api.GetSidsByArea( data )
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "AreaAddSid" {
		fmt.Println("AreaKickSid", data )
		data_buf = TrimX001( data_buf )
		var map_data map[string]string
		err_json := json.Unmarshal( data_buf ,&map_data )
		//data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub AreaAddSid json Unmarshal err:",err_json.Error())
			return
		}
		sid ,_ok1:= map_data["sid"]
		area_id ,_ok2:= map_data["area_id"]
		if( !_ok1 )  {
			LogError("Hub AreaAddSid json sid no found" )
			return
		}
		if( !_ok2 )  {
			LogError("Hub AreaAddSid json area_id no found"  )
			return
		}
		ret :=api.AreaAddSid(sid, area_id )
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}
	if cmd == "AreaKickSid" {

		data_buf = TrimX001( data_buf )
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub AreaKickSid json err:",err_json.Error())
			return
		}
		sid,err1 := data_json.GetString("sid")
		area_id,err2 := data_json.GetString("area_id")
		if( err1!=nil || err2!=nil )  {
			LogError("Hub AreaKickSid json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.AreaKickSid(sid, area_id )
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "Push" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub Push json err:",err_json.Error())
			return
		}
		to_sid,err2 := data_json.GetString("sid")
		if err2!=nil    {
			LogError("Hub Push json err:",err2.Error())
			return
		}
		fmt.Println( "hub recvice push:", string(data_buf) )
		ret := api.Push( to_sid ,sid, data_buf )

		str :="0"
		if ret {
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "BroadcastAll" {
		ret :=api.BroadcastAll(data_buf)
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}


	if cmd == "Broatcast" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub Broatcast json err:",err_json.Error())
			return
		}
		area_id,err2 := data_json.GetString("area_id")
		if(  err2!=nil )  {
			LogError("Hub data_json json err:",err2.Error() )
			return
		}
		ret := api.Broadcast( sid, area_id ,data_buf )
		str := "0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "UpdateSession" {

		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			LogError("Hub UpdateSession json err:",err_json.Error())
			return
		}
		sid,err1 := data_json.GetString("sid")
		to_data,err2 := data_json.GetString("data")
		if( err1!=nil || err2!=nil )  {
			LogError("Hub UpdateSession json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.UpdateSession(sid, to_data )
		str :="0"
		if ret{
			str = "1"
		}
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( str ) )
		conn.Write( write_buf )
		return
	}

	if cmd == "GetUserJoinedAreas" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			err_str :="Hub UpdateSession json err:"+err_json.Error()
			LogError( err_str )
			write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( "[]" ) )
			conn.Write( write_buf )
			return
		}
		sid, _ := data_json.GetString("sid")
		ret :=api.GetUserJoinedAreas(sid )

		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( ret) )
		conn.Write( write_buf )
		return

	}

	if cmd == "GetAllSession" {

		ret :=api.GetAllSession()
		write_buf,_:=protocol.HubPack( cmd,sid,seq,[]byte( ret) )
		conn.Write( write_buf )
		return

	}
}


