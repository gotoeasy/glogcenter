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
	Total       string                  `json:"total,omitempty"`       // 总件数（用10进制字符串形式以避免出现科学计数法）
	PageFirstId string                  `json:"pageFirstId,omitempty"` // 当前页第一条的文档ID或索引ID
	PageLastId  string                  `json:"pageLastId,omitempty"`  // 当前页最后一条的文档ID或索引ID
	Data        []*storage.LogDataModel `json:"data,omitempty"`        // 检索结果数据（日志文档数组）
}

// 单关键词浏览日志
func Search(storeName string, word string, pageSize int, currentId uint64, forward bool) *SearchResult {
	if word == "" {
		return searchLogData(storeName, pageSize, currentId, forward)
	}
	return searchWordIndex(storeName, word, pageSize, currentId, forward)
}

// 无关键词时走全量检索
func searchLogData(storeName string, pageSize int, currentId uint64, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	totalCount := storeLogData.TotalCount()                    // 总件数
	rs.Total = cmn.Uint64ToString(totalCount, 10)              // 返回的总件数用10进制字符串形式以避免出现科学计数法

	if totalCount == 0 {
		return rs
	}

	if currentId == 0 {
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
		rs.PageFirstId = cmn.Uint64ToString(max, 36)
		rs.PageLastId = cmn.Uint64ToString(min, 36)
	} else if forward {
		// 后一页
		if currentId > 1 {
			var min, max uint64
			if currentId > totalCount {
				max = totalCount
			} else {
				max = currentId - 1
			}
			if max > uint64(pageSize) {
				min = max - uint64(pageSize) + 1
			} else {
				min = 1
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel())
			}
			rs.PageFirstId = cmn.Uint64ToString(max, 36)
			rs.PageLastId = cmn.Uint64ToString(min, 36)
		}
	} else {
		// 前一页
		if totalCount > currentId {
			var min, max uint64
			min = currentId + 1
			max = min + uint64(pageSize) - 1
			if max > totalCount {
				max = totalCount
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i).ToLogDataModel())
			}
			rs.PageFirstId = cmn.Uint64ToString(max, 36)
			rs.PageLastId = cmn.Uint64ToString(min, 36)
		}
	}

	return rs
}

// 有关键词时走索引检索
func searchWordIndex(storeName string, word string, pageSize int, currentId uint64, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	storeIndex := storage.NewWordIndexStorage(storeName, word) // 索引
	totalCount := storeIndex.TotalCount()                      // 总件数
	rs.Total = cmn.Uint64ToString(totalCount, 10)              // 返回的总件数用10进制字符串形式以避免出现科学计数法

	if totalCount == 0 {
		return rs
	}

	if currentId == 0 {
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
		rs.PageFirstId = cmn.Uint64ToString(max, 36)
		rs.PageLastId = cmn.Uint64ToString(min, 36)
	} else if forward {
		// 后一页
		if currentId > 1 {
			var min, max uint64
			if currentId > totalCount {
				max = totalCount
			} else {
				max = currentId - 1
			}
			if max > uint64(pageSize) {
				min = max - uint64(pageSize) + 1
			} else {
				min = 1
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)).ToLogDataModel())
			}
			rs.PageFirstId = cmn.Uint64ToString(max, 36)
			rs.PageLastId = cmn.Uint64ToString(min, 36)
		}
	} else {
		// 前一页
		if totalCount > currentId {
			var min, max uint64
			min = currentId + 1
			max = min + uint64(pageSize) - 1
			if max > totalCount {
				max = totalCount
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)).ToLogDataModel())
			}
			rs.PageFirstId = cmn.Uint64ToString(max, 36)
			rs.PageLastId = cmn.Uint64ToString(min, 36)
		}
	}

	return rs
}
