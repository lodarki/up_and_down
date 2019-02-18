package handleUpAndDown

import (
	"ai_local/hardWare/rs485/rs485Constants"
	"ai_local/hardWare/rs485/rs485Helper"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tarm/serial"
	"time"
)

var port *serial.Port

func init() {
	comName := "/dev/ttyS0"
	baud := 38400
	config := serial.Config{Name: comName, Baud: baud}
	beego.Info("begin open port")
	p, e := serial.OpenPort(&config)
	for e != nil {
		beego.Error(e.Error())
		time.Sleep(time.Duration(5) * time.Second)
		p, e = serial.OpenPort(&config)
	}
	port = p

	go ReadFromPort()
}

func ReadFromPort() {
	for {
		var b = make([]byte, 1024)
		i, e := port.Read(b)
		if e != nil {
			beego.Error(e)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		beego.Debug(fmt.Sprintf("电机状态响应：%v", string(b[:i])))
	}
}

func UpOrDown() error {
	downStr := "9a:00:01:00:0a:ee:e5" //全下
	upStr := "9a:00:01:00:0a:dd:d6"   //全上

	bytes, e := rs485Helper.GetCommandBytesFromStr(downStr)
	if e != nil {
		return e
	}
	_, e = port.Write(bytes) //下
	if e != nil {
		return e
	}
	time.Sleep(time.Duration(1) * time.Minute)
	bytes, e = rs485Helper.GetCommandBytesFromStr(upStr)
	if e != nil {
		return e
	}
	_, e = port.Write(bytes) //上
	if e != nil {
		return e
	}

	return nil
}

func GetStatus() error {
	bytes, e := rs485Helper.GetCommandBytesFromStr(rs485Constants.QueryAll)
	if e != nil {
		return e
	}
	_, e = port.Write(bytes) //上
	if e != nil {
		return e
	}
	return nil
}