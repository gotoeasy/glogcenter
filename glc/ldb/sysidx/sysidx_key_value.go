/**
 * KV读写封装
 * 1）Key接口统一为string
 * 2）Value接口统一为SysidxData，并提供默认字段方便自行使用
 */
package sysidx

import (
	"bytes"
	"encoding/gob"
	"glc/cmn"
)

type SysidxData struct {
	Count   uint32
	Value   uint32
	Flag    bool
	Content string
}

func (s *SysidxStorage) SetSysidxData(key string, value *SysidxData) {
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

func (s *SysidxStorage) GetSysidxData(key string) *SysidxData {
	bs, err := s.Get(cmn.StringToBytes(key))
	if err != nil {
		return &SysidxData{
			Count: 0,
		}
	}

	rs := new(SysidxData)
	buffer := bytes.NewBuffer(bs)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(rs)
	if err != nil {
		panic(err.Error())
	}
	return rs
}
