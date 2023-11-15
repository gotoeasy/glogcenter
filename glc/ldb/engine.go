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
	for _, word := range kws {
		cnt := idxw.GetTotalCount(word)
		if (cmn.Startwiths(word, "~") && cond.OrgSystem != "" && cnt < 1) || (!cmn.Startwiths(word, "~") && cnt < 1) {
			// 系统名且是输入条件，或非系统名无索引件数
			cmn.Debug("关键词", word, "没有索引数据，直接返回空结果")
			rs := new(search.SearchResult)
			rs.Total = cmn.Uint32ToString(e.logStorage.TotalCount())
			rs.Count = cmn.Uint32ToString(0)
			return rs
		}
	}

	// 【快速检查2】，权限内的系统没有数据时，直接返回
	if cond.OrgSystem == "" && cond.OrgSystems[0] != "*" {
		var syscnt []string
		var allcnt uint32
		for _, word := range cond.OrgSystems {
			n := idxw.GetTotalCount(word)
			if n > 0 {
				syscnt = append(syscnt, word) // 这个系统有数据
			}
			allcnt += n // 累加总数
		}
		if allcnt < 1 {
			cmn.Debug("权限范围系统 ", cmn.Join(cond.OrgSystems, ","), " 在日志仓 ", e.storeName, " 中没有索引数据，直接返回空结果")
			rs := new(search.SearchResult)
			rs.Total = cmn.Uint32ToString(e.logStorage.TotalCount())
			rs.Count = cmn.Uint32ToString(0)
			return rs
		}

		// 重置调整优化系统条件
		cond.System = cond.OrgSystem
		cond.Systems = cond.OrgSystems

		if len(syscnt) == 1 {
			cond.System = syscnt[0] // 权限范围内仅这个系统是有数据的，直接调整作为系统条件提高查询性能
		} else {
			cond.Systems = syscnt // 过滤掉没有数据的系统，提高查询性能
		}
	} else {
		// 重置系统条件
		cond.System = cond.OrgSystem
		cond.Systems = cond.OrgSystems
	}

	cond.Kws = kws
	if len(cond.Kws) == 0 {
		// 无条件浏览模式（可能含多选条件）
		return search.SearchLogData(e.storeName, cond)
	}

	// 多关键词查询模式
	return search.SearchWordIndex(e.storeName, cond)
}

// 添加日志
func AddTextLog(md *logdata.LogDataModel) {
	engine := NewDefaultEngine()
	engine.AddTextLog(md.Date, md.Text, md.System)
}
