package cmn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"glc/conf"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"
)

// 字符串(10进制无符号整数形式)转int，超过int最大值会丢失精度
// 转换失败时返回默认值
func StringToInt(s string, defaultVal int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// 字符串(10进制无符号整数形式)转uint32，超过uint32最大值会丢失精度
// 转换失败时返回默认值
func StringToUint32(s string, defaultVal uint32) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return defaultVal
	}
	return uint32(v & 0xFFFFFFFF)
}

// // 字符串(指定进制无符号整数形式)转uint64，进制base范围为2~36
// // 参数错误或转换失败都返回默认值
// func StringToUint64(s string, base int, defaultVal uint64) uint64 {
// 	if s == "" || base < 2 || base > 36 {
// 		return defaultVal
// 	}

// 	v, err := strconv.ParseUint(s, base, 64)
// 	if err != nil {
// 		return defaultVal
// 	}
// 	return v
// }

// // Uint64转指定进制形式字符串
// func Uint64ToString(val uint64, base int) string {
// 	return strconv.FormatUint(val, base)
// }

// Uint32转字符串
func Uint32ToString(val uint32) string {
	return fmt.Sprintf("%d", val)
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
func StringToBool(s string, defaultVal bool) bool {
	lower := strings.ToLower(s)
	if lower == "true" {
		return true
	}
	if lower == "false" {
		return false
	}
	return defaultVal
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Uint32ToBytes(num uint32) []byte {
	bkey := make([]byte, 4)
	binary.BigEndian.PutUint32(bkey, num)
	return bkey
}

func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func Uint16ToBytes(num uint16) []byte {
	bkey := make([]byte, 2)
	binary.BigEndian.PutUint16(bkey, num)
	return bkey
}

// func Uint64ToBytes(num uint64) []byte {
// 	bkey := make([]byte, 8)
// 	binary.BigEndian.PutUint64(bkey, num)
// 	return bkey
// }

// func BytesToUint64(bytes []byte) uint64 {
// 	return binary.BigEndian.Uint64(bytes)
// }

func LenRune(str string) int {
	return utf8.RuneCountInString(str)
}

func LeftRune(str string, length int) string {
	if LenRune(str) <= length {
		return str
	}

	var rs string
	for i, s := range str {
		if i == length {
			break
		}
		rs = rs + string(s)
	}
	return rs
}

func PathSeparator() string {
	return string(os.PathSeparator)
}

// 字符串哈希处理后取模(余数)，返回值最大不超过mod值
func HashAndMod(str string, mod uint32) string {
	return fmt.Sprint(crc32.ChecksumIEEE(StringToBytes(str)) % mod)
}

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
		name = "default"
	}
	if conf.IsStoreNameAutoAddDate() {
		return fmt.Sprint(name, "-", time.Now().Format("20060102")) // name-yyyymmdd
	}
	return name
}

func EndwithsRune(str string, endstr string) bool {

	if endstr == "" || str == endstr {
		return true
	}

	strs := []rune(str)
	ends := []rune(endstr)
	lens := len(strs)
	lene := len(ends)
	if lens < lene {
		return false
	}

	dif := lens - lene
	for i := lene - 1; i >= 0; i-- {
		if strs[dif+i] != ends[i] {
			return false
		}
	}

	return true
}

func SubStringRune(str string, start int, end int) string {
	srune := []rune(str)
	slen := len(srune)
	if start >= slen || start >= end || start < 0 {
		return ""
	}

	rs := ""
	for i := start; i < slen && i < end; i++ {
		rs += string(srune[i])
	}
	return rs
}

func JoinBytes(bts ...[]byte) []byte {
	return bytes.Join(bts, []byte(""))
}

func GetStorageNames(path string, excludes ...string) []string {
	fileinf, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("读取目录失败", err)
		return []string{}
	}

	mapDir := make(map[string]string)
	for _, v := range fileinf {
		if v.IsDir() {
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
