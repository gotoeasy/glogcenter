package router

import (
	"glc/gweb"
	"glc/gweb/method"
	"glc/ldb/conf"
	"glc/web/controller"
	"glc/web/filter"
)

func Register() {

	// ContextPath
	contextPath := conf.GetContextPath()

	// 过滤器
	gweb.RegisterFilter(filter.ApiKeyFilter)

	// 控制器器（跳转）
	gweb.RegisterController(method.GET, "/", controller.RedirectToSearchController)
	gweb.RegisterController(method.GET, contextPath, controller.RedirectToSearchController)
	gweb.RegisterController(method.GET, contextPath+"/", controller.RedirectToSearchController)

	// 控制器
	gweb.RegisterController(method.POST, contextPath+"/search", controller.LogSearchController)
	gweb.RegisterController(method.POST, contextPath+"/add", controller.LogAddController)

}
