<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/logo.png"/>
</div>


# 缘起

日志中心，一直是用传统的三件套`ELK`，但终究还是不理想（定制安装并非简单、需要维护索引、页面打开初始化太慢，界面操作不习惯，资源占用太厉害甚至容易崩溃，等等诸如此类），替代品非常难找。<br>
<br>
终于，尝试用`go`写个日志中心，结果各种优异表现令人惊艳，故起名`glogcenter`，简称`GLC`，开仓建库<br>
<br>

[![Golang](https://img.shields.io/badge/golang-1.20-brightgreen.svg)](https://golang.google.cn)
[![Docker Pulls](https://img.shields.io/docker/pulls/gotoeasy/glc)](https://hub.docker.com/r/gotoeasy/glc)
[![GitHub release](https://img.shields.io/github/release/gotoeasy/glogcenter.svg)](https://github.com/gotoeasy/glogcenter/releases/latest)
[![License](https://img.shields.io/github/license/gotoeasy/glogcenter)](https://github.com/gotoeasy/glogcenter/blob/master/LICENSE)
<br>


## 特点
- [x] 使用`golang`实现，具备`go`的各种特性优势，关键是省资源、性能高
- [x] 借助`goleveldb`做数据保存，结合日志写多读少特点稍加设计，真是快
- [x] 关键词全文检索，支持中文分词（使用`jiebago`进行分词），毫秒级响应，自然流畅
- [x] 日志吞食量每秒近万条，闲时建索引速度每秒数千条，基本能满足大多项目需要
- [x] 支持个性化环境变量开关控制，支持日志仓自动化维护，灵活省心
- [x] 提供`docker`镜像，支持容器化部署，方便之极
- [x] 提供`java`项目日志收集包，`java`项目闭环支持
- [x] 提供`golang`项目日志收集包，`golang`项目闭环支持
- [x] 支持从`RabbitMQ`收取日志信息，满足更多闭环需求
- [x] 内置提供简洁的`VUE`实现的日志查询管理界面
- [x] 支持多服务集群模式部署，确保服务及数据的冗余性


<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/glogcenter.png"/>
</div>

<br>

## `docker`单机部署模式简易示例
```shell
# 简单示例
docker run -d -p 8080:8080 gotoeasy/glc

# 外挂数据目录
docker run -d -p 8080:8080 -v /glc:/glogcenter gotoeasy/glc
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
- [x] `GLC_SAVE_DAYS`日志仓按日存储自动维护时的保留天数(`0~180`)，`0`表示不自动删除，默认`180`天
- [x] `GLC_ENABLE_LOGIN`是否开启用户密码登录功能，默认`false`
- [x] `GLC_USERNAME`查询界面登录用的用户名，默认`glc`
- [x] `GLC_PASSWORD`查询界面登录用的密码，默认`GLogCenter100%666`
- [x] `GLC_ENABLE_SECURITY_KEY`日志添加的接口是否开启API秘钥校验，默认`false`
- [x] `GLC_HEADER_SECURITY_KEY`API秘钥的`header`键名，默认`X-GLC-AUTH`
- [x] `GLC_SECURITY_KEY`API秘钥，默认`glogcenter`
- [x] `GLC_ENABLE_AMQP_CONSUME`是否开启`rabbitMq`消费者接收日志，默认`false`
- [x] `GLC_AMQP_ADDR`消息队列`rabbitMq`连接地址，例：`amqp://user:password@ip:port/`，默认空白
- [x] `GLC_AMQP_JSON_FORMAT`消息队列`rabbitMq`消息文本是否为`json`格式，默认`true`
- [x] `GLC_CLUSTER_MODE`是否集群模式启动，默认`false`
- [x] `GLC_SERVER_URL`集群模式时的本节点服务地址，默认空白
- [x] `GLC_CLUSTER_URLS`集群模式时的关联节点服务地址，多个时`;`分隔，默认空白
- [x] `GLC_LOG_LEVEL`日志级别，可设定值为`debug/info/warn/error`，默认`info`
- [x] `GLC_GOMAXPROCS`使用最大CPU数量，值不在实际范围内时按最大值看待，默认最大值，常用于`docker`方式

## 命令行启动参数（适用`0.6.*`及以上版本）
- [x] 支持命令行参数`-v`查看版本
- [x] 在Linux系统下支持命令行参数`-d`以后台方式启动
- [x] 在Linux系统下支持命令行参数`stop`关闭程序
- [x] 在Linux系统下支持命令行参数`restart`重启程序
- [x] 服务端口可通过环境变量`GLC_SERVER_PORT`设定，默认`8080`
- [x] 其他环境变量同以上`docker`启动环境变量


## 接口
- [x] `/glc/v1/log/add`日志添加，`POST`，`application/json` <br>
      字段`system`： 字符串，对应页面的`分类` <br>
      字段`date`： 字符串，对应页面的`日期时间` <br>
      字段`text`： 字符串，对应页面的`日志内容` <br>

```shell
# 发送测试数据的参考脚本
curl -X POST -d '{"system":"demo", "date":"20230101 01:02:03.456","text":"demo log text"}' \
     -H "Content-Type:application/json" http://127.0.0.1:8080/glc/v1/log/add
```


## 使用`logback`的`java`项目，支持日志收集，确保主次版本和GLC版本一致
```xml
<!-- pom坐标 -->
<dependency>
    <groupId>top.gotoeasy</groupId>
    <artifactId>glc-logback-appender</artifactId>
    <version>0.9.0</version>
</dependency>
```

```xml
<!-- logback配置例子1，发送至 glogcenter -->
<appender name="GLC" class="top.gotoeasy.framework.glc.logback.appender.GlcHttpJsonAppender">
    <glcApiUrl>http://127.0.0.1:8080/glc/v1/log/add</glcApiUrl> <!-- 可通过环境变量 GLC_API_URL 设定 -->
    <glcApiKey>X-GLC-AUTH:glogcenter</glcApiKey>                <!-- 可通过环境变量 GLC_API_KEY 设定 -->
    <system>Demo</system>                                       <!-- 可通过环境变量 GLC_SYSTEM 设定 -->
    <layout>
        <pattern><![CDATA[%p %m %n]]></pattern>
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
        <pattern><![CDATA[%p %m %n]]></pattern>
    </layout>
</appender>

```

## 使用`golang`语言的项目，提供工具包，开箱即用
```shell
# 引入工具包
go get github.com/gotoeasy/glang

# 按需设定环境变量
export GLC_API_URL='http://127.0.0.1:8080/glc/v1/log/add'
export GLC_API_KEY='X-GLC-AUTH:glogcenter'
export GLC_SYSTEM=demo
export GLC_ENABLE=true
export GLC_LOG_LEVEL=debug # 日志级别（trace/debug/info/warn/error/fatal）
```

```golang
// 方式1： 通过 cmn.Debug(...)、cmn.Info(...)等方式，打印日志的同时发送至日志中心
// 方式2： 通过 cmn.NewGLogCenterClient()创建客户端对象后使用
//        更多内容详见文档 https://pkg.go.dev/github.com/gotoeasy/glang

import "github.com/gotoeasy/glang/cmn"

func main() {
    cmn.Info("启动WEB服务")
    err := cmn.NewFasthttpServer().Start()
    if err != nil {
        cmn.Fatalln("启动失败", err)
    }
}
```


## 更新履历

### 开发版`latest`

- [ ] 多语言
- [ ] 日志审计
- [ ] 集群支持动态删减节点（或是页面管理删除）


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
