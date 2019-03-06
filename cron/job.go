package main

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/robfig/cron"
	"log"
	"masterlab_socket/util"
	"os/exec"
)

// https://godoc.org/github.com/robfig/cron
//https://www.cnblogs.com/zuxingyu/p/6023919.html
func main() {
	c := cron.New()

	exampleJSON, err := util.ReadAll("C:/gopath/src/masterlab_socket/cron/cron.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	v, _ := jason.NewObjectFromBytes([]byte(exampleJSON))

	desc, _ := v.GetString("desc")
	exe_bin, _ := v.GetString("exe_bin")
	log.Println("desc:", desc)
	log.Println("exe_bin:", exe_bin)

	children, _ := v.GetObjectArray("schedule")
	for i, element := range children {
		log.Println(i, element)
		exp, err:= element.GetString("exp")
		if err!=nil{
			log.Println("exp:", err.Error())
		}
		file, err := element.GetString("file")
		if err!=nil{
			log.Println("exp:", err.Error())
		}
		arg, _ := element.GetString("arg")
		if err!=nil{
			log.Println("arg:", err.Error())
		}
		err = c.AddFunc(exp, func() {
			sh := fmt.Sprintf("%s %s %s", exe_bin, file, arg)
			log.Println(i, sh)
			out := Cmd(sh, false)
			log.Println(string(out))
		})
		if err!=nil{
			log.Println(err.Error())
		}
	}
	c.Start()

	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
	select {}
}

func Cmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	}
}
