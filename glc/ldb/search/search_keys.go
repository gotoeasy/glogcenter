/**
 * 日志检索
 * 1）无关键词时直接全量检索
 * 2）有关键词时检索索引
 * 3）支持指定相对ID及方向进行前后翻页检索
 */
package search

import (
	"glc/ldb/storage"
	"glc/ldb/storage/indexdoc"
	"glc/ldb/storage/indexword"
	"glc/ldb/storage/logdata"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

type SearchCondition struct {
	SearchKey        string   // 输入的检索文本
	StoreName        string   // 输入的日志仓条件
	System           string   // 输入的系统名条件
	DatetimeFrom     string   // 输入的日期范围（from）条件
	DatetimeTo       string   // 输入的日期范围（to）条件
	Loglevel         string   // 输入的日志级别（单选条件）条件【内部会过滤修改】
	Loglevels        []string // 输入的日志级别（多选条件）条件【内部会过滤修改】
	User             string   // 输入的用户条件
	CurrentStoreName string   // 隐藏条件，当前日志文档ID所属的日志仓
	CurrentId        uint32   // 隐藏条件，当前日志文档ID
	Forward          bool     // 隐藏条件，是否向前检索（玩下滚动查询）
	OldNearId        uint32   // 隐藏条件，相邻检索时的旧ID
	NewNearId        uint32   // 隐藏条件，相邻检索时的新ID
	NearStoreName    string   // 隐藏条件，相邻检索时新ID对应的日志仓
	Kws              []string // 【内部用】解析条件所得的检索关键词，非直接输入的检索文本
	SearchSize       int      // 【内部用】需要查询多少件（跨仓检索时可能多次检索，中间会内部调整）
	Systems          []string // 【内部用】有权限的系统名（一定有值，第一个元素是“*”时表示全部）
	OrgSystem        string   // 【内部用】保存输入的系统名条件
	OrgSystems       []string // 【内部用】保存有权限的系统名（一定有值，第一个元素是“*”时表示全部）
}

type SearchResult struct {
	Total         string                  `json:"total,omitempty"`         // 日志总量件数（用10进制字符串形式以避免出现科学计数法）
	Count         string                  `json:"count,omitempty"`         // 当前条件最多匹配件数（用10进制字符串形式以避免出现科学计数法）
	PageSize      string                  `json:"pagesize,omitempty"`      // 每次检索件数
	Data          []*logdata.LogDataModel `json:"data,omitempty"`          // 检索结果数据（日志文档数组）
	LastStoreName string                  `json:"laststorename,omitempty"` // 当前检索结果中，最后一条（最久远）日志所在日志仓
	TimeMessage   string                  `json:"timemessage,omitempty"`   // 查询耗时的文本消息表示，如：耗时30毫秒
}

type WidxStorage struct {
	word           string
	idxdocStorage  *indexdoc.DocIndexStorage
	idxwordStorage *indexword.WordIndexStorage
}

// 多关键词时计算关键词索引交集
func SearchWordIndex(storeName string, cond *SearchCondition) *SearchResult {
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据

	// 时间条件范围判断，默认全部，有检索条件时调整范围
	maxDocumentId := storeLogData.TotalCount()  // 时间范围条件内的最大文档ID
	minDocumentId := cmn.StringToUint32("1", 1) // 时间范围条件内的最小文档ID
	if !cmn.IsBlank(cond.DatetimeFrom) {
		minDocumentId = findMinDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, cond.DatetimeFrom) // 时间范围条件内的最小文档ID，找不到时返回0
		if minDocumentId == 0 {
			// 简单判断，无匹配时直接返回
			var rs = new(SearchResult)
			rs.Total = cmn.Uint32ToString(storeLogData.TotalCount())
			rs.Count = "0"
			return rs
		}
	}
	if !cmn.IsBlank(cond.DatetimeTo) {
		maxDocumentId = findMaxDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, cond.DatetimeTo) // 时间范围条件内的最大文档ID，找不到时返回0
		if maxDocumentId == 0 {
			// 简单判断，无匹配时直接返回
			var rs = new(SearchResult)
			rs.Total = cmn.Uint32ToString(storeLogData.TotalCount())
			rs.Count = "0"
			return rs
		}
	}
	if minDocumentId > maxDocumentId {
		// 简单判断，无匹配时直接返回
		var rs = new(SearchResult)
		rs.Total = cmn.Uint32ToString(storeLogData.TotalCount())
		rs.Count = "0"
		return rs
	}

	// 汇总索引进行关联查找
	var widxs []*WidxStorage
	for _, word := range cond.Kws {
		widxStorage := &WidxStorage{
			word:           word,
			idxdocStorage:  indexdoc.NewDocIndexStorage(storeName),
			idxwordStorage: indexword.NewWordIndexStorage(storeName),
		}
		widxs = append(widxs, widxStorage)
	}
	return findSame(cond, minDocumentId, maxDocumentId, storeLogData, widxs...)
}

// 按ID查取日志
func GetLogDataModelById(storeName string, id uint32) *logdata.LogDataModel {
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	ldm := storeLogData.GetLogDataModel(id)
	ldm.StoreName = storeName // 日志仓名称未序列化，前端要使用，这里补足赋值
	return ldm
}

// 无关键词时走全量检索
func SearchLogData(storeName string, cond *SearchCondition) *SearchResult {
	allloglevels := cmn.Join(cond.Loglevels, ",")              // 合并多选的级别条件
	noLogLevels := cmn.IsBlank(allloglevels)                   // 无多选条件
	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	totalCount := storeLogData.TotalCount()                    // 总件数
	rs.Total = cmn.Uint32ToString(totalCount)                  // 返回的日志总量件数，用10进制字符串形式以避免出现科学计数法
	rs.Count = cmn.Uint32ToString(totalCount)                  // 当前条件最多匹配件数
	rsCnt := 0                                                 // 已查到的件数

	if totalCount == 0 || cond.SearchSize == 0 {
		return rs // 无数据或不需要检索数据
	}

	// 时间条件范围判断，默认全部，有检索条件时调整范围
	maxDocumentId := totalCount                 // 时间范围条件内的最大文档ID
	minDocumentId := cmn.StringToUint32("1", 1) // 时间范围条件内的最小文档ID
	hasMin := !cmn.IsBlank(cond.DatetimeFrom)
	hasMax := !cmn.IsBlank(cond.DatetimeTo)
	if hasMin {
		minDocumentId = findMinDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, cond.DatetimeFrom) // 时间范围条件内的最小文档ID
		if minDocumentId == 0 {
			// 简单判断，无匹配时直接返回
			var rs = new(SearchResult)
			rs.Total = "0"
			rs.Count = "0"
			return rs
		}
	}
	if hasMax {
		maxDocumentId = findMaxDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, cond.DatetimeTo) // 时间范围条件内的最大文档ID
		if maxDocumentId == 0 {
			// 简单判断，无匹配时直接返回
			var rs = new(SearchResult)
			rs.Total = "0"
			rs.Count = "0"
			return rs
		}
	}
	if hasMax || hasMin {
		if minDocumentId > maxDocumentId {
			// 简单判断，无匹配时直接返回
			rs.Count = "0"
			return rs
		}
		rs.Count = cmn.Uint32ToString(maxDocumentId - minDocumentId + 1) // 估算的最大匹配件数
	}

	// 开始检索
	if cond.CurrentId == 0 {
		// 第一页
		var min, max uint32
		max = totalCount

		if max > maxDocumentId {
			max = maxDocumentId // 最大不超出时间范围限制内的最大文档ID
		}

		if min < minDocumentId {
			min = minDocumentId // 最小不超出时间范围限制内的最小文档ID
		}

		for i := max; i >= min; i-- {
			if cond.System == "" && cond.Systems[0] != "*" {
				// 权限范围系统内作过滤检查
				idxdocStorage := indexdoc.NewDocIndexStorage(storeName) // 判断系统权限用
				has := false
				for j, max2 := 0, len(cond.Systems); j < max2; j++ {
					if (idxdocStorage.GetWordDocSeq(cond.Systems[j], i)) > 0 {
						has = true
						break
					}
				}
				if !has {
					continue // 权限范围内的系统内找不到时，跳过该日志
				}
			}

			if cond.Loglevel == "" && len(cond.Loglevels) > 0 {
				// 日志级别范围内作过滤检查
				idxdocStorage := indexdoc.NewDocIndexStorage(storeName)
				has := false
				for j, max2 := 0, len(cond.Loglevels); j < max2; j++ {
					if (idxdocStorage.GetWordDocSeq("!"+cond.Loglevels[j], i)) > 0 {
						has = true
						break
					}
				}
				if !has {
					continue // 权限范围内的系统内找不到时，跳过该日志
				}
			}

			md := storeLogData.GetLogDataDocument(i).ToLogDataModel()
			md.StoreName = storeLogData.GetStoreName() // 日志仓名称未序列化，前端要使用，这里补足赋值
			if noLogLevels || cmn.ContainsIngoreCase(allloglevels, md.LogLevel) {
				rs.Data = append(rs.Data, md)
				rsCnt++
				if rsCnt >= cond.SearchSize {
					break
				}
			}
		}
	} else if cond.Forward {
		// 后一页（向下滚动触发检索）
		if cond.CurrentId > 1 {
			var min, max uint32
			if cond.CurrentId > totalCount {
				max = totalCount
			} else {
				max = cond.CurrentId - 1
			}

			if max > maxDocumentId {
				max = maxDocumentId // 最大不超出时间范围限制内的最大文档ID
			}

			if min < minDocumentId {
				min = minDocumentId // 最小不超出时间范围限制内的最小文档ID
			}

			for i := max; i >= min; i-- {
				if cond.System == "" && cond.Systems[0] != "*" {
					// 权限范围系统内作过滤检查
					idxdocStorage := indexdoc.NewDocIndexStorage(storeName) // 判断系统权限用
					has := false
					for j, max2 := 0, len(cond.Systems); j < max2; j++ {
						if (idxdocStorage.GetWordDocSeq(cond.Systems[j], i)) > 0 {
							has = true
							break
						}
					}
					if !has {
						continue // 权限范围内的系统内找不到时，跳过该日志
					}
				}

				if cond.Loglevel == "" && len(cond.Loglevels) > 0 {
					// 日志级别范围内作过滤检查
					idxdocStorage := indexdoc.NewDocIndexStorage(storeName)
					has := false
					for j, max2 := 0, len(cond.Loglevels); j < max2; j++ {
						if (idxdocStorage.GetWordDocSeq("!"+cond.Loglevels[j], i)) > 0 {
							has = true
							break
						}
					}
					if !has {
						continue // 权限范围内的系统内找不到时，跳过该日志
					}
				}

				md := storeLogData.GetLogDataDocument(i).ToLogDataModel()
				md.StoreName = storeLogData.GetStoreName() // 日志仓名称未序列化，前端要使用，这里补足赋值
				if noLogLevels || cmn.ContainsIngoreCase(allloglevels, md.LogLevel) {
					rs.Data = append(rs.Data, md)
					rsCnt++
					if rsCnt >= cond.SearchSize {
						break
					}
				}
			}
		}
	} else {
		// 前一页（向上滚动触发检索）【暂未使用】
		if totalCount > cond.CurrentId {
			var min, max uint32
			min = cond.CurrentId + 1

			if min < minDocumentId {
				min = minDocumentId // 最小不超出时间范围限制内的最小文档ID
			}

			if max > maxDocumentId {
				max = maxDocumentId // 最大不超出时间范围限制内的最大文档ID
			}

			for i := max; i >= min; i-- {
				if cond.System == "" && cond.Systems[0] != "*" {
					// 权限范围系统内作过滤检查
					idxdocStorage := indexdoc.NewDocIndexStorage(storeName) // 判断系统权限用
					has := false
					for j, max2 := 0, len(cond.Systems); j < max2; j++ {
						if (idxdocStorage.GetWordDocSeq(cond.Systems[j], i)) > 0 {
							has = true
							break
						}
					}
					if !has {
						continue // 权限范围内的系统内找不到时，跳过该日志
					}
				}

				if cond.Loglevel == "" && len(cond.Loglevels) > 0 {
					// 日志级别范围内作过滤检查
					idxdocStorage := indexdoc.NewDocIndexStorage(storeName)
					has := false
					for j, max2 := 0, len(cond.Loglevels); j < max2; j++ {
						if (idxdocStorage.GetWordDocSeq("!"+cond.Loglevels[j], i)) > 0 {
							has = true
							break
						}
					}
					if !has {
						continue // 权限范围内的系统内找不到时，跳过该日志
					}
				}

				md := storeLogData.GetLogDataDocument(i).ToLogDataModel()
				md.StoreName = storeLogData.GetStoreName() // 日志仓名称未序列化，前端要使用，这里补足赋值
				if noLogLevels || cmn.ContainsIngoreCase(allloglevels, md.LogLevel) {
					rs.Data = append(rs.Data, md)
					rsCnt++
					if rsCnt >= cond.SearchSize {
						break
					}
				}

			}
		}
	}

	return rs
}

// 参数widxs长度要求大于1，currentId不传就是查第一页
func findSame(cond *SearchCondition, minDocumentId uint32, maxDocumentId uint32, storeLogData *storage.LogDataStorageHandle, widxs ...*WidxStorage) *SearchResult {
	allloglevels := cmn.Join(cond.Loglevels, ",")            // 合并多选的级别条件
	noLogLevels := cmn.IsBlank(allloglevels)                 // 无多选条件
	var rs = new(SearchResult)                               // 查询结果
	rs.Total = cmn.Uint32ToString(storeLogData.TotalCount()) // 日志总量件数

	// 选个最短的索引
	cnt := len(widxs)
	minIdx := widxs[0]
	minCount := minIdx.idxwordStorage.GetTotalCount(minIdx.word)
	for i := 1; i < cnt; i++ {
		ctmp := widxs[i].idxwordStorage.GetTotalCount(widxs[i].word)
		if ctmp < minCount {
			minCount = ctmp
			minIdx = widxs[i]
		}
	}
	if minCount > maxDocumentId-minDocumentId+1 {
		minCount = maxDocumentId - minDocumentId + 1 // 最多匹配件数估算，不会超出时间条件范围，两者取其小
	}
	rs.Count = cmn.Uint32ToString(minCount) // 当前条件最多匹配件数

	if cond.SearchSize <= 0 {
		return rs // 只查关联件数，不查数据
	}

	// 简单检查排除没结果的情景
	totalCount := minIdx.idxwordStorage.GetTotalCount(minIdx.word)
	if totalCount == 0 || (cond.NewNearId == 0 && totalCount == 1 && cond.CurrentId > 0) {
		return rs // 索引件数0、或只有1条又还要跳过，都是找不到
	}

	// 找匹配位置并排除没结果的情景
	var pos uint32
	if cond.Forward {
		pos = totalCount // 向后检索更加旧的日志，游标从大开始递降
		if cond.NewNearId > 0 {
			// 相邻检索
			nearId := cond.NewNearId - 1                                  // 查找比定位日志小的旧日志
			pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, nearId) // 有相对文档ID时找相对位置
			for pos == 0 && nearId >= minDocumentId {
				nearId--
				pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, nearId)
			}
			if pos == 0 {
				return rs // 找不到
			}
			pos++ // 找到，后续会减1排除本条找到的日志，但相邻检索时本条日志是需要的，所以加1处理
		} else {
			// 普通检索
			if cond.CurrentId > 0 {
				pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, cond.CurrentId) // 有相对文档ID时找相对位置
				if pos <= 1 {
					return rs // 找不到、或最后条还要向后
				}
			}
		}
	} else {
		pos = 1 // 向前检索更加新的日志，游标从小开始递增
		if cond.NewNearId > 0 {
			// 相邻检索
			nearId := cond.NewNearId + 1                                  // 查找比定位日志大的新日志
			pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, nearId) // 有相对文档ID时找相对位置
			for pos == 0 && nearId <= maxDocumentId {
				nearId++
				pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, nearId)
			}
			if pos == 0 {
				return rs // 找不到
			}
			pos-- // 找到，后续会加1排除本条找到的日志，但相邻检索时本条日志是需要的，所以减1处理
		} else {
			// 普通检索
			if cond.CurrentId > 0 {
				pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, cond.CurrentId) // 有相对文档ID时找相对位置
				if pos == 0 || pos == totalCount {
					return rs //  找不到、或最前条还要向前
				}
			}
		}
	}

	// 位置就绪
	var rsCnt int = 0
	var flg bool
	if cond.CurrentId == 0 || (cond.CurrentId > 0 && cond.Forward) {
		// 无相对文档ID、或有且是后一页方向
		if cond.CurrentId > 0 {
			pos-- //  相对文档ID有的话才顺移
		}

		tmpMinPos := pos
		tmpMinIdx := minIdx
		for i := tmpMinPos; i > 0; {
			// 取值
			docId := minIdx.idxwordStorage.GetDocId(minIdx.word, i)
			if docId >= minDocumentId && docId <= maxDocumentId {
				// 在时间范围条件内时，继续查找比较
				flg = true
				for n := 0; n < cnt; n++ {
					if widxs[n] == minIdx {
						continue // 跳过比较自己
					}

					seq := widxs[n].idxdocStorage.GetWordDocSeq(widxs[n].word, docId)
					if seq == 0 {
						flg = false // 没找到
						break
					}
					if seq < tmpMinPos {
						tmpMinPos = seq
						tmpMinIdx = widxs[n] // 当前最短索引，存起来下回比较用
					}
				}

				if flg && cond.System == "" && cond.Systems[0] != "*" {
					// 权限范围系统内作过滤检查
					idxdocStorage := indexdoc.NewDocIndexStorage(storeLogData.GetStoreName()) // 判断系统权限用
					has := false
					for j, max2 := 0, len(cond.Systems); j < max2; j++ {
						if (idxdocStorage.GetWordDocSeq(cond.Systems[j], docId)) > 0 {
							has = true
							break
						}
					}
					if !has {
						flg = false // 没找到
					}
				}

				if flg && cond.Loglevel == "" && len(cond.Loglevels) > 0 {
					// 日志级别范围内作过滤检查
					idxdocStorage := indexdoc.NewDocIndexStorage(storeLogData.GetStoreName())
					has := false
					for j, max2 := 0, len(cond.Loglevels); j < max2; j++ {
						if (idxdocStorage.GetWordDocSeq("!"+cond.Loglevels[j], docId)) > 0 {
							has = true
							break
						}
					}
					if !has {
						flg = false // 没找到
					}
				}

				// 找到则加入结果
				if flg {
					md := storeLogData.GetLogDataModel(docId)
					md.StoreName = storeLogData.GetStoreName() // 日志仓名称未序列化，前端要使用，这里补足赋值
					if noLogLevels || cmn.ContainsIngoreCase(allloglevels, md.LogLevel) {
						rsCnt++
						rs.Data = append(rs.Data, md)
						if rsCnt >= cond.SearchSize {
							break // 最多找一页
						}
					}
				}
			}

			minIdx = tmpMinIdx // 当前最短索引
			i = tmpMinPos - 1  // 当前最短索引的后一个位置
			tmpMinPos--        // 当前最短索引可能不变，得正常减1，若变化则会被覆盖没有关系
		}
	} else {
		// 有相对文档ID且是前一页方向
		pos++
		var ary []*logdata.LogDataModel
		total := storeLogData.TotalCount() // 当前日志最大件数
		for i := pos; i <= total; i++ {
			// 取值
			docId := minIdx.idxwordStorage.GetDocId(minIdx.word, i)

			if docId < minDocumentId || docId > maxDocumentId {
				continue // 不在时间范围条件内，不匹配，跳过
			}

			// 比较
			flg = true
			for n := 0; n < cnt; n++ {
				if widxs[n] == minIdx {
					continue // 跳过比较自己
				}
				if widxs[n].idxdocStorage.GetWordDocSeq(widxs[n].word, docId) == 0 {
					flg = false // 没找到
					break
				}
			}

			if flg && cond.System == "" && cond.Systems[0] != "*" {
				// 权限范围系统内作过滤检查
				idxdocStorage := indexdoc.NewDocIndexStorage(storeLogData.GetStoreName()) // 判断系统权限用
				has := false
				for j, max2 := 0, len(cond.Systems); j < max2; j++ {
					if (idxdocStorage.GetWordDocSeq(cond.Systems[j], docId)) > 0 {
						has = true
						break
					}
				}
				if !has {
					flg = false // 没找到
				}
			}

			if flg && cond.Loglevel == "" && len(cond.Loglevels) > 0 {
				// 日志级别范围内作过滤检查
				idxdocStorage := indexdoc.NewDocIndexStorage(storeLogData.GetStoreName())
				has := false
				for j, max2 := 0, len(cond.Loglevels); j < max2; j++ {
					if (idxdocStorage.GetWordDocSeq("!"+cond.Loglevels[j], docId)) > 0 {
						has = true
						break
					}
				}
				if !has {
					flg = false // 没找到
				}
			}

			// 找到则加入结果
			if flg {
				md := storeLogData.GetLogDataModel(docId)
				md.StoreName = storeLogData.GetStoreName() // 日志仓名称未序列化，前端要使用，这里补足赋值
				if noLogLevels || cmn.ContainsIngoreCase(allloglevels, md.LogLevel) {
					rsCnt++
					ary = append(ary, md)
					if rsCnt >= cond.SearchSize {
						break // 最多找一页
					}
				}
			}
		}

		// 倒序放入结果
		for i := len(ary) - 1; i >= 0; i-- {
			rs.Data = append(rs.Data, ary[i])
		}
	}

	return rs
}

// 查找满足最小时间范围的最小文档id
func findMinDocumentIdByDatetime(storeLogData *storage.LogDataStorageHandle, uiMin uint32, uiMax uint32, minDatetime string) uint32 {
	if strings.Compare(minDatetime+".000", storeLogData.GetLogDataDocument(uiMin).ToLogDataModel().Date) <= 0 {
		return uiMin // 边界外输入条件常发生，特殊照顾确认边界，一定程度提高性能
	}

	rs := cmn.StringToUint32("0", 0)
	min := uiMin + 1 // 参数的最小已检查，跳过
	max := uiMax
	for min <= max {
		left, rigth, flg, target := findGE(storeLogData, min, max, minDatetime)
		min = left
		max = rigth
		if flg {
			rs = target
		}
	}
	return rs
}

func findGE(storeLogData *storage.LogDataStorageHandle, min uint32, max uint32, minDatetime string) (uint32, uint32, bool, uint32) {
	middle := (min + max) / 2
	if strings.Compare(minDatetime+".000", storeLogData.GetLogDataDocument(middle).ToLogDataModel().Date) <= 0 {
		return min, middle - 1, true, middle // 能匹配（middle的日时>=minDatetime），但不一定是最小匹配，继续返回下次待查找的范围
	}
	return middle + 1, max, false, 0 // 不匹配（middle的日时<minDatetime），返回下次待查找的范围
}

// 查找满足最大时间范围的最大文档id
func findMaxDocumentIdByDatetime(storeLogData *storage.LogDataStorageHandle, uiMin uint32, uiMax uint32, maxDatetime string) uint32 {
	if strings.Compare(storeLogData.GetLogDataDocument(uiMax).ToLogDataModel().Date, maxDatetime+".999") <= 0 {
		return uiMax // 能匹配（maxDatetime>=middle的日时），但不一定是最小匹配，继续返回下次待查找的范围
	}

	rs := cmn.StringToUint32("0", 0)
	min := uiMin
	max := uiMax - 1 // 参数的最大已检查，跳过
	for min <= max {
		left, rigth, flg, target := findLE(storeLogData, min, max, maxDatetime)
		min = left
		max = rigth
		if flg {
			rs = target
		}
	}
	return rs
}

func findLE(storeLogData *storage.LogDataStorageHandle, min uint32, max uint32, maxDatetime string) (uint32, uint32, bool, uint32) {
	middle := (min + max) / 2
	if strings.Compare(storeLogData.GetLogDataDocument(middle).ToLogDataModel().Date, maxDatetime+".999") <= 0 {
		return middle + 1, max, true, middle // 能匹配（maxDatetime>=middle的日时），但不一定是最小匹配，继续返回下次待查找的范围
	}
	return min, middle - 1, false, 0 // 不匹配（maxDatetime<middle的日时），返回下次待查找的范围
}
