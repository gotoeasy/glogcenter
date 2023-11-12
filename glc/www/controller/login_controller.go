package controller

import (
	"crypto/md5"
	"encoding/hex"
	"glc/conf"
	"glc/gweb"
	"glc/ldb/sysmnt"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

var catchLoginCheck *cmn.Cache // 缓存：登录失败次数检查
var catchSession *cmn.Cache    // 缓存：登录会话

func init() {
	if conf.IsEnableLogin() {
		catchLoginCheck = cmn.NewCache(time.Minute * 15)
		catchSession = cmn.NewCache(time.Hour * 12)
	}
}

func LoginController(req *gweb.HttpRequest) *gweb.HttpResult {

	if !InWhiteList(req) && InBlackList(req) {
		return gweb.Error403() // 黑名单，访问受限
	}
	if !conf.IsEnableLogin() {
		return gweb.Ok() // 登录相关变量没有初始化，不适合继续
	}

	username := req.GetFormParameter("username")
	password := req.GetFormParameter("password")
	key := getClientHash(req)
	val, find := catchLoginCheck.Get(key)
	cnt := 0
	if find {
		cnt = val.(int)
		if cnt >= 5 {
			catchLoginCheck.Set(key, cnt) // 还试，重新计算限制时间，再等15分钟吧
			return gweb.Error500("连续多次失败，当前已被限制登录")
		}
	}

	role := ""
	if username == conf.GetUsername() {
		// 管理员登录
		if password != conf.GetPassword() {
			cnt++
			catchLoginCheck.Set(key, cnt)
			return gweb.Error500("用户名或密码错误")
		}
		role = "admin"
	} else {
		// 一般用户登录
		user := sysmnt.NewSysmntStorage().GetSysUser(username)
		if user == nil || password != user.Password {
			cnt++
			catchLoginCheck.Set(key, cnt)
			return gweb.Error500("用户名或密码错误")
		}
	}

	token := createSessionid(username)
	catchSession.Set(token, username)

	catchLoginCheck.Delete(key)

	if conf.IsClusterMode() {
		user := &sysmnt.SysUser{Username: username, Password: password}
		go TransferGlc(conf.UserTransferLogin, user.ToJson()) // 转发其他GLC服务
	}

	return gweb.Result(cmn.OfMap("token", token, "role", role))
}

// 登录（来自数据转发）
func UserTransferLoginController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	loginuser := &sysmnt.SysUser{}
	req.BindJSON(loginuser)
	if loginuser.Username == conf.GetUsername() {
		// 管理员登录
		if loginuser.Password != conf.GetPassword() {
			return gweb.Error500("用户名或密码错误")
		}
	} else {
		// 一般用户登录
		user := sysmnt.NewSysmntStorage().GetSysUser(loginuser.Username)
		if user == nil || loginuser.Password != user.Password {
			return gweb.Error500("用户名或密码错误")
		}
	}

	token := createSessionid(loginuser.Username)
	catchSession.Set(token, loginuser.Username)
	return gweb.Ok()
}

func IsEnableLoginController(req *gweb.HttpRequest) *gweb.HttpResult {
	return gweb.Result(conf.IsEnableLogin())
}

func createSessionid(username string) string {
	ymd := cmn.Today()
	by1 := md5.Sum(cmn.StringToBytes(username + ymd))
	by2 := md5.Sum(cmn.StringToBytes(ymd + username + "添油"))
	by3 := md5.Sum(cmn.StringToBytes(ymd + username + "加醋" + conf.GetTokenSalt())) // 增加配置的令牌盐
	str1 := hex.EncodeToString(by1[:])
	str2 := hex.EncodeToString(by2[:])
	str3 := hex.EncodeToString(by3[:])
	return cmn.Right(str1, 15) + cmn.Left(str2, 15) + cmn.Left(str3, 30)
}

func GetUsernameByToken(token string) string {
	username, find := catchSession.Get(token)
	if find {
		return username.(string)
	}
	return ""
}

func getClientHash(req *gweb.HttpRequest) string {
	var ary []string
	ary = append(ary, req.GetHeader("Sec-Fetch-Site"))
	ary = append(ary, req.GetHeader("Sec-Fetch-Dest"))
	ary = append(ary, req.GetHeader("Sec-Ch-Ua-Mobile"))
	ary = append(ary, req.GetHeader("Accept-Language"))
	ary = append(ary, req.GetHeader("Accept-Encoding"))
	ary = append(ary, req.GetHeader("X-Forwarded-For"))
	ary = append(ary, req.GetHeader("Forwarded"))
	ary = append(ary, req.GetHeader("Sec-Ch-Ua-Platform"))
	ary = append(ary, req.GetHeader("User-Agent"))
	ary = append(ary, req.GetHeader("Sec-Fetch-Mode"))
	ary = append(ary, req.GetHeader("Sec-Ch-Ua"))
	ary = append(ary, req.GetHeader("Referer"))
	ary = append(ary, req.GinCtx.ClientIP())
	ary = append(ary, req.GetFormParameter("username"))
	return cmn.HashString(cmn.Join(ary, ","))
}

// 客户端IP是否在白名单中（内网地址总是在白名单中）
func InWhiteList(req *gweb.HttpRequest) bool {
	cityIp := cmn.GetCityIp(req.GinCtx.ClientIP())
	if cmn.Contains(cityIp, "内网") {
		return true
	}
	for i := 0; i < len(conf.GetWhiteList()); i++ {
		item := conf.GetWhiteList()[i]
		if item == "" {
			continue
		}
		if cmn.Endwiths(item, ".*") {
			item = cmn.ReplaceAll(item, "*", "") // 支持IP的最后一段使用通配符*
		}
		if cmn.Contains(cityIp, item) {
			return true
		}
	}
	return false
}

// 客户端IP是否在黑名单中（内网地址总是在白名单中）
func InBlackList(req *gweb.HttpRequest) bool {
	cityIp := cmn.GetCityIp(req.GinCtx.ClientIP())
	for i := 0; i < len(conf.GetBlackList()); i++ {
		item := conf.GetBlackList()[i]
		if item == "" {
			continue
		}
		if item == "*" {
			return true
		}
		if cmn.Endwiths(item, ".*") {
			item = cmn.ReplaceAll(item, "*", "") // 支持IP的最后一段使用通配符*
		}
		if cmn.Contains(cityIp, item) {
			return true
		}
	}
	return false
}
