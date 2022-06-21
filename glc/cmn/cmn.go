package cmn

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

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
	return fmt.Sprint(crc32.ChecksumIEEE([]byte(txt)) % mod)
}
