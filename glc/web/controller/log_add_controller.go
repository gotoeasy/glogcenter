package controller

import (
	"glc/gweb"
	"glc/ldb"
)

// 添加日志（表单提交方式）
func LogAddController(req *gweb.HttpRequest) *gweb.HttpResult {
	storeNmae := req.GetFormParameter("name")
	text := req.GetFormParameter("text")

	engine := ldb.NewEngine(storeNmae)
	engine.AddTextLog(text)
	return gweb.Ok()
}
