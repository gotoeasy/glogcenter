<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/logo.png"/>
</div>


# 缘起

昔者，吾日志中心之事，恒用传统三件套，曰`ELK`，然岁月流转，如行舟江河，难免潮涌沉浮，不维其盛，则不足以久存。故见诸多不如意，如定制之难，如索引之弗易，如初启之迟缓如蜗牛，如操作之陌生如他邦，如资甚之贪婪如饕餮之兽，如崩溃之险如山崩川竭，疑难重重，如堆积之沙石。<br>
<br>
终有一日，愈发执志，以`go`之巧工，铸造新日志中心，其表现多舛，实足以令人惊艳，是以，赋名曰`glogcenter`，亦称`GLC`，开仓建库。<br>
<br>
当下，架库之作已可窥见，与君共享。<br>
`（以上由GPT编辑）`
<br>

<p align="center">
    <a href="https://golang.google.cn"><img src="https://img.shields.io/badge/golang-1.23.1-brightgreen.svg"></a>
    <a href="https://hub.docker.com/r/gotoeasy/glc"><img src="https://img.shields.io/docker/pulls/gotoeasy/glc"></a>
    <a href="https://github.com/gotoeasy/glogcenter/releases/latest"><img src="https://img.shields.io/github/release/gotoeasy/glogcenter.svg"></a>
    <a href="https://github.com/gotoeasy/glogcenter/blob/master/LICENSE"><img src="https://img.shields.io/github/license/gotoeasy/glogcenter"></a>
<p>

<br>
<br>
国外仓库地址： https://github.com/gotoeasy/glogcenter <br>
国内（同步）： https://gitee.com/gotoeasy/glogcenter
<br>
<br>
演示地址(网络可能不稳定)： https://glc.gotoeasy.top
<br>
<br>

## 特点
- [x] 使用golang实现，具备go的各种特性优势，特别`【节省资源】`
- [x] 基于LSMT存储，融合日志写多读少的特点进行实现，毫秒级查询响应 `【性能卓越】`
- [x] 日志吞食量每秒近万条，闲时建索引每秒数千条，满足大多项目场景 `【广泛适用】`
- [x] 支持多关键词全文检索，支持多维度线索查询，支持定位相邻查询 `【功能丰富】`
- [x] 内置提供VUE实现的查询管理界面，页面简洁大方，操作习惯自然 `【体验优秀】`
- [x] 提供docker镜像，支持容器化部署，支持个性化环境变量配置 `【部署方便】`
- [x] 提供java/go/python等客户端工具包，日志收集信手拈来 `【集成简单】`
- [x] 支持登录验证，秘钥校验，权限控制，黑白名单等安全设定 `【安全可靠】`
- [x] 支持多服务集群模式部署，提供服务稳定性保障、数据冗余性保障 `【高可用保障】`
- [x] 上至央企大项目下至本地调试，已历经众多案例磨炼，表现出色 `【生产级别品质】`


## 概要图
<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/glcsummary.png"/>
</div>


## 部分截图
<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/glogcenter.png"/>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/storage.png"/>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/users.png"/>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/chatai.png"/>
</div>

<br>

## `docker`单机部署模式简易示例
```shell
# 快速体验（其中通过GLC_TEST_MODE=true开启测试模式，页面上会显示生成测试数据的按钮，方便测试或快速体验）
docker run -d -p 8080:8080 -e GLC_TEST_MODE=true gotoeasy/glc

# 外挂数据目录
docker run -d -p 8080:8080 -v /glc:/glogcenter gotoeasy/glc

# 【简易用法】
# 启动成功后即可按 http://ip:port 访问
# 支持多关键词检索，比如输入【key1、key2、key3】检索出同时满足这3个关键词的结果
# 默认每次检索100条，滚动到底部时自动检索后面100条
```

## `docker`集群部署模式简易示例
```shell
# 以下3台以集群方式启动，配置本节点地址及关联节点地址即可
# 采用“乐观集群”方式，简易选主（简单排序）+日志群发（忽略失败）+数据补偿（隔日同步历史数据）

# 服务1
docker run -d -p 8080:8080 -e GLC_CLUSTER_MODE=true -e GLC_SERVER_URL=http://172.27.59.51:8080 \
  -e GLC_CLUSTER_URLS=http://172.27.59.51:8080;http://172.27.59.52:8080;http://172.27.59.53:8080 \
  gotoeasy/glc
# 服务2
docker run -d -p 8080:8080 -e GLC_CLUSTER_MODE=true -e GLC_SERVER_URL=http://172.27.59.52:8080 \
  -e GLC_CLUSTER_URLS=http://172.27.59.51:8080;http://172.27.59.52:8080;http://172.27.59.53:8080 \
  gotoeasy/glc
# 服务3
docker run -d -p 8080:8080 -e GLC_CLUSTER_MODE=true -e GLC_SERVER_URL=http://172.27.59.53:8080 \
  -e GLC_CLUSTER_URLS=http://172.27.59.51:8080;http://172.27.59.52:8080;http://172.27.59.53:8080 \
  gotoeasy/glc
```


## `docker`启动环境变量
- [x] `GLC_STORE_NAME_AUTO_ADD_DATE`日志仓是否自动按日存储，默认`true`
- [x] `GLC_SAVE_DAYS`日志仓按日存储自动维护时的保留天数(有效范围`0~1200`)，`0`表示不自动删除，默认`180`天
- [x] `GLC_SEARCH_MULIT_LINE`，是否对日志列的全部行进行索引检索，默认`false`仅第一行
- [x] `GLC_ENABLE_LOGIN`是否开启用户密码登录功能，默认`false`
- [x] `GLC_USERNAME`管理员用户名，默认`glc`，从`0.13.0`版本开始，管理员有新增用户及权限管理功能，并且有全部系统的查询权限
- [x] `GLC_PASSWORD`管理员密码，默认`GLogCenter100%666`
- [x] `GLC_TOKEN_SALT`用以生成令牌的字符串令牌盐，开启登录功能时建议设定提高安全性，默认空白
- [x] `GLC_ENABLE_SECURITY_KEY`日志添加的接口是否开启API秘钥校验，默认`false`
- [x] `GLC_HEADER_SECURITY_KEY`API秘钥的`header`键名，默认`X-GLC-AUTH`
- [x] `GLC_SECURITY_KEY`API秘钥，默认`glogcenter`
- [x] `GLC_ENABLE_CORS`是否允许跨域，默认`false`
- [x] `GLC_PAGE_SIZE`每次检索件数，默认`100`（有效范围`1~1000`）
- [x] `GLC_ENABLE_WEB_GZIP`网页服务是否开启压缩，默认`false`
- [x] `GLC_ENABLE_AMQP_CONSUME`是否开启`rabbitMq`消费者接收日志，默认`false`
- [x] `GLC_AMQP_ADDR`消息队列`rabbitMq`连接地址，例：`amqp://user:password@ip:port/`，默认空白
- [x] `GLC_AMQP_JSON_FORMAT`消息队列`rabbitMq`消息文本是否为`json`格式，默认`true`
- [x] `GLC_CLUSTER_MODE`是否集群模式启动，默认`false`
- [x] `GLC_SERVER_URL`集群模式时的本节点服务地址，默认空白
- [x] `GLC_CLUSTER_URLS`集群模式时的关联节点服务地址，多个时`;`分隔，默认空白
- [x] `GLC_GOMAXPROCS`使用最大CPU数量，值不在实际范围内时按最大值看待，默认最大值，常用于`docker`方式
- [x] `GLC_TEST_MODE`是否开启测试模式，开启时显示生成测试数据的按钮，供测试或快速体验用，默认`false`
- [x] `GLC_WHITE_LIST`白名单，多个用逗号分隔，黑白名单冲突时白名单优先，默认空白。可设定IP，最后段支持通配符，如`1.2.3.*`，内网IP默认都是白名单不必设定，实验性质的支持区域名称（因为IP地域查询可能有误差），如`上海市,深圳市`
- [x] `GLC_BLACK_LIST`黑名单，多个用逗号分隔，黑白名单冲突时白名单优先，默认空白。可设定IP，最后段支持通配符，如`1.2.3.*`，也支持单个通配符`*`代表全部（也就是只允许内网或白名单指定使用），实验性质的支持区域名称（因为IP地域查询可能有误差）
- [x] `GLC_IP_ADD_CITY`对IP字段是否自动附加城市信息，默认`false`
- [x] `GLC_NEAR_SEARCH_SIZE`定位相邻检索时的检索件数，默认200，有效范围50-1000
- [x] `GLC_ENABLE_CHATAI`是否开启GLC智能助手，默认true，会在菜单栏显示


## 接口
- [x] `/glc/v1/log/add`日志添加，`POST`，`application/json` <br>
      字段`system`： 字符串，对应页面的`系统名` <br>
      字段`date`： 字符串，对应页面的`日期时间`，格式`yyyy-MM-dd HH:mm:ss.SSS` <br>
      字段`text`： 字符串，对应页面的`日志` <br>
      字段`servername`： 字符串，对应页面的`主机名` <br>
      字段`serverip`： 字符串，对应页面的`主机IP` <br>
      字段`loglevel`： 字符串，对应页面的`日志级别` <br>
      字段`traceid`： 字符串，对应页面的`追踪码` <br>
      字段`clientip`： 字符串，对应页面的`客户端IP` <br>
      字段`user`： 字符串，对应页面的`用户` <br>

```shell
# 发送测试数据的参考脚本
# 注意时间格式要一致，否则按时间范围检索可能无法得到预想结果
curl -X POST -d '{"system":"demo", "date":"2023-01-01 01:02:03.456","text":"demo log text"}' \
     -H "Content-Type:application/json" http://127.0.0.1:8080/glc/v1/log/add
```


## 使用`logback`的`java`项目，支持日志收集，确保主次版本和GLC版本一致
```xml
<!-- pom坐标 -->
<dependency>
    <groupId>top.gotoeasy</groupId>
    <artifactId>glc-logback-appender</artifactId>
    <version>0.17.1</version>
</dependency>
```

```xml
<!-- logback配置例子1，发送至 glogcenter -->
<appender name="GLC" class="top.gotoeasy.framework.glc.logback.appender.GlcHttpJsonAppender">
    <glcApiUrl>http://127.0.0.1:8080/</glcApiUrl> <!--可通过环境变量 GLC_API_URL 设定-->
    <glcApiKey>X-GLC-AUTH:glogcenter</glcApiKey>  <!--可通过环境变量 GLC_API_KEY 设定-->
    <system>demo</system>                         <!--可通过环境变量 GLC_SYSTEM 设定 -->
    <layout>
        <pattern><![CDATA[%m %n]]></pattern>
    </layout>
</appender>
```

```xml
<!-- logback配置例子2，发送至 rabbitmq -->
<appender name="GLC" class="top.gotoeasy.framework.glc.logback.appender.GlcAmqpAppender">
    <amqpHost>127.0.0.1</amqpHost>                <!-- 可通过环境变量 GLC_AMQP_HOST 设定 -->
    <amqpPort>5672</amqpPort>                     <!-- 可通过环境变量 GLC_AMQP_PORT 设定 -->
    <amqpUser>rabbitmqUsername</amqpUser>         <!-- 可通过环境变量 GLC_AMQP_USER 设定 -->
    <amqpPassword>rabbitmqPassword</amqpPassword> <!-- 可通过环境变量 GLC_AMQP_PASSWORD 设定 -->
    <system>Demo</system>                         <!-- 可通过环境变量 GLC_SYSTEM 设定 -->
    <layout>
        <pattern><![CDATA[%m %n]]></pattern>
    </layout>
</appender>
```

```xml
<!-- 一个简单的logback-spring.xml配置例子 -->
<?xml version="1.0" encoding="UTF-8"?>
<configuration debug="false">
    <appender name="CONSOLE"
        class="ch.qos.logback.core.ConsoleAppender">
        <Target>System.out</Target>
        <encoder>
            <pattern>%d-%c-%t-%5p: %m%n</pattern>
        </encoder>
    </appender>

    <appender name="GLC" class="top.gotoeasy.framework.glc.logback.appender.GlcHttpJsonAppender">
        <glcApiUrl>http://127.0.0.1:8080/glc/v1/log/add</glcApiUrl>
        <system>demo</system>
        <layout>
            <pattern><![CDATA[%m %n]]></pattern>
        </layout>
    </appender>

    <root level="DEBUG">
        <appender-ref ref="CONSOLE" />
        <appender-ref ref="GLC" />
    </root>
</configuration>
```

## 使用`golang`语言的项目，提供工具包，开箱即用
```shell
# 方式1）通过环境变量自动配置，程序直接使用cmn.Debug(...)写日志即可
export GLC_ENABLE=true # 此配置默认false，要发送日志中心必须配置为true
export GLC_ENABLE_CONSOLE_LOG=true # 默认true，控制台不打印时配置为false
export GLC_API_URL='http://127.0.0.1:8080/' # 未配置时将取消发送
export GLC_API_KEY='X-GLC-AUTH:glogcenter' # 这是默认值，按需修改
export GLC_SYSTEM=default  # 默认default，按需修改
export GLC_LOG_LEVEL=debug # 日志级别（debug/info/warn/error）
export GLC_TRACE_ID=12345  # 默认空，跨进程调用等一些特殊场景使用
export GLC_PRINT_SRC_LINE=true # 是否打印源码行号，go语言专用，默认false
```

```golang
// 方式2） 使用前通过程序cmn.SetGlcClient(...)手动配置初始化
import "github.com/gotoeasy/glang/cmn"

func main() {
    // 这里用手动初始化替代环境变量自动配置方式，更多选项详见GlcOptions字段说明
    cmn.SetGlcClient(cmn.NewGlcClient(&cmn.GlcOptions{
        ApiUrl:      "http://127.0.0.1:8080/",
        Enable:      "true",
    }))

    cmn.Debug("这是Debug级别日志")
    cmn.Info("这是Info级别日志", "多个参数", "会被拼接")
    gd := &cmn.GlcData{TraceId: "1234567890"} // 跟踪码相同的日志，传入该参数即可
    cmn.Warn("这里的GlcData类型参数都不会打印", "gd只起传值作用", gd)
    cmn.Error("gd参数顺序无关", gd, "用法如同log库，但对GlcData做了特殊的判断处理")
    cmn.WaitGlcFinish() // 停止接收新日志，等待日志都发送完成，常在退出前调用
}
```


## 使用`python`语言的项目，提供工具包，开箱即用
```shell
# 支持以下环境变量配置
export GLC_ENABLE=true # 默认false，要发送日志中心必须配置为true
export GLC_ENABLE_CONSOLE_LOG=true # 默认true，控制台不打印时配置为false
export GLC_API_URL='http://127.0.0.1:8080/' # 未配置时将取消发送
export GLC_API_KEY='X-GLC-AUTH:glogcenter' # 这是默认值，按需修改
export GLC_SYSTEM=default  # 默认default，按需修改
export GLC_LOG_LEVEL=debug # 日志级别（debug/info/warn/error），默认debug
export GLC_TRACE_ID=12345  # 默认空，跨进程调用等一些特殊场景使用
```

```python
# 安装
pip install glogcenter

# 使用
from glogcenter import glc

glc.debug("这是Debug级别日志")
glc.info("这是Info级别日志", "多个参数", "会被拼接")
gd = glc.GlcData()
gd.user = 'abcd'
glc.warn("这里的GlcData类型参数都不会打印", "gd只起传值作用", gd)
glc.error("gd参数顺序无关", gd, "用法如同log库，但对GlcData做了特殊的判断处理")
```


## 支持零侵入收集docker容器日志 (适用`0.17.0`及以上版本)
```shell
# 1) 使用 fluentd 收集日志（为啥？因为较高版本docker已默认支持）
# 本仓库中 fluent.conf 是简单配置示意，其中包含转发日志到GLC
# 官方镜像的时区不合适，懒得改可直接用 gotoeasy/fluentd:v1.17-1-zh 替代
docker run -d -p 24224:24224 -p 24224:24224/udp \
       -v ./fluent.conf:/fluentd/etc/fluent.conf fluentd:v1.17-1

# 2) 运行容器时指定日志驱动，指向 fluentd 服务端口
docker run -d -p --log-driver=fluentd --log-opt fluentd-address=192.168.169.170:24224 <你的镜像>

# 已经搞定啦，fluentd会把日志发到GLC （这个必开就不用说了）
# 接下来，去折腾 fluent.conf 就行，举一反三，但凡 fluentd 支持收集的东西都可以框进来

# 当然，这种日志在显示上有一定不足，但瑕不掩瑜，有时这么做还是很值得的
```



## 更新履历

### 开发版`latest`

- [ ] 日志审计、告警


### 版本`0.17.4`

- [x] 修复已知问题（#63）

### 版本`0.17.3`

- [x] 优化修复一些已知问题（#57 #58 等）

### 版本`0.17.2`

- [x] 修复：条件检索可能存在个别数据查不到
- [x] 取消用户词典参数`GLC_DICT_DIR`的支持，配置不当会影响分词影响检索结果令人困惑，得不偿失
- [x] 镜像缩小等优化


### 版本`0.17.1`

- [x] 支持自定义分词字典
- [x] 增加字典目录环境变量`GLC_DICT_DIR`，支持多个`*.txt`字典文件。比如环境变量设定为`/opt`，启动时使用`-v /your-dict-dir:/opt`映射好字典目录就行
- [x] 一些细节优化

### 版本`0.17.0`

- [x] 零侵入支持`docker`容器日志、文件等各种日志的收集
- [x] 增加接口 `/glc/v1/log/addBatch`，支持一次接收多条日志
- [x] 升级使用`Go1.23.1`进行编译


<details>
<summary><strong><mark>更多历史版本更新履历</mark></strong></summary> 

### 版本`0.16.0`

- [x] 分词优化
- [x] 大幅提升建索引速度（强烈推荐使用固态硬盘）

### 版本`0.15.2`

- [x] 支持#44： 方便ngnix目录方式代理
- [x] 部分瑕疵修复

### 版本`0.15.1`

- [x] 新增`GLC智能助手`，可以随时解答日志中心的相关问题
- [x] 修复#38：系统名大小写可能引起筛选不起作用


### 版本`0.15.0`

- [x] 新增`定位相邻检索`功能，非常适合快速定位进行上下文查询的场景。相应增加`GLC_NEAR_SEARCH_SIZE`参数，配置定位相邻检索时的检索件数，默认200，有效范围50-1000
- [x] 新增Python客户端工具包，方便Python语言的项目接入日志中心
- [x] 修复一些已知问题


### 版本`0.14.2`

- [x] 修复#32 隐蔽的权限控制错误问题

### 版本`0.14.1`

- [x] 新增`GLC_IP_ADD_CITY`环境变量，对IP字段是否自动附加城市信息，默认`false`
- [x] 工具包优化取主机IP地址逻辑（优先eth0网卡内网地址），`glc-logback-appender`同步升级

### 版本`0.14.0`

- [x] 日志新增用户字段，界面新增用户的精确检索条件，当要做特定用户维度的日志审计时，这会显得非常实用
- [x] 包`glc-logback-appender`同步升级，新增MDC存取用户的接口
- [x] 修复已知问题

### 版本`0.13.0`

- [x] 新增用户及系统权限管理，仅管理员能操作，可控制指定用户只能访问指定系统的日志，多系统共用且有数据安全需求时尤显重要
- [x] 升级使用`Go1.21.4`进行编译

### 版本`0.12.4`

- [x] 新增会话超时`GLC_SESSION_TIMEOUT`环境变量，单位为分钟，默认30分钟
- [x] 优化检索性能，部分多选日志级别的场景，性能改善明显

### 版本`0.12.3`

- [x] 新增跨仓检索支持，分仓模式下有时需要逐个检索日志仓进行确认，确实累人。现在只要清空日志仓条件再选择日期范围，就可以轻松的查取目标数据
- [x] 改善操作体验，修改页面检索条件后再滚动查询时仍旧是按原条件查询，避免新旧条件不同引发令人困惑的查询结果

### 版本`0.12.2`

- [x] 安全无小事，继续加固，连续登录失败5次后，锁定15分钟限制登录
- [x] 新加白名单`GLC_WHITE_LIST`环境变量，多个用逗号分隔，黑白名单冲突时白名单优先，默认空白。可设定IP，最后段支持通配符，如`1.2.3.*`，内网IP默认都是白名单不必设定，实验性质的支持区域名称（因为IP地域查询可能有误差），如`上海市,深圳市`
- [x] 新加黑名单`GLC_BLACK_LIST`环境变量，多个用逗号分隔，黑白名单冲突时白名单优先，默认空白。可设定IP，最后段支持通配符，如`1.2.3.*`，也支持单个通配符`*`代表全部（也就是只允许内网或白名单指定使用），实验性质的支持区域名称（因为IP地域查询可能有误差）

### 版本`0.12.1`

- [x] 新增`GLC_TOKEN_SALT`令牌盐环境变量，默认空串。如果日志内容比较敏感，应该修改用户密码开启登录功能，同时建议设定令牌盐，提高系统安全性
- [x] 新增`GLC_TEST_MODE`是否开启测试模式的开关，开启后将显示生成测试数据用的按钮，供测试或快速体验用，默认`false`
- [x] 优化改善

### 版本`0.12.0`

- [x] 增加配置开关`GLC_SEARCH_MULIT_LINE`，设定为`true`时，支持对日志列的全部行进行索引和检索，默认`false`。注意：不会对历史数据进行重新索引，也就是说，设定为`true`时，新加入的日志会做多行索引，但历史数据如果没有多行索引的仍旧没法进行多行检索
- [x] 同步升级`glc-logback-appender`，增加过滤器类`GlcFilter`用以生成客户端IP和跟踪码，可按需配置使用

### 版本`0.11.7`

- [x] 增加支持开始/停止自动查询，观察实时日志时实用

### 版本`0.11.6`

- [x] 升级使用`Go1.21.3`，一波安全更新
- [x] 增加`robots.txt`，拒绝爬虫爬取内容
- [x] 其他一些细节改善

### 版本`0.11.5`

- [x] 升级使用`Go1.21.1`
- [x] 更新升级依赖包，避免潜在问题
- [x] 前端页面细节改善

### 版本`0.11.4`

- [x] 支持日志级别多选条件，想排除某种级别进行检索时会很实用
- [x] 改善前端进一步提高使用体验

### 版本`0.11.3`

- [x] 保留页面中表格的配置，避免再次访问或刷新时自动重置

### 版本`0.11.2`

- [x] 支持下载保存当前检索结果

### 版本`0.11.1`

- [x] 升级使用`Go1.21.0`
- [x] 运行时基础镜像`alpine`升级至`3.18`
- [x] 页面检索件数提示信息改善

### 版本`0.11.0`

- [x] 前端全面重构改良，支持表格列宽、位置、显示隐藏等各种个性化设定
- [x] 新增`GLC_ENABLE_CORS`参数配置是否允许跨域，默认`false`，方便系统间对接
- [x] 新增`GLC_PAGE_SIZE`参数配置每次检索件数，默认`100`（有效范围`1~1000`）

### 版本`0.10.2`

- [x] 修复`issue #16`的查询BUG

### 版本`0.10.1`

- [x] 添加支持`日志级别`展示列及过滤条件，需同步使用`glc-logback-appender:0.10.1`
- [x] 添加支持`客户端IP`展示列，基于MDC实现，Java项目需参考使用MdcUtil类
- [x] 添加支持`TraceId`展示列，基于MDC实现，Java项目需参考使用MdcUtil类

### 版本`0.10.0`

- [x] 页面优化：系统名检索条件可选择输入，可以不用敲打了
- [x] 页面增加主机名、主机IP展示列，可配置是否显示，适用更多复杂使用场景
- [x] 同步使用`glc-logback-appender:0.10.0`，即可自动产生主机名、主机IP信息
- [x] 修复一些小瑕疵

### 版本`0.9.0`

- [x] 增加分类(系统)检索条件，支持多系统时准确筛选
- [x] 修复一些小瑕疵

### 版本`0.8.8`

- [x] 增加时间范围检索条件
- [x] 界面进一步简化优化

### 版本`0.8.7`

- [x] 修复：增加特殊字符转换处理，避免日志中的html标签字样无法显示
- [x] 修复：docker restart 失败问题

### 版本`0.8.6`

- [x] 升级使用Go1.20，更优秀的编译和运行时，进一步减少内存开销，进一步提高整体CPU性能

### 版本`0.8.5`

- [x] 后端配合前端路由，设定自动跳转改善使用体验

### 版本`0.8.4`

- [x] 代码整理优化，前后端升级依赖包
- [x] `golang`编译器升级至`1.19.4`
- [x] 运行时基础镜像`alpine`升级至`3.17`

### 版本`0.8.3`

- [x] 增加CPU使用限制的配置`GLC_GOMAXPROCS`，默认使用最大CPU数量，通常非`docker`启动方式使用
- [x] 代码优化，升级依赖包

### 版本`0.8.2`

- [x] 优化部分场景下的检索性能
- [x] 修复一些小问题

### 版本`0.8.1`

- [x] 代码重构改善
- [x] 支持日志级别配置`GLC_LOG_LEVEL`，可设定值为`debug/info/warn/error`，默认`info`

### 版本`0.8.0`

- [x] 集群支持动态扩增节点，日志自动转发
- [x] 集群支持自动选举Master
- [x] 隔日之前的历史日志仓自动检查同步

### 版本`0.7.0`

- [x] 增加日志转发功能，支持多服务集群模式部署，确保服务及数据保存的冗余性

### 版本`0.6.0`

- [x] 升级使用Go1.19
- [x] 优化执行文件体积，此版本考虑直接运行发布，以适用更多部署场景
- [x] 支持命令行参数使用`-v`查看版本
- [x] 在Linux系统下支持命令行参数使用`-d`以后台方式启动
- [x] 在Linux系统下支持命令行参数使用`stop`停止程序
- [x] 在Linux系统下支持命令行参数使用`restart`重启程序
- [x] `logback`用`jar`包，支持通过设定环境变量`GLC_ENABLE=false`关闭日志发送功能

### 版本`0.5.0`

- [x] 增加用户密码登录功能，可设定是否开启用户密码登录
- [x] 日志按日分仓存储时，默认自动维护保存最多180天，自动维护时不能手动删除日志仓
- [x] 改善日志仓管理页面的展示
- [x] 删除旧版接口`/glc/add`、`/glc/search`，`maven`公共仓库包同步修改并更新版本
- [x] `Docker`镜像设定默认时区`Asia/Shanghai`

### 版本`0.4.0`

- [x] 添加相应版本的`maven`公共仓库包，`java`项目日志可推至`RabbitMQ`
- [x] 添加`RabbitMQ`简单模式消费者，开启后能从`RabbitMQ`获取日志
- [x] 添加服务接口`/glc/v1/log/add`，接收`JSON`格式日志以便后续扩展
- [x] 添加日志仓管理功能，页面支持查看、删除等操作

### 版本`0.3.0`

- [x] 全面重构，不考虑旧版兼容
- [x] 控制索引文件数，避免大量日志时打开文件过多而崩溃
- [x] 降低索引文件的磁盘空间占用，优化索引创建速度
- [x] 检索页面，显示更友好的查询结果提示
- [x] `glc-logback-appender`的设定，可通过配置环境变量来覆盖

### 版本`0.2.0`

- [x] DIY了一个`logo`
- [x] 接口`/glc/add`添加`system`参数
- [x] 提供简洁的日志查询界面
- [x] 当前版设计为接收多个项目的日志，界面栏目为`分类`

### 初版`0.1.0`

- [x] 使用`golang`实现，就是快
- [x] 借助`goleveldb`做数据保存，结合日志写多读少特点稍加设计，真是快
- [x] 日志量虽大，却是真心节省内存资源
- [x] 常用的无条件查询最新日志，快到麻木无感
- [x] 关键词全文检索，支持中文分词，反向索引以空间换时间，快到麻木无感
- [x] 提供`docker`镜像支持容器化部署，方便之极
- [x] 提供`java`项目日志收集包，日志都发来发来发来
- [x] 服务接口`/glc/add`添加日志
- [x] 服务接口`/glc/search`查询日志

</details>
