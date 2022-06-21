package sysmnt

import "glc/ldb/storage"

// 管理用存储结构体
type SysmntStorage struct {
	storage *storage.LdbStorage // 存储器
}

var sysmntStorage *SysmntStorage

func GetSysmntStorage(storeName string) *SysmntStorage {
	if sysmntStorage != nil {
		if !sysmntStorage.storage.IsClose() {
			return sysmntStorage
		}
	}
	sysmntStorage = &SysmntStorage{
		storage: storage.GetStorage(storeName, "sysmnt", nil, fnSave, nil),
	}
	return sysmntStorage
}

// 检查指定关键词是否都有数据
func (s *SysmntStorage) HasKeyWord(kws []string) bool {
	for _, k := range kws {
		_, err := s.storage.Get([]byte(k))
		if err != nil {
			return false
		}
	}
	return true
}

// 添加关键词
func (s *SysmntStorage) AddKeyWords(kws []string) {
	for _, k := range kws {
		_, err := s.storage.Get([]byte(k))
		if err == nil {
			return
		}
		s.storage.Add(k)
	}
}

// 关键词作为key保存，值为空串
func fnSave(store *storage.LdbStorage, keyWord any) (*storage.LdbDocument, any) {
	store.Put([]byte(keyWord.(string)), []byte(""))
	return nil, nil
}
