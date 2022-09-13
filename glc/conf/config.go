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
var enableAmqpConsume bool
var enableWebGzip bool
var amqpAddr string
var amqpQueueName string
var amqpJsonFormat bool
var saveDays int
var enableLogin bool
var username string
var password string

func init() {
	UpdateConfigByEnv()

	// 自动判断创建目录
	_, err := os.Stat(storeRoot)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(storeRoot, 0766)
	}

}

func UpdateConfigByEnv() {
	// 读取环境变量初始化配置，各配置都有默认值
	storeRoot = Getenv("GLC_STORE_ROOT", "/glogcenter")                     // 存储根目录
	storeChanLength = GetenvInt("GLC_STORE_CHAN_LENGTH", 64)                // 存储通道长度
	maxIdleTime = GetenvInt("GLC_MAX_IDLE_TIME", 180)                       // 最大闲置时间（秒）,超过闲置时间将自动关闭，0时表示不关闭
	storeNameAutoAddDate = GetenvBool("GLC_STORE_NAME_AUTO_ADD_DATE", true) // 存储名是否自动添加日期（日志量大通常按日单位区分存储），默认true
	serverPort = GetenvInt("GLC_SERVER_PORT", 18080)                        // web服务端口，默认18080
	contextPath = Getenv("GLC_CONTEXT_PATH", "/glc")                        // web服务contextPath
	enableSecurityKey = GetenvBool("GLC_ENABLE_SECURITY_KEY", false)        // web服务是否开启API秘钥校验，默认false
	headerSecurityKey = Getenv("GLC_HEADER_SECURITY_KEY", "X-GLC-AUTH")     // web服务API秘钥的header键名
	securityKey = Getenv("GLC_SECURITY_KEY", "glogcenter")                  // web服务API秘钥
	enableWebGzip = GetenvBool("GLC_ENABLE_WEB_GZIP", true)                 // web服务是否开启Gzip
	enableAmqpConsume = GetenvBool("GLC_ENABLE_AMQP_CONSUME", false)        // 是否开启rabbitMq消费者接收日志
	amqpAddr = Getenv("GLC_AMQP_ADDR", "")                                  // rabbitMq连接地址，例："amqp://user:password@ip:port/"
	amqpQueueName = Getenv("GLC_AMQP_QUEUE_NAME", "glc-log-queue")          // rabbitMq队列名
	amqpJsonFormat = GetenvBool("GLC_AMQP_JSON_FORMAT", true)               // rabbitMq消息文本是否为json格式，默认true
	saveDays = GetenvInt("GLC_SAVE_DAYS", 180)                              // 日志分仓时的保留天数(0~180)，0表示不自动删除，默认180天
	enableLogin = GetenvBool("GLC_ENABLE_LOGIN", false)                     // 是否开启用户密码登录，默认“false”
	username = Getenv("GLC_USERNAME", "glc")                                // 登录用户名，默认“glc”
	password = Getenv("GLC_PASSWORD", "glogcenter")                         // 登录密码，默认“glogcenter”
}

// 取配置： 是否开启用户密码登录，可通过环境变量“GLC_ENABLE_LOGIN”设定，默认“false”
func IsEnableLogin() bool {
	return enableLogin
}

// 取配置： 登录用户名，可通过环境变量“GLC_USERNAME”设定，默认“glc”
func GetUsername() string {
	return username
}

// 取配置： 登录用户名，可通过环境变量“GLC_PASSWORD”设定，默认“glogcenter”
func GetPassword() string {
	return password
}

// 取配置： 日志分仓时的保留天数(0~180)，0表示不自动删除，可通过环境变量“GLC_SAVE_DAYS”设定，默认180天
func GetSaveDays() int {
	return saveDays
}

// 取配置： rabbitMq消息文本是否为json格式，可通过环境变量“GLC_AMQP_JSON_FORMAT”设定，默认值“true”
func IsAmqpJsonFormat() bool {
	return amqpJsonFormat
}

// 取配置： rabbitMq连接地址，可通过环境变量“GLC_AMQP_ADDR”设定，默认值“”
func GetAmqpQueueName() string {
	return amqpQueueName
}

// 取配置： rabbitMq连接地址，可通过环境变量“GLC_AMQP_ADDR”设定，默认值“”
func GetAmqpAddr() string {
	return amqpAddr
}

// 取配置： 是否开启rabbitMq消费者接收日志，可通过环境变量“GLC_ENABLE_AMQP_CONSUME”设定，默认值“false”
func IsEnableAmqpConsume() bool {
	return enableAmqpConsume
}

// 取配置： web服务API秘钥的header键名，可通过环境变量“GLC_HEADER_SECURITY_KEY”设定，默认值“X-GLC-AUTH”
func IsEnableWebGzip() bool {
	return enableWebGzip
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
	if strings.ToLower(s) == "false" {
		return false
	}
	return defaultValue
}
