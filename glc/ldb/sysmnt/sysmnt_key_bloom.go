/**
 * 利用leveldb简易实现布隆过滤器效果
 * （有必要时可直接使用布隆过滤器实现）
 */
package sysmnt

import (
	"glc/cmn"
)

const _PREFIX = "?"

// 检查指定关键词是否都有数据
func (s *SysmntStorage) ContainsKeyWord(kws []string) bool {
	for _, k := range kws {
		_, err := s.Get(cmn.StringToBytes(_PREFIX + k))
		if err != nil {
			return false
		}
	}
	return true
}

// 添加关键词
func (s *SysmntStorage) AddKeyWords(kws []string) {
	for _, k := range kws {
		s.Put(cmn.StringToBytes(_PREFIX+k), cmn.StringToBytes("")) // TODO 这个是否有性能问题？
		// log.Println("关键词已标记存在数据：", k)
	}
}
