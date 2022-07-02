/**
 * 系统管理存储器
 * 1）提供无序的KV形式读写功能，利用leveldb自动序列化存盘
 * 2）使用需自行控制避免发生Key的冲突问题
 */
package sysidx

import (
	"errors"
	"glc/cmn"
	"glc/ldb/conf"
	"glc/onexit"
	"log"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// 存储结构体
type SysidxStorage struct {
	storeName string      // 存储目录
	subPath   string      // 存储目录下的相对路径（存放数据）
	leveldb   *leveldb.DB // leveldb
	lastTime  int64       // 最后一次访问时间
	closing   bool        // 是否关闭中状态
}

var sdbMu sync.Mutex             // 锁
var sysidxStorage *SysidxStorage // 缓存用存储器

func init() {
	onexit.RegisterExitHandle(onExit) // 优雅退出
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func GetSysidxStorage(storeName string) *SysidxStorage { // 存储器，文档，自定义对象

	// 缓存有则取用
	subPath := "sysidx"
	cacheName := storeName + cmn.PathSeparator() + subPath
	if sysidxStorage != nil && !sysidxStorage.IsClose() { // 尝试用缓存实例存储器
		return sysidxStorage
	}

	// 缓存无则锁后创建返回并存缓存
	sdbMu.Lock()                                          // 上锁
	defer sdbMu.Unlock()                                  // 解锁
	if sysidxStorage != nil && !sysidxStorage.IsClose() { // 再次尝试用缓存实例存储器
		return sysidxStorage
	}

	store := new(SysidxStorage)
	store.storeName = storeName
	store.subPath = subPath
	store.closing = false
	store.lastTime = time.Now().Unix()

	dbPath := conf.GetStorageRoot() + cmn.PathSeparator() + cacheName
	db, err := leveldb.OpenFile(dbPath, nil) // 打开（在指定子目录中存放数据）
	if err != nil {
		log.Println("打开SysidxStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	sysidxStorage = store // 缓存起来

	// 逐秒判断，若闲置超时则自动关闭
	go autoCloseSysidxStorageWhenMaxIdle(store)

	log.Println("打开SysidxStorage：", cacheName)
	return store
}

func autoCloseSysidxStorageWhenMaxIdle(store *SysidxStorage) {
	if conf.GetMaxIdleTime() > 0 {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			if time.Now().Unix()-store.lastTime > int64(conf.GetMaxIdleTime()) {
				store.Close()
				ticker.Stop()
				break
			}
		}
	}
}

// 关闭Storage
func (s *SysidxStorage) Close() {
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
	sysidxStorage = nil

	log.Println("关闭SysidxStorage：", s.storeName+cmn.PathSeparator()+s.subPath)
}

// 直接存入数据到leveldb
func (s *SysidxStorage) Put(key []byte, value []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Put(key, value, nil)
}

// 直接从leveldb取数据
func (s *SysidxStorage) Get(key []byte) ([]byte, error) {
	if s.closing {
		return nil, errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Get(key, nil)
}

// 存储目录名
func (s *SysidxStorage) StoreName() string {
	return s.storeName
}

// 是否关闭中状态
func (s *SysidxStorage) IsClose() bool {
	return s.closing
}

func onExit() {
	if sysidxStorage != nil {
		sysidxStorage.Close()
	}
	log.Println("退出SysidxStorage")
}
