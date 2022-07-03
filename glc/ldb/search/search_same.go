/**
 * 反向索引求交集
 */
package search

import (
	"fmt"
	"glc/cmn"
	"glc/ldb/storage"
	"glc/ldb/storage/logdata"
)

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
