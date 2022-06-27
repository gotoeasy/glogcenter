package controller

import (
	"glc/cmn"
	"glc/gweb"
	"glc/ldb"
	"glc/ldb/conf"
	"glc/ldb/tokenizer"
)

// 日志检索（表单提交方式）
func LogSearchController(req *gweb.HttpRequest) *gweb.HttpResult {
	storeNmae := req.GetFormParameter("name")
	searchKey := tokenizer.GetSearchKey(req.GetFormParameter("searchKey"))
	pageSize := cmn.StringToInt(req.GetFormParameter("pageSize"), 20)
	currentId := cmn.StringToUint64(req.GetFormParameter("currentId"), 36, 0)
	forward := cmn.StringToBool(req.GetFormParameter("forward"), true)

	eng := ldb.NewEngine(storeNmae)
	rs := eng.Search(searchKey, pageSize, currentId, forward)
	return gweb.Result(rs)
}

// 重定向
func RedirectToSearchController(req *gweb.HttpRequest) *gweb.HttpResult {
	defer req.Redirect(conf.GetContextPath() + "/search")
	return nil
}
