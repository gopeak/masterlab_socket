package main

import (
	"bufio"
	"fmt"
	"masterlab_socket/protocol"
	"net"
	"time"
)

func main() {

	go ServerHub("0.0.0.0", 7004)
	time.Sleep(3 * time.Second)
	go client_hub()
	select {

	}
}


/**
 * 监听客户端连接
 */
func ServerHub(ip string, port int) {


	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		fmt.Println("ListenTCP Exception:", err.Error())
		return
	}
	// 初始化
	fmt.Println("  Server :", ip, port)
	for {
		conn, err := listen.AcceptTCP()
		defer conn.Close()
		if err != nil {
			fmt.Println("AcceptTCP Exception::", err.Error())
			continue
		}
		// 校验ip地址
		conn.SetKeepAlive(true)

		go handleMsg(conn)
	}
}
func handleMsg(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	fmt.Println("HandleConn client: ", conn.RemoteAddr() )
	for {

		_cmd,_sid, _seq,_payload, err := protocol.HubUnPack(reader)
		fmt.Println("server recvice : ",string(_cmd), string(_sid), string(_seq), string(_payload) )
		if err != nil {
			fmt.Println("HandleConn connection error: ", err.Error())
			break
		}
		resp_buf,err := protocol.HubPack( string(_cmd),string(_sid),string(_seq),_payload )
		if err != nil {
			fmt.Println("HandleConn protocol.EncodePacket error: ", err.Error())
			break
		}
		conn.Write( resp_buf )


	}
}

func client_hub() {

	// 客户端请求
	fmt.Println("  client_side " )
	service := "127.0.0.1:7004"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if ( err != nil ) {
		fmt.Println("  net.DialTCP error: ", err.Error())
		return
	}

	//fmt.Println( conn )
	reader := bufio.NewReader(conn)

	cmd := `GetUser`
	sid := `sid121`
	seq := `seq001`
	data := []byte(`{"user":"simarui","pass":123","data":"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww122"}`)

	buf,err := protocol.HubPack( cmd,sid,seq, data)
	if ( err != nil ) {
		fmt.Println("protocol.EncodePacket error: ", err.Error())
		return
	}
	fmt.Println("conn.Write:", string(buf) )
	conn.Write( buf )
	for {
		_cmd,_sid, _seq,_payload, err := protocol.HubUnPack(reader)
		if err != nil {
			fmt.Println(" connection error: ", err.Error())
			break
		}

		fmt.Println( "resp : ",string(_cmd), string(_sid), string(_seq), string(_payload) )
		break
	}
}