package controller

import (
	"glc/cmn"
	"glc/gweb"
	"glc/ldb"
	"glc/ldb/tokenizer"
)

func LogSearchController(req *gweb.HttpRequest) *gweb.HttpResult {
	storeNmae := req.GetUrlParameter("name")
	if storeNmae == "" {
		storeNmae = "default"
	}
	searchKey := tokenizer.GetSearchKey(req.GetUrlParameter("searchKey"))
	pageSize := cmn.StringToInt(req.GetUrlParameter("pageSize"), 20)
	currentId := cmn.StringToUint64(req.GetUrlParameter("currentId"), 36, 0)
	forward := cmn.StringToBool(req.GetUrlParameter("forward"), true)

	eng := ldb.NewEngine(storeNmae)
	rs := eng.Search(searchKey, pageSize, currentId, forward)
	return gweb.Result(rs)
}
