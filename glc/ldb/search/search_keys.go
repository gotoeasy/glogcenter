package search

import (
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
func Search(storeName string, keyWords []string, pageNo uint32, pageSize int) *SearchResult {

	var forward bool = true // 检索方向（true：后页方向，false：前页方向）
	var sps []*searchPos
	for i, max := 0, len(keyWords); i < max; i++ {
		sp := new(searchPos)
		sp.word = keyWords[i]
		sp.store = logindex.NewLogIndexStorage(storeName, sp.word)
		sp.pos = sp.store.TotalCount()
		sp.value = sp.store.Get(sp.pos)
		sps = append(sps, sp)
	}

	var results []uint32
	var curSp *searchPos = nil

	var rs = new(SearchResult)
	var finded bool
	for {
		curSp = getMinValueSp(sps) // 后页方向，拿最小值准备查找

		finded = true
		for _, sp := range sps {
			if curSp == sp { // if curSp.word == sp.word {
				continue // 跳过自身索引比较
			}

			if forward {
				if !findForward(curSp, sp) {
					finded = false
					break
				}
			} else {
				log.Println("// TODO")
			}
		}

		// 找到
		if finded {
			rs.TotalCount++
			if len(results) <= pageSize {
				results = append(results, curSp.value)
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
	for _, id := range results {
		rs.Result = append(rs.Result, storeLog.Get(id))
	}
	return rs
}

// 后页方向查找，同时移动位置。sp总是待比较位置、curSp值总是小于等于sp
func findForward(curSp *searchPos, sp *searchPos) bool {

	var startPos uint32 = sp.pos // 最初位置
	var step uint32 = sp.pos     // 初始设定最大查找范围(当前整长)

	for {
		// sp总是待比较位置，先直接比较，找到时移到下一个待比较位置后返回
		if curSp.value == sp.value {
			sp.pos--
			sp.value = sp.store.Get(sp.pos)
			return true
		}

		// 顺移到最后位置仍不匹配，且目标值更大，结束
		if sp.pos == 1 && curSp.value < sp.value {
			sp.pos-- // 位置0
			return false
		}

		// 步长是1时，直接判断处理后结束
		if step == 1 {
			if curSp.value < sp.value {
				// 当前位置更大时，丢弃。当前位置更小时，位置不变留用后续比较
				sp.pos-- // 瞬移1
				sp.value = sp.store.Get(sp.pos)
			}
			return false
		}

		// 默认二分查找，再根据索引差值调整优化
		step = (step + 1) / 2
		if curSp.value < sp.value {
			// 顺跳不够远
			min := (sp.value - curSp.value + 1) / 2
			if step > min {
				step = min // 优化：如果索引号差值比步长还小，按差值范围跳找。（索引号是递增的缘故，可按此进一步缩小查找范围）
			}
			sp.pos = sp.pos - step // 顺移跳找
		} else {
			// 顺跳太远了
			min := (curSp.value - sp.value + 1) / 2
			if step > min {
				step = min // 优化：如果索引号差值比步长还小，按差值范围跳找。（索引号是递增的缘故，可按此进一步缩小查找范围）
			}
			sp.pos = sp.pos + step // 倒移跳找

			if startPos == sp.pos {
				// 倒移后是最初位置时，结束，由于最初位置已经比较过是不匹配，直接瞬移1
				sp.pos--
				sp.value = sp.store.Get(sp.pos)
				return false
			}
		}
		sp.value = sp.store.Get(sp.pos)
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

// // 前页方向移动位置继续找
// func findBack(curSp *searchPos, sp *searchPos) bool {
// 	return false
// }
