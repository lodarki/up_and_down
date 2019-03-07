package controllers

import (
	"ai_lib/controllers"
	"up_and_down/handler/handleUpAndDown"
)

type TestController struct {
	controllers.BaseController
}

func (c *TestController) URLMapping() {
	c.Mapping("AllUp", c.AllUp)
	c.Mapping("AllDown", c.AllDown)
	c.Mapping("GetStatus", c.GetStatus)
}

// @router /all-up [get]
func (c *TestController) AllUp() {
	e := handleUpAndDown.AllUp()
	if e != nil {
		c.ServerError(e.Error())
		return
	}
	c.Success(nil)
}

// @router /all-down [get]
func (c *TestController) AllDown() {
	e := handleUpAndDown.AllDown()
	if e != nil {
		c.ServerError(e.Error())
		return
	}
	c.Success(nil)
}

// @title GetStatus
// @description 获取电机状态
// @router /get-status [get]
func (c *TestController) GetStatus() {

	position, status, err := handleUpAndDown.GetStatus()
	if err != nil {
		c.ServerError(err.Error())
		return
	}

	resultMap := make(map[string]interface{})
	resultMap["position"] = position
	resultMap["status"] = status
	c.Success(resultMap)
}
