package com

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"glc/conf"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

func ToBytes(data any) []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func GeyStoreNameByDate(name string) string {
	if name == "" {
		name = "logdata"
	}
	if conf.IsStoreNameAutoAddDate() {
		return fmt.Sprint(name, "-", time.Now().Format("20060102")) // name-yyyymmdd
	}
	return name
}

func JoinBytes(bts ...[]byte) []byte {
	return bytes.Join(bts, []byte(""))
}

// 取日志仓名列表，以“.”开头的默认忽略
func GetStorageNames(path string, excludes ...string) []string {
	fileinf, err := os.ReadDir(path)
	if err != nil {
		cmn.Error("读取目录失败", err)
		return []string{}
	}

	mapDir := make(map[string]string)
	for _, v := range fileinf {
		if v.IsDir() && !cmn.Startwiths(v.Name(), ".") {
			mapDir[v.Name()] = ""
		}
	}
	for i := 0; i < len(excludes); i++ {
		delete(mapDir, excludes[i])
	}

	var rs []string
	for k := range mapDir {
		rs = append(rs, k)
	}

	// 倒序
	sort.Slice(rs, func(i, j int) bool {
		return rs[i] > rs[j]
	})

	return rs
}

func GetDirInfo(path string) (uint32, int64, error) {
	var count uint32
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return err
	})
	return count, size, err
}

// 当前日期加减天数后的yyyymmdd格式
func GetYyyymmdd(days int) string {
	return time.Now().AddDate(0, 0, days).Format("20060102")
}

func Random() uint32 {
	rand.Seed(time.Now().UnixNano())
	for {
		v := rand.Uint32()
		if v != 0 {
			return v
		}
	}
}

func Unique(s []string) []string {
	m := make(map[string]struct{}, 0)
	newS := make([]string, 0)
	for _, i2 := range s {
		if _, ok := m[i2]; !ok {
			newS = append(newS, i2)
			m[i2] = struct{}{}
		}
	}
	return newS
}

func GetLocalGlcUrl() string {
	if conf.GetServerUrl() != "" {
		return conf.GetServerUrl()
	}

	if conf.GetServerIp() != "" {
		return "http://" + conf.GetServerIp() + ":" + conf.GetServerPort()
	}

	return "http://" + GetLocalIp() + ":" + conf.GetServerPort()

}

var localIp string

func GetLocalIp() string {
	if localIp != "" {
		return localIp
	}

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				localIp = ipnet.IP.String()
			}
		}
	}
	return localIp
}
