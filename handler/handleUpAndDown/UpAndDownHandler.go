package handleUpAndDown

import (
    "ai_local/hardWare/rs485/rs485Constants"
    "ai_local/hardWare/rs485/rs485Helper"
    "encoding/hex"
    "github.com/astaxie/beego"
    "github.com/pkg/errors"
    "github.com/tarm/serial"
    "time"
)

var Rs485Port *serial.Port

func InitRs485Port() {
    comName := "/dev/ttyS1"
    baud := 38400
    config := serial.Config{Name: comName, Baud: baud, ReadTimeout: time.Duration(2) * time.Second}
    beego.Info("begin open Rs485Port")
    p, e := serial.OpenPort(&config)
    if e != nil {
        beego.Error(e.Error())
        time.Sleep(time.Duration(5) * time.Second)
        p, e = serial.OpenPort(&config)
    }
    Rs485Port = p
}

func ReadFromPort() (result string, err error) {

    if Rs485Port == nil {
        err = errors.New("nil Rs485Port")
        return
    }

    var rb []byte
    var b = make([]byte, 1024)
    var i = 1
    var e error
    for i > 0 {
        i, e = Rs485Port.Read(b)
        if e != nil {
            err = e
            return
        }
        rb = append(rb, b[:i]...)
    }

    return hex.EncodeToString(rb), nil
}

func UpOrDown() error {
    downE := AllDown()
    if downE != nil {
        beego.Error(downE)
    }

    time.Sleep(time.Duration(1) * time.Minute)

    upE := AllUp()
    if upE != nil {
        beego.Error(upE)
    }

    return nil
}

func AllUp() error {

    var tryTimes = 0
    var result string
    var err error
    result, err = GetStatus()
    for err != nil && tryTimes < 3 {
        result, err = GetStatus()
        tryTimes++
        time.Sleep(time.Duration(1) * time.Second)
    }
    beego.Debug("result ", result)

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

    var tryTimes = 0
    var result string
    var err error
    result, err = GetStatus()
    for err != nil && tryTimes < 3 {
        result, err = GetStatus()
        tryTimes++
        time.Sleep(time.Duration(1) * time.Second)
    }
    beego.Debug("result ", result)

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

func GetStatus() (result string, err error) {
    if Rs485Port == nil {
        err = errors.New("nil Rs485Port")
        return
    }
    bytes, e := rs485Helper.GetCommandBytesFromStr(rs485Constants.QueryAll)
    if e != nil {
        err = e
        return
    }

    go func(result *string, err *error) {
        *result, *err = ReadFromPort()
    }(&result, &err)

    _, e = Rs485Port.Write(bytes) //上
    if e != nil {
        err = e
        return
    }

    time.Sleep(time.Duration(100) * time.Millisecond)
    return
}
