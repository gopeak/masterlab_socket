package cron

import (
	"flag"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/robfig/cron"
	"log"
	"masterlab_socket/util"
	"os/exec"
)

type Schedule struct {
}

// https://godoc.org/github.com/robfig/cron
//https://www.cnblogs.com/zuxingyu/p/6023919.html
func (this *Schedule) Run() {

	// "C:/gopath/src/masterlab_socket/cron/cron.json"
	var cron_file string
	flag.StringVar(&cron_file, "s", "cron.json", "cron.json's file path")
	fmt.Println("cron_file:", cron_file)
	cron_json, err := util.ReadAll(cron_file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	json_object, _ := jason.NewObjectFromBytes([]byte(cron_json))
	desc, _ := json_object.GetString("desc")
	log.Println("desc:", desc)

	children, _ := json_object.GetObjectArray("schedule")
	c := cron.New()
	for i, element := range children {
		log.Println(i, element)
		exe_bin, err := element.GetString("exe_bin")
		if err != nil {
			log.Println("exe_bin:", err.Error())
		}
		exp, err := element.GetString("exp")
		if err != nil {
			log.Println("exp:", err.Error())
		}
		file, err := element.GetString("file")
		if err != nil {
			log.Println("exp:", err.Error())
		}
		arg, _ := element.GetString("arg")
		if err != nil {
			log.Println("arg:", err.Error())
		}
		err = c.AddFunc(exp, func() {
			sh := fmt.Sprintf("%s %s %s", exe_bin, file, arg)
			log.Println(i, sh)
			out := this.Cmd(exe_bin, arg, file)
			log.Println(string(out))
		})
		if err != nil {
			log.Println(err.Error())
		}
	}
	c.Start()

	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
	select {}
}

func (this *Schedule) Cmd(exeBin, arg, cmd string) []byte {

	var out []byte
	var err  error
	if arg!=""{
		out, err  = exec.Command(exeBin, arg, cmd).Output()
	}else{
		out, err  = exec.Command(exeBin, cmd).Output()
	}
	if err != nil {
		fmt.Println("some error found:" + err.Error())
	}
	return out
}
