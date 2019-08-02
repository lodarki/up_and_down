package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"path"
	"runtime"
	"time"
)

func LogDebug(v ...interface{}) {
	beego.Debug(v...)
}

func LogInfo(v ...interface{}) {
	beego.Info(v...)
}

func LogWarn(v ...interface{}) {
	beego.Warn(v...)
}

func LogError(v ...interface{}) {
	beego.Error(v...)
}

func LogDev(msg string) {
	logMessage("V", msg, 2)
}

func LogCustom(prefix string, msg string) {
	logMessage(prefix, msg, 2)
}

func logMessage(prefix string, msg string, skip int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	m := fmt.Sprintf("%v [%v] [%v:%v] %v", DateTimeMStr(time.Now()), prefix, filename, line, msg)
	fmt.Println(m)
}
