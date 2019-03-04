// main loop

package golog

import (
	"fmt"
	"masterlab_socket/global"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Log *log.Logger
var SessionMongo *mgo.Session
var CollectionMongo *mgo.Collection

type MongoLog struct {
	Id_     bson.ObjectId `bson:"_id"`
	Name    string
	Level   string
	File    string
	Line    int
	Message string
	Time    int
}

// 初始化日志设置
func InitLogger() {

	if runtime.GOOS != "windows" {
		log.SetFormatter(&log.JSONFormatter{})

	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	//Log = logrus.New()

	fmt.Println("LogBehindType", global.Config.Log.LogBehindType)

	log.SetOutput(os.Stderr)

	// init logger
	loglevel := global.Config.Log.LogLevel
	if loglevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
	if loglevel == "error" {
		log.SetLevel(log.ErrorLevel)
	}
	if loglevel == "info" {
		log.SetLevel(log.InfoLevel)
	}
	if loglevel == "warn" {
		log.SetLevel(log.WarnLevel)
	}
	if loglevel == "fatal" {
		log.SetLevel(log.FatalLevel)
	}
	if loglevel == "panic" {
		log.SetLevel(log.PanicLevel)
	}
	fmt.Println("logger status : ", loglevel, runtime.GOOS)

}

func log2Mongo(level string, args ...interface{}) {

	SessionMongo, err := mgo.Dial(global.Config.Log.MongodbHost)
	if err != nil {
		panic(err)
	}
	defer SessionMongo.Close()

	// Optional. Switch the session to a monotonic behavior.
	SessionMongo.SetMode(mgo.Monotonic, true)
	CollectionMongo := SessionMongo.DB("gomore").C("logs")
	_, file, line, _ := runtime.Caller(2)
	fmt.Println("runtime.Caller", file, line)
	err = CollectionMongo.Insert(&MongoLog{bson.NewObjectId(), "", level, file, line, fmt.Sprint(args...), int(time.Now().Unix())})
	if err != nil {
		fmt.Println(err)
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {

	//log2Mongo("debug", args...)
	log.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {

	log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	//log2Mongo("info", args...)
	log.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	//log2Mongo("warn", args...)
	log.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	//log2Mongo("Warning", args...)
	log.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	//log2Mongo("error", args...)
	log.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	//log2Mongo("panic", args...)
	log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	//log2Mongo("fatal", args...)
	log.Fatal(args...)
}
