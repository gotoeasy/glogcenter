/**
 * TODO 待整理
 */
package tokenizer

import (
	"github.com/gotoeasy/glang/cmn"
	"github.com/wangbin/jiebago"
)

var seg jiebago.Segmenter // 参考 https://github.com/wangbin/jiebago
var ignoreWords []string  // 默认忽略的单词
// var simpleCutMode bool    // 是否使用简单分词模式

// 初始化装载字典
func init() {
	// 分词字典文件默认是“dict.txt”，支持通过环境变量“DICT_FILE”设定
	dictFile := cmn.GetEnvStr("DICT_FILE", "dict.txt") // TODO

	// 默认忽略的单字，主要是键盘能直接敲出的全角半角符号，字母数字不大好说不做排除
	ignoreWords = cmn.Split("`~!@# $%^&*()-_=+[{]}\\|;:'\",<.>/?，。《》；：‘　’“”、|】｝【｛＋－—（）×＆…％￥＃＠！～·\t\r\n", "")

	// 有些关键词是要按规则忽略的，支持通过环境变量“INGORE_WORDS”设定，以半角逗号分隔
	strWords := cmn.GetEnvStr("INGORE_WORDS", "") // TODO
	if strWords != "" {
		ignoreWords = append(ignoreWords, cmn.Split(strWords, ",")...)
	}

	// simpleCutMode = conf.GetenvBool("SIMPLE_CUT_MODE", true)

	seg.LoadDictionary(dictFile)
}

// 按搜索引擎模式进行分词后返回分词数组
func CutForSearch(text string) []string {
	return CutForSearchEx(text, nil, nil)
}

// 按搜索引擎模式进行分词后返回分词数组，可自定义添加或删除分词
func CutForSearchEx(text string, addWords []string, delWords []string) []string {

	if cmn.Trim(text) == "" {
		return []string{}
	}

	txt := cmn.ToLower(text)

	// 结巴分词
	sch := seg.CutForSearch(txt+" "+cmn.Join(addWords, " "), true) // TODO 暂且补丁
	var mapStr = make(map[string]string)
	tmp := ""
	for word := range sch {
		tmp = cmn.Trim(word)
		if tmp != "" {
			mapStr[tmp] = ""
		}
	}

	// // 简单分词
	// if simpleCutMode {
	// 	// 针对日志再保留特殊字符（【.】用于包名，【/】用工于路径或日期，【_】常用于表名，【-】常用于日期或连词）
	// 	txt = replaceByRegex(txt, "[,/;\"'?？，。!！=@#\\[\\]【】\\\\:]", " ")
	// 	//log.Println(txt)
	// 	keys := strings.Split(txt, " ")
	// 	for _, word := range keys {
	// 		tmp = strings.TrimSpace(word)
	// 		tmp = strings.TrimRight(strings.Trim(tmp, "."), "-") // 去除两边点和右边减号
	// 		if tmp != "" {
	// 			mapStr[tmp] = ""
	// 		}
	// 	}
	// }

	// 删除默认忽略的单词
	for _, word := range ignoreWords {
		delete(mapStr, word)
	}

	// 自定义添加分词
	for _, word := range addWords {
		tmp = cmn.Trim(cmn.ToLower(word))
		if tmp != "" {
			mapStr[cmn.Trim(cmn.ToLower(word))] = ""
		}
	}

	// 自定义删除分词
	for _, word := range delWords {
		tmp = cmn.Trim(cmn.ToLower(word))
		delete(mapStr, tmp)
	}

	// 遍历所有键存入数组返回，忽略大小写无重复
	var rs []string
	for k := range mapStr {
		rs = append(rs, k)
	}
	return rs
}

// func replaceByRegex(str string, rule string, replace string) string {
// 	reg, err := regexp.Compile(rule)
// 	if reg == nil || err != nil {
// 		panic("ssssssssssssssssssssssssssss")
// 		// return ""
// 	}
// 	return reg.ReplaceAllString(str, replace)
// }

// // 检索用文字进行分词，以及针对检索特殊场景的优化
// func GetSearchKey(searchKey string) string {
// 	if searchKey == "" {
// 		return ""
// 	}

// 	var mapKey = make(map[string]string)
// 	kws := CutForSearch(searchKey)

// 	for _, k := range kws {
// 		mapKey[k] = ""
// 	}

// 	for _, kw := range kws {
// 		ks := CutForSearch(kw)
// 		if len(ks) > 1 {
// 			for _, k := range ks {
// 				delete(mapKey, k)
// 			}
// 			mapKey[kw] = ""
// 		}
// 	}

// 	var rs []string
// 	for k := range mapKey {
// 		rs = append(rs, k)
// 	}

// 	// TODO
// 	log.Println("搜索关键词", kws, "优化后搜索", rs)
// 	return strings.Join(rs, " ")
// }
