package ldb

import (
	"glc/cmn"
	"glc/ldb/search"
	"glc/ldb/storage"
	"glc/ldb/sysmnt"
	"glc/ldb/tokenizer"
	"log"
)

type Engine struct {
	storeName  string
	logStorage *storage.LogDataStorageHandle // 日志存储控制器
	sysStorage *sysmnt.SysmntStorage         // 系统存储器
}

func NewEngine(storeName string) *Engine {
	storeName = cmn.GeyStoreNameByDate(storeName)
	return &Engine{
		storeName:  storeName,
		logStorage: storage.NewLogDataStorageHandle(storeName),
		sysStorage: sysmnt.GetSysmntStorage(storeName),
	}
}

func NewDefaultEngine() *Engine {
	var storeName string = cmn.GeyStoreNameByDate("default")
	return &Engine{
		storeName:  storeName,
		logStorage: storage.NewLogDataStorageHandle(storeName),
		sysStorage: sysmnt.GetSysmntStorage(storeName),
	}
}

// 添加日志
func (e *Engine) AddTextLog(logText string) {
	e.logStorage.AddTextLog(logText)
}

func (e *Engine) Search(searchKey string) *search.SearchResult {

	// 分词后检索
	kws := tokenizer.CutForSearch(searchKey) // TODO 检索用关键词处理

	log.Println("查询", kws)

	// 无数据判断，同布隆判断效果
	if !e.sysStorage.ContainsKeyWord(kws) {
		log.Println("检索条件布隆检查发现有不存在的关键词，直接返回空结果", kws)
		return new(search.SearchResult)
	}

	// 无条件浏览模式
	if len(kws) == 0 {
		return search.Search(e.storeName, "", 20, 0, true)
	}

	// 单关键词查询模式
	return search.Search(e.storeName, kws[0], 20, 11, false)
}
