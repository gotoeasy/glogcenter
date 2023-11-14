/**
 * 文档关键词索引序号存储器
 * 1）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 * 2）单线程调用设计，由日志存储器内部控制安全的调用，其他地方调用可能会有问题
 */
package indexdoc

import (
	"glc/com"
	"glc/conf"
	"glc/ldb/status"
	"sync"
	"time"

	"github.com/gotoeasy/glang/cmn"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DocIndexStorage struct {
	storeName string      // 存储目录
	subPath   string      // 存储目录下的相对路径（存放数据）
	leveldb   *leveldb.DB // leveldb
	lastTime  int64       // 最后一次访问时间
	closing   bool        // 是否关闭中状态
	mu        sync.Mutex  // 锁
}

var idxMu sync.Mutex
var mapStorage map[string](*DocIndexStorage)
var mapStorageMu sync.Mutex

func init() {
	mapStorage = make(map[string](*DocIndexStorage))
	cmn.OnExit(onExit) // 优雅退出
}

func getStorage(cacheName string) *DocIndexStorage {
	cacheStore := mapStorage[cacheName]
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore // 缓存中未关闭的存储对象
	}
	return nil
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func NewDocIndexStorage(storeName string) *DocIndexStorage { // 存储器，文档，自定义对象

	// 缓存有则取用
	subPath := "inverted" + cmn.PathSeparator() + "d"
	cacheName := storeName + cmn.PathSeparator() + subPath
	cacheStore := getStorage(cacheName)
	if cacheStore != nil {
		return cacheStore
	}

	// 缓存无则锁后创建返回并存缓存
	mapStorageMu.Lock()                // 缓存map锁
	defer mapStorageMu.Unlock()        // 缓存map解锁
	idxMu.Lock()                       // 上锁
	defer idxMu.Unlock()               // 解锁
	cacheStore = getStorage(cacheName) // 再次尝试取用缓存中存储器
	if cacheStore != nil {
		return cacheStore
	}

	store := new(DocIndexStorage)
	store.storeName = storeName
	store.subPath = subPath
	store.closing = false
	store.lastTime = time.Now().Unix()

	dbPath := conf.GetStorageRoot() + cmn.PathSeparator() + cacheName
	option := new(opt.Options)                  // leveldb选项
	option.Filter = filter.NewBloomFilter(10)   // 使用布隆过滤器
	db, err := leveldb.OpenFile(dbPath, option) // 打开（在指定子目录中存放数据）
	if err != nil {
		cmn.Error("打开DocIndexStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	status.UpdateStorageStatus(storeName, true) // 更新状态：当前日志仓打开
	mapStorage[cacheName] = store               // 缓存起来

	// 逐秒判断，若闲置超时则自动关闭
	go store.autoCloseWhenMaxIdle()

	cmn.Info("打开DocIndexStorage：", cacheName)
	return store
}

func (s *DocIndexStorage) autoCloseWhenMaxIdle() {
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

// 添加日志反向索引
func (s *DocIndexStorage) AddWordDocSeq(word string, docId uint32, seq uint32) error {

	s.lastTime = time.Now().Unix()
	err := s.leveldb.Put(com.JoinBytes(cmn.StringToBytes(word), cmn.Uint32ToBytes(docId)), cmn.Uint32ToBytes(seq), nil)
	if err != nil {
		cmn.Error("保存日志反向索引失败", err)
		return err
	}
	return nil
}

// 取日志所在关键词索引中的序号（返回0表示有问题）
func (s *DocIndexStorage) GetWordDocSeq(word string, docId uint32) uint32 {
	if s.closing {
		return 0
	}
	s.lastTime = time.Now().Unix()
	b, err := s.leveldb.Get(com.JoinBytes(cmn.StringToBytes(word), cmn.Uint32ToBytes(docId)), nil)
	if err != nil {
		return 0
	}
	return cmn.BytesToUint32(b)
}

// 关闭Storage
func (s *DocIndexStorage) Close() {
	if s == nil || s.closing { // 优雅退出时可能会正好nil，判断一下优雅点
		return
	}

	mapStorageMu.Lock()         // 缓存map锁
	defer mapStorageMu.Unlock() // 缓存map解锁
	s.mu.Lock()                 // 对象锁
	defer s.mu.Unlock()         // 对象解锁
	if s.closing {
		return
	}

	s.closing = true
	s.leveldb.Close()             // 走到这里时没有db操作了，可以关闭
	idxMu.Lock()                  // map锁
	defer idxMu.Unlock()          // map解锁
	mapStorage[s.storeName] = nil // 设空，下回GetStorage时自动再创建

	cmn.Info("关闭DocIndexStorage：", s.storeName+cmn.PathSeparator()+s.subPath)
}

// 存储目录名
func (s *DocIndexStorage) StoreName() string {
	return s.storeName
}

// 是否关闭中状态
func (s *DocIndexStorage) IsClose() bool {
	return s.closing
}

func onExit() {
	for k := range mapStorage {
		s := mapStorage[k]
		if s != nil {
			s.Close()
		}
	}
	cmn.Info("退出DocIndexStorage")
}
