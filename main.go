package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"up_and_down/handler/handleUpAndDown"
	_ "up_and_down/routers"
)

func main() {
	// 定时回收订单群聊
	autoUpAndDown := toolbox.NewTask("auto_up_and_down", "0 0/30 * * * *", func() error {
		beego.Info("do task")
		e := handleUpAndDown.UpOrDown()
		if e != nil {
			beego.Error(e)
		}
		return nil
	})

	// 定时任务组建初始化
	toolbox.AddTask("auto_up_and_down", autoUpAndDown)
	toolbox.StartTask()
	beego.Run()
}
