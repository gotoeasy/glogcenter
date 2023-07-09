package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb"

	"github.com/gotoeasy/glang/cmn"
)

// 日志检索（表单提交方式）
func LogSearchController(req *gweb.HttpRequest) *gweb.HttpResult {

	if conf.IsEnableLogin() && req.GetFormParameter("token") != GetSessionid() {
		return gweb.Error403() // 登录检查
	}

	storeName := req.GetFormParameter("storeName")
	//searchKey := tokenizer.GetSearchKey(req.GetFormParameter("searchKey"))
	searchKey := req.GetFormParameter("searchKey")
	pageSize := cmn.StringToInt(req.GetFormParameter("pageSize"), 20)
	currentId := cmn.StringToUint32(req.GetFormParameter("currentId"), 0)
	forward := cmn.StringToBool(req.GetFormParameter("forward"), true)
	datetimeFrom := req.GetFormParameter("datetimeFrom")
	datetimeTo := req.GetFormParameter("datetimeTo")
	system := req.GetFormParameter("system")
	loglevel := req.GetFormParameter("loglevel")

	if !cmn.IsBlank(system) {
		system = "~" + cmn.Trim(system)
	}
	if !cmn.IsBlank(loglevel) {
		loglevel = "!" + cmn.Trim(loglevel)
	}

	eng := ldb.NewEngine(storeName)
	rs := eng.Search(searchKey, system, datetimeFrom, datetimeTo, loglevel, pageSize, currentId, forward)
	return gweb.Result(rs)
}
