package controllers

import (
	"ai_lib/controllers"
	"ai_local/hardWare/rs485/rs485Helper"
	"github.com/astaxie/beego"
	"up_and_down/handler/handleUpAndDown"
)

type TestController struct {
	controllers.BaseController
}

func (c *TestController) URLMapping() {
	c.Mapping("AllUp", c.AllUp)
	c.Mapping("AllDown", c.AllDown)
}

// @Router /all-up [get]
func (c *TestController) AllUp() {
	upStr := "9a:00:01:00:0a:dd:d6"
	beego.Debug("全上")
	bytes, e := rs485Helper.GetCommandBytesFromStr(upStr)
	if e != nil {
		c.ServerError(e.Error())
		return
	}
	_, e = handleUpAndDown.Rs485Port.Write(bytes) //上
	if e != nil {
		c.ServerError(e.Error())
		return
	}

	c.Success(nil)
}

// @Router /all-down [get]
func (c *TestController) AllDown() {

	downStr := "9a:00:01:00:0a:ee:e5" //全下
	beego.Debug("全下")
	bytes, e := rs485Helper.GetCommandBytesFromStr(downStr)
	if e != nil {
		c.ServerError(e.Error())
		return
	}
	_, e = handleUpAndDown.Rs485Port.Write(bytes) //下
	if e != nil {
		c.ServerError(e.Error())
		return
	}

	c.Success(nil)
}
