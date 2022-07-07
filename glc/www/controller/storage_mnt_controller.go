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
	name := req.GetFormParameter("name")
	err := sysmnt.DeleteStorage(name)
	if err != nil {
		log.Println("日志仓", name, "删除失败", err)
		return gweb.Error500("删除失败")
	}
	return gweb.Ok()
}
