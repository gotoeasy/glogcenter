package cmn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"glc/conf"
	"hash/crc32"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
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

func RightRune(str string, length int) string {
	lenr := LenRune(str)
	if lenr <= length {
		return str
	}

	var rs string
	start := lenr - length
	for i, s := range str {
		if i >= start {
			rs = rs + string(s)
		}
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
		name = "logdata"
	}
	if conf.IsStoreNameAutoAddDate() {
		return fmt.Sprint(name, "-", time.Now().Format("20060102")) // name-yyyymmdd
		// return fmt.Sprint(name, "-", time.Now().Format("200601021504")) // name-yyyymmddHHMM
	}
	return name
}

func StartwithsRune(str string, startstr string) bool {

	if startstr == "" || str == startstr {
		return true
	}

	strs := []rune(str)
	tmps := []rune(startstr)
	lens := len(strs)
	lentmp := len([]rune(tmps))
	if lens < lentmp {
		return false
	}

	for i := 0; i < lentmp; i++ {
		if tmps[i] != strs[i] {
			return false
		}
	}

	return true
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
	fileinf, err := os.ReadDir(path)
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

// 按M或G单位显示
func GetSizeInfo(size uint64) string {
	if size > 1024*1024*1024 {
		return fmt.Sprintf("%.1fG", float64(size)/1024/1024/1024)
	}
	return fmt.Sprintf("%.1fM", float64(size)/1024/1024)
}

// 当前日期加减天数后的yyyymmdd格式
func GetYyyymmdd(days int) string {
	return time.Now().AddDate(0, 0, days).Format("20060102")
}

// 判断文件是否存在
func IsExistFile(file string) bool {
	s, err := os.Stat(file)
	if err == nil {
		return !s.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return !s.IsDir()
}

// 判断文件夹是否存在
func IsExistDir(dir string) bool {
	s, err := os.Stat(dir)
	if err == nil {
		return s.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return s.IsDir()
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
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
