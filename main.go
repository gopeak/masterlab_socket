//
//  main
//

package main

import (
	"masterlab_socket/global"
	"masterlab_socket/lib/syncmap"
	_ "net/http/pprof"
	"runtime"
)


// 所有的场景名称列表
var Areas = make([]string, 0, 1000)

// 场景集合
var AreasMap *syncmap.SyncMap

// 一个全局的场景
var GlobalArea   *AreaType

// 所有的用户连接对象
var AllConns *syncmap.SyncMap
var AllWsConns *syncmap.SyncMap

// 用户加入过的场景列表
var UserJoinedAreas *syncmap.SyncMap


// 初始化全局变量
func initGlobal() {

	global.SumConnections = 0
	global.Qps = 0

	// 先在global声明,再使用make函数创建一个非nil的map，nil map不能赋值
	global.AuthCmds = make([]string, 0)
	global.UserSessions = syncmap.New()
	global.SingleMode = global.Config.SingleMode
	global.AuthCmds = global.Config.Connector.AuthCcmds
	UserJoinedAreas = syncmap.New()
	global.InitWorkerAddr()
}

/**
 * 框架启动
 */
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置和全局变量
	global.InitConfig()
	InitLogger()
	initGlobal()

	// 前端的socket服务
	go SocketConnector("", global.Config.Connector.SocketPort)
	go WebsocketConnector("", global.Config.Connector.WebsocketPort)

	// 开启hub服务器
	hub := new(Hub)
	go hub.Server()

	// 预创建多个场景
	go AreaInitConfig()

	// 计划任务
	schedule := new(Schedule)
	go schedule.Run()

	LogInfo("Server started!")

	// C:\gopath\mongodb\bin\mongod.exe --dbpath=C:\gopath\mongodb\data
	// D:\soft\MongoDB\bin\mongod.exe --dbpath=D:\soft\MongoDB\data
	select {}

}
