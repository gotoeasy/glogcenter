/**
 * KV读写封装
 * 1）Key接口统一为string
 * 2）Value接口统一为SysmntData，并提供默认字段方便自行使用
 */
package sysmnt

import (
	"bytes"
	"encoding/gob"

	"github.com/gotoeasy/glang/cmn"
)

type SysmntData struct {
	Count   uint32
	Value   uint32
	Flag    bool
	Content string
}

func (s *SysmntStorage) SetSysmntData(key string, value *SysmntData) {
	k := cmn.StringToBytes(key)
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(value)
	if err != nil {
		panic(err)
	}
	v := buffer.Bytes()
	s.Put(k, v)
}

func (s *SysmntStorage) GetSysmntData(key string) *SysmntData {
	bs, err := s.Get(cmn.StringToBytes(key))
	if err != nil {
		return &SysmntData{
			Count: 0,
		}
	}

	rs := new(SysmntData)
	buffer := bytes.NewBuffer(bs)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(rs)
	if err != nil {
		panic(err.Error())
	}
	return rs
}
