// automatically generated, do not modify

package protocol

import (

	"encoding/binary"
	"bytes"
	"bufio"
	"errors"
	"fmt"
)


func HubPack(  cmd  ,sid  ,callback_seq string , payload []byte) ( []byte ,error ){
	// len(total)+ len(header) == 12
	var pkg *bytes.Buffer = new(bytes.Buffer)
	cmd_len := uint32(len(cmd))
	sid_len := uint32(len(sid))
	seq_len := uint32(len(callback_seq))
	payload_len := uint32(len(string(payload)))
	totalsize := cmd_len + sid_len + seq_len +  payload_len+12

	// set totalsize
	err:=binary.Write( pkg , binary.LittleEndian, totalsize)
	if err != nil {
		return nil, err
	}
	err=binary.Write( pkg , binary.LittleEndian, cmd_len)
	if err != nil {
		return nil, err
	}
	err=binary.Write( pkg , binary.LittleEndian, sid_len)
	if err != nil {
		return nil, err
	}
	err=binary.Write( pkg , binary.LittleEndian, seq_len)
	if err != nil {
		return nil, err
	}

	err = binary.Write( pkg, binary.LittleEndian, []byte(cmd))
	if err != nil {
		return nil, err
	}
	err = binary.Write( pkg, binary.LittleEndian, []byte(sid))
	if err != nil {
		return nil, err
	}
	err = binary.Write( pkg, binary.LittleEndian, []byte(callback_seq))
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


func HubUnPack(r *bufio.Reader) ( []byte, []byte,  []byte, []byte, error) {
	var totalsize  uint32

	lengthByte, _ := r.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	err := binary.Read(lengthBuff, binary.LittleEndian, &totalsize)
	if err != nil {
		return nil,nil,nil,nil,errors.New("read total size error") // errors.Annotate(err, "read total size")
	}
	//fmt.Println( "totalsize" , totalsize)
	if totalsize < 12 {
		return nil,nil, nil,nil,errors.New( fmt.Sprintf("bad packet. totalsize:%d", totalsize))
	}

	pack := make([]byte, int(totalsize)+4)
	_, err = r.Read(pack)
	if err != nil {
		return nil,nil,nil,nil, err
	}
	if len(pack)<int(totalsize)+4 {
		return nil,nil,nil,nil, errors.New("read headersize error")
	}
	//fmt.Println( "pack"  ,string(pack))

	cmd_size :=   uint32( pack[4] )
	sid_size :=   uint32( pack[8] )
	seq_size :=   uint32( pack[12] )
	//fmt.Println( "size"  ,cmd_size,sid_size, seq_size )
	cmd :=  pack[16:cmd_size+16]
	sid :=  pack[16+cmd_size:sid_size+16+cmd_size]
	callback_seq :=  pack[16+cmd_size+sid_size:16+cmd_size+sid_size+seq_size]
	payload := pack[(16+cmd_size+sid_size+seq_size):int(totalsize)+4]
	/*
	payload_str := string(payload)
	callback_seq_str := string(callback_seq)
	sid_str := string(sid)
	cmd_str := string( cmd )
	fmt.Println( cmd_str,sid_str,callback_seq_str,payload_str)
	*/
	return  cmd , sid , callback_seq ,payload, nil
}



