package worker

import (
	"fmt"
	_ "fmt"
	"github.com/antonholmquist/jason"
	"gopkg.in/gomail.v2"
	"strconv"
)

func (this TaskType) Mail() ReturnType {

	//sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	// 获取数据
	fmt.Println("Mail this.Data:", string(this.Data))
	json_obj, err := jason.NewObjectFromBytes(this.Data)
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:" + err.Error()}
		return ret
	}
	host, err := json_obj.GetString("host")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:host not found"}
		return ret
	}
	port, err := json_obj.GetString("port")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:port not found"}
		return ret
	}
	user, err := json_obj.GetString("user")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:user not found"}
		return ret
	}
	password, err := json_obj.GetString("password")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:password not found"}
		return ret
	}
	from, err := json_obj.GetString("from")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:from not found"}
		return ret
	}
	to, err := json_obj.GetString("to")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:to not found"}
		return ret
	}
	cc, err := json_obj.GetString("cc")
	if err != nil {
		cc = ""
	}
	cc_name, err := json_obj.GetString("cc_name")
	if err != nil {
		cc_name = ""
	}
	subject, err := json_obj.GetString("subject")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:subject not found"}
		return ret
	}
	body, err := json_obj.GetString("body")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:body not found"}
		return ret
	}
	fmt.Println(host, port, user, password, from, to, subject, body)
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	if cc != "" {
		m.SetAddressHeader("Cc", cc, cc_name)
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.Attach("D:/timg.jpg")

	port_int, err := strconv.Atoi(port)
	if err != nil {
		ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, "port error"}
		return ret
	}

	d := gomail.NewDialer(host, port_int, user, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, err.Error()}
		return ret
	}
	ret := ReturnType{"ok", "mail", this.ReqHeader.Sid, ""}
	return ret
}
