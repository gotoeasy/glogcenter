/**
 * 日志存储器控制
 * 1）添加日志的入口
 * 2）优先响应保存日志，闲时才建索引
 * 3）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 */
package storage

import (
	"glc/ldb/storage/logdata"
	"sync"

	"github.com/gotoeasy/glang/cmn"
)

var mapStorageHandle sync.Map
var mapStorageMu sync.Mutex

// 日志存储结构体
type LogDataStorageHandle struct {
	storage *logdata.LogDataStorage // 存储器
}

// 创建日志存储
func NewLogDataStorageHandle(storeName string) *LogDataStorageHandle {
	value, ok := mapStorageHandle.Load(storeName)
	if ok && value != nil {
		cacheStore := value.(*LogDataStorageHandle)
		if !cacheStore.storage.IsClose() {
			return cacheStore // 缓存中未关闭的存储对象
		}
	}

	// 删除缓存中已关闭的存储对象
	mapStorageMu.Lock()
	defer mapStorageMu.Unlock()
	var delKeys []any
	mapStorageHandle.Range(func(key, value any) bool {
		s := value.(*LogDataStorageHandle)
		if s.storage.IsClose() {
			delKeys = append(delKeys, key)
		}
		return true
	})
	for _, key := range delKeys {
		mapStorageHandle.Delete(key)
	}

	store := &LogDataStorageHandle{
		storage: logdata.NewLogDataStorage(storeName),
	}
	mapStorageHandle.Store(storeName, store)
	return store
}

// 添加日志（参数是普通文本日志）
func (s *LogDataStorageHandle) AddTextLog(date string, logText string, system string) {
	txt := cmn.Trim(logText)
	ary := cmn.Split(txt, "\n")

	d := new(logdata.LogDataModel)
	d.Text = cmn.Trim(ary[0])
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
		cmn.Error("竟然失败，再来一次", s.storage.IsClose(), err)
		if s.storage.IsClose() {
			s.storage = logdata.NewLogDataStorage(s.storage.StoreName())
		}
		s.storage.Add(d)
	}
}

// 添加日志（参数LogDataModel）
func (s *LogDataStorageHandle) AddLogDataModel(data *logdata.LogDataModel) {
	ary := cmn.Split(data.Text, "\n")
	if len(ary) > 1 {
		data.Detail = data.Text
		data.Text = ary[0]
	}

	if s.storage.IsClose() {
		s.storage = logdata.NewLogDataStorage(s.storage.StoreName())
	}
	err := s.storage.Add(data)
	if err != nil {
		cmn.Error("竟然失败，再来一次", s.storage.IsClose(), err)
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

func (s *LogDataStorageHandle) GetStoreName() string {
	return s.storage.StoreName()
}
