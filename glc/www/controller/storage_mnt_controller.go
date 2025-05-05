package controller

import (
	"fmt"
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"glc/ldb/status"
	"glc/ldb/sysmnt"
	"glc/ver"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

var glcLatest string = ""
var glcOrigin string = ""

// 查询是否测试模式
func TestModeController(req *gweb.HttpRequest) *gweb.HttpResult {
	return gweb.Result(conf.IsTestMode())
}

// 查询版本信息
func VersionController(req *gweb.HttpRequest) *gweb.HttpResult {
	rs := cmn.OfMap("version", ver.VERSION, "latest", glcLatest) // version当前版本号，latest最新版本号
	return gweb.Result(rs)
}

func SetOrigin(req *gweb.HttpRequest) {
	origin := req.GinCtx.GetHeader("Origin")
	if origin != "" {
		glcOrigin = origin
	}
}

// 查询日志仓名称列表
func StorageNamesController(req *gweb.HttpRequest) *gweb.HttpResult {
	if (!InWhiteList(req) && InBlackList(req)) || (conf.IsEnableLogin() && GetUsernameByToken(req.GetToken()) == "") {
		return gweb.Error403() // 黑名单检查、登录检查
	}

	rs := com.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	return gweb.Result(rs)
}

// 查询系统名称列表
func SystemNamesController(req *gweb.HttpRequest) *gweb.HttpResult {
	if (!InWhiteList(req) && InBlackList(req)) || (conf.IsEnableLogin() && GetUsernameByToken(req.GetToken()) == "") {
		return gweb.Error403() // 黑名单检查、登录检查
	}

	if conf.IsEnableLogin() {
		username := GetUsernameByToken(req.GetToken())
		mnt := sysmnt.NewSysmntStorage()
		if username == conf.GetUsername() {
			// 管理员
			names := mnt.GetSysUsernames()
			var all []string
			var m map[string]bool = make(map[string]bool)
			for i := 0; i < len(names); i++ {
				user := mnt.GetSysUser(names[i])
				if user != nil && user.Systems != "*" {
					ary := cmn.Split(user.Systems, ",")
					for j := 0; j < len(ary); j++ {
						if !m[cmn.ToLower(ary[j])] {
							m[cmn.ToLower(ary[j])] = true
							all = append(all, ary[j])
						}
					}
				}
			}
			if len(all) == 0 {
				all = GetAllSystemNames() // 所有系统
			}
			return gweb.Result(all)
		} else {
			// 非管理员
			user := mnt.GetSysUser(username)
			if user != nil {
				if user.Systems == "*" {
					return gweb.Result(GetAllSystemNames()) // 所有系统
				}
				return gweb.Result(cmn.Split(user.Systems, ",")) // 按设定的系统
			}
		}
	} else {
		all := GetAllSystemNames()
		if len(all) > 0 {
			return gweb.Result(all)
		}
	}

	return gweb.Ok() // 都有权限，不返回结果
}

// 查询日志仓信息列表
func StorageListController(req *gweb.HttpRequest) *gweb.HttpResult {
	token := req.GetToken()
	if (!InWhiteList(req) && InBlackList(req)) || (conf.IsEnableLogin() && GetUsernameByToken(token) == "") {
		return gweb.Error403() // 黑名单检查、登录检查
	}
	if conf.IsEnableLogin() {
		catchSession.Set(token, GetUsernameByToken(token)) // 会话重新计时
	}

	rs := sysmnt.GetStorageList()
	return gweb.Result(rs)
}

// 删除指定日志仓
func StorageDeleteController(req *gweb.HttpRequest) *gweb.HttpResult {
	if (!InWhiteList(req) && InBlackList(req)) || (conf.IsEnableLogin() && GetUsernameByToken(req.GetToken()) == "") {
		return gweb.Error403() // 黑名单检查、登录检查
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

	cacheTime = time.Now().Add(-1 * time.Hour) // 让检索时不用缓存名，避免查询不存在的日志仓

	return gweb.Ok()
}

// 尝试查询最新版本号（注：服务不一定总是可用，每小时查取一次）
func init() {
	go func() {
		url := "https://glc.gotoeasy.top/glogcenter/current/version.json?v=" + ver.VERSION + "&h=" + cmn.Base62Encode(cmn.StringToBytes(glcOrigin))
		v := cmn.GetGlcLatestVersion(url)
		glcLatest = cmn.IifStr(v != "", v, glcLatest)
		ticker := time.NewTicker(time.Hour)
		for range ticker.C {
			url = "https://glc.gotoeasy.top/glogcenter/current/version.json?v=" + ver.VERSION + "&h=" + cmn.Base62Encode(cmn.StringToBytes(glcOrigin))
			v = cmn.GetGlcLatestVersion(url)
			glcLatest = cmn.IifStr(v != "", v, glcLatest)
		}
	}()
}
