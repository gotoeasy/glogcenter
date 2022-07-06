/**
 * 日志检索
 * 1）无关键词时直接全量检索
 * 2）有关键词时检索索引
 * 3）支持指定相对ID及方向进行前后翻页检索
 */
package search

import (
	"fmt"
	"glc/cmn"
	"glc/ldb/storage"
	"glc/ldb/storage/indexdoc"
	"glc/ldb/storage/indexword"
	"glc/ldb/storage/logdata"
)

type SearchResult struct {
	Total string                  `json:"total,omitempty"` // 日志总量件数（用10进制字符串形式以避免出现科学计数法）
	Count string                  `json:"count,omitempty"` // 当前条件最多匹配件数（用10进制字符串形式以避免出现科学计数法）
	Data  []*logdata.LogDataModel `json:"data,omitempty"`  // 检索结果数据（日志文档数组）
}

type WidxStorage struct {
	word           string
	idxdocStorage  *indexdoc.DocIndexStorage
	idxwordStorage *indexword.WordIndexStorage
}

// 多关键词时计算关键词索引交集
func SearchWordIndex(storeName string, kws []string, pageSize int, currentDocId uint32, forward bool) *SearchResult {
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	var widxs []*WidxStorage
	for _, word := range kws {
		widxStorage := &WidxStorage{
			word:           word,
			idxdocStorage:  indexdoc.NewDocIndexStorage(storeName),
			idxwordStorage: indexword.NewWordIndexStorage(storeName),
		}
		widxs = append(widxs, widxStorage)
	}
	return findSame(pageSize, currentDocId, forward, storeLogData, widxs...)
}

// 无关键词时走全量检索
func SearchLogData(storeName string, pageSize int, currentDocId uint32, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	totalCount := storeLogData.TotalCount()                    // 总件数
	rs.Total = cmn.Uint32ToString(totalCount)                  // 返回的总件数用10进制字符串形式以避免出现科学计数法
	rs.Count = cmn.Uint32ToString(totalCount)                  // 当前条件最多匹配件数

	if totalCount == 0 {
		return rs
	}

	if currentDocId == 0 {
		// 第一页
		var min, max uint32
		max = totalCount
		if max > uint32(pageSize) {
			min = max - uint32(pageSize) + 1
		} else {
			min = 1
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
			if max > uint32(pageSize) {
				min = max - uint32(pageSize) + 1
			} else {
				min = 1
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
			max = min + uint32(pageSize) - 1
			if max > totalCount {
				max = totalCount
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel())
			}
		}
	}

	return rs
}

// 参数widxs长度要求大于1，currentDocId不传就是查第一页
func findSame(pageSize int, currentDocId uint32, forward bool, storeLogData *storage.LogDataStorageHandle, widxs ...*WidxStorage) *SearchResult {

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

		for i := pos; i > 0; i-- {
			// 取值
			docId := minIdx.idxwordStorage.GetDocId(minIdx.word, i)
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
				rs.Data = append(rs.Data, storeLogData.GetLogDataModel(docId))
				if rsCnt >= pageSize {
					break // 最多找一页
				}
			}
		}
	} else {
		// 有相对文档ID且是前一页方向
		pos++
		var ary []*logdata.LogDataModel
		for i := pos; i <= totalCount; i++ {
			// 取值
			docId := minIdx.idxwordStorage.GetDocId(minIdx.word, i)
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
				if rsCnt >= pageSize {
					break // 最多找一页
				}
			}
		}

		// 倒序放入结果
		for i := len(ary) - 1; i >= 0; i-- {
			rs.Data = append(rs.Data, ary[i])
		}
	}

	rs.Total = fmt.Sprintf("%d", storeLogData.TotalCount())
	return rs
}
