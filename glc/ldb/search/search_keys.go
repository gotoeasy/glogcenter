/**
 * 日志检索
 * 1）无关键词时直接全量检索
 * 2）有关键词时检索索引
 * 3）支持指定相对ID及方向进行前后翻页检索
 */
package search

import (
	"glc/conf"
	"glc/ldb/storage"
	"glc/ldb/storage/indexdoc"
	"glc/ldb/storage/indexword"
	"glc/ldb/storage/logdata"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

type SearchResult struct {
	Total    string                  `json:"total,omitempty"`    // 日志总量件数（用10进制字符串形式以避免出现科学计数法）
	Count    string                  `json:"count,omitempty"`    // 当前条件最多匹配件数（用10进制字符串形式以避免出现科学计数法）
	PageSize string                  `json:"pagesize,omitempty"` // 每次检索件数
	Data     []*logdata.LogDataModel `json:"data,omitempty"`     // 检索结果数据（日志文档数组）
}

type WidxStorage struct {
	word           string
	idxdocStorage  *indexdoc.DocIndexStorage
	idxwordStorage *indexword.WordIndexStorage
}

// 多关键词时计算关键词索引交集
func SearchWordIndex(storeName string, kws []string, currentDocId uint32, forward bool, minDatetime string, maxDatetime string) *SearchResult {
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据

	// 时间条件范围判断，默认全部，有检索条件时调整范围
	maxDocumentId := storeLogData.TotalCount()  // 时间范围条件内的最大文档ID
	minDocumentId := cmn.StringToUint32("1", 1) // 时间范围条件内的最小文档ID
	if !cmn.IsBlank(minDatetime) {
		minDocumentId = findMinDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, minDatetime) // 时间范围条件内的最小文档ID，找不到时返回0
		if minDocumentId == 0 {
			// 简单判断，无匹配时直接返回
			var rs = new(SearchResult)
			rs.Total = cmn.Uint32ToString(storeLogData.TotalCount())
			rs.Count = "0"
			return rs
		}
	}
	if !cmn.IsBlank(maxDatetime) {
		maxDocumentId = findMaxDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, maxDatetime) // 时间范围条件内的最大文档ID，找不到时返回0
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
	for _, word := range kws {
		widxStorage := &WidxStorage{
			word:           word,
			idxdocStorage:  indexdoc.NewDocIndexStorage(storeName),
			idxwordStorage: indexword.NewWordIndexStorage(storeName),
		}
		widxs = append(widxs, widxStorage)
	}
	return findSame(currentDocId, forward, minDocumentId, maxDocumentId, storeLogData, widxs...)
}

// 无关键词时走全量检索
func SearchLogData(storeName string, currentDocId uint32, forward bool, minDatetime string, maxDatetime string) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	totalCount := storeLogData.TotalCount()                    // 总件数
	rs.Total = cmn.Uint32ToString(totalCount)                  // 返回的日志总量件数，用10进制字符串形式以避免出现科学计数法
	rs.Count = cmn.Uint32ToString(totalCount)                  // 当前条件最多匹配件数

	if totalCount == 0 {
		return rs
	}

	// 时间条件范围判断，默认全部，有检索条件时调整范围
	maxDocumentId := totalCount                 // 时间范围条件内的最大文档ID
	minDocumentId := cmn.StringToUint32("1", 1) // 时间范围条件内的最小文档ID
	hasMin := !cmn.IsBlank(minDatetime)
	hasMax := !cmn.IsBlank(maxDatetime)
	if hasMin {
		minDocumentId = findMinDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, minDatetime) // 时间范围条件内的最小文档ID
		if minDocumentId == 0 {
			// 简单判断，无匹配时直接返回
			var rs = new(SearchResult)
			rs.Total = "0"
			rs.Count = "0"
			return rs
		}
	}
	if hasMax {
		maxDocumentId = findMaxDocumentIdByDatetime(storeLogData, minDocumentId, maxDocumentId, maxDatetime) // 时间范围条件内的最大文档ID
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
	if currentDocId == 0 {
		// 第一页
		var min, max uint32
		max = totalCount

		if max > maxDocumentId {
			max = maxDocumentId // 最大不超出时间范围限制内的最大文档ID
		}

		if max > uint32(conf.GetPageSize()) {
			min = max - uint32(conf.GetPageSize()) + 1
		} else {
			min = 1
		}

		if min < minDocumentId {
			min = minDocumentId // 最小不超出时间范围限制内的最小文档ID
		}

		for i := max; i >= min; i-- {
			rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel()) // 件数等同日志文档ID
		}
	} else if forward {
		// 后一页
		if currentDocId > 1 {
			var min, max uint32
			if currentDocId > totalCount {
				max = totalCount
			} else {
				max = currentDocId - 1
			}

			if max > maxDocumentId {
				max = maxDocumentId // 最大不超出时间范围限制内的最大文档ID
			}

			if max > uint32(conf.GetPageSize()) {
				min = max - uint32(conf.GetPageSize()) + 1
			} else {
				min = 1
			}

			if min < minDocumentId {
				min = minDocumentId // 最小不超出时间范围限制内的最小文档ID
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel())
			}
		}
	} else {
		// 前一页
		if totalCount > currentDocId {
			var min, max uint32
			min = currentDocId + 1

			if min < minDocumentId {
				min = minDocumentId // 最小不超出时间范围限制内的最小文档ID
			}

			max = min + uint32(conf.GetPageSize()) - 1
			if max > totalCount {
				max = totalCount
			}

			if max > maxDocumentId {
				max = maxDocumentId // 最大不超出时间范围限制内的最大文档ID
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel())
			}
		}
	}

	return rs
}

// 参数widxs长度要求大于1，currentDocId不传就是查第一页
func findSame(currentDocId uint32, forward bool, minDocumentId uint32, maxDocumentId uint32, storeLogData *storage.LogDataStorageHandle, widxs ...*WidxStorage) *SearchResult {

	var rs = new(SearchResult)
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

	// 简单检查排除没结果的情景
	totalCount := minIdx.idxwordStorage.GetTotalCount(minIdx.word)
	if totalCount == 0 || (totalCount == 1 && currentDocId > 0) {
		return rs // 索引件数0、或只有1条又还要跳过，都是找不到
	}

	// 找匹配位置并排除没结果的情景
	pos := totalCount // 默认检索最新第一页
	if currentDocId > 0 {
		pos = minIdx.idxdocStorage.GetWordDocSeq(minIdx.word, currentDocId) // 有相对文档ID时找相对位置
		if pos == 0 || (pos == 1 && forward) || (pos == totalCount && !forward) {
			return rs // 找不到、或最后条还要向后、或最前条还要向前，都是找不到
		}
	}

	// 位置就绪
	var rsCnt int = 0
	var flg bool
	if currentDocId == 0 || currentDocId > 0 && forward {
		// 无相对文档ID、或有且是后一页方向
		if currentDocId > 0 {
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
				// 找到则加入结果
				if flg {
					rsCnt++
					rs.Data = append(rs.Data, storeLogData.GetLogDataModel(docId))
					if rsCnt >= conf.GetPageSize() {
						break // 最多找一页
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
		for i := pos; i <= totalCount; i++ {
			// 取值
			docId := minIdx.idxwordStorage.GetDocId(minIdx.word, i)

			if docId < minDocumentId || docId > maxDocumentId {
				continue // 不在时间范围条件内，不匹配，跳过
			}

			// 比较
			flg = true
			for i := 0; i < cnt; i++ {
				if widxs[i] == minIdx {
					continue // 跳过比较自己
				}
				if widxs[i].idxdocStorage.GetWordDocSeq(widxs[i].word, docId) == 0 {
					flg = false // 没找到
					break
				}
			}
			// 找到则加入结果
			if flg {
				rsCnt++
				ary = append(ary, storeLogData.GetLogDataModel(docId))
				if rsCnt >= conf.GetPageSize() {
					break // 最多找一页
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
