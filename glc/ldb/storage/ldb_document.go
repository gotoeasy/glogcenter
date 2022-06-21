package storage

import (
	"bytes"
	"encoding/gob"
)

// 日志的文档索引
type LdbDocument struct {
	Id      uint32 `json:"id,omitempty"`      // 文档ID，从1开始递增
	Content string `json:"content,omitempty"` // 文档内容，内容格式自行定义
}

func (d *LdbDocument) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(d)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func ParseBytes(data []byte) *LdbDocument {
	d := new(LdbDocument)
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(d)
	if err != nil {
		return nil
	}
	return d
}

// func (d *LdbDocument) ToJson() string {
// 	bt, _ := json.Marshal(d)
// 	return string(bt)
// }

// func ParseJson(jsonstr string) *LdbDocument {
// 	d := new(LdbDocument)
// 	json.Unmarshal([]byte(jsonstr), d)
// 	return d
// }
