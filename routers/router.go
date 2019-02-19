package routers

import (
	"github.com/astaxie/beego"
	"up_and_down/controllers"
)

func init() {
	ns :=
		beego.NewNamespace("/test",
			beego.NSInclude(
				&controllers.TestController{},
			),
		)
	beego.AddNamespace(ns)
}
