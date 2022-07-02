/**
 * 日志文档
 * 1）面向leveldb存储接口
 */
package logdata

import (
	"bytes"
	"encoding/gob"
)

// 日志的文档索引
type LogDataDocument struct {
	Id      uint32 `json:"id,omitempty"`      // 文档ID，从1开始递增
	Content string `json:"content,omitempty"` // 文档内容，内容格式自行定义
}

func (d *LogDataDocument) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(d)
	return buffer.Bytes()
}

func (d *LogDataDocument) LoadBytes(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(d)
}

func (d *LogDataDocument) ToLogDataModel() *LogDataModel {
	rs := new(LogDataModel)
	rs.LoadJson(d.Content)
	return rs
}
