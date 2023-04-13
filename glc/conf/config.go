/**
 * 系统配置
 * 1）统一管理系统的全部配置信息
 * 2）所有配置都有默认值以便直接使用
 * 3）所有配置都可以通过环境变量设定覆盖，方便自定义配置，方便容器化部署
 */

package conf

import (
	"gopkg.in/yaml.v2"
	"os"
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
var userList []User
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

type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Config struct {
	StoreRoot            string   `json:"storeRoot" yaml:"storeRoot"`                       //  存储根目录
	StoreChanLength      int      `json:"storeChanLength" yaml:"storeChanLength"`           //  存储通道长度
	MaxIdleTime          int      `json:"maxIdleTime" yaml:"maxIdleTime"`                   //  最大闲置时间（秒）,超过闲置时间将自动关闭，0时表示不关闭
	StoreNameAutoAddDate bool     `json:"storeNameAutoAddDate" yaml:"storeNameAutoAddDate"` //  存储名是否自动添加日期（日志量大通常按日单位区分存储），默认true
	ServerUrl            string   `json:"serverUrl" yaml:"serverUrl"`                       //   服务URL，默认“”，集群配置时自动获取地址可能不对，可通过这个设定
	ServerIp             string   `json:"serverIp" yaml:"serverIp"`                         //  服务IP，默认“”，当“”时会自动获取
	ServerPort           string   `json:"serverPort" yaml:"serverPort"`                     //   web服务端口，默认“8080”
	ContextPath          string   `json:"contextPath" yaml:"contextPath"`                   //   web服务contextPath
	EnableSecurityKey    bool     `json:"enableSecurityKey" yaml:"enableSecurityKey"`       //    web服务是否开启API秘钥校验，默认false
	SecurityKey          string   `json:"securityKey" yaml:"securityKey"`                   //    web服务API秘钥的header键名
	HeaderSecurityKey    string   `json:"headerSecurityKey" yaml:"headerSecurityKey"`       //     web服务API秘钥
	EnableAmqpConsume    bool     `json:"enableAmqpConsume" yaml:"enableAmqpConsume"`       //     是否开启rabbitMq消费者接收日志
	EnableWebGzip        bool     `json:"enableWebGzip" yaml:"enableWebGzip"`               //     web服务是否开启Gzip
	AmqpAddr             string   `json:"amqpAddr" yaml:"amqpAddr"`                         //     rabbitMq连接地址，例："amqp://user:password@ip:port/"
	AmqpQueueName        string   `json:"amqpQueueName" yaml:"amqpQueueName"`               //     rabbitMq队列名
	AmqpJsonFormat       bool     `json:"amqpJsonFormat" yaml:"amqpJsonFormat"`             // rabbitMq消息文本是否为json格式，默认true
	SaveDays             int      `json:"saveDays" yaml:"saveDays"`                         // 日志分仓时的保留天数(0~180)，0表示不自动删除，默认180天
	EnableLogin          bool     `json:"enableLogin" yaml:"enableLogin"`                   // 是否开启用户密码登录，默认“false”
	User                 []User   `json:"user" yaml:"user"`                                 // 登录用户名密码，默认“admin，admin@@2023”
	ClusterMode          bool     `json:"clusterMode" yaml:"clusterMode"`                   // 是否开启集群模式，默认false
	ClusterUrls          []string `json:"clusterUrls" yaml:"clusterUrls"`                   // 从服务器地址，多个时逗号分开，默认“”
	EnableBackup         bool     `json:"enableBackup" yaml:"enableBackup"`                 // 是否开启备份，默认false
	GlcGroup             string   `json:"glcGroup" yaml:"glcGroup"`                         // 日志中心分组名，默认“default”
	MinioUrl             string   `json:"minioUrl" yaml:"minioUrl"`                         // MINIO地址，默认“”
	MinioUser            string   `json:"minioUser" yaml:"minioUser"`                       // MINIO用户名，默认“”
	MinioPassword        string   `json:"minioPassword" yaml:"minioPassword"`               // MINIO密码，默认“”
	MinioBucket          string   `json:"minioBucket" yaml:"minioBucket"`                   // MINIO桶名，默认“”
	EnableUploadMinio    bool     `json:"enableUploadMinio" yaml:"enableUploadMinio"`       // 是否开启上传备份至MINIO服务器，默认false
	GoMaxProcess         int      `json:"goMaxProcess" yaml:"goMaxProcess"`                 // 使用的最大CPU数量，默认是最大CPU数量（设定值不在实际数量范围是按最大看待）
	LogLevel             string   `json:"logLevel" yaml:"logLevel"`                         // 日志级别 默认INFO
}

func init() {
	var isNew bool
	// 默认INFO级别日志
	cmn.SetLogLevel("INFO")
	var setting Config
	config, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		cmn.Error("读取配置文件失败", err)
		isNew = true
	}
	var yamlerr interface{} = yaml.Unmarshal(config, &setting)
	if yamlerr != nil {
		cmn.Error("读取配置文件失败", err)
	}

	if setting.StoreRoot == "" {
		setting.StoreRoot = "/glogcenter"
		isNew = true
	}
	if setting.StoreChanLength == 0 {
		setting.StoreChanLength = 64
		isNew = true
	}
	if setting.MaxIdleTime == 0 {
		setting.MaxIdleTime = 180
		isNew = true
	}
	if setting.StoreNameAutoAddDate == false {
		setting.StoreNameAutoAddDate = true
		isNew = true
	}
	if setting.ServerPort == "" {
		setting.ServerPort = "8080"
		isNew = true
	}
	if setting.ContextPath == "" {
		setting.ContextPath = "/glc"
		isNew = true
	}
	if setting.HeaderSecurityKey == "" {
		setting.HeaderSecurityKey = "X-GLC-AUTH"
		isNew = true
	}
	if setting.SecurityKey == "" {
		setting.SecurityKey = "glogcenter"
		isNew = true
	}
	if setting.AmqpQueueName == "" {
		setting.AmqpQueueName = "glc-log-queue"
		isNew = true
	}
	if setting.AmqpJsonFormat == false {
		setting.AmqpJsonFormat = true
		isNew = true
	}
	if setting.SaveDays == 0 {
		setting.SaveDays = 180
		isNew = true
	}
	if setting.User == nil {
		setting.User = []User{
			{
				Username: "admin",
				Password: "admin@@2023",
			},
		}
		username = "admin"
		password = "admin@@2023"
		isNew = true
	}
	if setting.GlcGroup == "" {
		setting.GlcGroup = "default"
		isNew = true
	}
	if setting.GoMaxProcess == 0 || setting.GoMaxProcess > getGoMaxProcessConf(-1) {
		setting.GoMaxProcess = getGoMaxProcessConf(-1)
		isNew = true
	}
	if setting.LogLevel == "" {
		setting.LogLevel = "INFO"
		isNew = true
	}
	if isNew {
		data, err := yaml.Marshal(setting)
		if err != nil {
			cmn.Error("读取配置文件失败", err)
		}
		err = os.WriteFile("./config/config.yaml", data, 0777)
		if err != nil {
			cmn.Error("写入配置文件失败", err)
		}
	}
	cmn.SetLogLevel(setting.LogLevel)
	storeRoot = setting.StoreRoot
	storeChanLength = setting.StoreChanLength
	maxIdleTime = setting.MaxIdleTime
	storeNameAutoAddDate = setting.StoreNameAutoAddDate
	serverUrl = setting.ServerUrl
	serverIp = setting.ServerIp
	serverPort = setting.ServerPort
	contextPath = setting.ContextPath
	enableSecurityKey = setting.EnableSecurityKey
	securityKey = setting.SecurityKey
	headerSecurityKey = setting.HeaderSecurityKey
	enableAmqpConsume = setting.EnableAmqpConsume
	enableWebGzip = setting.EnableWebGzip
	amqpAddr = setting.AmqpAddr
	amqpQueueName = setting.AmqpQueueName
	amqpJsonFormat = setting.AmqpJsonFormat
	saveDays = setting.SaveDays
	enableLogin = setting.EnableLogin
	clusterMode = setting.ClusterMode
	clusterUrls = setting.ClusterUrls
	enableBackup = setting.EnableBackup
	glcGroup = setting.GlcGroup
	minioUrl = setting.MinioUrl
	minioUser = setting.MinioUser
	minioPassword = setting.MinioPassword
	minioBucket = setting.MinioBucket
	enableUploadMinio = setting.EnableUploadMinio
	goMaxProcess = setting.GoMaxProcess
	userList = setting.User
	for _, user := range userList {
		username = user.Username
		password = user.Password
	}
}

// GetGoMaxProcess 取配置： 使用的最大CPU数量，可通过环境变量“GLC_GOMAXPROCS”设定，默认最大CPU数量
func GetGoMaxProcess() int {
	return goMaxProcess
}
func getGoMaxProcessConf(n int) int {
	max := runtime.NumCPU()
	if n < 1 || n > max {
		n = max
	}
	return n
}

// GetServerUrl 取配置： 服务URL，集群配置时自动获取地址可能不对，可通过环境变量“GLC_ENABLE_BACKUP”设定，默认“”
func GetServerUrl() string {
	return serverUrl
}

// IsEnableBackup 取配置： 是否开启MINIO备份，可通过环境变量“GLC_ENABLE_BACKUP”设定，默认false
func IsEnableBackup() bool {
	return enableBackup
}

// GetGlcGroup 取配置： 日志中心分组名，可通过环境变量“GLC_GROUP”设定，默认“default”
func GetGlcGroup() string {
	return glcGroup
}

// GetMinioUrl 取配置： MINIO地址，可通过环境变量“GLC_MINIO_URL”设定，默认“”
func GetMinioUrl() string {
	return minioUrl
}

// GetMinioUser 取配置： MINIO用户名，可通过环境变量“GLC_MINIO_USER”设定，默认“”
func GetMinioUser() string {
	return minioUser
}

// GetMinioPassword 取配置： MINIO密码，可通过环境变量“GLC_MINIO_PASS”设定，默认“”
func GetMinioPassword() string {
	return minioPassword
}

// GetMinioBucket 取配置： MINIO桶名，可通过环境变量“GLC_MINIO_BUCKET”设定，默认“”
func GetMinioBucket() string {
	return minioBucket
}

// IsEnableUploadMinio 取配置： 是否开启上传备份至MINIO服务器，可通过环境变量“GLC_ENABLE_UPLOAD_MINIO”设定，默认false
func IsEnableUploadMinio() bool {
	return enableUploadMinio
}

// IsClusterMode 取配置： 是否开启转发日志到其他GLC服务，可通过环境变量“GLC_CLUSTER_MODE”设定，默认false
func IsClusterMode() bool {
	return clusterMode
}

// GetClusterUrls 取配置： 从服务器地址，可通过环境变量“GLC_SLAVE_HOSTS”设定，默认“”
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

// IsEnableLogin 取配置： 是否开启用户密码登录，可通过环境变量“GLC_ENABLE_LOGIN”设定，默认“false”
func IsEnableLogin() bool {
	return enableLogin
}

// GetUsername 取配置： 登录用户名，可通过环境变量“GLC_USERNAME”设定，默认“glc”
func GetUsername() string {
	return username
}

// GetPassword 取配置： 登录用户名，可通过环境变量“GLC_PASSWORD”设定，默认“glogcenter”
func GetPassword() string {
	return password
}

// GetSaveDays 取配置： 日志分仓时的保留天数(0~180)，0表示不自动删除，可通过环境变量“GLC_SAVE_DAYS”设定，默认180天
func GetSaveDays() int {
	return saveDays
}

// IsAmqpJsonFormat 取配置： rabbitMq消息文本是否为json格式，可通过环境变量“GLC_AMQP_JSON_FORMAT”设定，默认值“true”
func IsAmqpJsonFormat() bool {
	return amqpJsonFormat
}

// GetAmqpQueueName 取配置： rabbitMq连接地址，可通过环境变量“GLC_AMQP_ADDR”设定，默认值“”
func GetAmqpQueueName() string {
	return amqpQueueName
}

// GetAmqpAddr 取配置： rabbitMq连接地址，可通过环境变量“GLC_AMQP_ADDR”设定，默认值“”
func GetAmqpAddr() string {
	return amqpAddr
}

// IsEnableAmqpConsume 取配置： 是否开启rabbitMq消费者接收日志，可通过环境变量“GLC_ENABLE_AMQP_CONSUME”设定，默认值“false”
func IsEnableAmqpConsume() bool {
	return enableAmqpConsume
}

// IsEnableWebGzip 取配置： web服务API秘钥的header键名，可通过环境变量“GLC_HEADER_SECURITY_KEY”设定，默认值“X-GLC-AUTH”
func IsEnableWebGzip() bool {
	return enableWebGzip
}

// IsEnableSecurityKey 取配置： web服务API秘钥的header键名，可通过环境变量“GLC_HEADER_SECURITY_KEY”设定，默认值“X-GLC-AUTH”
func IsEnableSecurityKey() bool {
	return enableSecurityKey
}

// GetHeaderSecurityKey 取配置： web服务API秘钥的header键名，可通过环境变量“GLC_HEADER_SECURITY_KEY”设定，默认值“X-GLC-AUTH”
func GetHeaderSecurityKey() string {
	return headerSecurityKey
}

// GetSecurityKey 取配置： web服务API秘钥，可通过环境变量“GLC_SECURITY_KEY”设定，默认值“glogcenter”
func GetSecurityKey() string {
	return securityKey
}

// GetContextPath 取配置： web服务端口，可通过环境变量“GLC_CONTEXT_PATH”设定，默认值“8080”
func GetContextPath() string {
	return contextPath
}

// GetServerIp 取配置： 服务IP，可通过环境变量“GLC_SERVER_IP”设定，默认值“”，自动获取
func GetServerIp() string {
	return serverIp
}

// GetServerPort 取配置： web服务端口，可通过环境变量“GLC_SERVER_PORT”设定，默认值“8080”
func GetServerPort() string {
	return serverPort
}

// GetStorageRoot 取配置：存储根目录，可通过环境变量“GLC_STORE_ROOT”设定，默认值“/glogcenter”
func GetStorageRoot() string {
	return storeRoot
}

// GetStoreChanLength 取配置：存储通道长度，可通过环境变量“GLC_STORE_CHAN_LENGTH”设定，默认值“64”
func GetStoreChanLength() int {
	return storeChanLength
}

// GetMaxIdleTime 取配置：最大闲置时间（秒），可通过环境变量“GLC_MAX_IDLE_TIME”设定，默认值“180”，超过闲置时间将自动关闭存储器，0时表示不关闭
func GetMaxIdleTime() int {
	return maxIdleTime
}

// IsStoreNameAutoAddDate 取配置：存储名是否自动添加日期（日志量大通常按日单位区分存储），可通过环境变量“GLC_STORE_NAME_AUTO_ADD_DATE”设定，默认值“true”
func IsStoreNameAutoAddDate() bool {
	return storeNameAutoAddDate
}

// GetUserList 用户列表
func GetUserList() []User {
	return userList
}
