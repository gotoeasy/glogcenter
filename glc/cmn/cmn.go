package cmn

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// 字符串(10进制无符号整数形式)转uint32，超过uint32最大值会丢失精度
// 转换失败时返回默认值
func StringToUint32(s string, defaultVal uint32) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return defaultVal
	}
	return uint32(v & 0xFFFFFFFF)
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
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
	return uint32(binary.BigEndian.Uint32(bytes))
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

func Getenv(name string, defaultValue string) string {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}
	return s
}

func GetenvInt(name string, defaultValue int) int {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}

func GetenvBool(name string, defaultValue bool) bool {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}

	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}

// 字符串哈希处理后取模(余数)，返回值最大不超过mod值
func HashMod(str string, mod uint32) string {
	txt := "添油" + str + "加醋"
	return fmt.Sprint(crc32.ChecksumIEEE(StringToBytes(txt)) % mod)
}
