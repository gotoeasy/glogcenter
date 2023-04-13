package controller

import (
	"fmt"
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"glc/ldb/status"
	"glc/ldb/sysmnt"

	"github.com/gotoeasy/glang/cmn"
)

// StorageNamesController 查询日志仓名称列表
func StorageNamesController(req *gweb.HttpRequest) *gweb.HttpResult {
	for _, s := range GetSessionid() {
		if conf.IsEnableLogin() && req.GetFormParameter("token") == s["sessionid"] {
			rs := com.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
			return gweb.Result(rs)

		}
	}
	return gweb.Error403() // 登录检查

}

// StorageListController 查询日志仓信息列表
func StorageListController(req *gweb.HttpRequest) *gweb.HttpResult {
	for _, s := range GetSessionid() {
		if conf.IsEnableLogin() && req.GetFormParameter("token") == s["sessionid"] {
			rs := sysmnt.GetStorageList()
			return gweb.Result(rs)

		}
	}
	return gweb.Error403() // 登录检查
}

// StorageDeleteController 删除指定日志仓
func StorageDeleteController(req *gweb.HttpRequest) *gweb.HttpResult {
	for _, s := range GetSessionid() {
		if conf.IsEnableLogin() && req.GetFormParameter("token") == s["sessionid"] {
			name := req.GetFormParameter("storeName")
			if name == ".sysmnt" {
				return gweb.Error500("不能删除 .sysmnt")
			} else if conf.IsStoreNameAutoAddDate() {
				if conf.GetSaveDays() > 0 {
					ymd := cmn.Right(name, 8)
					if cmn.Len(ymd) == 8 && cmn.Startwiths(ymd, "20") {
						msg := fmt.Sprintf("当前是日志仓自动维护模式，最多保存 %d 天，不能手动删除", conf.GetSaveDays())
						return gweb.Error500(msg)
					}
				}
			} else if name == "logdata" {
				return gweb.Error500("不能删除当前使用的唯一日志仓 " + "logdata")
			}

			if status.IsStorageOpening(name) {
				return gweb.Error500("日志仓 " + name + " 正在使用，不能删除")
			}

			err := sysmnt.DeleteStorage(name)
			if err != nil {
				cmn.Error("日志仓", name, "删除失败", err)
				return gweb.Error500("日志仓 " + name + " 正在使用，不能删除")
			}
			return gweb.Ok()

		}
	}
	return gweb.Error403() // 登录检查
}
