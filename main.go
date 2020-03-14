//
//  main
//
package main

import (
	"masterlab_socket/area"
	"masterlab_socket/cmd"
	"masterlab_socket/connector"
	"masterlab_socket/cron"
	"masterlab_socket/global"
	"masterlab_socket/golog"
	"masterlab_socket/hub"
	"masterlab_socket/lib/syncmap"
	_ "net/http/pprof"
	"runtime"
)



// 初始化全局变量
func initGlobal() {

	global.SumConnections = 0
	global.Qps = 0

	// 先在global声明,再使用make函数创建一个非nil的map，nil map不能赋值
	global.AuthCmds = make([]string, 0)
	global.UserSessions = syncmap.New()
	global.SingleMode = global.Config.SingleMode
	global.AuthCmds = global.Config.Connector.AuthCcmds
	area.UserJoinedAreas = syncmap.New()
	global.InitWorkerAddr()
}

/**
 * 框架启动
 */
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd.Execute();

	// 初始化配置和全局变量
	global.InitConfig()
	golog.InitLogger()
	initGlobal()

	// 前端的socket服务
	frontSocket := new(connector.Connector)
	go frontSocket.Socket("", global.Config.Connector.SocketPort)
	go frontSocket.Websocket("", global.Config.Connector.WebsocketPort)

	// 开启hub服务器
	hubObj := new(hub.Hub)
	go hubObj.Server()

	// 预创建多个场景
	go area.InitConfig()

	// 计划任务
	schedule := new(cron.Schedule)
	go schedule.Run()

	golog.Info("Server started!")

	//go build;
	select {}

}
