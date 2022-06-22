package search

import (
	"glc/cmn"
	"glc/ldb/logdata"
	"glc/ldb/logindex"
	"glc/ldb/storage"
	"glc/ldb/tokenizer"
	"log"
)

type SearchResult struct {
	TotalCount uint32
	Result     []*storage.LdbDocument
}

type searchPos struct {
	word  string                    // 关键词
	pos   uint32                    // 游标
	value uint32                    // 游标位置的索引值
	store *logindex.LogIndexStorage // 反向索引存储器
}

// 无条件浏览日志
func BrowsePage(storeName string, pageNo uint32, pageSize int) *SearchResult {

	var rs = new(SearchResult)

	store := logdata.NewLogDataStorage(storeName)
	rs.TotalCount = store.TotalCount()
	for i, min := rs.TotalCount-pageNo*uint32(pageSize), rs.TotalCount-pageNo*uint32(pageSize)-uint32(pageSize); i > min && i > 0; i-- {
		rs.Result = append(rs.Result, store.Get(i)) // 件数就是日志文档ID
	}
	return rs
}

// 单关键词浏览日志
func BrowsePageByOneWord(storeName string, word string, pageNo uint32, pageSize int) *SearchResult {

	var rs = new(SearchResult)

	storeLog := logdata.NewLogDataStorage(storeName)
	storeIndex := logindex.NewLogIndexStorage(storeName, word)
	rs.TotalCount = storeIndex.TotalCount()
	for i, min := rs.TotalCount-pageNo*uint32(pageSize), rs.TotalCount-pageNo*uint32(pageSize)-uint32(pageSize); i > min && i > 0; i-- {
		rs.Result = append(rs.Result, storeLog.Get(storeIndex.Get(i))) // 索引存放日志文档ID
	}
	return rs
}

// 检索用文字进行分词，以及针对检索特殊场景的优化
func CutSearchKey(searchKey string) []string {
	var mapKey = make(map[string]string)
	kws := tokenizer.CutForSearch(searchKey)

	for _, k := range kws {
		mapKey[k] = ""
	}

	for _, kw := range kws {
		ks := tokenizer.CutForSearch(kw)
		if len(ks) > 1 {
			for _, k := range ks {
				delete(mapKey, k)
			}
			mapKey[kw] = ""
		}
	}

	var rs []string
	for k := range mapKey {
		rs = append(rs, k)
	}

	log.Println("搜索关键词", kws, "优化后搜索", rs)
	return rs
}

// 多关键词检索（求交集）
func Search(storeName string, keyWords []string, pageSize int, tipId string, forward bool) *SearchResult {
	var sps []*searchPos
	for i, max := 0, len(keyWords); i < max; i++ {
		sp := new(searchPos)
		sp.word = keyWords[i]
		sp.store = logindex.NewLogIndexStorage(storeName, sp.word)
		sp.pos = sp.store.TotalCount()
		sp.value = sp.store.Get(sp.pos)
		sps = append(sps, sp)
	}

	// 先快速顺移到tipId提示位置，tipId是要丢弃的
	if forward {
		movePosForwardByTid(cmn.StringToUint32(tipId, 0)-1, sps)
	} else {
		movePosBackByTid(cmn.StringToUint32(tipId, 0)+1, sps)
	}

	var results []uint32
	var curSp *searchPos = nil

	var rs = new(SearchResult)
	var finded bool
	for {
		curSp = getMinValueSp(sps) // 拿最小值准备查找

		finded = true
		for _, sp := range sps {
			if curSp == sp {
				continue // 跳过自身索引比较
			}

			if forward {
				if !findForward(curSp.value, sp, sp.pos) {
					finded = false
					break
				}
			} else {
				if !findBack(curSp.value, sp, sp.pos) {
					finded = false
					break
				}
			}
		}

		// 找到
		if finded {
			rs.TotalCount++
			if len(results) < pageSize {
				results = append(results, curSp.value)
			}

			if !forward {
				if curSp.pos >= curSp.store.TotalCount() {
					curSp.pos = 0
				} else {
					curSp.pos++ // 倒移找到时，补上curSp的位置移动
					curSp.value = curSp.store.Get(curSp.pos)
				}
			}
		}

		if rs.TotalCount >= 1000 {
			break
		}

		var finish bool = false
		for _, sp := range sps {
			if sp.pos == 0 {
				finish = true
			}
		}
		if finish {
			break
		}

	}

	storeLog := logdata.NewLogDataStorage(storeName)
	if forward {
		for _, id := range results {
			rs.Result = append(rs.Result, storeLog.Get(id))
		}
	} else {
		// 倒找时，当页数据要排为倒序
		for i := len(results) - 1; i >= 0; i-- {
			rs.Result = append(rs.Result, storeLog.Get(results[i]))
		}
	}
	return rs
}

// 后页方向查找，同时移动位置。sp总是待比较位置、curValue值总是小于等于sp
func findForward(curValue uint32, sp *searchPos, startPos uint32) bool {
	endPos := uint32(1)                 // 最后位置
	var step uint32 = endPos - startPos // 最大查找范围

	for {
		// sp总是待比较位置，先直接比较，找到时移到下一个待比较位置后返回
		if curValue == sp.value {
			sp.pos--
			sp.value = sp.store.Get(sp.pos)
			return true
		}

		// 顺移到最后位置仍不匹配，且目标值更大，结束
		if sp.pos <= endPos && curValue < sp.value {
			sp.pos = 0 // 位置0
			return false
		}

		// 步长是1时，直接判断处理后结束
		if step == 1 {
			if curValue < sp.value {
				// 当前位置更大时，丢弃。当前位置更小时，位置不变留用后续比较
				sp.pos-- // 顺移1
				sp.value = sp.store.Get(sp.pos)
			}
			return false
		}

		// 默认二分查找，再根据索引差值调整优化
		step = (step + 1) / 2
		if curValue < sp.value {
			// 顺跳不够远
			min := (sp.value - curValue + 1) / 2
			if step > min {
				step = min // 优化：如果索引号差值比步长还小，按差值范围跳找。（索引号是递增的缘故，可按此进一步缩小查找范围）
			}
			sp.pos = sp.pos - step // 顺移跳找
		} else {
			// 顺跳太远了
			min := (curValue - sp.value + 1) / 2
			if step > min {
				step = min // 优化：如果索引号差值比步长还小，按差值范围跳找。（索引号是递增的缘故，可按此进一步缩小查找范围）
			}
			sp.pos = sp.pos + step // 倒移跳找

			if startPos == sp.pos {
				// 倒移后是最初位置时，结束，由于最初位置已经比较过是不匹配，直接顺移1
				sp.pos--
				sp.value = sp.store.Get(sp.pos)
				return false
			}
		}
		sp.value = sp.store.Get(sp.pos)
	}

}

// 后页方向查找，同时移动位置。sp总是待比较位置、curValue值总是小于等于sp
func findBack(curValue uint32, sp *searchPos, startPos uint32) bool {
	endPos := sp.store.TotalCount()     // 最后位置
	var step uint32 = endPos - startPos // 最大查找范围

	for {
		// sp总是待比较位置，先直接比较，找到时移到下一个待比较位置后返回
		if curValue == sp.value {
			if sp.pos < endPos {
				sp.pos++
				sp.value = sp.store.Get(sp.pos)
			} else {
				sp.pos = 0
			}
			return true
		}

		// 倒移到最大位置仍不匹配，且目标值更小，结束
		if sp.pos >= endPos && curValue > sp.value {
			sp.pos = 0 // 位置0表示可以结束查找
			return false
		}

		// 步长是1时，直接判断处理后结束
		if step == 1 {
			if curValue > sp.value {
				// 当前位置更小时，丢弃。当前位置更大时，位置不变留用后续比较
				sp.pos++ // 倒移1
				sp.value = sp.store.Get(sp.pos)
			}
			return false
		}

		// 默认二分查找，再根据索引差值调整优化
		step = (step + 1) / 2
		if curValue > sp.value {
			// 倒跳不够远
			min := (curValue - sp.value + 1) / 2
			if step > min {
				step = min // 优化：如果索引号差值比步长还小，按差值范围跳找。（索引号是递增的缘故，可按此进一步缩小查找范围）
			}
			sp.pos = sp.pos + step // 倒移跳找
		} else {
			// 倒跳太远了
			min := (sp.value - curValue + 1) / 2
			if step > min {
				step = min // 优化：如果索引号差值比步长还小，按差值范围跳找。（索引号是递增的缘故，可按此进一步缩小查找范围）
			}
			sp.pos = sp.pos - step // 顺移跳找

			if sp.pos == startPos {
				// 顺移后是最初位置时，结束，由于最初位置已经比较过是不匹配，直接倒移1
				sp.pos++
				sp.value = sp.store.Get(sp.pos)
				return false
			}

		}
		sp.value = sp.store.Get(sp.pos)
	}

}

// 顺移到第一个比该值小的位置
func movePosForwardByTid(tid uint32, sps []*searchPos) {
	if tid == 0 {
		return // 不会有0值id，默认检索第一页，不用处理
	}
	for _, sp := range sps {
		if tid < sp.value {
			findForward(tid, sp, sp.pos) // sp很长，tid是其中某个值，顺移比较直到比tid还小
		}
	}
}

// 倒移到第一个比该值小的位置
func movePosBackByTid(tid uint32, sps []*searchPos) {

	// 默认从1开始全范围找
	for _, sp := range sps {
		sp.pos = 1
		sp.value = sp.store.Get(sp.pos)
	}

	if tid == 0 {
		return // 不会有0值id，0表示没有提示，默认从最后页开始找
	}

	// 调整到tid的位置找
	for _, sp := range sps {
		b := findBack(tid, sp, 1) // sp很长，tid是其中某个值，倒移比较直到大于tid
		if b {
			sp.pos++ // 相同时回退一位
		} else if sp.pos < 1 {
			sp.pos = sp.store.TotalCount() // 最小1，找不到时可能被置为0，改成最后位置供后面查找处理
			sp.value = sp.store.Get(sp.pos)
		}
	}
}

func getMinValueSp(sps []*searchPos) *searchPos {
	var curSp *searchPos = nil
	for _, sp := range sps {
		if curSp == nil || sp.value < curSp.value {
			curSp = sp
			continue
		}
	}
	return curSp
}
