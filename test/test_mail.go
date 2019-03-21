package main
import (
    "fmt"
    "net/smtp"
    "strings"
)
//如果go语言的版本为1.9.2，出现错误提示:“unencrypted connection”，因为此版本需要加密认证，可采用LOGIN认证，则需要增加以下内容：
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
    msg := []byte("To: " + to_address + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\nReply-To: " + replyToAddress + "\r\nCc: " + cc_address + "\r\nBcc: " + bcc_address + "\r\n" + content_type + "\r\n\r\n" + body)

    send_to := MergeSlice(to, cc)
    send_to = MergeSlice(send_to, bcc)
    err := smtp.SendMail(host, auth, user, send_to, msg)
    return err
}
func main() {
    user := "sender@smtp.masterlab.vip"
    password := "MasterLab123Pwd"
    host := "smtpdm.aliyun.com:25"
    to := []string{"121642038@qq.com","weichaoduo@163.com"}
    cc := []string{"79720699@qq.com"}
    bcc := []string{"79720699@qq.com","79720699@qq.com"}
    subject := "test Golang to sendmail"
    mailtype :="html"
    replyToAddress:="sender@smtp.masterlab.vip"
    body := `
        <html>
        <body>
        <h3>
        "Test send to email"
        </h3>
        </body>
        </html>
        `
    fmt.Println("send email")
    err := SendToMail(user, password, host, subject, body, mailtype, replyToAddress, to, cc, bcc)
    if err != nil {
        fmt.Println("Send mail error!")
        fmt.Println(err)
    } else {
        fmt.Println("Send mail success!")
    }
}
func MergeSlice(s1 []string, s2 []string) []string {
    slice := make([]string, len(s1)+len(s2))
    copy(slice, s1)
    copy(slice[len(s1):], s2)
    return slice
}