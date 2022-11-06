package service

import (
	"encoding/json"
	"glc/ldb/sysmnt"

	"github.com/gotoeasy/glang/cmn"
)

type KeyValue struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Version string `json:"version,omitempty"`
}

func GetSysmntItem(key string) (*KeyValue, error) {
	sysmnt := sysmnt.NewSysmntStorage()
	by, err := sysmnt.Get(cmn.StringToBytes(key))
	if err != nil {
		return nil, err
	}

	dkv := &KeyValue{}
	dkv.LoadJson(cmn.BytesToString(by))

	return dkv, nil
}

func DelSysmntItem(key string) error {
	sysmnt := sysmnt.NewSysmntStorage()
	err := sysmnt.Del(cmn.StringToBytes(key))
	if err != nil {
		return err
	}

	return nil
}

func SetSysmntItem(kv *KeyValue) (*KeyValue, error) {
	dkv, _ := GetSysmntItem(kv.Key)
	if dkv == nil {
		dkv = &KeyValue{}
	}

	// if dkv.Version != kv.Version {
	// 	return nil, errors.New("数据版本不是最新")
	// }

	dkv.Key = kv.Key
	dkv.Value = kv.Value
	dkv.Version = cmn.Uint32ToString(cmn.StringToUint32(dkv.Version, 0) + 1)

	sysmnt := sysmnt.NewSysmntStorage()
	sysmnt.Put(cmn.StringToBytes(dkv.Key), cmn.StringToBytes(dkv.ToJson()))

	return dkv, nil
}

func (d *KeyValue) ToJson() string {
	bt, _ := json.Marshal(d)
	return cmn.BytesToString(bt)
}

func (d *KeyValue) LoadJson(jsonstr string) error {
	if jsonstr == "" {
		return nil
	}
	return json.Unmarshal(cmn.StringToBytes(jsonstr), d)
}
