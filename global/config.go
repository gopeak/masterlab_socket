package global

import (
	"fmt"
	"math/rand"
	"masterlab_socket/lib/BurntSushi/toml"
	"flag"
)


type Appconfig struct {
	Enable int64 //  Listen to clients
	Status string
}


type configType struct {
	Name         string
	Enable       bool
	Status       string
	Version      string
	Loglevel     string
	SingleMode   bool	  `toml:"single_mode"`
	Log          log          `toml:"log"`
	Admin        admin        `toml:"admin"`
	Connector    connector    `toml:"connector"`
	Object       object       `toml:"object"`
	ToWorker     toWorker 	  `toml:"worker"`
	Hub          hub          `toml:"hub"`
	Area         area         `toml:"area"`
}

type log struct {
	LogLevel      string `toml:"log_level"`
	LogBehindType string `toml:"log_behind_type"`
	MongodbHost   string `toml:"mongodb_host"`
	MongodbPort   string `toml:"mongodb_port"`
}
type admin struct {
	HttpPort string `toml:"http_port"`
}

type connector struct {
	WebsocketPort     int `toml:"websocket_port"`
	SocketPort        int `toml:"socket_port"`
	MaxConections     int `toml:"max_conections"`
	MaxConntionsIp    int `toml:"max_conntions_ip"`
	MaxPacketRate     int `toml:"max_packet_rate"`
	MaxPacketRateUnit int `toml:"max_packet_rate_unit"`
	AuthCcmds	[]string `toml:"auth_cmds"`
}

type object struct {
	DataType      string `toml:"data_type"`
	RedisHost     string `toml:"redis_host"`
	RedisPort     string `toml:"redis_port"`
	RedisPassword string `toml:"redis_password"`
}

type toWorker struct {
	Servers [][]string `toml:"to_servers"`
}

type hub struct {
	Hub_host string `toml:"hub_host"`
	Hub_port string `toml:"hub_port"`
}

type area struct {
	Init_area []string
}


func InitConfig() {

	var filepath string
	flag.StringVar(&filepath,"c", "config.toml", "config.toml's file path")
	fmt.Println( "filepath:", filepath )
	if _, err := toml.DecodeFile( filepath, &Config); err != nil {
		fmt.Println("toml.DecodeFile error:", err)
		return
	}
}

func GetRandWorkerAddr() string  {
	rand_index := rand.Intn(len(WorkerServers))
	return  WorkerServers[rand_index]
}

func InitWorkerAddr()   {

	for _,data := range Config.ToWorker.Servers {
		worker_host  := data[0]
		worker_port_str  := data[1]
		WorkerServers = append( WorkerServers ,worker_host + ":" + worker_port_str )
	}
}

