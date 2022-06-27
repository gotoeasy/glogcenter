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
var serverPort int
var contextPath string
var enableSecurityKey bool
var securityKey string
var headerSecurityKey string

func init() {
	UpdateConfigByEnv()
}

func UpdateConfigByEnv() {
	// 读取环境变量初始化配置，各配置都有默认值
	storeRoot = Getenv("GLC_STORE_ROOT", "/glogcenter")                     // 存储根目录
	storeChanLength = GetenvInt("GLC_STORE_CHAN_LENGTH", 64)                // 存储通道长度
	maxIdleTime = GetenvInt("GLC_MAX_IDLE_TIME", 180)                       // 最大闲置时间（秒）,超过闲置时间将自动关闭，0时表示不关闭
	storeNameAutoAddDate = GetenvBool("GLC_STORE_NAME_AUTO_ADD_DATE", true) // 存储名是否自动添加日期（日志量大通常按日单位区分存储），默认true
	serverPort = GetenvInt("GLC_SERVER_PORT", 8080)                         // web服务端口
	contextPath = Getenv("GLC_CONTEXT_PATH", "/glc")                        // web服务contextPath
	enableSecurityKey = GetenvBool("GLC_ENABLE_SECURITY_KEY", false)        // web服务是否开启API秘钥校验，默认false
	headerSecurityKey = Getenv("GLC_HEADER_SECURITY_KEY", "X-GLC-AUTH")     // web服务API秘钥的header键名
	securityKey = Getenv("GLC_SECURITY_KEY", "glogcenter")                  // web服务API秘钥
}

// 取配置： web服务API秘钥的header键名，可通过环境变量“GLC_HEADER_SECURITY_KEY”设定，默认值“X-GLC-AUTH”
func IsEnableSecurityKey() bool {
	return enableSecurityKey
}

// 取配置： web服务API秘钥的header键名，可通过环境变量“GLC_HEADER_SECURITY_KEY”设定，默认值“X-GLC-AUTH”
func GetHeaderSecurityKey() string {
	return headerSecurityKey
}

// 取配置： web服务API秘钥，可通过环境变量“GLC_SECURITY_KEY”设定，默认值“glogcenter”
func GetSecurityKey() string {
	return securityKey
}

// 取配置： web服务端口，可通过环境变量“GLC_CONTEXT_PATH”设定，默认值“8080”
func GetContextPath() string {
	return contextPath
}

// 取配置： web服务端口，可通过环境变量“GLC_SERVER_PORT”设定，默认值“8080”
func GetServerPort() int {
	return serverPort
}

// 取配置：存储根目录，可通过环境变量“GLC_STORE_ROOT”设定，默认值“/glogcenter”
func GetStorageRoot() string {
	return storeRoot
}

// 取配置：存储通道长度，可通过环境变量“GLC_STORE_CHAN_LENGTH”设定，默认值“64”
func GetStoreChanLength() int {
	return storeChanLength
}

// 取配置：最大闲置时间（秒），可通过环境变量“GLC_MAX_IDLE_TIME”设定，默认值“180”，超过闲置时间将自动关闭存储器，0时表示不关闭
func GetMaxIdleTime() int {
	return maxIdleTime
}

// 取配置：存储名是否自动添加日期（日志量大通常按日单位区分存储），可通过环境变量“GLC_STORE_NAME_AUTO_ADD_DATE”设定，默认值“true”
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
