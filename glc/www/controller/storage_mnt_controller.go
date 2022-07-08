package controller

import (
	"glc/cmn"
	"glc/conf"
	"glc/gweb"
	"glc/ldb/status"
	"glc/ldb/sysmnt"
	"log"
)

// 查询日志仓名称列表
func StorageNamesController(req *gweb.HttpRequest) *gweb.HttpResult {
	rs := cmn.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	return gweb.Result(rs)
}

// 查询日志仓信息列表
func StorageListController(req *gweb.HttpRequest) *gweb.HttpResult {
	rs := sysmnt.GetStorageList()
	return gweb.Result(rs)
}

// 删除指定日志仓
func StorageDeleteController(req *gweb.HttpRequest) *gweb.HttpResult {
	name := req.GetFormParameter("storeName")
	if status.IsStorageOpening(name) {
		return gweb.Error500("日志仓 " + name + " 正在使用，不能删除")
	}

	err := sysmnt.DeleteStorage(name)
	if err != nil {
		msg := err.Error()
		log.Println("日志仓", name, "删除失败", msg)
		return gweb.Error500("日志仓 " + name + " 正在使用，不能删除")
	}
	return gweb.Ok()
}
