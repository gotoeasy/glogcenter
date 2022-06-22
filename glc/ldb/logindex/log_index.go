package logindex

import (
	"glc/cmn"
	"glc/ldb/storage"
	"log"
	"math"
)

var mapStorage map[string](*LogIndexStorage)

// 单关键词反向索引存储结构体
type LogIndexStorage struct {
	storage *storage.LdbStorage // 存储器
}

func init() {
	mapStorage = make(map[string](*LogIndexStorage))
}

// 单关键词反向索引存储器
func NewLogIndexStorage(storeName string, word string) *LogIndexStorage {

	var subPath string = "inverted" + cmn.PathSeparator() + cmn.HashMod(word, 100) + cmn.PathSeparator() + "k_" + cmn.HashMod(word, math.MaxUint32)
	cacheName := storeName + ":" + subPath
	cacheStore := mapStorage[cacheName] // 缓存中的存储对象
	if cacheStore != nil {
		if !cacheStore.storage.IsClose() {
			return cacheStore
		}
	}

	store := &LogIndexStorage{
		storage: storage.GetStorage(storeName, subPath, nil, fnSave, nil),
	}
	mapStorage[cacheName] = store
	return store
}

// 添加日志id
func (s *LogIndexStorage) Add(id uint32) {
	s.storage.Add(id)
}

func fnSave(store *storage.LdbStorage, id any) (*storage.LdbDocument, any) {
	// key递增(此时TotalCount已递增完成)，value是文档id
	store.Put(cmn.Uint32ToBytes(store.TotalCount()), cmn.Uint32ToBytes(id.(uint32)))
	log.Println("保存关键字索引 ", id)
	return nil, nil
}

// 取日志id
func (s *LogIndexStorage) Get(id uint32) uint32 {
	bytes, err := s.storage.Get(cmn.Uint32ToBytes(id))
	if err != nil {
		log.Println("id=", id, err)
	}
	return cmn.BytesToUint32(bytes)
}

// 总件数
func (s *LogIndexStorage) TotalCount() uint32 {
	return s.storage.TotalCount()
}
