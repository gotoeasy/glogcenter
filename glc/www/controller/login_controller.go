package controller

import (
	"crypto/md5"
	"encoding/hex"
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

var sessionid []map[string]string

func init() {
	if conf.IsEnableLogin() {
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
	userList := conf.GetUserList()
	for _, user := range userList {
		if username == user.Username && password == user.Password {
			for _, s := range sessionid {
				if s["username"] == username {
					return gweb.Result(s["sessionid"])
				}
			}
		}
	}

	return gweb.Error500("用户名或密码错误")
}

func IsEnableLoginController(req *gweb.HttpRequest) *gweb.HttpResult {
	return gweb.Result(conf.IsEnableLogin())
}

func createSessionid() []map[string]string {
	userList := conf.GetUserList()
	for _, user := range userList {
		ymd := com.GetYyyymmdd(0)
		by1 := md5.Sum(cmn.StringToBytes(user.Username + ymd))
		by2 := md5.Sum(cmn.StringToBytes(ymd + conf.GetPassword()))
		by3 := md5.Sum(cmn.StringToBytes(ymd + "添油" + conf.GetUsername() + "加醋" + conf.GetPassword()))
		str1 := hex.EncodeToString(by1[:])
		str2 := hex.EncodeToString(by2[:])
		str3 := hex.EncodeToString(by3[:])
		sessionid = append(sessionid, map[string]string{"username": user.Username, "sessionid": cmn.Right(str1, 15) + cmn.Left(str2, 15) + cmn.Left(str3, 30)})
	}
	return sessionid
}

func GetSessionid() []map[string]string {
	return sessionid
}
