package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)


type Util struct {


}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}


	return ioutil.ReadAll(f)
}
func saveFile(str string, n int) {
	f, err := os.Create("./output" + strconv.Itoa(n) + ".txt") //创建文件

	if err != nil {
		LogError("os.Create Error:", err.Error())
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	_, errw := w.WriteString(str)
	if errw != nil {
		LogError("WriteString Error:", errw.Error())
		return
	}
	//fmt.Printf("写入 %d 个字节\n", n4)
	w.Flush()
	f.Close()

}

//  转义json字符串
func EncodeJsonStr(str string) string {
	str = strings.Replace(str, `"`, `\"`, -1)
	return str
}

// 反解json字符串
func DecodeJsonStr(str string) string {
	str = strings.Replace(str, `\"`, `"`, -1)
	return str
}

func TrimStr(str string) string {
	str = strings.Replace(str, " ", "",  -1)
	str = strings.Replace(str, "\n", "",  -1)
	return str
}


func TrimX001(data_buf []byte) []byte {
	return data_buf
	for i, ch := range data_buf {

		switch {
		case ch > '~':   data_buf[i] = ' '
		case ch == '\r':
		case ch == '\n':
		case ch == '\t':
		case ch < ' ':   data_buf[i] = ' '
		}
	}
	return data_buf
}



func Int2String( from int ) string{
	str := strconv.Itoa(from)
	return str
}

func IntFormat2String( from int ) string{
	str := fmt.Sprintf("%d",from)
	return str
}


func RandInt64(min,max int64) int64{
	maxBigInt:=big.NewInt(max)
	i,_:=rand.Int(rand.Reader,maxBigInt)
	if i.Int64()<min{
		RandInt64(min,max)
	}
	return i.Int64()
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}


func Convert2Byte( invoker_ret interface{}) []byte{

	var data_buf []byte

	switch invoker_ret.(type) {      //多选语句switch
	case string:
		data_buf = []byte( invoker_ret.(string) )
	case int:
		data_buf = []byte(strconv.Itoa( invoker_ret.(int) ))
	case int64:
		data_buf  = Int64ToBytes( invoker_ret.(int64) )
	case float32:
		data_buf  = Float32ToByte( invoker_ret.(float32) )
	case float64:
		data_buf  = Float64ToByte( invoker_ret.(float64) )
	case []string:
		data_buf,_ = json.Marshal( invoker_ret.([]string) )
	case map[string]string:
		tmp,err := json.Marshal( invoker_ret.(map[string]string) )
		if err!=nil{
			fmt.Println( "json.Marshal( invoker_ret err :", err.Error())
			return data_buf
		}
		data_buf = tmp
	case map[string]interface{}:
		tmp,err := json.Marshal( invoker_ret.(map[string]interface{}) )
		if err!=nil{
			fmt.Println( "json.Marshal( invoker_ret err :", err.Error())
			return data_buf
		}
		data_buf = tmp
	}
	return data_buf
}


func GetJsonChildObj( str string ,key string) ( string, error){

	reg := regexp.MustCompile(  `"`+key+`":\{([\w\W]*)\}\s*,`)
	reg_ret := reg.FindAllStringSubmatch(str, -1)
	if len(reg_ret)>0 {
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			return "{"+reg_ret[0][1]+"}",nil
		}
	}else{
		reg := regexp.MustCompile(`"`+key+`":\{([\w\W]*)\}\s*\}$`)
		reg_ret = reg.FindAllStringSubmatch(str, -1)
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			return "{"+reg_ret[0][1]+"}",nil
		}
	}
	return "{}", errors.New("match result no found")
}

func GetJsonChildArray( str string ,key string) ( string, error){

 	reg := regexp.MustCompile(  `"`+key+`":\[([\w\W]*)\]\s*,`)
	reg_ret := reg.FindAllStringSubmatch(str, -1)
	if len(reg_ret)>0 {
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			return "["+reg_ret[0][1]+"]",nil
		}
	}else{
		reg := regexp.MustCompile(`"`+key+`":\[([\w\W]*)\]\s*\}$`)
		reg_ret = reg.FindAllStringSubmatch(str, -1)
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			return "["+reg_ret[0][1]+"]",nil
		}
	}
	return "", errors.New("match result no found")
}

func GetJsonChildStr( str string ,key string) ( string, error){

 	reg := regexp.MustCompile(  `"`+key+`":"([\w\W]*)"\s*,`)
	reg_ret := reg.FindAllStringSubmatch(str, -1)
	if len(reg_ret)>0 {
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			return `"`+reg_ret[0][1]+`"`,nil
		}
	}else{
		reg := regexp.MustCompile(`"`+key+`":"([\w\W]*)"\s*\}$`)
		reg_ret = reg.FindAllStringSubmatch(str, -1)
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			return `"`+reg_ret[0][1]+`"`,nil
		}
	}
	return "[]", errors.New("match result no found")
}


func GetJsonChildBool( str string ,key string) ( bool, error){

 	reg := regexp.MustCompile(  `"`+key+`":([\w\W]*)\s*,`)
	reg_ret := reg.FindAllStringSubmatch(str, -1)
	if len(reg_ret)>0 {
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			if strings.ToLower(reg_ret[0][1])=="false"{
				return false,nil
			}else{
				return true,nil
			}
		}
	}else{
		reg := regexp.MustCompile(`"`+key+`":([\w\W]*)\s*\}$`)
		reg_ret = reg.FindAllStringSubmatch(str, -1)
		if( len(reg_ret[0])>1 ) {
			fmt.Printf("reg::::::::::::::::::%s\n",  reg_ret[0][1])
			if strings.ToLower(reg_ret[0][1])=="false"{
				return false,nil
			}else{
				return true,nil
			}
		}
	}
	return false, errors.New("match result no found")
}






