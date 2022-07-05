package onstart

import (
	"glc/conf"
	"glc/gweb"
	"glc/gweb/http"
	"glc/gweb/method"
	"glc/ldb"
	"glc/rabbitmq"
	"glc/www/controller"
	"glc/www/filter"
	"glc/www/html"
)

func Run() {

	http.StartHttpServer(func() {

		contextPath := conf.GetContextPath() // ContextPath

		// 过滤器
		gweb.RegisterFilter(filter.CrossFilter)

		// 控制器（跳转）
		gweb.RegisterController(method.GET, "/", html.RedirectToHomeController)
		gweb.RegisterController(method.GET, contextPath, html.RedirectToHomeController)

		// Html静态文件
		gweb.RegisterController(method.GET, contextPath+"/", html.HomeIndexHtmlController)
		gweb.RegisterController(method.GET, "/**/*.html", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.css", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.js", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.ico", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.png", html.StaticFileController)

		// 控制器
		gweb.RegisterController(method.POST, contextPath+"/search", controller.LogSearchController) // Deprecated
		gweb.RegisterController(method.POST, contextPath+"/add", controller.LogAddController)       // Deprecated

		gweb.RegisterController(method.POST, contextPath+"/v1/log/add", controller.JsonLogAddController)
		gweb.RegisterController(method.POST, contextPath+"/v1/log/search", controller.LogSearchController)

		// 默认引擎空转一下，触发未建索引继续建
		go ldb.NewDefaultEngine().AddTextLog("", "", "")

		rabbitmq.Start()
	})

}
