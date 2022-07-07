package controller

import (
	"glc/gweb"
	"glc/ldb/sysmnt"
	"log"
)

// 查询日志仓列表
func StorageListController(req *gweb.HttpRequest) *gweb.HttpResult {
	rs := sysmnt.GetStorageList()
	return gweb.Result(rs)
}

// 删除指定日志仓
func StorageDeleteController(req *gweb.HttpRequest) *gweb.HttpResult {
	name := req.GetFormParameter("storeName")
	err := sysmnt.DeleteStorage(name)
	if err != nil {
		msg := err.Error()
		log.Println("日志仓", name, "删除失败", msg)
		return gweb.Error500("日志仓 " + name + " 正在使用中，无法删除")
	}
	return gweb.Ok()
}
