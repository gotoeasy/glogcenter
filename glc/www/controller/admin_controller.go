package controller

import (
	"glc/conf"
	"glc/gweb"
)

// 前端检索页面
func AdminController(req *gweb.HttpRequest) *gweb.HttpResult {
	defer req.Redirect(conf.GetContextPath() + "/search")
	return nil
}
