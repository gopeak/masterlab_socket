package main

import (
	"bufio"
	"fmt"
	"masterlab_socket/protocol"
	"net"
	"time"
)

func main() {

	go Server("0.0.0.0", 7003)
	time.Sleep(3 * time.Second)
	go client_side()
	select {

	}
}


/**
 * 监听客户端连接
 */
func Server(ip string, port int) {


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

		go handleClientMsg(conn)
	}
}
func handleClientMsg(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	fmt.Println("HandleConn client: ", conn.RemoteAddr() )
	for {

		_,header, data,_, err := protocol.DecodePacket( reader )
		fmt.Println("server recvice header: ", string(header), " data:", string(data))
		if err != nil {
			fmt.Println("HandleConn connection error: ", err.Error())
			break
		}
		resp_buf,err := protocol.EncodePacket(  protocol.TypeReply , header,data)
		if err != nil {
			fmt.Println("HandleConn protocol.EncodePacket error: ", err.Error())
			break
		}
		conn.Write( resp_buf )


	}
}

func client_side() {

	// 客户端请求
	fmt.Println("  client_side " )
	service := "127.0.0.1:7003"
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

	header := []byte( `{"cmd":"Auth","sid":"1234516","ver":"1.2","seq":12123,"token":"sssssssssss121"}`)
	data := []byte(`{"user":"simarui","pass":123","data":"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww122"}`)
	buf,err := protocol.EncodePacket( protocol.TypeReq , header, data)
	if ( err != nil ) {
		fmt.Println("protocol.EncodePacket error: ", err.Error())
		return
	}
	fmt.Println("conn.Write:", string(buf) )
	conn.Write( buf )
	for {
		_type,resp_header, resp_data,_, err := protocol.DecodePacket(reader)
		if err != nil {
			fmt.Println(" connection error: ", err.Error())
			break
		}

		fmt.Println( "resp _type: ", _type )
		fmt.Println( "resp header: ", string(resp_header), " data:", string(resp_data) )
		break
	}
}