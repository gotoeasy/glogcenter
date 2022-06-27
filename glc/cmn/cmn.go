package cmn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"glc/ldb/conf"
	"hash/crc32"
	"os"
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

// // 字符串(10进制无符号整数形式)转uint32，超过uint32最大值会丢失精度
// // 转换失败时返回默认值
// func StringToUint32(s string, defaultVal uint32) uint32 {
// 	v, err := strconv.ParseUint(s, 10, 32)
// 	if err != nil {
// 		return defaultVal
// 	}
// 	return uint32(v & 0xFFFFFFFF)
// }

// 字符串(指定进制无符号整数形式)转uint64，进制base范围为2~36
// 参数错误或转换失败都返回默认值
func StringToUint64(s string, base int, defaultVal uint64) uint64 {
	if s == "" || base < 2 || base > 36 {
		return defaultVal
	}

	v, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// Uint64转指定进制形式字符串
func Uint64ToString(val uint64, base int) string {
	return strconv.FormatUint(val, base)
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

// func Uint32ToBytes(num uint32) []byte {
// 	bkey := make([]byte, 4)
// 	binary.BigEndian.PutUint32(bkey, num)
// 	return bkey
// }

// func BytesToUint32(bytes []byte) uint32 {
// 	return uint32(binary.BigEndian.Uint32(bytes))
// }

func Uint64ToBytes(num uint64) []byte {
	bkey := make([]byte, 8)
	binary.BigEndian.PutUint64(bkey, num)
	return bkey
}

func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

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
	txt := "添油" + str + "加醋"
	return fmt.Sprint(crc32.ChecksumIEEE(StringToBytes(txt)) % mod)
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
