//
//  main
//

package masterlab_socket

import (
	"masterlab_socket/area"
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

	// 初始化配置和全局变量
	global.InitConfig()
	golog.InitLogger()
	initGlobal()

	// 前端的socket服务
	go SocketConnector("", global.Config.Connector.SocketPort)
	go WebsocketConnector("", global.Config.Connector.WebsocketPort)

	// 开启hub服务器
	 go hub.HubServer()

	// 预创建多个场景
	go area.InitConfig()

	// 计划任务
	schedule := new(Schedule)
	go schedule.Run()

	golog.Info("Server started!")

	// C:\gopath\mongodb\bin\mongod.exe --dbpath=C:\gopath\mongodb\data
	// D:\soft\MongoDB\bin\mongod.exe --dbpath=D:\soft\MongoDB\data
	select {}

}
