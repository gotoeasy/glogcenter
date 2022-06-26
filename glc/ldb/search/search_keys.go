/**
 * 日志检索
 * 1）无关键词时直接全量检索
 * 2）有关键词时检索索引
 * 3）支持指定相对ID及方向进行前后翻页检索
 */
package search

import (
	"glc/ldb/storage"
)

type SearchResult struct {
	TotalCount  uint64                     // 总件数
	PageFirstId uint64                     // 当前页第一条的文档ID或索引ID
	PageLastId  uint64                     // 当前页最后一条的文档ID或索引ID
	Data        []*storage.LogDataDocument // 检索结果数据（日志文档数组）
}

// 单关键词浏览日志
func Search(storeName string, word string, pageSize int, currentId uint64, forward bool) *SearchResult {

	// 检查修正pageSize
	if pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	if word == "" {
		return searchLogData(storeName, pageSize, currentId, forward)
	}
	return searchWordIndex(storeName, word, pageSize, currentId, forward)
}

// 无关键词时走全量检索
func searchLogData(storeName string, pageSize int, currentId uint64, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	rs.TotalCount = storeLogData.TotalCount()                  // 总件数

	if rs.TotalCount == 0 {
		return rs
	}

	if currentId == 0 {
		// 第一页
		var min, max uint64
		max = rs.TotalCount
		if max > uint64(pageSize) {
			min = max - uint64(pageSize) + 1
		} else {
			min = 1
		}

		for i := max; i >= min; i-- {
			rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i)) // 件数等同日志文档ID
		}
		rs.PageFirstId = max
		rs.PageLastId = min
	} else if forward {
		// 后一页
		if currentId > 1 {
			var min, max uint64
			max = currentId - 1
			if max > uint64(pageSize) {
				min = max - uint64(pageSize) + 1
			} else {
				min = 1
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i))
			}
			rs.PageFirstId = max
			rs.PageLastId = min
		}
	} else {
		// 前一页
		if rs.TotalCount > currentId {
			var min, max uint64
			min = currentId + 1
			max = min + uint64(pageSize) - 1
			if max > rs.TotalCount {
				max = rs.TotalCount
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(i))
			}
			rs.PageFirstId = max
			rs.PageLastId = min
		}
	}

	return rs
}

// 有关键词时走索引检索
func searchWordIndex(storeName string, word string, pageSize int, currentId uint64, forward bool) *SearchResult {

	var rs = new(SearchResult)                                 // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	storeIndex := storage.NewWordIndexStorage(storeName, word) // 索引
	rs.TotalCount = storeIndex.TotalCount()                    // 总件数

	if rs.TotalCount == 0 {
		return rs
	}

	if currentId == 0 {
		// 第一页
		var min, max uint64
		max = rs.TotalCount
		if max > uint64(pageSize) {
			min = max - uint64(pageSize) + 1
		} else {
			min = 1
		}

		for i := max; i >= min; i-- {
			rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i))) // 经索引取日志文档ID
		}
		rs.PageFirstId = max
		rs.PageLastId = min
	} else if forward {
		// 后一页
		if currentId > 1 {
			var min, max uint64
			max = currentId - 1
			if max > uint64(pageSize) {
				min = max - uint64(pageSize) + 1
			} else {
				min = 1
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)))
			}
			rs.PageFirstId = max
			rs.PageLastId = min
		}
	} else {
		// 前一页
		if rs.TotalCount > currentId {
			var min, max uint64
			min = currentId + 1
			max = min + uint64(pageSize) - 1
			if max > rs.TotalCount {
				max = rs.TotalCount
			}

			for i := max; i >= min; i-- {
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.Get(i)))
			}
			rs.PageFirstId = max
			rs.PageLastId = min
		}
	}

	return rs
}
