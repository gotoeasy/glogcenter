/**
 * TODO 待整理
 */
package tokenizer

import (
	"glc/conf"

	"github.com/gotoeasy/glang/cmn"
)

var sego *cmn.TokenizerSego

// 初始化装载字典
func init() {

	// 加载用户词典
	dicFiles, _ := cmn.GetFiles(conf.GetDictDir(), ".txt")
	sego = cmn.NewTokenizerSego(dicFiles...)
}

// 按搜索引擎模式进行分词后返回分词数组
func CutForSearch(text string) []string {
	return sego.CutForSearch(text)
}

// 按搜索引擎模式进行分词后返回分词数组，可自定义添加或删除分词
func CutForSearchEx(text string, addWords []string, delWords []string) []string {
	return sego.CutForSearchEx(text, addWords, delWords)
}
