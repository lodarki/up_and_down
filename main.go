package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"up_and_down/handler/handleUpAndDown"
	_ "up_and_down/routers"
)

func main() {
	// 定时回收订单群聊
	autoUpAndDown := toolbox.NewTask("auto_up_and_down", "0 0 * * * *", func() error {
		beego.Info("do task")
		e := handleUpAndDown.UpOrDown()
		if e != nil {
			beego.Error(e)
		}
		return nil
	})

	autoStatus := toolbox.NewTask("auto_status", "0/5 * * * * *", func() error {
		e := handleUpAndDown.GetStatus()
		if e != nil {
			beego.Error(e)
		}
		return nil
	})

	// 定时任务组建初始化
	toolbox.AddTask("auto_up_and_down", autoUpAndDown)
	toolbox.AddTask("auto_status", autoStatus)
	toolbox.StartTask()
	beego.Run()
}
