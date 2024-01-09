package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb"
	"glc/ldb/storage/logdata"

	"github.com/gotoeasy/glang/cmn"
)

// 添加日志（JSON提交方式）
func JsonLogAddController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	md := &logdata.LogDataModel{}
	err := req.BindJSON(md)
	if err != nil {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	md.Text = cmn.Trim(md.Text)
	addDataModelLog(md)

	if conf.IsClusterMode() {
		go TransferGlc(conf.LogTransferAdd, md.ToJson()) // 转发其他GLC服务
	}

	return gweb.Ok()
}

// 添加日志（来自数据转发）
func JsonLogTransferAddController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	md := &logdata.LogDataModel{}
	err := req.BindJSON(md)
	if err != nil {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	md.Text = cmn.Trim(md.Text)
	addDataModelLog(md)
	// addTextLog(md)
	return gweb.Ok()
}

// 添加日志
func addDataModelLog(data *logdata.LogDataModel) {
	engine := ldb.NewDefaultEngine()

	// 按配置要求在IP字段上附加城市信息（当IP含空格时认为已附加过）
	if conf.IsIpAddCity() {
		if data.ClientIp != "" && !cmn.Contains(data.ClientIp, " ") {
			data.ClientIp = cmn.GetCityIp(data.ClientIp)
		}
		if data.ServerIp != "" && !cmn.Contains(data.ServerIp, " ") {
			data.ServerIp = cmn.GetCityIp(data.ServerIp)
		}
	}

	engine.AddLogDataModel(data)
}
