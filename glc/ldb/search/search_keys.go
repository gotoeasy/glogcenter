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
	"glc/ldb/storage/indexdoc"
	"glc/ldb/storage/indexword"
	"glc/ldb/storage/logdata"
)

type SearchResult struct {
	Total string                  `json:"total,omitempty"` // 总件数（用10进制字符串形式以避免出现科学计数法）
	Data  []*logdata.LogDataModel `json:"data,omitempty"`  // 检索结果数据（日志文档数组）
}

type WidxStorage struct {
	word           string
	idxdocStorage  *indexdoc.DocIndexStorage
	idxwordStorage *indexword.WordIndexStorage
}

// 多关键词浏览日志
func Search(storeName string, kws []string, pageSize int, currentDocId uint32, forward bool) *SearchResult {
	storeLogData := storage.NewLogDataStorageHandle(storeName) // 数据
	var widxs []*WidxStorage
	for _, word := range kws {
		widxStorage := &WidxStorage{
			word:           word,
			idxdocStorage:  indexdoc.NewDocIndexStorage(storeName, word),
			idxwordStorage: indexword.NewWordIndexStorage(storeName, word),
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

// 有关键词时走索引检索
func SearchWordIndex(storeName string, word string, pageSize int, currentDocId uint32, forward bool) *SearchResult {

	var rs = new(SearchResult)                                   // 检索结果
	storeLogData := storage.NewLogDataStorageHandle(storeName)   // 数据
	storeIndex := indexword.NewWordIndexStorage(storeName, word) // 索引
	totalCount := storeIndex.GetTotalCount(word)                 // 总件数
	rs.Total = cmn.Uint32ToString(storeLogData.TotalCount())     // 返回的总件数用10进制字符串形式以避免出现科学计数法

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
			rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.GetDocId(word, i)).ToLogDataModel()) // 经索引取日志文档ID
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
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.GetDocId(word, i)).ToLogDataModel())
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
				rs.Data = append(rs.Data, storeLogData.GetLogDataDocument(storeIndex.GetDocId(word, i)).ToLogDataModel())
			}
		}
	}

	return rs
}
