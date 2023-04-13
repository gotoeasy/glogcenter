package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb"

	"github.com/gotoeasy/glang/cmn"
)

// 日志检索（表单提交方式）
func LogSearchController(req *gweb.HttpRequest) *gweb.HttpResult {
	for _, s := range GetSessionid() {
		if conf.IsEnableLogin() && req.GetFormParameter("token") == s["sessionid"] {
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
	}

	return gweb.Error403() // 登录检查
}
