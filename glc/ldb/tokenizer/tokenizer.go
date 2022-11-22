/**
 * TODO 待整理
 */
package tokenizer

import (
	"path/filepath"
	"runtime"

	"github.com/gotoeasy/glang/cmn"
)

var jiebago *cmn.TokenizerJiebago

// 初始化装载字典
func init() {
	_, filename, _, _ := runtime.Caller(0) // 当前go文件所在路径
	dictfile := filepath.Join(filepath.Dir(filename), "dict.txt")
	jiebago = cmn.NewTokenizerJiebago(dictfile)
}

// 按搜索引擎模式进行分词后返回分词数组
func CutForSearch(text string) []string {
	return jiebago.CutForSearch(text)
}

// 按搜索引擎模式进行分词后返回分词数组，可自定义添加或删除分词
func CutForSearchEx(text string, addWords []string, delWords []string) []string {
	return jiebago.CutForSearchEx(text, addWords, delWords)
}
