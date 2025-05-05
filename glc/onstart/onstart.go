package onstart

import (
	"glc/conf"
	"glc/gweb"
	"glc/gweb/http"
	"glc/gweb/method"
	"glc/ldb"
	"glc/rabbitmq"
	"glc/www/cluster"
	"glc/www/controller"
	"glc/www/html"
)

func Run() {

	http.StartHttpServer(func() {

		contextPath := conf.GetContextPath() // ContextPath

		// Html静态文件
		gweb.RegisterController(method.GET, contextPath+"/", html.HomeIndexHtmlController) // [响应/glc/]
		gweb.RegisterController(method.GET, "/**/*.html", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.css", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.js", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.txt", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.ico", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.png", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.jpg", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.jpeg", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.gif", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.svg", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.json", html.StaticFileController)
		gweb.RegisterController(method.GET, "/**/*.xml", html.StaticFileController)

		// 控制器
		gweb.RegisterController(method.POST, contextPath+"/v1/log/add", controller.JsonLogAddController)                         // 添加日志
		gweb.RegisterController(method.POST, contextPath+"/v1/log/addBatch", controller.JsonLogAddBatchController)               // 添加日志数组
		gweb.RegisterController(method.POST, contextPath, controller.JsonLogAddBatchController)                                  // 添加日志数组（响应简化路径如 http://host:port/glc）
		gweb.RegisterController(method.POST, "/", controller.JsonLogAddBatchController)                                          // 添加日志数组（响应简化路径如 http://host:port/）
		gweb.RegisterController(method.POST, contextPath+"/v1/log/addTestData", controller.JsonLogAddTestDataController)         // 添加测试日志（仅测试模式有效）
		gweb.RegisterController(method.POST, contextPath+conf.LogTransferAdd, controller.JsonLogTransferAddController)           // 日志数据转发添加日志
		gweb.RegisterController(method.POST, contextPath+"/v1/log/search", controller.LogSearchController)                       // 查询日志
		gweb.RegisterController(method.POST, contextPath+"/v1/store/names", controller.StorageNamesController)                   // 查询日志仓名称列表
		gweb.RegisterController(method.POST, contextPath+"/v1/store/list", controller.StorageListController)                     // 查询日志仓信息列表
		gweb.RegisterController(method.POST, contextPath+"/v1/store/delete", controller.StorageDeleteController)                 // 删除日志仓
		gweb.RegisterController(method.POST, contextPath+"/v1/store/systems", controller.SystemNamesController)                  // 查询系统名列表
		gweb.RegisterController(method.POST, contextPath+"/v1/store/mode", controller.TestModeController)                        // 查询是否测试模式
		gweb.RegisterController(method.POST, contextPath+"/v1/user/enableLogin", controller.IsEnableLoginController)             // 查询是否开启用户密码登录功能
		gweb.RegisterController(method.POST, contextPath+"/v1/user/login", controller.LoginController)                           // Login
		gweb.RegisterController(method.POST, contextPath+conf.UserTransferLogin, controller.UserTransferLoginController)         // 转发Login
		gweb.RegisterController(method.POST, contextPath+"/v1/version/info", controller.VersionController)                       // 查询版本信息
		gweb.RegisterController(method.POST, contextPath+"/v1/sysuser/list", controller.UserListController)                      // [用户]列表查询
		gweb.RegisterController(method.POST, contextPath+"/v1/sysuser/save", controller.UserSaveController)                      // [用户]保存
		gweb.RegisterController(method.POST, contextPath+"/v1/sysuser/del", controller.UserDelController)                        // [用户]删除
		gweb.RegisterController(method.POST, contextPath+"/v1/sysuser/changePsw", controller.UserChangePswController)            // [用户]修改自己密码
		gweb.RegisterController(method.POST, contextPath+conf.SysUserTransferChgPsw, controller.UserTransferChangePswController) // [用户]转发修改自己密码
		gweb.RegisterController(method.POST, contextPath+conf.SysUserTransferSave, controller.UserTransferSaveController)        // [用户]转发保存
		gweb.RegisterController(method.POST, contextPath+conf.SysUserTransferDel, controller.UserTransferDelController)          // [用户]转发删除

		// 集群操作接口
		gweb.RegisterController(method.POST, contextPath+"/sys/cluster/info", controller.ClusterGetClusterInfoController)   // 获取集群信息
		gweb.RegisterController(method.POST, contextPath+"/sys/cluster/save", controller.ClusterMasterSaveKvDataController) // 保存集群信息
		gweb.RegisterController(method.POST, contextPath+"/sys/cluster/async", controller.ClusterMasterAsyncDataController) // 保存Master发来的集群信息
		gweb.RegisterController(method.GET, contextPath+"/sys/cluster/down", controller.ClusterDownloadStoreDataController) // 打包下载指定日志仓数据
		gweb.RegisterController(method.POST, contextPath+"/sys/item/get", controller.ClusterGetItemController)              // KV获取
		gweb.RegisterController(method.POST, contextPath+"/sys/item/set", controller.ClusterSetItemController)              // KV设定
		gweb.RegisterController(method.POST, contextPath+"/sys/item/del", controller.ClusterDelItemController)              // KV删除

		// 默认引擎空转一下，触发未建索引继续建
		go ldb.NewDefaultEngine().AddTextLog("", "", "")

		cluster.Start() // 显式调用触发加入集群或初始化集群
		rabbitmq.Start()
	})

}
