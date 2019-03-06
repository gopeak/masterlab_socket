package worker

import (
	"github.com/antonholmquist/jason"
	"masterlab_socket/golog"
	"fmt" 
	"strconv"
)




func (this TaskType)Message(   ) string {

	sdk:=new(Sdk).Init( this.ReqType,this.ReqHeader,this.Data   )

	data_json ,err_json:= jason.NewObjectFromBytes( this.Data )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}

	_,err1 := data_json.GetString("type")
	_,err2 := data_json.GetString("message")
	sid,err3 := data_json.GetString("id")
	if( err1!=nil || err2!=nil || err3!=nil ){
		//golog.Error("todpole message json err:",err1.Error()+err2.Error()+err3.Error())
		return ""
	}
	//broatcast_msg := fmt.Sprintf(`{"type":"message","message":"%s","id":"%s" }`,message,sid)
	sdk.Broatcast( sid,"area-global", this.Data   )
	//json_ret := fmt.Sprintf(`{"type":"messageresp","id":"%s" }`,sid)
	return "";


}

func (this TaskType)Update(   ) string {

	sdk:=new(Sdk).Init( this.ReqType,this.ReqHeader,this.Data   )

	data_json ,err_json:= jason.NewObjectFromBytes( this.Data )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}

	type_str,_ := data_json.GetString("type")
	angle_str,_ := data_json.GetString("angle")
	_id,_ := data_json.GetString("id")
	momentum_str,_ := data_json.GetString("momentum")
	x_str,_ := data_json.GetString("x")
	y_str,_ := data_json.GetString("y")
	name,err_name := data_json.GetString("name")
	if( err_name!=nil ) {
		name = "Guest."+ _id ;
	}
	angle,_ := strconv.ParseFloat(angle_str, 32)
	momentum ,_:= strconv.ParseFloat(momentum_str, 32)
	x,_ := strconv.ParseFloat(x_str, 32)
	y ,_:= strconv.ParseFloat(y_str, 32)
	broatcast_data := fmt.Sprintf(`{"type":"%s","id":"%s","angle":%.3f,"momentum":%.3f,"x":%.3f,"y":%.3f,"life":1,"name":"%s","authorized":%s}`,
		type_str,_id,float32(angle),float32(momentum),float32(x),float32(y),name,"false" )

	sdk.Broatcast( this.ReqHeader.Sid,"area-global",[]byte(broatcast_data) )
	return ""
	json_ret := fmt.Sprintf(`{"type":"%s","id":"%s" }`,"none",_id)
	return json_ret;


}



