package main

import (
	"bufio"
	"fmt"
	"net"
	"masterlab_socket/protocol"
)



func main() {

	// 客户端请求
	fmt.Println("  client_side " )
	service := "127.0.0.1:9002"
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

	header := []byte( `{"cmd":"Mail","sid":"1234516","ver":"1.2","seq":12123,"token":"sssssssssss121"}`)
	data := []byte(`{"host":"smtpdm.aliyun.com","port":"465","user":"sender@smtp.masterlab.vip","password":"MasterLab123Pwd","from":"sender@smtp.masterlab.vip","to":"121642038@qq.com","subject":"Hello","body":"hello world"}`)
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