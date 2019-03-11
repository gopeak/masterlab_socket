package main

import (
	"bufio"
	"fmt"
	"masterlab_socket/protocol"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
	"github.com/antonholmquist/jason"
)

var Conns   []*net.TCPConn
var Sids []string
var Tokens []string


func createReqConns(num int64)  {

	Conns = make([]*net.TCPConn, 0)
	Sids = make([]string, 0)
	for i := 0; i < int(num); i++ {
		service := os.Args[1]
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			fmt.Println(err.Error())
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		//defer conn.Close()
		Conns = append(Conns, conn)
		time.Sleep(10 * time.Millisecond)
		//srcData :=  []byte( strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(i)  )
		//str:=md5.Sum([]byte(srcData))

		protocolPack:= new(protocol.Pack)
		protocolPack.Init()
		//sid := strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(i)
		//token := ""
		//data :=  strconv.FormatInt(int64(time.Now().Unix()), 10)
		//buf,err :=  protocolPack.WrapReq( "Auth", sid, token, i, []byte(data) )

		header := []byte( `{"cmd":"Auth","sid":"1234516","ver":"1.2","seq":12123,"token":"sssssssssss121"}`)
		data := []byte(`{"user":"simarui","pass":123","data":"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww122"}`)
		buf,err := protocol.EncodePacket( protocol.TypeReq , header, data)

		fmt.Println(" protocolPack.WrapReq:", string(buf))
		if err != nil {
			fmt.Println(" protocolPack.WrapReq err: ", err.Error() )
			continue
		}
		conn.Write( buf )
		r := bufio.NewReader(conn)
		for {
			_, _,resp_data,_,err :=  protocol.DecodePacket( r )
			fmt.Println( "Auth:",   string(resp_data) )
			json,_ := jason.NewObjectFromBytes( resp_data )

			sid ,err:= json.GetString("sid")
			if err != nil {
				fmt.Println("json.GetString(`sid`) err: ", err.Error() )
				continue
			}
			token ,err:= json.GetString("msg")
			if err != nil {
				fmt.Println("json.GetString(`token`) err: ", err.Error() )
				continue
			}

			buf,_ =  protocolPack.WrapReq( "JoinArea", sid, token, 0, []byte("area-global") )
			conn.Write([]byte( buf ))

			//fmt.Println( "", sid  )
			Sids  = append(  Sids, sid )
			Tokens =  append(  Tokens,token )


			break
		}
	}


} //

func hanleConnResp( conn *net.TCPConn ,times int64, conn_num int64 ,i int ){

	//fmt.Println( conn )
	reader := bufio.NewReader(conn)
	var success int64
	success = 0
	req_sid := Sids[i]
	fmt.Println("req_sid:",req_sid)
	token := Tokens[i]
	//data :=  protocol.WrapReqStr("GetUserSession",req_sid,0,req_sid )
	protocolPack:= new(protocol.Pack)
	protocolPack.Init()
	buf,_ :=  protocolPack.WrapReq( "GetUserSession", req_sid, token, 0, []byte(req_sid) )
	n,err := conn.Write([]byte( buf ))
	if( n<=0 ) {
		fmt.Println( " GetUserSession write size err:",n )
	}
	if( err!=nil ) {
		fmt.Println( " GetUserSession write err:",err.Error() )
	}
	have_lava_area :=false
	recvice_br_msg:=0
	defer conn.Close()
	for {
		ptype, resp_header,resp_data,_,err :=  protocol.DecodePacket( reader )
		//fmt.Println( "protocol.DecodePacket:",ptype,  string(resp_header), string(resp_data) )
		if err != nil {
			//fmt.Println("HandleConn connection error: ", err.Error())
			conn.Close()
			return
		}
		_type := fmt.Sprintf("%d",ptype)
		success++

		if( _type==protocol.TypeResp ){
			resp_header_obj,msg_err := protocolPack.GetRespHeaderObj( resp_header )
			if msg_err != nil {
				fmt.Println("msg error: ", msg_err.Error() )
				continue
			}
			// 登录认证,然后获取用户信息
			req_id := resp_header_obj.SeqId
			// 获取当前信息后 发送点对点信息
			if resp_header_obj.Cmd=="GetUserSession"  {

				fmt.Println("GetUserSession:",string(resp_data))
				// 发送点对点消息
				go func() {
					to_sid_index := i - 1
					if to_sid_index < 0 {
						to_sid_index = 0
					}
					to_sid := Sids[to_sid_index]
					push_data := fmt.Sprintf(`{"sid":"%s","data":"%s"}`, to_sid, "md55555555555")
					buf, _ := protocolPack.WrapReq("Push", req_sid, token, 0, []byte(push_data))
					conn.Write([]byte( buf ))
				}()
				// 获取场景列表
				//go func() {
					buf, _ := protocolPack.WrapReq("GetAreas", req_sid, token, 0, []byte(""))
					conn.Write([]byte( buf ))
				//}()


			}
			if resp_header_obj.Cmd=="GetAreas"  {
				fmt.Println("GetAreas Revcie:",string(resp_data))
			}

			if resp_header_obj.Cmd=="JoinArea"  {
				//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				push_data := fmt.Sprintf(`{"area_id":"area-global","data":"%s"}`,"md56666666666")
				buf,_ =  protocolPack.WrapReq( "Broadcast", req_sid, token, req_id+1, []byte(push_data) )
				conn.Write([]byte( buf ))
				fmt.Println("Broadcast sended:",i , string(buf))
			}

			if resp_header_obj.Cmd=="LeaveChannel"  {
				buf,_ =  protocolPack.WrapReq( "KickSelf", req_sid, token, req_id+1, []byte(req_sid) )
				conn.Write([]byte( buf ))
			}
			if resp_header_obj.Cmd=="KickSelf"  {
				conn.Close()
				return
			}

		}
		if _type==protocol.TypePush  {
			fmt.Println("Push Revcie:",string(resp_data))
		}

		// 发送广播
		if _type==protocol.TypeBroatcast  {
			recvice_br_msg++
			if(  recvice_br_msg>=int(conn_num) ){
				fmt.Println("Broadcast Revcie:",recvice_br_msg)
				if !have_lava_area {
					time.Sleep(200 * time.Millisecond)
					have_lava_area = true
					buf,_ =  protocolPack.WrapReq( "LeaveChannel", req_sid, token, 0, []byte("area-global") )
					conn.Write([]byte( buf ))
				}
			}



		}
	}
}



func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//start := time.Now().Unix()
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage:%s  host:port connections send_times packet_type ", os.Args[0])
		os.Exit(1)
	}

	times, _ := strconv.ParseInt(os.Args[3], 10, 32)
	conn_num, _ := strconv.ParseInt(os.Args[2], 10, 32)
	fmt.Println("Connections and  times:", conn_num, times)

	createReqConns(conn_num)
	//ch_success := make(chan int64, 0)
	time.Sleep(2 * time.Second)
	var i int64
	for i = 0; i < conn_num; i++ {
		conn := Conns[i]
		go hanleConnResp(  conn,  times, conn_num , int(i) )
	}
	select {

	}
	/*
	var qps int64
	var recv_times int64
	qps = 0
	recv_times = 0
	for {
		select {
		case r := <-ch_success:
			recv_times++
		//fmt.Println("recv_times:", recv_times)
			qps = qps + r
			if recv_times == conn_num-1 {
				fmt.Printf(".")
				end := time.Now().Unix()
				els_time := end - start
				//fmt.Println("time:", els_time, qps)
				fmt.Printf("\nels_time:%d recv_times:%d qps:%d", els_time, recv_times, qps)
				return
			}
		default:
			fmt.Printf(".")
			time.Sleep(100 * time.Millisecond)
		}
	}
	*/

}

