/**
 * 日志检索
 * 1）无关键词时直接全量检索
 * 2）有关键词时检索索引
 * 3）支持指定相对ID及方向进行前后翻页检索
 */
package search

import (
	"glc/cmn"
	"glc/ldb/storage"
)

type SearchResult struct {
	Total string                  `json:"total,omitempty"` // 总件数（用10进制字符串形式以避免出现科学计数法）
	Data  []*storage.LogDataModel `json:"data,omitempty"`  // 检索结果数据（日志文档数组）
}

// 多关键词浏览日志
func Search(storeName string, kws []string, pageSize int, currentDocId uint64, forward bool) *SearchResult {
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	var widxs []*storage.WordIndexStorage
	for _, word := range kws {
		widxs = append(widxs, storage.NewWordIndexStorage(storeName, word))
	}
	return findSame(pageSize, currentDocId, forward, storeLogData, widxs...)
}

// 无关键词时走全量检索
func SearchLogData(storeName string, pageSize int, currentDocId uint64, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	totalCount := storeLogData.TotalCount()                    // 总件数
	rs.Total = cmn.Uint64ToString(totalCount, 10)              // 返回的总件数用10进制字符串形式以避免出现科学计数法

	if totalCount == 0 {
		return rs
	}

	if currentDocId == 0 {
		// 第一页
		var min, max uint64
		max = totalCount
		if max > uint64(pageSize) {
			min = max - uint64(pageSize) + 1
		} else {
			min = 1
		}

		for i := max; i >= min; i-- {
			rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel()) // 件数等同日志文档ID
		}
	} else if forward {
		// 后一页
		if currentDocId > 1 {
			var min, max uint64
			if currentDocId > totalCount {
				max = totalCount
			} else {
				max = currentDocId - 1
			}
			if max > uint64(pageSize) {
				min = max - uint64(pageSize) + 1
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
			var min, max uint64
			min = currentDocId + 1
			max = min + uint64(pageSize) - 1
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

// 有关键词时走索引检索
func SearchWordIndex(storeName string, word string, pageSize int, currentDocId uint64, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	storeIndex := storage.NewWordIndexStorage(storeName, word) // 索引
	totalCount := storeIndex.TotalCount()                      // 总件数
	rs.Total = cmn.Uint64ToString(totalCount, 10)              // 返回的总件数用10进制字符串形式以避免出现科学计数法

	if totalCount == 0 {
		return rs
	}

	if currentDocId == 0 {
		// 第一页
		var min, max uint64
		max = totalCount
		if max > uint64(pageSize) {
			min = max - uint64(pageSize) + 1
		} else {
			min = 1
		}

		for i := max; i >= min; i-- {
			rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)).ToLogDataModel()) // 经索引取日志文档ID
		}
	} else if forward {
		// 后一页
		if currentDocId > 1 {
			var min, max uint64
			if currentDocId > totalCount {
				max = totalCount
			} else {
				max = currentDocId - 1
			}
			if max > uint64(pageSize) {
				min = max - uint64(pageSize) + 1
			} else {
				min = 1
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)).ToLogDataModel())
			}
		}
	} else {
		// 前一页
		if totalCount > currentDocId {
			var min, max uint64
			min = currentDocId + 1
			max = min + uint64(pageSize) - 1
			if max > totalCount {
				max = totalCount
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)).ToLogDataModel())
			}
		}
	}

	return rs
}
