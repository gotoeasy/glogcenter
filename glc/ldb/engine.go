package ldb

import (
	"glc/ldb/logdata"
	"glc/ldb/search"
	"glc/ldb/sysmnt"
	"glc/ldb/tokenizer"
)

type Engine struct {
	storeName  string
	logStorage *logdata.LogDataStorage // 日志存储器
	sysStorage *sysmnt.SysmntStorage   // 系统存储器
}

func NewEngine(storeName string) *Engine {
	return &Engine{
		storeName:  storeName,
		logStorage: logdata.NewLogDataStorage(storeName),
		sysStorage: sysmnt.GetSysmntStorage(storeName),
	}
}

func NewDefaultEngine() *Engine {
	var storeName string = "default"
	return &Engine{
		storeName:  storeName,
		logStorage: logdata.NewLogDataStorage(storeName),
		sysStorage: sysmnt.GetSysmntStorage(storeName),
	}
}

// 添加日志
func (e *Engine) AddTextLog(logText string) {
	e.logStorage.AddTextLog(logText)
}

func (e *Engine) Search(searchKey string) *search.SearchResult {

	// 分词后检索
	kws := tokenizer.CutForSearch(searchKey)

	// 无数据判断，同布隆判断效果
	if !e.sysStorage.HasKeyWord(kws) {
		return new(search.SearchResult)
	}

	// 无条件浏览模式
	if len(kws) == 0 {
		return search.BrowsePage(e.storeName, 1, 50)
	}

	// 单关键词查询模式
	if len(kws) == 1 {
		return search.BrowsePageByOneWord(e.storeName, kws[0], 1, 50)
	}

	// 关键词检索模式
	return search.Search(e.storeName, tokenizer.CutForSearch(searchKey), 1, 20)
}
