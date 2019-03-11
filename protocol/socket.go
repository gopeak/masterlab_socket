//通讯协议处理，主要处理封包和解包的过程
package protocol


import (

	"encoding/binary"
	"errors"
	"fmt"
	"bufio"
	"bytes"
	"encoding/json"
	"strconv"
	"masterlab_socket/util"
)


//|总长度4|type4|头长度4|header|data|checksum
type ClientPacket struct {
	TotalSize uint32    //4
	Type      uint8
	HeaderSize uint32   //4
	Header   []byte
	Payload   []byte
	Checksum  uint32    //4
}


func EncodePacket(  _type string, header []byte,  payload []byte) ( []byte ,error ){
	// len(totaol)+ len(header) + len(Checksum) == 12
	var pkg *bytes.Buffer = new(bytes.Buffer)
	totalsize := uint32(len(string(header)) +  len(string(payload)) ) +4

	// set totalsize
	err:=binary.Write( pkg , binary.LittleEndian, totalsize)
	if err != nil {
		return nil, err
	}
	//fmt.Println( "set totalsize" , totalsize )
	_type_int,err:=strconv.Atoi(_type)
	// set type
	err = binary.Write( pkg, binary.LittleEndian, uint32(_type_int) )
	if err != nil {
		return nil, err
	}
	// set headersize
	headersize := uint32(len(string(header)))
	err = binary.Write( pkg, binary.LittleEndian, headersize)
	if err != nil {
		return nil, err
	}
	//fmt.Println( "set headersize" , headersize )

	// set header
	err = binary.Write( pkg, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	// set payload
	err = binary.Write(pkg, binary.LittleEndian, payload )
	if err != nil {
		return nil, err
	}
	// write checksum
	return  pkg.Bytes(),nil
}


func DecodePacket(r *bufio.Reader) ( uint32, []byte,  []byte, []byte, error) {
	var totalsize , headersize uint32

	lengthByte, _ := r.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	err := binary.Read(lengthBuff, binary.LittleEndian, &totalsize)
	if err != nil {
		return 0,nil,nil,nil,err// errors.Annotate(err, "read total size")
	}
	fmt.Println( "totalsize" , totalsize)
	if totalsize < 12 {
		return 0,nil, nil,nil,errors.New( fmt.Sprintf("bad packet. totalsize:%d", totalsize))
	}

	pack := make([]byte, int(totalsize)+4+4)
	_, err = r.Read(pack)
	if err != nil {
		return 0,nil,nil,nil, err
	}
	if len(pack)<4 {
		return 0,nil,nil,nil, errors.New("read headersize error")
	}

	_type :=   uint32(pack[4] )
	fmt.Println( "type:" ,_type )

	headersize =   uint32(pack[8] )
	fmt.Println( "headersize"  ,headersize )
	if len(pack)< int(headersize+12) {
		fmt.Println( "pack:"  ,string(pack) )
		return 0,nil,nil,nil, errors.New("headersize error")
	}
	header :=  pack[12:headersize+12]

	payload := pack[(12+headersize):(totalsize+4+4)]
	fmt.Println( "header:" , string(header))
	fmt.Println( "payload:" ,  string(payload) )
	return _type,header,payload,pack, nil
}


type Pack struct {
	ProtocolObj ProtocolType
	Data        []byte
}

func (this *Pack) Init() *Pack {

	this.ProtocolObj = ProtocolType{}
	this.ProtocolObj.ReqObj = ReqRoot{}
	this.ProtocolObj.RespObj = ResponseRoot{}
	this.ProtocolObj.BroatcastObj = BroatcastRoot{}
	this.ProtocolObj.PushObj = PushRoot{}
	return this
}

func (this *Pack) GetReqObjByReader( reader *bufio.Reader ) (*ReqRoot, error) {

	stb := &ReqRoot{}
	_type,header,data,_,err := DecodePacket( reader )
	if err!=nil {
		return stb, err
	}
	return this.GetReqObj( _type,header,data )

}

func (this *Pack) GetReqHeaderObj(  header []byte) (*ReqHeader, error) {

	stb := &ReqHeader{}
	err := json.Unmarshal( header, stb )
	return stb, err
}

func (this *Pack) GetReqObj( _type uint32 ,header []byte, data []byte ) (*ReqRoot, error) {

	var req_header ReqHeader
	stb := &ReqRoot{}
	header = util.TrimX001(header)
	stb.Type = fmt.Sprintf( "%d", _type )
	err :=json.Unmarshal(header, &req_header)
	if err!=nil {
		return stb, err
	}
	stb.Header = req_header
	stb.Data = data
	//this.ProtocolObj.ReqObj = stb
	return stb, err
}


func (this *Pack) GetRespHeaderObj(  header []byte) (*RespHeader, error) {

	stb := &RespHeader{}
	err := json.Unmarshal( header, stb )
	return stb, err
}

func (this *Pack) GetRespObj( _type uint32,  header []byte,  data []byte) (*ResponseRoot, error) {

	var resp_header RespHeader
	stb := &ResponseRoot{}
	err :=json.Unmarshal(header, &resp_header)
	if err!=nil {
		return stb, err
	}
	stb.Type = fmt.Sprintf( "%d",_type )
	stb.Header = resp_header
	stb.Data = data
	return stb, err
}


func (this *Pack) GetBroatcastHeaderObj(  header []byte) (*BroatcastHeader, error) {

	stb := &BroatcastHeader{}
	err := json.Unmarshal( header, stb )
	return stb, err
}


func (this *Pack) GetBroatcastObj(data []byte) (*BroatcastRoot, error) {
	this.Data = data
	stb := &BroatcastRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.BroatcastObj = stb
	return stb, err
}

func (this *Pack) GetPushObj(data []byte) (*PushRoot, error) {
	this.Data = data
	stb := &PushRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.PushObj = stb
	return stb, err
}

func (this *Pack) WrapReq( cmd ,sid ,token string, seq int, data []byte ) ([]byte, error) {
	req_obj_header := &ReqHeader{}
	req_obj_header.Cmd = cmd
	req_obj_header.Sid = sid
	req_obj_header.Token = token
	req_obj_header.SeqId = seq
	header_buf ,_ := json.Marshal( req_obj_header )
	fmt.Println( "header_buf:", string(header_buf) )
	return  EncodePacket( TypeReq, header_buf, data  )
}

func (this *Pack) WrapReqWithHeader( req_header *ReqHeader, data []byte ) ([]byte, error) {

	header_buf ,_ := json.Marshal( req_header )
	//fmt.Println( "header_buf:", string(header_buf) )
	return  EncodePacket( TypeReq, header_buf, data  )
}

func (this *Pack) WrapResp( cmd, req_sid string, seq int, status int,  data []byte ) ( []byte,error ) {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = cmd
	resp_header_obj.Sid = req_sid
	resp_header_obj.SeqId = seq
	resp_header_obj.Status = status
	header_buf ,_ := json.Marshal( resp_header_obj )
	return  EncodePacket( TypeResp, header_buf, data  )

}


/**
 * 封包返回客户端错误的消息
 */
func (this *Pack) WrapRespErr( err string) ( []byte,error ) {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = "WrapRespErr"
	resp_header_obj.Sid = ""
	resp_header_obj.SeqId = 0
	resp_header_obj.Status = 500

	header_buf ,_ := json.Marshal( resp_header_obj )
	return  EncodePacket(  TypeError,header_buf, []byte(err)  )
}


func (this *Pack) WrapRespObj( req_obj *ReqRoot, invoker_ret []byte, status int ) ResponseRoot {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = req_obj.Header.Cmd
	resp_header_obj.SeqId = req_obj.Header.SeqId
	resp_header_obj.Gzip = req_obj.Header.Gzip
	resp_header_obj.Sid = req_obj.Header.Sid
	resp_header_obj.Status = status
	this.ProtocolObj.RespObj.Header =resp_header_obj
	this.ProtocolObj.RespObj.Data = invoker_ret
	this.ProtocolObj.RespObj.Type = "2"

	return this.ProtocolObj.RespObj
}

func (this *Pack) WrapPushRespObj(to_sid string, from_sid string , data[]byte ) PushRoot {

	push_header_obj := PushHeader{}
	push_header_obj.Sid = from_sid

	push_obj := PushRoot{}
	push_obj.Header =push_header_obj
	push_obj.Data  = data
	push_obj.Type  = "push"

	return push_obj
}

func (this *Pack) WrapPushResp(to_sid string, from_sid string , data_buf []byte ) ([]byte,error) {

	push_header_obj := PushHeader{}
	push_header_obj.Sid = from_sid
	header_buf,_ := json.Marshal( push_header_obj )
	return  EncodePacket(  TypePush ,header_buf, data_buf )

}


func (this *Pack) WrapBroatcastRespObj( area_id , from_sid string , data []byte ) BroatcastRoot {

	broatcast_header_obj := BroatcastHeader{}
	broatcast_header_obj.Sid = from_sid
	broatcast_header_obj.AreaId = area_id

	broatcast_obj := BroatcastRoot{}
	broatcast_obj.Header =broatcast_header_obj
	broatcast_obj.Data  = data
	broatcast_obj.Type  = TypeBroatcast

	return broatcast_obj
}

func (this *Pack) WrapBroatcastResp( area_id, from_sid string,  data []byte ) ( []byte,error ) {

	broatcast_header_obj := BroatcastHeader{}
	broatcast_header_obj.Sid = from_sid
	broatcast_header_obj.AreaId = area_id

	header_buf ,_ := json.Marshal( broatcast_header_obj )
	return  EncodePacket(  TypeBroatcast ,header_buf, data  )

}



