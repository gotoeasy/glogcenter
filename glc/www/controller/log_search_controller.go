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
	searchKey := req.GetFormParameter("searchKey")
	currentId := cmn.StringToUint32(req.GetFormParameter("currentId"), 0)
	forward := cmn.StringToBool(req.GetFormParameter("forward"), true)
	datetimeFrom := req.GetFormParameter("datetimeFrom")
	datetimeTo := req.GetFormParameter("datetimeTo")
	system := req.GetFormParameter("system")
	loglevel := req.GetFormParameter("loglevel") // 单选条件
	loglevels := cmn.Split(loglevel, ",")        // 多选条件
	if len(loglevels) <= 1 || len(loglevels) >= 4 {
		loglevels = make([]string, 0) // 多选的单选或全选，都清空（单选走loglevel索引，全选等于没选）
	}

	if !cmn.IsBlank(system) {
		system = "~" + cmn.Trim(system)
	}
	if !cmn.IsBlank(loglevel) && len(loglevels) == 0 {
		loglevel = "!" + cmn.Trim(loglevel) // 单个条件时作为索引条件
	} else {
		loglevel = "" // 多选条件时不使用，改用loglevels
	}

	eng := ldb.NewEngine(storeName)
	rs := eng.Search(searchKey, system, datetimeFrom, datetimeTo, loglevel, loglevels, currentId, forward)

	// 检索结果后处理
	rs.PageSize = cmn.IntToString(conf.GetPageSize())
	return gweb.Result(rs)
}
