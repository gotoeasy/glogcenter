package controller

import (
	"encoding/json"
	"fmt"
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"glc/ldb/status"
	"glc/ldb/sysmnt"
	"glc/ver"

	"github.com/gotoeasy/glang/cmn"
)

// 查询是否测试模式
func TestModeController(req *gweb.HttpRequest) *gweb.HttpResult {
	return gweb.Result(conf.IsTestMode())
}

// 查询版本信息
func VersionController(req *gweb.HttpRequest) *gweb.HttpResult {
	rs := cmn.OfMap("version", ver.VERSION, "latest", getLatestVersion()) // version当前版本号，latest最新版本号
	return gweb.Result(rs)
}

// 查询日志仓名称列表
func StorageNamesController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableLogin() && req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	rs := com.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	return gweb.Result(rs)
}

// 查询日志仓信息列表
func StorageListController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableLogin() && req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	rs := sysmnt.GetStorageList()
	return gweb.Result(rs)
}

// 删除指定日志仓
func StorageDeleteController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableLogin() && req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	name := req.GetFormParameter("storeName")
	if name == ".sysmnt" {
		return gweb.Error500("不能删除 .sysmnt")
	} else if conf.IsStoreNameAutoAddDate() {
		if conf.GetSaveDays() > 0 {
			ymd := cmn.Right(name, 8)
			if cmn.Len(ymd) == 8 && cmn.Startwiths(ymd, "20") {
				msg := fmt.Sprintf("当前是日志仓自动维护模式，最多保存 %d 天，不支持手动删除", conf.GetSaveDays())
				return gweb.Error500(msg)
			}
		}
	} else if name == "logdata" {
		return gweb.Error500("日志仓 " + name + " 正在使用，不能删除")
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

// 尝试查询最新版本号（注：这里不能保证服务一定可用），查不到返回空串
func getLatestVersion() string {
	// {"version":"v0.12.0"}
	bts, err := cmn.HttpGetJson("https://glc.gotoeasy.top/glogcenter/current/version.json?v="+ver.VERSION, "Auth:glc") // 取最新版本号
	if err != nil {
		return ""
	}

	var data struct {
		Version string `json:"version,omitempty"`
	}
	if err := json.Unmarshal(bts, &data); err != nil {
		return ""
	}
	return data.Version
}
