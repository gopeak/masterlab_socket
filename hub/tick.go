// main loop

package hub

import (
	"fmt"
	"time"
	"net"
	"masterlab_socket/global"
	"masterlab_socket/lib/syncmap"
	"masterlab_socket/golog"
	"masterlab_socket/area"
	"github.com/garyburd/redigo/redis"
)

func tick() {
	timer := time.Tick(100 * time.Millisecond)
	for now := range timer {
		// entity updates (you could use now for physic engine calculs)
		// this is called every 100 millisecondes
		// playerFactory.Update()
		fmt.Println("now", now)
	}
}

func TickSyncSession() {

	redisc, err := redis.Dial("tcp", global.Config.Object.RedisHost+`:`+string(global.Config.Object.RedisPort))
	//defer redisc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	timer := time.Tick(1 * time.Second)
	var LastSessions *syncmap.SyncMap

	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		/*var UserSessions = map[string]*area.Session{}
		for item := range global.SyncUserSessions.IterItems() {
			UserSessions[item.Key] = item.Value.(*area.Session)
		}
		js1, _ := json2.Marshal(UserSessions)
		*/
		if LastSessions != global.UserSessions {
			redisc.Do("Set", "masterlab_socket/user_session", global.UserSessions)
			redisc.Flush()
			LastSessions = global.UserSessions
		}

	}
}

func LoadSessionFromRedis() {

	redisc, err := redis.Dial("tcp", global.Config.Object.RedisHost+`:`+string(global.Config.Object.RedisPort))
	//defer redisc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	reply, err_get := redisc.Do("Get", "masterlab_socket/user_session")
	if err_get != nil {
		fmt.Println(err_get)
		return
	}
	fmt.Println("GET morego/user_session ", reply)
	if reply != nil {
		global.UserSessions = reply.(*syncmap.SyncMap)
		var UserSessions = map[string]*area.Session{}
		for item := range global.UserSessions.IterItems() {
			UserSessions[item.Key] = item.Value.(*area.Session)
			fmt.Println(UserSessions[item.Key].Sid)
		}
	}
}

func TickWorkerServer() {
	// 先暂停10秒
	time.Sleep(5 * time.Second)
	timer := time.Tick(10 * time.Second)
	for now := range timer {
		//fmt.Println("now", now)
		ch_success := make(chan string, 0)
		for _, data := range global.Config.ToWorker.Servers {
			go func(data []string) {
				worker_host := data[0]
				worker_port_str := data[1]
				ip_port := worker_host + ":" + worker_port_str

				//fmt.Println("tcpAddr: ",index," ", ip_port)
				conn, err_req := net.DialTimeout("tcp", ip_port, 5*time.Second)
				if err_req != nil {
					golog.Error("检测到 workerserver:", ip_port, " 连接异常!", now)
					for i, addr := range global.WorkerServers {
						if addr == ip_port {
							global.WorkerServers = append(global.WorkerServers[:i], global.WorkerServers[i+1:]...)
						}
					}
					ch_success <- ip_port + err_req.Error()
				} else {
					exist := false
					for _, addr := range global.WorkerServers {
						if addr == ip_port {
							exist = true
							break
						}
					}
					if !exist {
						global.WorkerServers = append(global.WorkerServers, ip_port)
					}
					ch_success <- ip_port + "ok"
				}
				//fmt.Println("result: ", ip_port, " ok")
				//req_str:= fmt.Sprintf("%d||%s||%s||%d||%s\n", protocol.TypePing, "Ping", "", 0, "")
				//conn.Write([]byte(req_str))
				conn.Close()
			}(data)
		}
		sum := 0
		for i := 0; i < len(global.Config.ToWorker.Servers)+1; i++ {
			select {
			case <-ch_success:
				//fmt.Println("recv_result:", r)
				sum++
				if sum == len(global.Config.ToWorker.Servers) {
					break
				}

			default:
				//fmt.Printf(".")
				time.Sleep(10 * time.Millisecond)
			}
		}
		//fmt.Println("sum:", sum)

	}
}
