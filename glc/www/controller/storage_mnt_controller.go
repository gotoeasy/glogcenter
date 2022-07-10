package controller

import (
	"fmt"
	"glc/cmn"
	"glc/conf"
	"glc/gweb"
	"glc/ldb/status"
	"glc/ldb/sysmnt"
	"log"
)

// 查询日志仓名称列表
func StorageNamesController(req *gweb.HttpRequest) *gweb.HttpResult {
	if req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	rs := cmn.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	return gweb.Result(rs)
}

// 查询日志仓信息列表
func StorageListController(req *gweb.HttpRequest) *gweb.HttpResult {
	if req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	rs := sysmnt.GetStorageList()
	return gweb.Result(rs)
}

// 删除指定日志仓
func StorageDeleteController(req *gweb.HttpRequest) *gweb.HttpResult {
	if req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	name := req.GetFormParameter("storeName")

	if conf.IsStoreNameAutoAddDate() && conf.GetSaveDays() > 0 {
		msg := fmt.Sprintf("当前是日志仓自动维护模式，最多保存 %d 天，不能手动删除", conf.GetSaveDays())
		return gweb.Error500(msg)
	}

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
