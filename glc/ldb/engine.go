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

func (e *Engine) Search(searchKey string, minDatetime string, maxDatetime string, pageSize int, currentDocId uint32, forward bool) *search.SearchResult {

	// 检查修正pageSize
	if pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	// 分词后检索
	kws := tokenizer.CutForSearch(searchKey) // TODO 检索用关键词处理

	if searchKey == "" {
		cmn.Debug("无条件查询", "currentDocId =", currentDocId)
	} else {
		cmn.Debug("查询", searchKey, "，分词后检索", kws, "currentDocId =", currentDocId)
	}

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

	if len(kws) == 0 {
		// 无条件浏览模式
		return search.SearchLogData(e.storeName, pageSize, currentDocId, forward, minDatetime, maxDatetime)
	}

	// 多关键词查询模式
	return search.SearchWordIndex(e.storeName, kws, pageSize, currentDocId, forward, minDatetime, maxDatetime)
}

// 添加日志
func AddTextLog(md *logdata.LogDataModel) {
	engine := NewDefaultEngine()
	engine.AddTextLog(md.Date, md.Text, md.System)
}
