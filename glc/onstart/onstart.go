package onstart

import (
	"glc/gweb"
	"glc/gweb/http"
	"glc/gweb/method"
	"glc/ldb"
	"glc/ldb/conf"
	"glc/www/controller"
	"glc/www/filter"
	"glc/www/html"
)

func Run() {

	http.StartHttpServer(func() {

		contextPath := conf.GetContextPath() // ContextPath

		// 过滤器
		gweb.RegisterFilter(filter.ApiKeyFilter)
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
		gweb.RegisterController(method.POST, contextPath+"/search", controller.LogSearchController)
		gweb.RegisterController(method.POST, contextPath+"/add", controller.LogAddController)

		// 默认引擎空转一下，触发未建索引继续建
		go ldb.NewDefaultEngine().AddTextLog("", "", "")
	})

}
