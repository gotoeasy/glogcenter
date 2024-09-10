/**
 * TODO 待整理
 */
package tokenizer

import (
	"embed"
	"glc/conf"

	"github.com/gotoeasy/glang/cmn"
)

var (
	//go:embed dict.zip
	dictionary embed.FS
)

var sego *cmn.TokenizerSego

// 初始化装载字典
func init() {

	// 默认字典，不存在时尝试从go:embed复制
	defaultDictFile := "/glogcenter/.dictionary/dict.txt"
	if !cmn.IsExistFile(defaultDictFile) {
		bts, err := dictionary.ReadFile("dict.zip")
		if err == nil {
			cmn.UnZipBytes(bts, defaultDictFile)
		}
	}

	// 加载字典，自定义字典+默认字典
	dicFiles, _ := cmn.GetFiles(conf.GetDictDir(), ".txt")
	dicFiles = append(dicFiles, defaultDictFile)

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
