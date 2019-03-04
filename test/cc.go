package main

import (
	 
   "fmt"
   "crypto/rand"
   "math/big"
   "time" 
   "encoding/json" 
   "masterlab_socket/protocol"
   "github.com/antonholmquist/jason"
   "masterlab_socket/global"
   //"io/ioutil"
   "net"
   "os"
   "os/signal"
   "strconv"
   "runtime" 
   "bufio"
  // "log"
)


var reqs = []string{ }
 

type SocketResponse struct {
  data      string
  response  string
  err       error
}

func main2() {

    runtime.GOMAXPROCS(runtime.NumCPU()) 
    start:=time.Now().Unix()
    go end_hook( start )
    global.PackSplitType = "breakline"
    num, _ := strconv.ParseInt( os.Args[2], 10, 32)  
	for i:=0;i<int(num);i++{
		 
		max := big.NewInt(1000)
		rand, _ := rand.Int(rand.Reader, max) 
		reqs  = append( reqs,fmt.Sprintf( `{"token":"%s", "cmd":"socket.user_login","params":{"user":"admin_xbd","password":"258369"}}`  ,rand) )		
	}
    //golog.Println( "urls:" ,urls  )
	results := asyncReq( reqs )
  
    i:=0
	for _, result := range results {
        i++
		fmt.Printf( "%d result length: %d\n",i,  len(result.data) )
	}
    
    fmt.Printf( " result num: %d\n",  len(results) )

    select{} 
}

func end_hook( start int64  ) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)

    s := <-c
    end:=time.Now().Unix()
    els_time := end-start
    fmt.Println("Got signal:", s )
    fmt.Println("need time :", els_time )
    os.Exit( 1 )
    
}

func asyncReq( datas []string ) []*SocketResponse {

	//ch := make( chan *SocketResponse )
	responses  := []*SocketResponse{}
	
	for _, data:= range datas {
		go func(data string) {
		
			//fmt.Printf("Req %s \n", data)
			times, _ := strconv.ParseInt( os.Args[3], 10, 32)  
			if len(os.Args) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
				os.Exit(1)
			}
			service := os.Args[1]
			tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
			checkError(err)
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
            defer conn.Close()
			checkError(err) 
			 reqWithBufferio( conn, data, int(times) )
			 return  
            //time.Sleep(10 * time.Millisecond)
            if len(os.Args) >=5 {
                req_type :=  string(os.Args[4])
                if req_type=="json" {
                    reqWithJson( conn, data, int(times) )
                } 
            }else{
                reqWithBufferio( conn, data, int(times) )
            }
			
		}(data)
	} 
 
	return responses
}

 

func reqWithJson(  conn *net.TCPConn , data string ,times int  ){
    
    _, err := conn.Write([]byte( data+"\n" ))
	checkError(err)
	var i int
    i = 0
    sid:=""
    str :=`{ "cmd":"socket.getSession"}`
	d := json.NewDecoder(conn)   
    for {
		
		var msg interface{}
     
        err := d.Decode(&msg) 
        if  err != nil { 
            
            conn.Close() 
            fmt.Println( "d.Decode(&msg) ", err.Error()  )  
            break 
        }
		json_encode,err_encode := json.Marshal( msg ) 
        if err_encode!=nil { 
            fmt.Println( "json.Marshal error:",err_encode.Error() )  
            conn.Close()                   
            break
        } 
		response_str := string( json_encode )
		 
		//fmt.Printf( " response: %s\n", response_str )
		msg_json, errjson := jason.NewObjectFromBytes( []byte(response_str) ) 
		checkError(errjson)
		cmd,  _ := msg_json.GetString("cmd") 
		//fmt.Printf( " cmd: %s\n", cmd )                    
		if cmd=="socket.user_login" { 
			sid,  _ = msg_json.GetString("data","sid")   
			//fmt.Printf( " sid: %s\n", sid )  
			str=fmt.Sprintf( `{ "cmd":"socket.getSession","params":{"sid":"%s"} }` ,sid )
			str =  str+"\n" 
			//fmt.Println( " post : ", str )  
			_,err = conn.Write([]byte( str )) 
			//time.Sleep(10 * time.Millisecond )
			checkError(err)  
            
		}
		if cmd=="socket.getSession"  {  
		   
			_,err = conn.Write([]byte(  str+"\n" )) 
			checkError(err) 
			i++  
			time.Sleep(10 * time.Millisecond )
			if( i>=times ){
				//conn.Close()
                fmt.Println( " i : ", i )  
				fmt.Println( " conn close! " ,i ,"\n" )
				break
			}
		}
        //ch <- &SocketResponse{ data, string(response_str), err} 
    }
    //conn.Close()
    
}



func reqWithBufferio(  conn *net.TCPConn , data string ,times int  ){
    
    
    req_ready_byte,_ := protocol.Packet ( data ) 
	conn.Write( req_ready_byte )
	
	var i int
    i = 0
    sid:=""
    str :=`{ "cmd":"socket.getSession"}`
	
    reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString( '\n' )
        if err != nil {
            fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
            conn.Close()    
			return 
        } 
       
        if( msg=="" ) {
            continue
        }
     
        msg_json, errjson := jason.NewObjectFromBytes( []byte(msg) ) 
		if errjson != nil {
		    //fmt.Println( " errjson  :",  errjson.Error() )
            continue		     
		}
		cmd,  _ := msg_json.GetString("cmd") 
		//fmt.Printf( " msg: %s\n", msg )                    
		if cmd=="socket.user_login" { 
			sid,  _ = msg_json.GetString("data" )   
			//fmt.Printf( " sid: %s\n", sid )  
			str=fmt.Sprintf( `{ "cmd":"socket.getSession","params":{"sid":"%s"} }` ,sid )
			session_req_byte,_ := protocol.Packet ( str ) 
	        conn.Write( session_req_byte )
			//time.Sleep(10 * time.Millisecond )
			checkError(err)  
            
		}
		if cmd=="socket.getSession"  && sid!=""  {  
		   
		    // session_str,  _ := msg_json.GetString("data" ) 
		    //session_json, _ := jason.NewObjectFromBytes( []byte(session_str) )   
			//sid ,_= session_json.GetString("msg","Sid")
			//fmt.Printf( " sid: %s\n", sid )  
			str=fmt.Sprintf( `{ "cmd":"socket.getSession","params":{"sid":"%s"} }` ,sid )
			session_req_byte,_ :=protocol.Packet ( str )
	        conn.Write( session_req_byte  )
			checkError(err) 
			i++  
			time.Sleep(10 * time.Millisecond )
			if( i>times ){
				conn.Close()
                fmt.Println( " i : ", i )  
				fmt.Println( " conn close! " ,i ,"\n" )
				break
			}
		}
		 
	 }
   
    conn.Close()
    
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

 