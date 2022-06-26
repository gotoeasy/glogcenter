package router

import (
	"glc/gweb"
	"glc/gweb/method"
	"glc/web/controller"
)

func Register() {
	gweb.RegisterController(method.GET, "/glc/list", controller.LogSearchController)
}
