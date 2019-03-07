package worker

import (
	"fmt"
	_ "fmt"
	"github.com/antonholmquist/jason"
	main "masterlab_socket"
	"net"
)

func (this TaskType) Auth() ReturnType {

	//sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	fmt.Println( "Auth this.Data:",string(this.Data) )
	sid := main.AreaCreateSid()
	if (sid != "") {
		ret := ReturnType{"ok", "welcome", sid, ""}
		return ret
	} else {
		ret := ReturnType{"failed", "failed", sid, ""}
		return ret
	}
}

func (this TaskType) Push() string {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)

	from_sid := this.ReqHeader.Sid
	json_obj, err := jason.NewObjectFromBytes(this.Data)
	if err != nil {
		return ""
	}
	to_sid, _ := json_obj.GetString("sid")
	sdk.Push(from_sid, to_sid, this.Data)

	return ""
}

func (this TaskType) Broadcast() string {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)

	from_sid := this.ReqHeader.Sid
	json_obj, err := jason.NewObjectFromBytes(this.Data)
	if err != nil {
		return ""
	}
	area_id, _ := json_obj.GetString("area_id")

	if (area_id == "global") {
		main.LogError("broatcast global failed")
		return ""
	} else {
		sdk.Broatcast(from_sid, area_id, this.Data)
	}
	return ""
}

func (this TaskType) GetUserSession() string {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)

	return sdk.GetSession(this.ReqHeader.Sid)

}

func (this TaskType) GetAreas() string {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)

	return sdk.GetAreasStr()

}

func (this TaskType) JoinArea() string {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)
	//fmt.Println( "JoinChannel",this.Data  )
	if (sdk.AreaAddSid(this.ReqHeader.Sid, string(this.Data))) {
		return "ok"
	} else {
		return "failed"
	}

}

func (this TaskType) LeaveChannel() interface{} {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)
	//fmt.Println( "LeaveChannel header:",this.ReqHeader  )
	if (sdk.AreaKickSid(this.ReqHeader.Sid, string(this.Data))) {
		return "ok"
	} else {
		return "failed"
	}

}

func (this TaskType) KickSelf() interface{} {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)

	if (sdk.Kick(this.ReqHeader.Sid)) {
		return "ok"
	} else {
		return "failed"
	}

}

func (this TaskType) GetBase(conn *net.TCPConn, cmd string, req_sid string, req_id int, req_data string) string {

	sdk := new(Sdk).Init(this.ReqType, this.ReqHeader, this.Data)
	return sdk.GetBase()

}
