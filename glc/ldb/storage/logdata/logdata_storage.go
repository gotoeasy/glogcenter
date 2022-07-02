/**
 * 日志存储器
 * 1）以通道控制接收日志，ID自动递增作为键有序保存
 * 2）优先响应保存日志，闲时创建关键词反向索引
 * 3）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 */
package logdata

import (
	"errors"
	"glc/cmn"
	"glc/ldb/conf"
	"glc/ldb/storage/indexword"
	"glc/ldb/sysmnt"
	"glc/ldb/tokenizer"
	"glc/onexit"
	"log"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// 存储结构体
type LogDataStorage struct {
	storeName    string             // 存储目录
	subPath      string             // 存储目录下的相对路径（存放数据）
	storeChan    chan *LogDataModel // 存储通道
	leveldb      *leveldb.DB        // leveldb
	currentCount uint32             // 当前件数
	lastTime     int64              // 最后一次访问时间
	closing      bool               // 是否关闭中状态
	mu           sync.Mutex         // 锁
	wg           sync.WaitGroup     // 计数
}

var ldbMu sync.Mutex
var mapStorage map[string](*LogDataStorage)

func init() {
	mapStorage = make(map[string](*LogDataStorage))
	onexit.RegisterExitHandle(onExit) // 优雅退出
}

func getCacheStore(cacheName string) *LogDataStorage {
	cacheStore := mapStorage[cacheName]
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore // 缓存中未关闭的存储对象
	}
	return nil
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func NewLogDataStorage(storeName string, subPath string) *LogDataStorage { // 存储器，文档，自定义对象

	// 缓存有则取用
	cacheName := storeName + cmn.PathSeparator() + subPath
	cacheStore := getCacheStore(cacheName)
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore
	}

	// 缓存无则锁后创建返回并存缓存
	ldbMu.Lock()                          // 上锁
	defer ldbMu.Unlock()                  // 解锁
	cacheStore = getCacheStore(cacheName) // 再次尝试取用缓存中存储器
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore
	}

	store := new(LogDataStorage)
	store.storeName = storeName
	store.subPath = subPath
	store.closing = false
	store.lastTime = time.Now().Unix()
	store.storeChan = make(chan *LogDataModel, conf.GetStoreChanLength()) // 初始化管道，设定缓冲

	dbPath := conf.GetStorageRoot() + cmn.PathSeparator() + cacheName
	db, err := leveldb.OpenFile(dbPath, nil) // 打开（在指定子目录中存放数据）
	if err != nil {
		log.Println("打开LogDataStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	store.currentCount = store.loadTotalCount() // 读取总件数
	mapStorage[cacheName] = store               // 缓存起来

	// 消费就绪
	go readyGo(store)

	// 逐秒判断，若闲置超时则自动关闭
	go autoCloseWhenMaxIdle(store)

	log.Println("打开LogDataStorage：", cacheName)
	return store
}

// 等待接收日志，优先响应保存日志，空时再生成索引
func readyGo(store *LogDataStorage) {
	for {
		select {
		case data := <-store.storeChan:
			store.wg.Done()
			// 优先响应保存日志
			if data == nil {
				if !store.IsClose() {
					close(store.storeChan) // 关闭通道
				}
				break
			}
			saveLogData(store, data) // 保存日志数据
		default:
			// 空时再生成索引，一次一条日志，有空则生成直到全部完成
			n := createInvertedIndex(store) // 生成反向索引

			// 索引生成完成后，等待接收保存日志
			if n < 1 {
				log.Println("空闲等待接收日志")
				data := <-store.storeChan // 没有索引可生成时，等待storeChan
				store.wg.Done()
				if data == nil {
					if !store.IsClose() {
						close(store.storeChan) // 关闭通道
					}
					break
				}
				saveLogData(store, data) // 保存日志数据
			}
		}
	}
}

func saveLogData(store *LogDataStorage, model *LogDataModel) {
	//store.wg.Done()
	store.currentCount++                              // ID递增
	doc := new(LogDataDocument)                       // 文档
	doc.Id = store.currentCount                       // 已递增好的值
	model.Id = cmn.Uint32ToString(store.currentCount) // 模型数据要转Json存，也得更新ID,ID用36进制字符串形式表示
	doc.Content = model.ToJson()                      // 转json作为内容(含Id)

	// 保存
	store.put(cmn.Uint32ToBytes(doc.Id), doc.ToBytes())                                 // 日志数据
	store.leveldb.Put(cmn.Uint32ToBytes(0), cmn.Uint32ToBytes(store.currentCount), nil) // 保存日志总件数
	log.Println("保存日志数据 ", doc.Id)
}

// 创建日志索引（一次建一条日志的索引）,没有可建索引时返回false
func createInvertedIndex(s *LogDataStorage) int {

	// 索引信息和日志数量相互比较，判断是否继续创建索引
	mntKey := "INDEX:" + s.StoreName()
	sysStorage := sysmnt.GetSysmntStorage(s.StoreName())
	sysmntData := sysStorage.GetSysmntData(mntKey)
	if s.TotalCount() == 0 || sysmntData.Count >= s.TotalCount() {
		return 0 // 没有新的日志需要建索引
	}

	sysmntData.Count++                               // 下一条要建索引的日志id
	docm, err := s.GetLogDataModel(sysmntData.Count) // 取出日志模型数据
	if err != nil {
		log.Println("取日志模型数据失败：", sysmntData.Count, err)
		return 2
	}

	// 整理关键词
	adds := docm.Keywords
	adds = append(adds, docm.Tags...)
	adds = append(adds, docm.Client, docm.Server, docm.System, docm.User)
	kws := tokenizer.CutForSearchEx(docm.Text, adds, docm.Sensitives) // 两数组参数的元素可以重复或空白，会被判断整理
	//	log.Println("GetLogDataModel=", docm.ToJson(), kws)

	// 每个关键词都创建反向索引
	for _, word := range kws {
		idxw := indexword.NewWordIndexStorage(s.StoreName(), word)
		idxw.Add(word, cmn.StringToUint32(docm.Id, 0)) // 日志ID加入索引
	}
	log.Println("创建日志索引：", cmn.StringToUint32(docm.Id, 0))

	// 保存当前创建了多少索引
	sysStorage.SetSysmntData(mntKey, sysmntData)

	return 1
}

func autoCloseWhenMaxIdle(store *LogDataStorage) {
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

// 直接存入数据到leveldb
func (s *LogDataStorage) put(key []byte, value []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务，调用方需自行重取连接等处理
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Put(key, value, nil)
}

// 直接从leveldb取数据
func (s *LogDataStorage) Get(key []byte) ([]byte, error) {
	if s.closing {
		return nil, errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Get(key, nil)
}

// 直接从leveldb取数据并转换为LogDataModel
func (s *LogDataStorage) GetLogDataModel(id uint32) (*LogDataModel, error) {
	if s.closing {
		return nil, errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	bts, err := s.leveldb.Get(cmn.Uint32ToBytes(id), nil)
	if err != nil {
		return nil, err
	}

	doc := new(LogDataDocument)
	doc.LoadBytes(bts)
	return doc.ToLogDataModel(), nil
}

// 添加数据，经通道后由自定义的保存函数处理保存
func (s *LogDataStorage) Add(model *LogDataModel) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	if model == nil {
		return errors.New("ingore data nil") // 拒绝接收nil（nil作通道的关闭信号用，避免冲突）
	}
	s.lastTime = time.Now().Unix()
	s.wg.Add(1)
	s.storeChan <- model // 把文本发送到消息通道
	return nil
}

// 关闭Storage
func (s *LogDataStorage) Close() {
	if s == nil || s.closing { // 优雅退出时可能会正好nil，判断一下优雅点
		return
	}

	s.mu.Lock()         // 锁
	defer s.mu.Unlock() // 解锁
	if s.closing {
		return
	}

	s.closing = true
	s.wg.Wait()                   // 等待通道清空
	s.wg.Add(1)                   // 通道消息计数
	s.storeChan <- nil            // 通道正在在阻塞等待接收，给个nil让它接收后关闭
	s.leveldb.Close()             // 走到这里时没有db操作了，可以关闭
	mapStorage[s.storeName] = nil // 设空，下回GetStorage时自动再创建

	log.Println("关闭LogDataStorage：", s.storeName+cmn.PathSeparator()+s.subPath)
}

func (s *LogDataStorage) loadTotalCount() uint32 {
	bytes, err := s.leveldb.Get(cmn.Uint32ToBytes(0), nil)
	if err != nil || bytes == nil {
		return 0
	}
	return cmn.BytesToUint32(bytes)
}

// 总件数
func (s *LogDataStorage) TotalCount() uint32 {
	return s.currentCount
}

// 存储目录名
func (s *LogDataStorage) StoreName() string {
	return s.storeName
}

// 是否关闭中状态
func (s *LogDataStorage) IsClose() bool {
	return s.closing
}

func onExit() {
	for k := range mapStorage {
		mapStorage[k].Close()
	}
	log.Println("退出LogDataStorage")
}
