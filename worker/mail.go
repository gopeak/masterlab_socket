package worker

import (
	"fmt"
	_ "fmt"
	"github.com/antonholmquist/jason"
	"gopkg.in/gomail.v2"
	"masterlab_socket/lib"
	"masterlab_socket/util"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type loginAuth struct {
	username, password string
}
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// return "LOGIN", []byte{}, nil
	return "LOGIN", []byte(a.username), nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

func SendToMail(user, password, host, subject, body, mailtype, replyToAddress string, to, cc, bcc []string) error {
	// hp := strings.Split(host, ":")
	//auth := smtp.PlainAuth("", user, password, hp[0])
	auth := LoginAuth(user, password)
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	cc_address := strings.Join(cc, ";")
	bcc_address := strings.Join(bcc, ";")
	to_address := strings.Join(to, ";")
	msg := []byte("To: " + to_address + "\r\nFrom: Masterlab<" + user + ">\r\nSubject: " + subject + "\r\nReply-To: " + replyToAddress + "\r\nCc: " + cc_address + "\r\nBcc: " + bcc_address + "\r\n" + content_type + "\r\n\r\n" + body)

	send_to := MergeSlice(to, cc)
	send_to = MergeSlice(send_to, bcc)
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
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
	cc, err := json_obj.GetStringArray("cc")
	if err != nil {
		cc = nil
	}
	bcc, err := json_obj.GetStringArray("bcc")
	if err != nil {
		bcc = nil
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
	m.SetHeader("From", from, "Masterlab")
	toStr := ""
	for _, to := range toArr {
		//m.SetHeader("To", to)
		m.SetAddressHeader("To", to, to)
		toStr =   fmt.Sprintf("%s;%s", toStr, to)
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
	now := time.Now().Unix()
	timestamp := strconv.FormatInt(now,10)
	create_time_nano := fmt.Sprintf("%v", time.Now().UnixNano());

	var seq string
	seq, err = json_obj.GetString("seq")
	if err != nil {
		seq = create_time_nano
	}
	if port=="465" || port=="995"{
		d := gomail.NewDialer(host, port_int, user, password)
		// Send the email toArr Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, err.Error()}
			_, err  = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
				seq, subject, toStr, "error",timestamp, err.Error())
			return ret
		}
	}else{
		m.Reset();
		noSslHost := fmt.Sprintf("%s:%s", host, port)
		err := SendToMail(user, password, noSslHost, subject, body, "html", "", toArr, cc, bcc)
		if err != nil {
			fmt.Println("Send mail error!")
			ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, err.Error()}
			_, err  = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
				seq, subject, toStr, "error",timestamp, err.Error())
			return ret
		}
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