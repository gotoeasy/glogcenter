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
	loglevel := req.GetFormParameter("loglevel")

	if !cmn.IsBlank(system) {
		system = "~" + cmn.Trim(system)
	}
	if !cmn.IsBlank(loglevel) {
		loglevel = "!" + cmn.Trim(loglevel)
	}

	eng := ldb.NewEngine(storeName)
	rs := eng.Search(searchKey, system, datetimeFrom, datetimeTo, loglevel, currentId, forward)

	// 检索结果后处理
	rs.PageSize = cmn.IntToString(conf.GetPageSize())
	if !forward {
		// 修复最多匹配件数：检索（非检索更多）数据量少于1页时，最多匹配件数=检索结果件数，避免个别特殊场景下两者不一致
		size := len(rs.Data)
		if size < conf.GetPageSize() {
			rs.Count = cmn.IntToString(size)
		}
	}

	return gweb.Result(rs)
}
