package worker

import (
	"fmt"
	_ "fmt"
	"github.com/json-iterator/go"
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

type MailContent struct {
	Seq         string   `json:"seq"`
	Host        string   `json:"host"`
	Port        string   `json:"port"`
	User        string   `json:"user"`
	Password    string   `json:"password"`
	From        string   `json:"from"`
	FromName    string   `json:"from_name"`
	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Bcc         []string `json:"bcc"`
	ContentType string   `json:"content_type"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attach      string   `json:"attach"`
}

func (this *MailContent) Init() *MailContent {

	this.Seq = ""
	this.Host = ""
	this.Port = ""
	this.User = ""
	this.Password = ""
	this.From = ""
	this.FromName = "Masterlab"
	this.To = []string{}
	this.Cc = []string{}
	this.Bcc = []string{}
	this.ContentType = "html"
	this.Subject = ""
	this.Body = ""
	this.Attach = ""

	return this
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

// 通过非加密方式发送
func sendByNoSSL(mailContent *MailContent) error {
	replyToAddress := ""
	auth := LoginAuth(mailContent.User, mailContent.Password)
	content_type := "Content-Type: text/" + mailContent.ContentType + "; charset=UTF-8"
	noSslHost := fmt.Sprintf("%s:%s", mailContent.Host, mailContent.Port)
	cc_address := strings.Join(mailContent.Cc, ";")
	bcc_address := strings.Join(mailContent.Bcc, ";")
	to_address := strings.Join(mailContent.To, ";")
	msg := []byte("To: " + to_address + "\r\nFrom: Masterlab<" + mailContent.From + ">\r\nSubject: " + mailContent.Subject + "\r\nReply-To: " + replyToAddress + "\r\nCc: " + cc_address + "\r\nBcc: " + bcc_address + "\r\n" + content_type + "\r\n\r\n" + mailContent.Body)

	send_to := MergeSlice(mailContent.To, mailContent.Cc)
	send_to = MergeSlice(send_to, mailContent.Bcc)
	err := smtp.SendMail(noSslHost, auth, mailContent.From, send_to, msg)
	return err
}

// 通过gomail库SSL加密发送
func sendBySSL(mailContent *MailContent) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailContent.FromName+"<"+mailContent.From+">")
	m.SetHeader("To", mailContent.To...)
	m.SetHeader("Subject", mailContent.Subject)
	if mailContent.ContentType == "html" {
		m.SetBody("text/html", mailContent.Body)
	} else {
		m.SetBody("text/plain", mailContent.Body)
	}
	if mailContent.Attach != "" && util.Exists(mailContent.Attach) {
		m.Attach(mailContent.Attach)
	}
	port_int, err := strconv.Atoi(mailContent.Port)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(mailContent.Host, port_int, mailContent.User, mailContent.Password)
	// Send the email toArr Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (this TaskType) Mail() ReturnType {

	//sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	// 获取数据
	//fmt.Println("Mail this.Data:", string(this.Data))

	mailContent := new(MailContent).Init()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json_err := json.Unmarshal(this.Data, &mailContent)
	if json_err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "json err:" + json_err.Error()}
		return ret
	}
	//fmt.Println(mailContent)
	to_address := strings.Join(mailContent.To, ";")
	db := new(lib.Mysql)
	_, err := db.ShortConnect()
	if err != nil {
		ret := ReturnType{"failed", "failed", this.ReqHeader.Sid, "数据库连接失败:" + err.Error()}
		return ret
	}
	// 获取时间戳
	now := time.Now().Unix()
	timestamp := strconv.FormatInt(now, 10)
	create_time_nano := fmt.Sprintf("%v", time.Now().UnixNano())
	if mailContent.Seq == "" {
		mailContent.Seq = create_time_nano
	}
	if mailContent.Port == "465" || mailContent.Port == "995" {
		// Send the email toArr Bob, Cora and Dan.
		if err := sendBySSL(mailContent); err != nil {
			fmt.Println("Send mail error!", err.Error())
			ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, err.Error()}
			_, err = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
				mailContent.Seq, mailContent.Subject, to_address, "error", timestamp, err.Error())
			return ret
		}
	} else {
		err := sendByNoSSL(mailContent)
		if err != nil {
			fmt.Println("Send mail error!", err.Error())
			ret := ReturnType{"filed", "mail", this.ReqHeader.Sid, err.Error()}
			_, err = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
				mailContent.Seq, mailContent.Subject, to_address, "error", timestamp, err.Error())
			return ret
		}
	}

	_, err = db.Insert("REPLACE INTO   `main_mail_queue` (seq, `title`, `address`, `status`, `create_time`, `error`) VALUES ( ?,?,?,?,?,?)",
		mailContent.Seq, mailContent.Subject, to_address, "done", timestamp, "")
	if err != nil {
		ret := ReturnType{"failed", "mail", this.ReqHeader.Sid, "sql replace into error"}
		return ret
	}
	ret := ReturnType{"ok", "mail", this.ReqHeader.Sid, ""}
	return ret
}
