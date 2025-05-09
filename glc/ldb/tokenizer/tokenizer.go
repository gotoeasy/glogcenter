/**
 * TODO 待整理
 */
package tokenizer

import (
	"github.com/gotoeasy/glang/cmn"
)

// 分词器
var sego *cmn.TokenizerSego

// 初始化分词器
func init() {
	sego = cmn.NewTokenizerSego()
}

// 按搜索引擎模式进行分词后返回分词数组
func CutForSearch(text string) []string {
	return sego.CutForSearch(text)
}

// 按搜索引擎模式进行分词后返回分词数组，可自定义添加或删除分词
func CutForSearchEx(text string, addWords []string, delWords []string) []string {
	return sego.CutForSearchEx(text, addWords, delWords)
}
