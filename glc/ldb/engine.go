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
	adds = append(adds, cond.System, cond.Loglevel)
	kws := tokenizer.CutForSearchEx(cond.SearchKey, adds, nil) // 检索用关键词处理

	// 简单检查，存在无索引数据的关键词时，直接返回
	for _, word := range kws {
		idxw := indexword.NewWordIndexStorage(e.storeName)
		if idxw.GetTotalCount(word) < 1 {
			cmn.Debug("关键词", word, "没有索引数据，直接返回空结果")
			rs := new(search.SearchResult)
			rs.Total = cmn.Uint32ToString(e.logStorage.TotalCount())
			rs.Count = cmn.Uint32ToString(0)
			return rs
		}
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
