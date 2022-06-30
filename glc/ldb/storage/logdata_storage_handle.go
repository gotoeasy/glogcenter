/**
 * 日志存储器控制
 * 1）添加日志的入口
 * 2）优先响应保存日志，闲时创建关键词反向索引
 * 3）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 */
package storage

import (
	"glc/cmn"
	"log"
	"strings"
)

var mapLogDataStorageHandle map[string](*LogDataStorageHandle)

// 日志存储结构体
type LogDataStorageHandle struct {
	storage *LogDataStorage // 存储器
}

func init() {
	mapLogDataStorageHandle = make(map[string](*LogDataStorageHandle))
}

// 创建日志存储
func NewLogDataStorageHandle(storeName string) *LogDataStorageHandle {

	cacheStore := mapLogDataStorageHandle[storeName] // 缓存中的存储对象
	if cacheStore != nil {
		if !cacheStore.storage.IsClose() {
			return cacheStore
		}
	}

	store := &LogDataStorageHandle{
		storage: NewLogDataStorage(storeName, "data"),
	}
	mapLogDataStorageHandle[storeName] = store
	return store
}

// 添加日志（参数是普通文本日志）
func (s *LogDataStorageHandle) AddTextLog(date string, logText string, system string) {
	txt := strings.TrimSpace(logText)
	if txt == "" {
		return
	}
	ary := strings.Split(txt, "\n")

	d := new(LogDataModel)
	d.Text = strings.TrimSpace(ary[0])
	if len(ary) > 1 {
		d.Detail = txt
	}
	d.Date = date
	d.System = system

	if s.storage.IsClose() {
		s.storage = NewLogDataStorage(s.storage.storeName, "data")
	}
	err := s.storage.Add(d)
	if err != nil {
		log.Println("竟然失败，再来一次", s.storage.IsClose(), err)
		if s.storage.IsClose() {
			s.storage = NewLogDataStorage(s.storage.storeName, "data")
		}
		s.storage.Add(d)
	}
}

// // 添加日志（参数是LogDataModel形式的json字符串）
// func (s *LogDataStorage) AddJsonLog(logJson string) {
// 	s.storage.Add(logJson)
// }

// 取日志（文档）
func (s *LogDataStorageHandle) GetLogDataDocument(id uint64) *LogDataDocument {
	bytes, _ := s.storage.Get(cmn.Uint64ToBytes(id))
	doc := new(LogDataDocument)
	doc.LoadBytes(bytes)
	return doc
}

// 取日志（模型）
func (s *LogDataStorageHandle) GetLogDataModel(id uint64) *LogDataModel {
	d := s.GetLogDataDocument(id)
	m := new(LogDataModel)
	m.LoadJson(d.Content)
	return m
}

// 总件数
func (s *LogDataStorageHandle) TotalCount() uint64 {
	return s.storage.TotalCount()
}
