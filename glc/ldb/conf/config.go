/**
 * 系统配置
 * 1）统一管理系统的全部配置信息
 * 2）所有配置都有默认值以便直接使用
 * 3）所有配置都可以通过环境变量设定覆盖，方便自定义配置，方便容器化部署
 */
package conf

import (
	"os"
	"strconv"
	"strings"
)

var storeRoot string
var storeChanLength int
var maxIdleTime int
var storeNameAutoAddDate bool

func init() {
	// 读取环境变量初始化配置，各配置都有默认值
	storeRoot = Getenv("STORE_ROOT", "e:/222")                          // 存储根目录
	storeChanLength = GetenvInt("STORE_CHAN_LENGTH", 64)                // 存储通道长度
	maxIdleTime = GetenvInt("MAX_IDLE_TIME", 10)                        // 最大闲置时间（秒）,超过闲置时间将自动关闭，0时表示不关闭
	storeNameAutoAddDate = GetenvBool("STORE_NAME_AUTO_ADD_DATE", true) // 存储名是否自动添加日期（日志量大通常按日单位区分存储），默认true
}

// 取配置：存储根目录，可通过环境变量“STORE_ROOT”设定，默认值“/glogcenter”
func GetStorageRoot() string {
	return storeRoot
}

// 取配置：存储通道长度，可通过环境变量“STORE_CHAN_LENGTH”设定，默认值“64”
func GetStoreChanLength() int {
	return storeChanLength
}

// 取配置：最大闲置时间（秒），可通过环境变量“MAX_IDLE_TIME”设定，默认值“300”，超过闲置时间将自动关闭存储器，0时表示不关闭
func GetMaxIdleTime() int {
	return maxIdleTime
}

// 取配置：存储名是否自动添加日期（日志量大通常按日单位区分存储），可通过环境变量“STORE_NAME_AUTO_ADD_DATE”设定，默认值“true”
func IsStoreNameAutoAddDate() bool {
	return storeNameAutoAddDate
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
