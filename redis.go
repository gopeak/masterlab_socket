package main

import (
	"fmt"
	"masterlab_socket/global"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Pool      *redis.Pool
	RedisConn *redis.Conn
)

func RedisInit() {

	fmt.Println(global.Config.Object.RedisPort)
	RedisConn, err := redis.Dial("tcp", global.Config.Object.RedisHost+`:`+string(global.Config.Object.RedisPort))
	//defer RedisConn.Close()
	if err != nil {
		fmt.Println("Redis server connect failed,", err.Error())
		return
	}

	//n, err := RedisConn.Do("Set", "aaa", "vvvvvvvvv")
	//fmt.Println(n, err)

	data := &Session{
		``,
		"{}",
		true,  // 登录成功
		false, // 是否被踢出
		``,
		time.Now().Unix(), //加入时间
		time.Now().Unix(),
	}
	RedisConn.Do("Set", "111", data)
	v1, err := redis.String(RedisConn.Do("Get", "111"))
	fmt.Println(v1, err)

	v, err := redis.String(RedisConn.Do("Get", "aaa"))
	fmt.Println(v, err)

	Pool = NewPool(global.Config.Object.RedisHost+":"+string(global.Config.Object.RedisPort), global.Config.Object.RedisPassword)
	cc := Pool.Get()
	fmt.Println("cc:", cc)

	v2, err := redis.String(cc.Do("Get", "aaa"))
	fmt.Println(v2, err)

}

func NewPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Set(key string, args ...interface{}) (bool, error) {

	Pool = NewPool(global.Config.Object.RedisHost+":"+string(global.Config.Object.RedisPort), global.Config.Object.RedisPassword)
	cc := Pool.Get()
	return redis.Bool(cc.Do("Set", `ueli/`+key, args))

}

func Get(key string) (string, error) {

	Pool = NewPool(global.Config.Object.RedisHost+":"+string(global.Config.Object.RedisPort), global.Config.Object.RedisPassword)
	cc := Pool.Get()
	return redis.String(cc.Do("Get", `ueli/`+key))

}

func _Delete(key string) (bool, error) {

	Pool = NewPool(global.Config.Object.RedisHost+":"+string(global.Config.Object.RedisPort), global.Config.Object.RedisPassword)
	cc := Pool.Get()
	return redis.Bool(cc.Do("Delete", `ueli/`+key))

}
