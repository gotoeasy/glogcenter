package controller

import (
	"glc/gweb"
	"glc/ldb/conf"
)

// 前端检索页面
func AdminController(req *gweb.HttpRequest) *gweb.HttpResult {
	defer req.Redirect(conf.GetContextPath() + "/search")
	return nil
}
