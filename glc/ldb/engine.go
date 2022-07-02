package ldb

import (
	"glc/cmn"
	"glc/ldb/search"
	"glc/ldb/storage"
	"glc/ldb/sysidx"
	"glc/ldb/tokenizer"
	"log"
)

type Engine struct {
	storeName  string
	logStorage *storage.LogDataStorageHandle // 日志存储控制器
	sysStorage *sysidx.SysidxStorage         // 系统存储器
}

func NewEngine(storeName string) *Engine {
	storeName = cmn.GeyStoreNameByDate(storeName)
	return &Engine{
		storeName:  storeName,
		logStorage: storage.NewLogDataStorageHandle(storeName),
		sysStorage: sysidx.GetSysidxStorage(storeName),
	}
}

func NewDefaultEngine() *Engine {
	var storeName string = cmn.GeyStoreNameByDate("default")
	return &Engine{
		storeName:  storeName,
		logStorage: storage.NewLogDataStorageHandle(storeName),
		sysStorage: sysidx.GetSysidxStorage(storeName),
	}
}

// 添加日志
func (e *Engine) AddTextLog(date string, logText string, system string) {
	e.logStorage.AddTextLog(date, logText, system)
}

func (e *Engine) Search(searchKey string, pageSize int, currentDocId uint32, forward bool) *search.SearchResult {

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
		log.Println("无条件查询")
	} else {
		log.Println("查询", searchKey, "，分词后检索", kws)
	}

	// 无数据判断，同布隆判断效果
	if !e.sysStorage.ContainsKeyWord(kws) {
		log.Println("检索条件布隆检查发现有不存在的关键词，直接返回空结果", kws)
		return new(search.SearchResult)
	}

	if len(kws) == 0 {
		// 无条件浏览模式
		return search.SearchLogData(e.storeName, pageSize, currentDocId, forward)
	} else if len(kws) == 1 {
		// 单关键词查询模式
		return search.SearchWordIndex(e.storeName, kws[0], pageSize, currentDocId, forward)
	} else {
		// 多关键词查询模式
		return search.Search(e.storeName, kws, pageSize, currentDocId, forward)
	}

}
