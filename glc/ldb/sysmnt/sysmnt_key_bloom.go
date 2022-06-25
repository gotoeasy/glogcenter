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
		_, err := s.Get(cmn.StringToBytes(_PREFIX + k))
		if err == nil {
			return
		}
		s.Put(cmn.StringToBytes(_PREFIX+k), cmn.StringToBytes(""))
	}
}
