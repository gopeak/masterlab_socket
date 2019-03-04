package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/weekface/mgorus"
)

func main() {

	log := golog.New()
	hooker, err := mgorus.NewHooker("localhost:27017", "db", "collection")
	if err == nil {
		golog.Hooks.Add(hooker)
	} else {
		fmt.Println("mongodb err:", err)
	}

	golog.WithFields(golog.Fields{
		"name": "zhangsan1215555555551155555",
		"age":  28225552,
	}).Info("Hello world!")

	golog.Warn("2222222222221111122222")
}
