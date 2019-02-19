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
