package filter

import (
	"glc/conf"
	"glc/gweb"
	"net/http"
)

// 允许跨域
func CrossFilter(req *gweb.HttpRequest) *gweb.HttpResult {

	if conf.IsEnableCross() {
		req.SetHeader("Access-Control-Allow-Origin", "*")
		req.SetHeader("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		req.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		req.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		req.SetHeader("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if req.GetMethod() == "OPTIONS" {
			req.AbortWithStatus(http.StatusOK)
		}
	}

	return nil
}
