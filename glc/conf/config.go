/**
 * 系统配置
 * 1）统一管理系统的全部配置信息
 * 2）所有配置都有默认值以便直接使用
 * 3）所有配置都可以通过环境变量设定覆盖，方便自定义配置，方便容器化部署
 */
package conf

import (
	"runtime"
	"sort"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

var storeRoot string
var storeChanLength int
var maxIdleTime int
var storeNameAutoAddDate bool
var serverUrl string
var serverIp string
var serverPort string
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
var clusterMode bool
var clusterUrls []string
var enableBackup bool
var glcGroup string
var minioUrl string
var minioUser string
var minioPassword string
var minioBucket string
var enableUploadMinio bool
var goMaxProcess int

func init() {
	cmn.SetLogLevel(cmn.GetEnvStr("GLC_LOG_LEVEL", "INFO")) // 默认INFO级别日志
	UpdateConfigByEnv()

	// 在这个地方建目录，如果创建失败就比较难看，比如仅命令行查看版本的情景
	// // 自动判断创建目录
	// _, err := os.Stat(storeRoot)
	// if err != nil && os.IsNotExist(err) {
	// 	os.MkdirAll(storeRoot, 0766)
	// }

}

func UpdateConfigByEnv() {
	// 读取环境变量初始化配置，各配置都有默认值
	storeRoot = cmn.GetEnvStr("GLC_STORE_ROOT", "/glogcenter")                              // 存储根目录
	storeChanLength = cmn.GetEnvInt("GLC_STORE_CHAN_LENGTH", 64)                            // 存储通道长度
	maxIdleTime = cmn.GetEnvInt("GLC_MAX_IDLE_TIME", 180)                                   // 最大闲置时间（秒）,超过闲置时间将自动关闭，0时表示不关闭
	storeNameAutoAddDate = cmn.GetEnvBool("GLC_STORE_NAME_AUTO_ADD_DATE", true)             // 存储名是否自动添加日期（日志量大通常按日单位区分存储），默认true
	serverUrl = cmn.GetEnvStr("GLC_SERVER_URL", "")                                         // 服务URL，默认“”，集群配置时自动获取地址可能不对，可通过这个设定
	serverIp = cmn.GetEnvStr("GLC_SERVER_IP", "")                                           // 服务IP，默认“”，当“”时会自动获取
	serverPort = cmn.GetEnvStr("GLC_SERVER_PORT", "8080")                                   // web服务端口，默认“8080”
	contextPath = cmn.GetEnvStr("GLC_CONTEXT_PATH", "/glc")                                 // web服务contextPath
	enableSecurityKey = cmn.GetEnvBool("GLC_ENABLE_SECURITY_KEY", false)                    // web服务是否开启API秘钥校验，默认false
	headerSecurityKey = cmn.GetEnvStr("GLC_HEADER_SECURITY_KEY", "X-GLC-AUTH")              // web服务API秘钥的header键名
	securityKey = cmn.GetEnvStr("GLC_SECURITY_KEY", "glogcenter")                           // web服务API秘钥
	enableWebGzip = cmn.GetEnvBool("GLC_ENABLE_WEB_GZIP", false)                            // web服务是否开启Gzip
	enableAmqpConsume = cmn.GetEnvBool("GLC_ENABLE_AMQP_CONSUME", false)                    // 是否开启rabbitMq消费者接收日志
	amqpAddr = cmn.GetEnvStr("GLC_AMQP_ADDR", "")                                           // rabbitMq连接地址，例："amqp://user:password@ip:port/"
	amqpQueueName = cmn.GetEnvStr("GLC_AMQP_QUEUE_NAME", "glc-log-queue")                   // rabbitMq队列名
	amqpJsonFormat = cmn.GetEnvBool("GLC_AMQP_JSON_FORMAT", true)                           // rabbitMq消息文本是否为json格式，默认true
	saveDays = cmn.GetEnvInt("GLC_SAVE_DAYS", 180)                                          // 日志分仓时的保留天数(0~180)，0表示不自动删除，默认180天
	enableLogin = cmn.GetEnvBool("GLC_ENABLE_LOGIN", false)                                 // 是否开启用户密码登录，默认“false”
	username = cmn.GetEnvStr("GLC_USERNAME", "glc")                                         // 登录用户名，默认“glc”
	password = cmn.GetEnvStr("GLC_PASSWORD", "GLogCenter100%666")                           // 登录密码，默认“GLogCenter100%666”
	clusterMode = cmn.GetEnvBool("GLC_CLUSTER_MODE", false)                                 // 是否开启集群模式，默认false
	splitUrls(cmn.GetEnvStr("GLC_CLUSTER_URLS", ""))                                        // 从服务器地址，多个时逗号分开，默认“”
	enableBackup = cmn.GetEnvBool("GLC_ENABLE_BACKUP", false)                               // 是否开启备份，默认false
	glcGroup = cmn.GetEnvStr("GLC_GROUP", "default")                                        // 日志中心分组名，默认“default”
	minioUrl = cmn.GetEnvStr("GLC_MINIO_URL", "")                                           // MINIO地址，默认“”
	minioUser = cmn.GetEnvStr("GLC_MINIO_USER", "")                                         // MINIO用户名，默认“”
	minioPassword = cmn.GetEnvStr("GLC_MINIO_PASS", "")                                     // MINIO密码，默认“”
	minioBucket = cmn.GetEnvStr("GLC_MINIO_BUCKET", "")                                     // MINIO桶名，默认“”
	enableUploadMinio = cmn.GetEnvBool("GLC_ENABLE_UPLOAD_MINIO", false)                    // 是否开启上传备份至MINIO服务器，默认false
	goMaxProcess = getGoMaxProcessConf(cmn.GetEnvInt("GLC_GOMAXPROCS", runtime.NumCPU()-1)) // 是否开启上传备份至MINIO服务器，默认false
}

// 取配置： 服务URL，集群配置时自动获取地址可能不对，可通过环境变量“GLC_ENABLE_BACKUP”设定，默认“”
func GetGoMaxProcess() int {
	return goMaxProcess
}
func getGoMaxProcessConf(n int) int {
	if n < 1 {
		n = 1
	}
	if n > runtime.NumCPU() {
		n = runtime.NumCPU()
	}
	return n
}

// 取配置： 服务URL，集群配置时自动获取地址可能不对，可通过环境变量“GLC_ENABLE_BACKUP”设定，默认“”
func GetServerUrl() string {
	return serverUrl
}

// 取配置： 是否开启MINIO备份，可通过环境变量“GLC_ENABLE_BACKUP”设定，默认false
func IsEnableBackup() bool {
	return enableBackup
}

// 取配置： 日志中心分组名，可通过环境变量“GLC_GROUP”设定，默认“default”
func GetGlcGroup() string {
	return glcGroup
}

// 取配置： MINIO地址，可通过环境变量“GLC_MINIO_URL”设定，默认“”
func GetMinioUrl() string {
	return minioUrl
}

// 取配置： MINIO用户名，可通过环境变量“GLC_MINIO_USER”设定，默认“”
func GetMinioUser() string {
	return minioUser
}

// 取配置： MINIO密码，可通过环境变量“GLC_MINIO_PASS”设定，默认“”
func GetMinioPassword() string {
	return minioPassword
}

// 取配置： MINIO桶名，可通过环境变量“GLC_MINIO_BUCKET”设定，默认“”
func GetMinioBucket() string {
	return minioBucket
}

// 取配置： 是否开启上传备份至MINIO服务器，可通过环境变量“GLC_ENABLE_UPLOAD_MINIO”设定，默认false
func IsEnableUploadMinio() bool {
	return enableUploadMinio
}

// 取配置： 是否开启转发日志到其他GLC服务，可通过环境变量“GLC_CLUSTER_MODE”设定，默认false
func IsClusterMode() bool {
	return clusterMode
}

// 取配置： 从服务器地址，可通过环境变量“GLC_SLAVE_HOSTS”设定，默认“”
func GetClusterUrls() []string {
	return clusterUrls
}

func splitUrls(str string) {
	hosts := strings.Split(str, ";")
	for i := 0; i < len(hosts); i++ {
		url := strings.TrimSpace(hosts[i])
		if url != "" {
			clusterUrls = append(clusterUrls, url)
		}
	}

	// 倒序
	sort.Slice(clusterUrls, func(i, j int) bool {
		return clusterUrls[i] > clusterUrls[j]
	})
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

// 取配置： 服务IP，可通过环境变量“GLC_SERVER_IP”设定，默认值“”，自动获取
func GetServerIp() string {
	return serverIp
}

// 取配置： web服务端口，可通过环境变量“GLC_SERVER_PORT”设定，默认值“8080”
func GetServerPort() string {
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
