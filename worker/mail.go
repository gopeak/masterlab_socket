package worker

import (
	"fmt"
	_ "fmt"
	"github.com/antonholmquist/jason"
	"gopkg.in/gomail.v2"
	"masterlab_socket/lib"
	"masterlab_socket/util"
	"strconv"
	"time"
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
	toArr, err := json_obj.GetStringArray("to")
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:toArr not found"}
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

	fmt.Println(host, port, user, password, from, toArr, subject, body)
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	toStr := ""
	for _, to := range toArr {
		m.SetHeader("To", to)
		toStr =   fmt.Sprintf("%s,%s", toStr, to)
	}

	if cc != "" {
		m.SetAddressHeader("Cc", cc, cc_name)
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.Attach("D:/timg.jpg")
	var attach string
	attach, err = json_obj.GetString("attach")
	if err != nil {
		attach = ""
	}
	if attach!="" && util.Exists(attach){
		m.Attach(attach)
	}
	port_int, err := strconv.Atoi(port)
	if err != nil {
		ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, "port error"}
		return ret
	}


	
	db := new(lib.Mysql)
	_, err = db.ShortConnect()
	if err != nil {
		ret := ReturnType{ "failed","failed" ,this.ReqHeader.Sid, "数据库连接失败:" + err.Error() }
		return ret
	}
	// 获取时间戳
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	create_time_nano := timestamp[:10]

	var seq string
	seq, err = json_obj.GetString("seq")
	if err != nil {
		seq = create_time_nano
	}

	d := gomail.NewDialer(host, port_int, user, password)
	// Send the email toArr Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, err.Error()}
		_, err  = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
			seq, subject, toStr, "error",timestamp, err.Error())
		return ret
	}

	_, err  = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
		seq, subject, toStr, "done",timestamp, "")
	if err != nil {
		ret := ReturnType{"failed", "mail", this.ReqHeader.Sid, "sql replace into error"}
		return ret
	}

	ret := ReturnType{"ok", "mail", this.ReqHeader.Sid, ""}
	return ret
}
