package controller

import (
	"glc/gweb"
	"glc/ldb"
)

// 添加日志（表单提交方式）
func LogAddController(req *gweb.HttpRequest) *gweb.HttpResult {
	storeNmae := req.GetFormParameter("name")
	text := req.GetFormParameter("text")
	date := req.GetFormParameter("date")
	system := req.GetFormParameter("system")

	engine := ldb.NewEngine(storeNmae)
	engine.AddTextLog(date, text, system)
	return gweb.Ok()
}
