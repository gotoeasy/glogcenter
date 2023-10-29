package controller

import (
	"crypto/md5"
	"encoding/hex"
	"glc/conf"
	"glc/gweb"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

var sessionid string
var catch *cmn.Cache

func init() {
	if conf.IsEnableLogin() {
		catch = cmn.NewCache(time.Minute * 15)
		sessionid = createSessionid()
		go func() {
			ticker := time.NewTicker(time.Hour) // 一小时更新一次
			for {
				<-ticker.C
				sessionid = createSessionid()
			}
		}()
	}
}

func LoginController(req *gweb.HttpRequest) *gweb.HttpResult {
	username := req.GetFormParameter("username")
	password := req.GetFormParameter("password")
	key := getClientHash(req)
	val, find := catch.Get(key)
	cnt := 0
	if find {
		cnt = val.(int)
		if cnt >= 5 {
			catch.Set(key, cnt) // 还试，重新计算限制时间，再等15分钟吧
			return gweb.Error500("连续多次失败，当前已被限制登录")
		}
	}
	if username != conf.GetUsername() || password != conf.GetPassword() {
		cnt++
		catch.Set(key, cnt)
		return gweb.Error500("用户名或密码错误")
	}

	catch.Delete(key)
	return gweb.Result(sessionid)
}

func IsEnableLoginController(req *gweb.HttpRequest) *gweb.HttpResult {
	return gweb.Result(conf.IsEnableLogin())
}

func createSessionid() string {
	ymd := cmn.Today()
	by1 := md5.Sum(cmn.StringToBytes(conf.GetUsername() + ymd))
	by2 := md5.Sum(cmn.StringToBytes(ymd + conf.GetPassword()))
	by3 := md5.Sum(cmn.StringToBytes(ymd + "添油" + conf.GetUsername() + "加醋" + conf.GetPassword() + conf.GetTokenSalt())) // 增加配置的令牌盐
	str1 := hex.EncodeToString(by1[:])
	str2 := hex.EncodeToString(by2[:])
	str3 := hex.EncodeToString(by3[:])
	return cmn.Right(str1, 15) + cmn.Left(str2, 15) + cmn.Left(str3, 30)
}

func GetSessionid() string {
	return sessionid
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
	return cmn.HashString(cmn.Join(ary, ","))
}
