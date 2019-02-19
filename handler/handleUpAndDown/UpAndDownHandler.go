package handleUpAndDown

import (
	"ai_local/hardWare/rs485/rs485Constants"
	"ai_local/hardWare/rs485/rs485Helper"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tarm/serial"
	"time"
)

var Rs485Port *serial.Port

func init() {
	comName := "/dev/ttyS1"
	baud := 38400
	config := serial.Config{Name: comName, Baud: baud}
	beego.Info("begin open Rs485Port")
	p, e := serial.OpenPort(&config)
	if e != nil {
		beego.Error(e.Error())
		time.Sleep(time.Duration(5) * time.Second)
		p, e = serial.OpenPort(&config)
	}
	Rs485Port = p

	go ReadFromPort()
}

func ReadFromPort() {
	for {
		if Rs485Port == nil {
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		var b = make([]byte, 1024)
		i, e := Rs485Port.Read(b)
		if e != nil {
			beego.Error(e)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		beego.Debug(fmt.Sprintf("电机状态响应：%v", hex.EncodeToString(b[:i])))
	}
}

func UpOrDown() error {
	upE := AllUp()
	if upE != nil {
		beego.Error(upE)
	}
	go func() {
		time.Sleep(time.Duration(5) * time.Second)
		GetStatus()
	}()
	time.Sleep(time.Duration(1) * time.Minute)
	downE := AllDown()
	if downE != nil {
		beego.Error(downE)
	}
	go func() {
		time.Sleep(time.Duration(5) * time.Second)
		GetStatus()
	}()
	return nil
}

func AllUp() error {
	upStr := "9a:00:01:00:0a:dd:d6"
	beego.Debug("全上")
	bytes, e := rs485Helper.GetCommandBytesFromStr(upStr)
	if e != nil {
		beego.Error(e)
		return e
	}
	_, e = Rs485Port.Write(bytes) //上
	if e != nil {
		beego.Error(e)
		return e
	}
	return nil
}

func AllDown() error {
	downStr := "9a:00:01:00:0a:ee:e5" //全下
	beego.Debug("全下")
	bytes, e := rs485Helper.GetCommandBytesFromStr(downStr)
	if e != nil {
		beego.Error(e)
		return e
	}
	_, e = Rs485Port.Write(bytes) //下
	if e != nil {
		beego.Error(e)
		return e
	}
	return nil
}

func GetStatus() error {
	if Rs485Port == nil {
		return nil
	}
	bytes, e := rs485Helper.GetCommandBytesFromStr(rs485Constants.QueryAll)
	if e != nil {
		return e
	}
	_, e = Rs485Port.Write(bytes) //上
	if e != nil {
		return e
	}
	return nil
}
