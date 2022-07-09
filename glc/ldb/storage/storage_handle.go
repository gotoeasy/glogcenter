/**
 * 日志存储器控制
 * 1）添加日志的入口
 * 2）优先响应保存日志，闲时才建索引
 * 3）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 */
package storage

import (
	"glc/cmn"
	"glc/ldb/storage/logdata"
	"log"
	"strings"
)

var mapStorageHandle map[string](*LogDataStorageHandle)

// 日志存储结构体
type LogDataStorageHandle struct {
	storage *logdata.LogDataStorage // 存储器
}

func init() {
	mapStorageHandle = make(map[string](*LogDataStorageHandle))
}

// 创建日志存储
func NewLogDataStorageHandle(storeName string) *LogDataStorageHandle {

	cacheStore := mapStorageHandle[storeName] // 缓存中的存储对象
	if cacheStore != nil {
		if !cacheStore.storage.IsClose() {
			return cacheStore
		}
	}

	store := &LogDataStorageHandle{
		storage: logdata.NewLogDataStorage(storeName),
	}
	mapStorageHandle[storeName] = store
	return store
}

// 添加日志（参数是普通文本日志）
func (s *LogDataStorageHandle) AddTextLog(date string, logText string, system string) {
	txt := strings.TrimSpace(logText)
	if txt == "" {
		return
	}
	ary := strings.Split(txt, "\n")

	d := new(logdata.LogDataModel)
	d.Text = strings.TrimSpace(ary[0])
	if len(ary) > 1 {
		d.Detail = txt
	}
	d.Date = date
	d.System = system

	if s.storage.IsClose() {
		s.storage = logdata.NewLogDataStorage(s.storage.StoreName())
	}
	err := s.storage.Add(d)
	if err != nil {
		log.Println("竟然失败，再来一次", s.storage.IsClose(), err)
		if s.storage.IsClose() {
			s.storage = logdata.NewLogDataStorage(s.storage.StoreName())
		}
		s.storage.Add(d)
	}
}

// 添加日志（参数LogDataModel）
func (s *LogDataStorageHandle) AddLogDataModel(data *logdata.LogDataModel) {
	ary := strings.Split(data.Text, "\n")
	data.Text = strings.TrimSpace(ary[0])
	if len(ary) > 1 {
		data.Detail = data.Text
	}

	if s.storage.IsClose() {
		s.storage = logdata.NewLogDataStorage(s.storage.StoreName())
	}
	err := s.storage.Add(data)
	if err != nil {
		log.Println("竟然失败，再来一次", s.storage.IsClose(), err)
		if s.storage.IsClose() {
			s.storage = logdata.NewLogDataStorage(s.storage.StoreName())
		}
		s.storage.Add(data)
	}
}

// 取日志（文档）
func (s *LogDataStorageHandle) GetLogDataDocument(id uint32) *logdata.LogDataDocument {
	bytes, _ := s.storage.Get(cmn.Uint32ToBytes(id))
	doc := new(logdata.LogDataDocument)
	doc.LoadBytes(bytes)
	return doc
}

// 取日志（模型）
func (s *LogDataStorageHandle) GetLogDataModel(id uint32) *logdata.LogDataModel {
	d := s.GetLogDataDocument(id)
	m := new(logdata.LogDataModel)
	m.LoadJson(d.Content)
	return m
}

// 总件数
func (s *LogDataStorageHandle) TotalCount() uint32 {
	return s.storage.TotalCount()
}
