package router

import (
	"glc/gweb"
	"glc/gweb/method"
	"glc/ldb/conf"
	"glc/web/controller"
)

func Register() {

	contextPath := conf.GetContextPath()

	gweb.RegisterController(method.GET, "/", controller.RedirectToSearchController)
	gweb.RegisterController(method.GET, contextPath, controller.RedirectToSearchController)
	gweb.RegisterController(method.GET, contextPath+"/", controller.RedirectToSearchController)

	gweb.RegisterController(method.POST, contextPath+"/search", controller.LogSearchController)
	gweb.RegisterController(method.POST, contextPath+"/add", controller.LogAddController)

}
