/**
 * 系统管理存储器
 * 1）提供无序的KV形式读写功能，利用leveldb自动序列化存盘
 * 2）使用需自行控制避免发生Key的冲突问题
 */
package sysmnt

import (
	"errors"
	"glc/cmn"
	"glc/conf"
	"glc/onexit"
	"log"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// 存储结构体
type SysmntStorage struct {
	subPath  string      // 存储目录下的相对路径（存放数据）
	leveldb  *leveldb.DB // leveldb
	lastTime int64       // 最后一次访问时间
	closing  bool        // 是否关闭中状态
}

var sdbMu sync.Mutex             // 锁
var sysmntStorage *SysmntStorage // 缓存用存储器

func init() {
	onexit.RegisterExitHandle(onExit) // 优雅退出
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func NewSysmntStorage() *SysmntStorage { // 存储器，文档，自定义对象

	// 缓存有则取用
	subPath := ".sysmnt"
	cacheName := subPath
	if sysmntStorage != nil && !sysmntStorage.IsClose() { // 尝试用缓存实例存储器
		return sysmntStorage
	}

	// 缓存无则锁后创建返回并存缓存
	sdbMu.Lock()                                          // 上锁
	defer sdbMu.Unlock()                                  // 解锁
	if sysmntStorage != nil && !sysmntStorage.IsClose() { // 再次尝试用缓存实例存储器
		return sysmntStorage
	}

	store := new(SysmntStorage)
	store.subPath = subPath
	store.closing = false
	store.lastTime = time.Now().Unix()

	dbPath := conf.GetStorageRoot() + cmn.PathSeparator() + cacheName
	db, err := leveldb.OpenFile(dbPath, nil) // 打开（在指定子目录中存放数据）
	if err != nil {
		log.Println("打开SysmntStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	sysmntStorage = store // 缓存起来

	// 逐秒判断，若闲置超时则自动关闭
	go store.autoCloseWhenMaxIdle()

	log.Println("打开SysmntStorage：", cacheName)
	return store
}

func (s *SysmntStorage) autoCloseWhenMaxIdle() {
	if conf.GetMaxIdleTime() > 0 {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			if time.Now().Unix()-s.lastTime > int64(conf.GetMaxIdleTime()) {
				s.Close()
				ticker.Stop()
				break
			}
		}
	}
}

// 关闭Storage
func (s *SysmntStorage) Close() {
	if s == nil || s.closing { // 优雅退出时可能会正好nil，判断一下优雅点
		return
	}

	sdbMu.Lock()         // 锁
	defer sdbMu.Unlock() // 解锁
	if s.closing {
		return
	}

	s.closing = true
	s.leveldb.Close()
	sysmntStorage = nil

	log.Println("关闭SysmntStorage：", s.subPath)
}

func (s *SysmntStorage) GetStorageDataCount(storeName string) uint32 {
	bt, err := s.Get(cmn.StringToBytes("data:" + storeName))
	if err != nil {
		return 0
	}
	return cmn.BytesToUint32(bt)
}

func (s *SysmntStorage) SetStorageDataCount(storeName string, count uint32) {
	s.Put(cmn.StringToBytes("data:"+storeName), cmn.Uint32ToBytes(count))
}

func (s *SysmntStorage) GetStorageIndexCount(storeName string) uint32 {
	bt, err := s.Get(cmn.StringToBytes("index:" + storeName))
	if err != nil {
		return 0
	}
	return cmn.BytesToUint32(bt)
}

func (s *SysmntStorage) SetStorageIndexCount(storeName string, count uint32) {
	s.Put(cmn.StringToBytes("index:"+storeName), cmn.Uint32ToBytes(count))
}

func (s *SysmntStorage) DeleteStorageInfo(storeName string) error {
	err := s.leveldb.Delete(cmn.StringToBytes("data:"+storeName), nil)
	if err != nil {
		return err
	}
	err = s.leveldb.Delete(cmn.StringToBytes("index:"+storeName), nil)
	if err != nil {
		return err
	}
	return nil
}

// 直接存入数据到leveldb
func (s *SysmntStorage) Put(key []byte, value []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Put(key, value, nil)
}

// 直接从leveldb取数据
func (s *SysmntStorage) Get(key []byte) ([]byte, error) {
	if s.closing {
		return nil, errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Get(key, nil)
}

// 直接从leveldb取数据
func (s *SysmntStorage) Del(key []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Delete(key, nil)
}

// 是否关闭中状态
func (s *SysmntStorage) IsClose() bool {
	return s.closing
}

func onExit() {
	if sysmntStorage != nil {
		sysmntStorage.Close()
	}
	log.Println("退出SysmntStorage")
}
