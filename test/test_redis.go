package main

import (
	"flag"
	"fmt"
	"time"
	"github.com/garyburd/redigo/redis"
)

var (
	pool          *redis.Pool
	redisServer   = flag.String("redisServer", "127.0.0.1:6379", "")
	redisPassword = flag.String("redisPassword", "", "")
)

func main3() {

	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	n, err := c.Do("Set", "aaa", "vvvvvvvvv")
	fmt.Println(n, err)
	v, err := redis.String(c.Do("Get", "aaa"))
	fmt.Println(v, err)

	flag.Parse()
	pool = newPool(*redisServer, *redisPassword)
	cc := pool.Get()
	fmt.Println("cc:", cc)

	v2, err := redis.String(cc.Do("Get", "aaa"))
	fmt.Println(v2, err)

}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			//     c.Close()
			//     return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
