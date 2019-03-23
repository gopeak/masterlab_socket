// main loop

package golog

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"masterlab_socket/global"
	"os"
	"runtime"
)

var Log *log.Logger


// 初始化日志设置
func InitLogger() {

	if runtime.GOOS != "windows" {
		log.SetFormatter(&log.JSONFormatter{})

	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

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

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {

	log.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {

	log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	log.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
