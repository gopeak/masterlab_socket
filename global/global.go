/**
 *  定义全局变量
 *
 */

package global

import (
	"github.com/robfig/cron"
	"masterlab_socket/lib/syncmap"
)

const (
	ERROR_PACKET_RATES     = `Packet rate limit`
	ERROR_MAX_CONNECTIONS  = `Max connection limit`
	ERROR_RESPONSE          = `RecvMessage error`
	DISBALE_RESPONSE        = `Server has been stopped!`
)
// 全局配置变量
var Config configType

// 服务器当前状态
var AppConfig = & Appconfig{}

var WorkerServers = make([]string, 0, 1000)

var SumConnections int32

var Qps int64

//  用户会话对象
var  UserSessions *syncmap.SyncMap

// 是否为单机运行模式
var SingleMode bool

// 用户认证的命令
var AuthCmds []string

// 定时任务
var Crons = map[string]*cron.Cron{}


func IsAuthCmd( cmd string ) bool {

	//fmt.Println( "global.AuthCmds:",AuthCmds )
	for _,c:= range AuthCmds{
		if( c==cmd ){
			return true
		}
	}
	return false

}