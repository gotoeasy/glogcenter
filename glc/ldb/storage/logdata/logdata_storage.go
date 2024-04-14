/**
 * 日志存储器
 * 1）以通道控制接收日志，ID自动递增作为键有序保存
 * 2）优先响应保存日志，闲时创建关键词反向索引
 * 3）获取存储对象线程安全，带缓存无则创建有则直取，空闲超时自动关闭leveldb，再次获取时自动打开
 */
package logdata

import (
	"errors"
	"glc/conf"
	"glc/ldb/status"
	"glc/ldb/storage/indexword"
	"glc/ldb/sysmnt"
	"glc/ldb/tokenizer"
	"sync"
	"time"

	"github.com/gotoeasy/glang/cmn"
	"github.com/syndtr/goleveldb/leveldb"
)

// 存储结构体
type LogDataStorage struct {
	storeName         string             // 存储目录
	subPath           string             // 存储目录下的相对路径（存放数据）
	storeChan         chan *LogDataModel // 存储通道
	leveldb           *leveldb.DB        // leveldb
	currentCount      uint32             // 当前件数
	savedCurrentCount uint32             // 已保存的当前件数
	indexedCount      uint32             // 已创建的索引件数
	savedIndexedCount uint32             // 已保存的索引件数(定时保存indexedCount，存起来以便下次启动继续建索引)
	lastTime          int64              // 最后一次访问时间
	closing           bool               // 是否关闭中状态
	mu                sync.Mutex         // 锁
	wg                sync.WaitGroup     // 计数
}

var zeroUint32Bytes []byte = cmn.Uint32ToBytes(0)
var ldbMu sync.Mutex
var mapStorage map[string](*LogDataStorage)
var mapStorageMu sync.Mutex

func init() {
	mapStorage = make(map[string](*LogDataStorage))
	cmn.OnExit(onExit) // 优雅退出
}

func getCacheStore(cacheName string) *LogDataStorage {
	cacheStore := mapStorage[cacheName]
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore // 缓存中未关闭的存储对象
	}
	return nil
}

// 获取存储对象，线程安全（带缓存无则创建有则直取）
func NewLogDataStorage(storeName string) *LogDataStorage { // 存储器，文档，自定义对象

	subPath := "data"
	// 缓存有则取用
	cacheName := storeName + cmn.PathSeparator() + subPath
	cacheStore := getCacheStore(cacheName)
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore
	}

	// 缓存无则锁后创建返回并存缓存
	mapStorageMu.Lock()                   // 缓存map锁
	defer mapStorageMu.Unlock()           // 缓存map解锁
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
		cmn.Error("打开LogDataStorage失败：", dbPath)
		panic(err)
	}
	store.leveldb = db
	store.loadMetaData()                        // 初始化件数等信息
	status.UpdateStorageStatus(storeName, true) // 更新状态：当前日志仓打开
	mapStorage[cacheName] = store               // 缓存起来

	// 消费就绪
	go store.readyGo()

	// 定时判断保存总件数，避免每次保存以提高性能
	go store.readySaveMetaDate()

	// 逐秒判断，若闲置超时则自动关闭
	go store.autoCloseWhenMaxIdle()

	cmn.Info("打开LogDataStorage：", cacheName)
	return store
}

// 定时调用保存件数信息，避免每次都存levledb
func (s *LogDataStorage) readySaveMetaDate() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		<-ticker.C
		if s.IsClose() {
			ticker.Stop()
			break
		}
		s.saveMetaData()
	}
}

// 等待接收日志，优先响应保存日志，空时再生成索引
func (s *LogDataStorage) readyGo() {
	for {
		select {
		case data := <-s.storeChan:
			s.wg.Done()
			// 优先响应保存日志
			if data == nil {
				if !s.IsClose() {
					close(s.storeChan) // 关闭通道
				}
				break
			}
			s.saveLogData(data) // 保存日志数据
		default:
			// 空时再生成索引，一次一条日志，有空则生成直到全部完成
			n := s.createInvertedIndex() // 生成反向索引

			// 索引生成完成后，等待接收保存日志
			if n < 1 {
				cmn.Info("空闲等待接收日志")
				data := <-s.storeChan // 没有索引可生成时，等待storeChan
				s.wg.Done()
				if data == nil {
					if !s.IsClose() {
						close(s.storeChan) // 关闭通道
					}
					break
				}
				s.saveLogData(data) // 保存日志数据
			}
		}
	}
}

func (s *LogDataStorage) saveLogData(model *LogDataModel) {
	if model.Text == "" {
		return // Text没有内容的话就不保存了
	}

	//store.wg.Done()
	s.currentCount++                              // ID递增
	doc := new(LogDataDocument)                   // 文档
	doc.Id = s.currentCount                       // 已递增好的值
	model.Id = cmn.Uint32ToString(s.currentCount) // 模型数据要转Json存，也得更新ID,ID用36进制字符串形式表示
	doc.Content = model.ToJson()                  // 转json作为内容(含Id)

	// 保存
	s.put(cmn.Uint32ToBytes(doc.Id), doc.ToBytes()) // 日志数据
	cmn.Debug("保存日志数据 ", doc.Id)
}

// 创建日志索引（一次建一条日志的索引）,没有可建索引时返回0
func (s *LogDataStorage) createInvertedIndex() int {

	// 索引信息和日志数量相互比较，判断是否继续创建索引
	if s.TotalCount() == 0 || s.indexedCount >= s.TotalCount() {
		return 0 // 没有新的日志需要建索引
	}

	s.indexedCount++                               // 下一条要建索引的日志id
	docm, err := s.GetLogDataModel(s.indexedCount) // 取出日志模型数据
	if err != nil {
		cmn.Error("取日志模型数据失败：", s.indexedCount, err)
		return 2
	}

	// 整理生成关键词
	var adds []string
	if docm.System != "" {
		adds = append(adds, "~"+cmn.ToLower(docm.System))
	}
	if docm.LogLevel != "" {
		adds = append(adds, "!"+cmn.ToLower(docm.LogLevel))
	}
	if docm.User != "" {
		adds = append(adds, "@"+cmn.ToLower(docm.User))
	}

	tgtStr := docm.System + " " + docm.ServerName + " " + docm.ServerIp + " " + docm.ClientIp + " " + docm.TraceId + " " + docm.LogLevel + " " + docm.User
	if docm.Detail != "" && conf.IsMulitLineSearch() {
		tgtStr = tgtStr + " " + docm.Detail // 支持日志列全部行作为索引检索对象
	} else {
		tgtStr = tgtStr + " " + docm.Text // 日志列仅第一行作为索引检索对象
	}
	kws := tokenizer.CutForSearchEx(tgtStr, adds, nil) // 两数组参数的元素可以重复或空白，会被判断整理

	// 每个关键词都创建反向索引
	for _, word := range kws {
		idxw := indexword.NewWordIndexStorage(s.StoreName())
		idxw.Add(word, cmn.StringToUint32(docm.Id, 0)) // 日志ID加入索引
	}
	// cmn.Debug("创建日志索引：", cmn.StringToUint32(docm.Id, 0))

	return 1
}

func (s *LogDataStorage) autoCloseWhenMaxIdle() {
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

	mapStorageMu.Lock()         // 缓存map锁
	defer mapStorageMu.Unlock() // 缓存map解锁
	s.mu.Lock()                 // 锁
	defer s.mu.Unlock()         // 解锁
	if s.closing {
		return
	}

	s.closing = true
	s.wg.Wait()                                    // 等待通道清空
	s.saveMetaData()                               // 保存件数等元信息
	s.wg.Add(1)                                    // 通道消息计数
	s.storeChan <- nil                             // 通道正在在阻塞等待接收，给个nil让它接收后关闭
	s.leveldb.Close()                              // 走到这里时没有db操作了，可以关闭
	mapStorage[s.storeName] = nil                  // 设空，下回GetStorage时自动再创建
	status.UpdateStorageStatus(s.storeName, false) // 更新状态：当前日志仓关闭

	cmn.Info("关闭LogDataStorage：", s.storeName+cmn.PathSeparator()+s.subPath)
}

func (s *LogDataStorage) loadMetaData() {

	// 初始化：当前日志件数
	bytes, err := s.leveldb.Get(zeroUint32Bytes, nil)
	if err != nil || bytes == nil {
		s.currentCount = 0
	} else {
		s.currentCount = cmn.BytesToUint32(bytes)
		s.savedCurrentCount = s.currentCount
	}

	// 初始化：已建索引件数
	idxw := indexword.NewWordIndexStorage(s.StoreName())
	s.indexedCount = idxw.GetIndexedCount()
	s.savedIndexedCount = s.indexedCount

	// 检查更新系统管理存储器中的日志总件数
	sysmntStore := sysmnt.NewSysmntStorage() // 系统管理存储器
	if sysmntStore.GetStorageDataCount(s.storeName) != s.currentCount {
		sysmntStore.SetStorageDataCount(s.storeName, s.currentCount)
	}
	if sysmntStore.GetStorageIndexCount(s.storeName) != s.currentCount {
		sysmntStore.SetStorageIndexCount(s.storeName, s.indexedCount)
	}
}

func (s *LogDataStorage) saveMetaData() {

	if s.savedCurrentCount < s.currentCount {
		s.savedCurrentCount = s.currentCount
		s.leveldb.Put(zeroUint32Bytes, cmn.Uint32ToBytes(s.savedCurrentCount), nil) // 保存日志总件数
		sysmntStore := sysmnt.NewSysmntStorage()                                    // 系统管理存储器
		sysmntStore.SetStorageDataCount(s.storeName, s.savedCurrentCount)           // 保存日志总件数
		cmn.Info("保存LogDataStorage件数:", s.savedCurrentCount)
	}

	if s.savedIndexedCount < s.indexedCount {
		s.savedIndexedCount = s.indexedCount
		idxw := indexword.NewWordIndexStorage(s.StoreName())
		idxw.SaveIndexedCount(s.savedIndexedCount)                         // 保存索引总件数
		sysmntStore := sysmnt.NewSysmntStorage()                           // 系统管理存储器
		sysmntStore.SetStorageIndexCount(s.storeName, s.savedCurrentCount) // 保存索引总件数
		cmn.Info("保存LogDataStorage已建索引件数:", s.savedIndexedCount)
	}
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
		s := mapStorage[k]
		if s != nil {
			s.Close()
		}
	}
	cmn.Info("退出LogDataStorage")
}
