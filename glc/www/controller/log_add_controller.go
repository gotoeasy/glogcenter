package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb"
	"glc/ldb/storage/logdata"
	"log"
)

// 添加日志（表单提交方式）
func LogAddController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() {
		auth := req.GetHeader(conf.GetHeaderSecurityKey())
		if auth != conf.GetSecurityKey() {
			return gweb.Error(403, "未经授权的访问，拒绝服务")
		}
	}

	storeNmae := req.GetFormParameter("name")
	text := req.GetFormParameter("text")
	date := req.GetFormParameter("date")
	system := req.GetFormParameter("system")

	engine := ldb.NewEngine(storeNmae)
	engine.AddTextLog(date, text, system)
	return gweb.Ok()
}

// 添加日志（JSON提交方式）
func JsonLogAddController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() {
		auth := req.GetHeader(conf.GetHeaderSecurityKey())
		if auth != conf.GetSecurityKey() {
			return gweb.Error(403, "未经授权的访问，拒绝服务")
		}
	}

	md := &logdata.LogDataModel{}
	err := req.BindJSON(md)
	if err != nil {
		log.Println("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	engine := ldb.NewDefaultEngine()
	engine.AddTextLog(md.Date, md.Text, md.System)
	return gweb.Ok()
}
