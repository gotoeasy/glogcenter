package controller

import (
	"glc/cmn"
	"glc/gweb"
	"glc/ldb"
)

// 日志检索（表单提交方式）
func LogSearchController(req *gweb.HttpRequest) *gweb.HttpResult {

	if req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	storeName := req.GetFormParameter("storeName")
	//searchKey := tokenizer.GetSearchKey(req.GetFormParameter("searchKey"))
	searchKey := req.GetFormParameter("searchKey")
	pageSize := cmn.StringToInt(req.GetFormParameter("pageSize"), 20)
	currentId := cmn.StringToUint32(req.GetFormParameter("currentId"), 0)
	forward := cmn.StringToBool(req.GetFormParameter("forward"), true)

	eng := ldb.NewEngine(storeName)
	rs := eng.Search(searchKey, pageSize, currentId, forward)
	return gweb.Result(rs)
}
