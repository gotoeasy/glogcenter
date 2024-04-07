/**
 * 关键词反向索引存储器
 * 1）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 * 2）单线程调用设计，由日志存储器内部控制安全的调用，其他地方调用可能会有问题
 */
package indexword

import (
	"glc/com"
	"glc/conf"
	"glc/ldb/status"
	"glc/ldb/storage/indexdoc"
	"sync"
	"time"

	"github.com/gotoeasy/glang/cmn"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type WordIndexStorage struct {
	storeName    string      // 存储目录
	subPath      string      // 存储目录下的相对路径（存放数据）
	leveldb      *leveldb.DB // leveldb
	lastTime     int64       // 最后一次访问时间
	indexedCount uint32      // 已建索引件数
	closing      bool        // 是否关闭中状态
	mu           sync.Mutex  // 锁
}

var zeroUint32Bytes []byte = cmn.Uint32ToBytes(0)
var zeroUint16Bytes []byte = cmn.Uint16ToBytes(0) // 索引件数的key

var idxMu sync.Mutex
var mapStorage map[string](*WordIndexStorage)
var mapStorageMu sync.Mutex

func init() {
	mapStorage = make(map[string](*WordIndexStorage))
	cmn.OnExit(onExit) // 优雅退出
}

func getStorage(cacheName string) *WordIndexStorage {
	cacheStore := mapStorage[cacheName]
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore // 缓存中未关闭的存储对象
	}
	return nil
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func NewWordIndexStorage(storeName string) *WordIndexStorage { // 存储器，文档，自定义对象

	// 缓存有则取用
	subPath := "inverted" + cmn.PathSeparator() + "k"
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

	store := new(WordIndexStorage)
	store.storeName = storeName
	store.subPath = subPath
	store.closing = false
	store.lastTime = time.Now().Unix()

	dbPath := conf.GetStorageRoot() + cmn.PathSeparator() + cacheName
	option := new(opt.Options)                  // leveldb选项
	option.Filter = filter.NewBloomFilter(10)   // 使用布隆过滤器
	db, err := leveldb.OpenFile(dbPath, option) // 打开（在指定子目录中存放数据）
	if err != nil {
		cmn.Error("打开WordIndexStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	store.loadIndexedCount()                    // 加载已建索引件数
	status.UpdateStorageStatus(storeName, true) // 更新状态：当前日志仓打开
	mapStorage[cacheName] = store               // 缓存起来

	// 逐秒判断，若闲置超时则自动关闭
	go store.autoCloseWhenMaxIdle()

	cmn.Info("打开WordIndexStorage：", cacheName)
	return store
}

func (s *WordIndexStorage) autoCloseWhenMaxIdle() {
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

// 取已建索引件数
func (s *WordIndexStorage) GetIndexedCount() uint32 {
	return s.indexedCount
}

// 初期加载已建索引件数
func (s *WordIndexStorage) loadIndexedCount() {
	bt, err := s.leveldb.Get(zeroUint16Bytes, nil)
	if err == nil {
		s.indexedCount = cmn.BytesToUint32(bt)
	}
}

// 保存已建索引件数
func (s *WordIndexStorage) SaveIndexedCount(count uint32) error {
	s.indexedCount = count
	return s.leveldb.Put(zeroUint16Bytes, cmn.Uint32ToBytes(count), nil)
}

// 取关键词索引当前的文档数
func (s *WordIndexStorage) GetTotalCount(word string) uint32 {
	wordBytes := cmn.StringToBytes(word)
	bt, err := s.leveldb.Get(com.JoinBytes(wordBytes, zeroUint32Bytes), nil) // TODO 是否有性能问题?
	if err != nil {
		return 0
	}
	return cmn.BytesToUint32(bt)
}

// 存关键词索引当前的文档数
func (s *WordIndexStorage) setTotalCount(word string, cnt uint32) error {
	return s.leveldb.Put(com.JoinBytes(cmn.StringToBytes(word), zeroUint32Bytes), cmn.Uint32ToBytes(cnt), nil)
}

// 添加关键词反向索引
func (s *WordIndexStorage) Add(word string, docId uint32) error {

	// 加关键词反向索引
	s.lastTime = time.Now().Unix()
	seq := s.GetTotalCount(word)
	seq++
	err := s.leveldb.Put(com.JoinBytes(cmn.StringToBytes(word), cmn.Uint32ToBytes(seq)), cmn.Uint32ToBytes(docId), nil)
	if err != nil {
		cmn.Error("保存关键词反向索引失败", err)
		return err
	}

	// 添加文档反向索引
	diStorage := indexdoc.NewDocIndexStorage(s.storeName)
	err = diStorage.AddWordDocSeq(word, docId, seq)
	if err != nil {
		cmn.Error("保存日志反向索引失败", err)
		return err
	}

	// 保存建好的索引数
	err = s.setTotalCount(word, seq)
	if err != nil {
		cmn.Error("保存关键词反向索引件数失败", err)
		return err // 忽略事务问题，可下回重建
	}
	// cmn.Debug("创建日志索引：", docId, "，关键词：", word)
	return nil
}

// 取日志ID（返回0表示有问题）
func (s *WordIndexStorage) GetDocId(word string, seq uint32) uint32 {
	if s.closing {
		return 0
	}
	s.lastTime = time.Now().Unix()
	b, err := s.leveldb.Get(com.JoinBytes(cmn.StringToBytes(word), cmn.Uint32ToBytes(seq)), nil)
	if err != nil {
		return 0
	}
	return cmn.BytesToUint32(b)
}

// 关闭Storage
func (s *WordIndexStorage) Close() {
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

	cmn.Info("关闭WordIndexStorage：", s.storeName+cmn.PathSeparator()+s.subPath)
}

// 存储目录名
func (s *WordIndexStorage) StoreName() string {
	return s.storeName
}

// 是否关闭中状态
func (s *WordIndexStorage) IsClose() bool {
	return s.closing
}

func onExit() {
	for k := range mapStorage {
		s := mapStorage[k]
		if s != nil {
			s.Close()
		}
	}
	cmn.Info("退出WordIndexStorage")
}
