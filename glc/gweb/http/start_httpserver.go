package http

import "glc/gweb"

// StartHttpServer 注册控制器，启动web服务
func StartHttpServer(fnRegister func()) {

	// 注册控制器
	fnRegister()

	// 启动web服务
	gweb.Run()
}
