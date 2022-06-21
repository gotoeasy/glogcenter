package logdata

import (
	"glc/cmn"
	"glc/ldb/logindex"
	"glc/ldb/storage"
	"glc/ldb/sysmnt"
	"glc/ldb/tokenizer"
	"log"
	"strings"
)

var mapStorage map[string](*LogDataStorage)

// 日志存储结构体
type LogDataStorage struct {
	storage *storage.LdbStorage // 存储器
}

func init() {
	mapStorage = make(map[string](*LogDataStorage))
}

// 创建日志存储
func NewLogDataStorage(storeName string) *LogDataStorage {

	cacheStore := mapStorage[storeName] // 缓存中的存储对象
	if cacheStore != nil {
		if !cacheStore.storage.IsClose() {
			return cacheStore
		}
	}

	store := &LogDataStorage{
		storage: storage.GetStorage(storeName, "data", fnBeforeSave, fnSave, fnAfterSave),
	}
	mapStorage[storeName] = store
	return store
}

// 添加日志（参数是普通文本日志）
func (s *LogDataStorage) AddTextLog(logText string) {
	txt := strings.TrimSpace(logText)
	if txt == "" {
		return
	}
	ary := strings.Split(txt, "\n")

	d := new(LogDataModel)
	d.Text = ary[0]
	if len(ary) > 1 {
		d.Detail = txt
	}

	s.storage.Add(d)
}

// // 添加日志（参数是LogDataModel形式的json字符串）
// func (s *LogDataStorage) AddJsonLog(logJson string) {
// 	s.storage.Add(logJson)
// }

// 日志预处理
func fnBeforeSave(store *storage.LdbStorage, logDataModel any) any {
	// TODO 预处理
	return logDataModel
}

// 日志保存处理
func fnSave(store *storage.LdbStorage, logDataModel any) (*storage.LdbDocument, any) {
	model := logDataModel.(*LogDataModel)
	d := new(storage.LdbDocument)
	d.Id = store.TotalCount()                       // 已递增好的值
	d.Content = model.ToJson()                      // 转json作为内容
	store.Put(cmn.Uint32ToBytes(d.Id), d.ToBytes()) // 保存
	log.Println("保存日志数据 ", d.Id)
	return d, logDataModel
}

// 日志后处理
func fnAfterSave(store *storage.LdbStorage, doc *storage.LdbDocument, logDataModel any) {

	model := logDataModel.(*LogDataModel)
	kws := tokenizer.CutForSearch(model.Text)
	mnt := sysmnt.GetSysmntStorage(store.StoreName())
	mnt.AddKeyWords(kws)
	log.Println("日志后处理 ", doc.Id)

	for _, key := range kws {
		idx := logindex.NewLogIndexStorage(store.StoreName(), key)
		idx.Add(doc.Id) // 反向索引
	}
}

// 取日志
func (s *LogDataStorage) Get(id uint32) *storage.LdbDocument {
	bytes, _ := s.storage.Get(cmn.Uint32ToBytes(id))
	return storage.ParseBytes(bytes)
}

// 取日志（Json字符串形式）
func (s *LogDataStorage) GetJsonLog(id uint32) string {
	bytes, _ := s.storage.Get(cmn.Uint32ToBytes(id))
	doc := storage.ParseBytes(bytes)
	return doc.Content
}

// 总件数
func (s *LogDataStorage) TotalCount() uint32 {
	return s.storage.TotalCount()
}
