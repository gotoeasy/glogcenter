package controller

import (
	"crypto/md5"
	"encoding/hex"
	"glc/cmn"
	"glc/conf"
	"glc/gweb"
	"time"
)

var sessionid string

func init() {
	sessionid = createSessionid()
	go func() {
		ticker := time.NewTicker(time.Hour) // 一小时更新一次
		for {
			<-ticker.C
			sessionid = createSessionid()
		}
	}()
}

func LoginController(req *gweb.HttpRequest) *gweb.HttpResult {
	username := req.GetFormParameter("username")
	password := req.GetFormParameter("password")
	if username != conf.GetUsername() || password != conf.GetPassword() {
		return gweb.Error500("用户名或密码错误")
	}

	return gweb.Result(sessionid)
}

func createSessionid() string {
	ymd := cmn.GetYyyymmdd(0)
	by1 := md5.Sum(cmn.StringToBytes(conf.GetUsername() + ymd))
	by2 := md5.Sum(cmn.StringToBytes(ymd + conf.GetPassword()))
	by3 := md5.Sum(cmn.StringToBytes(ymd + "添油" + conf.GetUsername() + "加醋" + conf.GetPassword()))
	str1 := hex.EncodeToString(by1[:])
	str2 := hex.EncodeToString(by2[:])
	str3 := hex.EncodeToString(by3[:])
	return cmn.RightRune(str1, 15) + cmn.LeftRune(str2, 15) + cmn.LeftRune(str3, 30)
}

func GetSessionid() string {
	return sessionid
}
