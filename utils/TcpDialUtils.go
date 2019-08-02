package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

var tcpConnMap map[string]connExpire
var tcpDialOnce sync.Once
var tcpLock sync.Mutex

type connExpire struct {
	Co         net.Conn
	ExpireTime time.Time
}

func init() {
	tcpDialOnce.Do(func() {
		tcpConnMap = make(map[string]connExpire)
	})
}

func getTcpConn(host, port string) (net.Conn, error) {
	k := fmt.Sprintf("%v<>%v", host, port)

	tcpLock.Lock()
	if cep, ok := tcpConnMap[k]; ok {
		exp := time.Now().Add(time.Duration(30) * time.Second)
		cep.ExpireTime = exp
		tcpConnMap[k] = cep
		tcpLock.Unlock()
		return cep.Co, nil
	}
	tcpLock.Unlock()

	conn, e := net.DialTimeout("tcp", host+":"+port, time.Duration(5)*time.Second)
	if e != nil {
		return nil, e
	}

	tcpLock.Lock()
	tcpConnMap[k] = connExpire{Co: conn, ExpireTime: time.Now().Add(time.Duration(30) * time.Second)}
	tcpLock.Unlock()

	go func(c net.Conn, key string) {
		// 延时回收TCP连接
		for {
			time.Sleep(time.Duration(5) * time.Second)
			tcpLock.Lock()
			if !AfterNow(tcpConnMap[k].ExpireTime) {
				cE := c.Close()
				if cE != nil {
					beego.Error(cE)
				}
				delete(tcpConnMap, k)
				tcpLock.Unlock()
				break
			}
			tcpLock.Unlock()
		}
	}(conn, k)

	return conn, nil
}

func TcpSendWithResponse(content []byte, host, port string) ([]byte, error) {
	//beego.Info(fmt.Sprintf(">>>>>>>>>>>>>> Tcp下发 ip:%s, port:%s, 开始 >>>>>>>>>>>>", host, port))
	conn, err := getTcpConn(host, port)
	//beego.Debug("tcp 拨号结束！")
	if err != nil {
		beego.Warn(err.Error())
		time.Sleep(time.Duration(100) * time.Millisecond)
		conn, err = net.DialTimeout("tcp", host+":"+port, time.Duration(5)*time.Second)
		i := 0
		for err != nil && i <= 2 {
			beego.Error(err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			conn, err = net.DialTimeout("tcp", host+":"+port, time.Duration(5)*time.Second)
			i++
		}
	}

	if conn == nil {
		if err != nil {
			return []byte{}, err
		}
		return []byte{}, errors.New("empty conn")
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(3) * time.Second))
	if err != nil {
		beego.Error(err)
	}

	var results = make([]byte, 1024)
	_, e := conn.Write(content)
	if e != nil {
		return nil, e
	}
	i, e2 := conn.Read(results)
	if e2 != nil {
		return nil, e2
	}
	//beego.Info(fmt.Sprintf(">>>>>>>>>>>>>> Tcp下发 ip:%s, port:%s, 结束 >>>>>>>>>>>>", host, port))
	return results[:i], nil
}

func TcpSendOnly(content []byte, host, port string, timeOutSecondsClose int) error {
	beego.Info(fmt.Sprintf(">>>>>>>>>>>>>> Tcp下发 ip:%s, port:%s, 开始 >>>>>>>>>>>>", host, port))
	addr := GetTcpAddr(host, port)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, e := conn.Write(content)
	if e != nil {
		return e
	}
	if timeOutSecondsClose > 0 {
		time.Sleep(time.Duration(timeOutSecondsClose) * time.Millisecond)
	}
	beego.Info(fmt.Sprintf(">>>>>>>>>>>>>> Tcp下发 ip:%s, port:%s, 结束 >>>>>>>>>>>>", host, port))
	return nil
}

func GetTcpAddr(host, port string) *net.TCPAddr {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+port)
	if err != nil {
		beego.Error(fmt.Sprintf("Fatal error: %s", err.Error()))
		// 重连
		time.Sleep(time.Duration(30) * time.Second)
		return GetTcpAddr(host, port)
	} else {
		return tcpAddr
	}
}

func ReadFromTcp(host, port string, cb func([]byte)) {
	addr := GetTcpAddr(host, port)
	conn, e := net.DialTCP("tcp", nil, addr)
	if e != nil {
		beego.Error(e)
		return
	}
	for {
		var results = make([]byte, 1024)
		i, e2 := conn.Read(results)
		if e2 != nil {
			beego.Error(e2)
			break
		}
		cb(results[:i])
	}
}
