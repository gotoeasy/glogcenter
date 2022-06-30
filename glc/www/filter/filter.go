package filter

import (
	"glc/gweb"
	"glc/ldb/conf"
	"net/http"
)

// 校验HEADER的API秘钥
func ApiKeyFilter(req *gweb.HttpRequest) *gweb.HttpResult {

	//log.Println("================================", req.RequestURI())
	// 开启API秘钥校验时才检查
	if !conf.IsEnableSecurityKey() {
		return nil
	}

	auth := req.GetHeader(conf.GetHeaderSecurityKey())
	if auth != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}
	return nil // 返回nil表示正常过滤成功
}

// 校验HEADER的API秘钥
func CrossFilter(req *gweb.HttpRequest) *gweb.HttpResult {

	req.SetHeader("Access-Control-Allow-Origin", "*")
	req.SetHeader("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	req.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	req.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	req.SetHeader("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if req.GetMethod() == "OPTIONS" {
		req.AbortWithStatus(http.StatusNoContent)
	}
	return nil
}
