package storage

import (
	"errors"
	"glc/cmn"
	"log"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// 存储结构体
type LdbStorage struct {
	storeName    string      // 存储目录名
	subPath      string      // 存储目录名下的目录
	storeChan    chan any    // 存储通道
	leveldb      *leveldb.DB // leveldb
	currentCount uint32      // 当前件数
	lastTime     int64       // 最后一次访问时间
	closing      bool        // 是否关闭中状态
	mu           sync.Mutex  // 锁
}

var storePath string    // 存储根目录
var storeChanBuffer int // 存储器通道缓冲数，可通过环境变量‘STORE_CHAN_BUFFER’设定
var timeout int         // 默认在5分钟内没有操作则关闭相应的Storage，可通过环境变量‘TIMEOUT’设定

var mu sync.Mutex
var wg sync.WaitGroup
var mapStorage map[string](*LdbStorage)

func init() {
	mapStorage = make(map[string](*LdbStorage))

	// 检查环境变量设定配置
	timeout = cmn.GetenvInt("TIMEOUT", 300)                  // 默认秒
	storeChanBuffer = cmn.GetenvInt("STORE_CHAN_BUFFER", 64) // 默认秒
	storePath = cmn.Getenv("STORE_PATH", "e:\\222")          // 默认‘/glogcenter’
}

func getCacheStore(cacheName string) *LdbStorage {
	cacheStore := mapStorage[cacheName] // 缓存中的存储对象
	if cacheStore != nil {
		cacheStore.lastTime = time.Now().Unix()
		return cacheStore
	}
	return nil
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func GetStorage(storeName string, subPath string,
	fnBeforeSave func(*LdbStorage, any) any, // 存储器，原始入参，预备对象
	fnSave func(*LdbStorage, any) (*LdbDocument, any), // 存储器，预备对象，文档，自定义对象
	fnAfterSave func(*LdbStorage, *LdbDocument, any)) *LdbStorage { // 存储器，文档，自定义对象

	// 缓存有则取用
	cacheName := storeName + cmn.PathSeparator() + subPath
	cacheStore := getCacheStore(cacheName)
	if cacheStore != nil {
		return cacheStore
	}

	// 缓存无则锁后创建返回并存缓存
	mu.Lock()
	cacheStore = getCacheStore(cacheName) // 再次尝试取用缓存中存储器
	if cacheStore != nil {
		return cacheStore
	}

	store := new(LdbStorage)
	store.storeName = storeName
	store.subPath = subPath
	store.closing = false
	store.storeChan = make(chan any, storeChanBuffer) // 初始化管道，设定缓冲

	// 打开（在指定子目录中存放数据）
	dbPath := storePath + cmn.PathSeparator() + cacheName
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Println("打开LdbStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	store.currentCount = store.loadTotalCount() // 读取总件数
	mapStorage[cacheName] = store               // 缓存起来

	// 消费就绪
	go func() {
		for {
			data := <-store.storeChan
			wg.Done()
			if data == nil {
				close(store.storeChan) // 关闭通道
				break
			}
			if fnBeforeSave != nil {
				obj := fnBeforeSave(store, data) // 自定义预处理
				if fnSave != nil {
					store.currentCount++
					doc, obj := fnSave(store, obj) // 自定义保存处理
					if fnAfterSave != nil {
						fnAfterSave(store, doc, obj) // 自定义后处理
					}
					store.leveldb.Put(cmn.Uint32ToBytes(0), cmn.Uint32ToBytes(store.currentCount), nil) // 保存总件数
				}
			} else if fnSave != nil {
				store.currentCount++
				doc, obj := fnSave(store, data) // 自定义保存处理
				if fnAfterSave != nil {
					fnAfterSave(store, doc, obj) // 自定义后处理
				}
				store.leveldb.Put(cmn.Uint32ToBytes(0), cmn.Uint32ToBytes(store.currentCount), nil) // 保存总件数
			}

		}
	}()

	// 逐秒判断，若超时则自动关闭
	if timeout > 0 {
		go func() {
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				if time.Now().Unix()-store.lastTime > int64(timeout) {
					store.Close()
					ticker.Stop()
					break
				}
			}
		}()
	}

	mu.Unlock()
	log.Println("打开LdbStorage：", cacheName)
	return store
}

// 直接存入数据到leveldb
func (s *LdbStorage) Put(key []byte, value []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Put(key, value, nil)
}

// 直接从leveldb取数据
func (s *LdbStorage) Get(key []byte) ([]byte, error) {
	if s.closing {
		return nil, errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Get(key, nil)
}

// 关闭Storage
func (s *LdbStorage) Close() {
	if s.closing {
		return
	}
	s.closing = true
	wg.Wait()                     // 等待通道清空
	s.mu.Lock()                   // 锁
	wg.Add(1)                     // 通道消息计数
	s.storeChan <- nil            // 通道正在在阻塞等待接收，给个nil让它接收后关闭
	s.leveldb.Close()             // 走到这里时没有db操作了，可以关闭
	mapStorage[s.storeName] = nil // 设空，下回GetStorage时自动再创建
	s.mu.Unlock()                 // 解锁

	log.Println("关闭LdbStorage：", s.storeName)
}

// 添加数据，经通道后由自定义的保存函数处理保存
func (s *LdbStorage) Add(data any) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	if data == nil {
		return errors.New("ingore data nil") // 拒绝接收nil（nil作通道的关闭信号用，避免冲突）
	}

	s.lastTime = time.Now().Unix()
	s.storeChan <- data // 把文本发送到消息通道
	wg.Add(1)
	return nil
}

func (s *LdbStorage) loadTotalCount() uint32 {
	bytes, err := s.leveldb.Get(cmn.Uint32ToBytes(0), nil)
	if err != nil || bytes == nil {
		return 0
	}
	return cmn.BytesToUint32(bytes)
}

// 存储根目录
func StorePath() string {
	return storePath
}

// 总件数
func (s *LdbStorage) TotalCount() uint32 {
	return s.currentCount
}

// 存储目录名
func (s *LdbStorage) StoreName() string {
	return s.storeName
}

// 是否关闭中状态
func (s *LdbStorage) IsClose() bool {
	s.lastTime = time.Now().Unix()
	return s.closing
}
