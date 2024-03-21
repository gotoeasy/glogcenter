package ldb

import (
	"glc/com"
	"glc/ldb/search"
	"glc/ldb/storage"
	"glc/ldb/storage/indexword"
	"glc/ldb/storage/logdata"
	"glc/ldb/tokenizer"

	"github.com/gotoeasy/glang/cmn"
)

type Engine struct {
	storeName   string
	logStorage  *storage.LogDataStorageHandle // 日志存储控制器
	idxwStorage *indexword.WordIndexStorage   // 关键词反向索引存储器
}

func NewEngine(storeName string) *Engine {
	if storeName == "" {
		storeName = com.GeyStoreNameByDate("")
	}

	return &Engine{
		storeName:   storeName,
		logStorage:  storage.NewLogDataStorageHandle(storeName),
		idxwStorage: indexword.NewWordIndexStorage(storeName),
	}
}

func NewDefaultEngine() *Engine {
	var storeName string = com.GeyStoreNameByDate("")
	return &Engine{
		storeName:   storeName,
		logStorage:  storage.NewLogDataStorageHandle(storeName),
		idxwStorage: indexword.NewWordIndexStorage(storeName),
	}
}

// 添加日志
func (e *Engine) AddTextLog(date string, logText string, system string) {
	e.logStorage.AddTextLog(date, logText, system)
}

// 添加日志
func (e *Engine) AddLogDataModel(data *logdata.LogDataModel) {
	e.logStorage.AddLogDataModel(data)
}

func (e *Engine) Search(cond *search.SearchCondition) *search.SearchResult {

	// 分词后检索
	var adds []string
	adds = append(adds, cond.OrgSystem, cond.Loglevel, cond.User)
	kws := tokenizer.CutForSearchEx(cond.SearchKey, adds, nil) // 检索用关键词处理

	// 【快速检查1】，存在无索引数据的关键词时，直接返回
	idxw := indexword.NewWordIndexStorage(e.storeName)
	totalCount := e.logStorage.TotalCount() // 当前日志仓总件数
	maxMatchCount := totalCount             // 最大匹配件数
	for _, word := range kws {
		cnt := idxw.GetTotalCount(word)
		if cnt < 1 {
			cmn.Debug("关键词", word, "没有索引数据，直接返回空结果")
			rs := new(search.SearchResult)
			rs.Total = cmn.Uint32ToString(totalCount)
			rs.Count = "0"
			return rs
		}
		if cnt < maxMatchCount {
			maxMatchCount = cnt // 取最小件数，尽量减少最大匹配件数的误差
		}
	}

	// 【快速检查2】，日志级别多选时，合计下件数，没有数据时直接返回
	if len(cond.Loglevels) > 0 {
		var allcnt uint32
		for _, word := range cond.Loglevels {
			n := idxw.GetTotalCount("!" + word)
			allcnt += n // 累加总数
		}
		if allcnt < 1 {
			cmn.Debug("日志级别范围 ", cmn.Join(cond.Loglevels, ","), " 在日志仓 ", e.storeName, " 中没有索引数据，直接返回空结果")
			rs := new(search.SearchResult)
			rs.Total = cmn.Uint32ToString(totalCount)
			rs.Count = "0"
			return rs
		}

		if allcnt < maxMatchCount {
			maxMatchCount = allcnt // 取其小，尽量减少最大匹配件数的误差
		}
	}

	// 【快速检查3】，权限内的系统没有数据时，直接返回（场景：虽然没有输入系统条件，但当前日志仓没有权限范围系统的数据）
	if cond.OrgSystem == "" && cond.OrgSystems[0] != "*" {
		var sysnames []string
		var allcnt uint32
		for _, word := range cond.OrgSystems {
			n := idxw.GetTotalCount(word)
			if n > 0 {
				sysnames = append(sysnames, word) // 这个系统有数据
			}
			allcnt += n // 累加总数
		}
		if allcnt < 1 {
			cmn.Debug("权限范围系统 ", cmn.Join(cond.OrgSystems, ","), " 在日志仓 ", e.storeName, " 中没有索引数据，直接返回空结果")
			rs := new(search.SearchResult)
			rs.Total = cmn.Uint32ToString(totalCount)
			rs.Count = "0"
			return rs
		}

		if allcnt < maxMatchCount {
			maxMatchCount = allcnt // 取其小，尽量减少最大匹配件数的误差
		}

		// 重置调整优化系统条件
		cond.System = cond.OrgSystem
		cond.Systems = cond.OrgSystems

		if len(sysnames) == 1 {
			cond.System = sysnames[0] // 权限范围内仅这个系统是有数据的，直接调整作为系统条件提高查询性能
		} else {
			cond.Systems = sysnames // 过滤掉没有数据的系统，提高查询性能
		}
	} else {
		// 重置系统条件
		cond.System = cond.OrgSystem
		cond.Systems = cond.OrgSystems
	}

	if cond.System != "" {
		// 可能内部优化调整增加了系统条件，需要确保加入检索关键字
		hasSystemCondition := false
		for _, w := range kws {
			if w == cond.System {
				hasSystemCondition = true
				break
			}
		}
		if !hasSystemCondition {
			kws = append(kws, cond.System)
		}
	}

	cond.Kws = kws
	var rs *search.SearchResult
	if len(cond.Kws) == 0 {
		// 无条件浏览模式（可能含多选条件）
		rs = search.SearchLogData(e.storeName, cond)
	} else {
		// 多关键词查询模式
		rs = search.SearchWordIndex(e.storeName, cond)
	}

	if maxMatchCount < cmn.StringToUint32(rs.Count, 0) {
		rs.Count = cmn.Uint32ToString(maxMatchCount) // 取其小，尽量减少最大匹配件数的误差
	}
	return rs
}

// 添加日志
func AddTextLog(md *logdata.LogDataModel) {
	engine := NewDefaultEngine()
	engine.AddTextLog(md.Date, md.Text, md.System)
}
