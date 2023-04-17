/**
 * 日志模型
 * 1）面向日志接口，设定常用属性方便扩充
 */

package logdata

import (
	"encoding/json"

	"github.com/gotoeasy/glang/cmn"
)

// LogDataModel 日志模型
// Text是必须有的日志内容，Id自增，内置其他属性可选
// 其中Tags是空格分隔的标签，日期外各属性值会按空格分词
// 对应的json属性统一全小写
type LogDataModel struct {
	Id         string   `json:"id,omitempty"`         // 从1开始递增(36进制字符串)
	Text       string   `json:"text,omitempty"`       // 【必须】日志内容，多行时仅为首行，直接显示用，是全文检索对象
	Detail     string   `json:"detail,omitempty"`     // 多行时的详细日志信息，通常是包含错误堆栈等的日志内容
	Date       string   `json:"date,omitempty"`       // 日期（格式自定义）
	Tags       []string `json:"tags,omitempty"`       // 自定义标签，都作为关键词看待处理
	Server     string   `json:"server,omitempty"`     // 服务器
	Client     string   `json:"client,omitempty"`     // 客户端
	User       string   `json:"user,omitempty"`       // 用户
	System     string   `json:"system,omitempty"`     // 所属系统
	TraceId    string   `json:"traceid,omitempty"`    // 跟踪ID
	Keywords   []string `json:"keywords,omitempty"`   // 自定义的关键词
	Sensitives []string `json:"sensitives,omitempty"` // 要删除的敏感词
}

func (d *LogDataModel) ToJson() string {
	bt, _ := json.Marshal(d)
	return cmn.BytesToString(bt)
}

func (d *LogDataModel) LoadJson(jsonstr string) error {
	return json.Unmarshal(cmn.StringToBytes(jsonstr), d)
}
