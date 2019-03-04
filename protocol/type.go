package protocol

//  websocket 通讯协议
//  请求数据包 { "header":{ "cmd":"", "seq_id":0,  "sid":"" , "token":"", "version":"1.0" ,"gzip":true}  , "type":"req", "data":{}  }
//  请求响应数据包  { "header":{ "cmd":"", "seq_id":0,  "sid":"", "gzip":true  }  , "type":"response",  "status":0, "msg":"",  "data":{}  }
//  接收广播数据包 { "header":{ "area_id":"", "sid":""   }  , "type":"broatcast",   "data":{}  }
//  接收点对点数据包 { "header":{  "sid":""  }  , "type":"push",   "data":{}  }

const  TypeReq  	 = "1"
const  TypeResp  	 = "2"
const  TypeBroatcast  = "3"
const  TypePush 	 = "4"
const  TypeError  	 = "5"
const  TypeReply	 = "6"
const  TypePing	 = "7"

type ProtocolType struct {

	ReqObj ReqRoot
	RespObj ResponseRoot
	BroatcastObj BroatcastRoot
	PushObj  PushRoot
	Init func()

}

type BaseRoot struct {
	Type string             `json:"type"`
	Header  interface{}    `json:"header"`
	Data []byte         	 `json:"data"`
}

type ReqRoot struct {
	Type string            	 `json:"type"`
	Header  ReqHeader        `json:"header"`
	WsData interface{}	 `json:"data"`
	Data []byte

}


type ReqHeader struct {
	Cmd    string      	`json:"cmd"`
	SeqId  int          	`json:"seq_id"`
	Sid    string     	`json:"sid"`
	NoResp bool         	`json:"no_resp"`
	Token    string     	`json:"token"`
	Version    string     	`json:"version"`
	Gzip    bool     	`json:"gzip"`
}


type ResponseRoot struct {
	Type string             `json:"type"`
	Header  RespHeader     	`json:"header"`
	Data []byte             `json:"data"`
}

type RespHeader struct {
	Cmd    string        `json:"cmd"`
	SeqId  int           `json:"seq_id"`
	Sid    string        `json:"sid"`
	Gzip    bool         `json:"gzip"`
	Status int 	     `json:"status"`
}

type BroatcastRoot struct {
	Type string             	`json:"type"`
	Header  BroatcastHeader          `json:"header"`
	Data    interface{}   		`json:"data"`
}

type BroatcastHeader struct {
	AreaId  string      `json:"area_id"`
	Sid    string          `json:"sid"`
}

type PushRoot struct {
	Type string             `json:"type"`
	Header  PushHeader      `json:"header"`
	Data interface{}    	`json:"data"`
}

type PushHeader struct {
	Sid    string          `json:"sid"`
}
