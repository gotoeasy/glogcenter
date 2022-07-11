<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/logo.png"/>
</div>


# 缘起

日志中心，一直是用传统的三件套`ELK`，但终究还是不理想（定制安装并非简单、页面打开初始化太慢，界面操作不习惯，最主要的是资源占用太厉害，甚至隔断时间就会崩溃），替代品总是找不到<br>
终于，在又一次的挂掉之后，忍不住用`go`试写一个`logcenter`，结果各种优异表现确实是惊艳到了自己，故起名`glogcenter`，简称`GLC`，开仓建库，目标是逐步替换线上的`ELK`
<br>

[![Docker Pulls](https://img.shields.io/docker/pulls/gotoeasy/glc)](https://hub.docker.com/r/gotoeasy/glc)
[![GitHub release](https://img.shields.io/github/release/gotoeasy/glogcenter.svg)](https://github.com/gotoeasy/glogcenter/releases/latest)
<br>


## 特点
- [x] 使用`golang`实现，具备`go`的各种特性优势，关键是省资源、性能高
- [x] 借助`goleveldb`做数据保存，结合日志写多读少特点稍加设计，真是快
- [x] 关键词全文检索，支持中文分词（使用`jiebago`进行分词），毫秒级响应，自然流畅
- [x] 日志吞食量每秒近万条，闲时建索引速度每秒数千条，基本能满足大多项目需要
- [x] 提供`docker`镜像，支持容器化部署，方便之极
- [x] 提供`java`项目日志收集包，`java`项目闭环支持
- [x] 支持从`RabbitMQ`收取日志信息，满足更多闭环需求
- [x] 内置提供简洁的`VUE`实现的日志查询管理界面

<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/glogcenter.png"/>
</div>

<br>

## `docker`运行
```
// 简单示例
docker run -d -p 8080:8080 gotoeasy/glc

// 外挂数据目录
docker run -d -p 8080:8080 -v /glc:/glogcenter gotoeasy/glc
```


## `docker`启动环境变量
- [x] `GLC_STORE_NAME_AUTO_ADD_DATE`日志仓是否自动按日存储，默认`true`
- [x] `GLC_SAVE_DAYS`日志仓按日存储自动维护时的保留天数(`0~180`)，`0`表示不自动删除，默认`180`天
- [x] `GLC_ENABLE_LOGIN`是否开启用户密码登录功能，默认`false`
- [x] `GLC_USERNAME`查询界面登录用的用户名，默认`glc`
- [x] `GLC_PASSWORD`查询界面登录用的密码，默认`glogcenter`
- [x] `GLC_ENABLE_SECURITY_KEY`日志添加的接口是否开启API秘钥校验，默认`false`
- [x] `GLC_HEADER_SECURITY_KEY`API秘钥的`header`键名，默认`X-GLC-AUTH`
- [x] `GLC_SECURITY_KEY`API秘钥，默认`glogcenter`
- [x] `GLC_ENABLE_AMQP_CONSUME`是否开启`rabbitMq`消费者接收日志，默认`false`
- [x] `GLC_AMQP_ADDR`消息队列`rabbitMq`连接地址，例：`amqp://user:password@ip:port/`，默认空白
- [x] `GLC_AMQP_JSON_FORMAT`消息队列`rabbitMq`消息文本是否为`json`格式，默认`true`



## 接口
- [x] `/glc/v1/log/add`日志添加，`POST`，`application/json` <br>
      字段`system`： 字符串，对应页面的`分类` <br>
      字段`date`： 字符串，对应页面的`日期时间` <br>
      字段`text`： 字符串，对应页面的`日志内容` <br>



## 使用`logback`的`java`项目，支持日志收集
```xml
<!-- pom坐标 -->
<dependency>
    <groupId>top.gotoeasy</groupId>
    <artifactId>glc-logback-appender</artifactId>
    <version>0.5.0</version>
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


## 更新履历

### 开发版`latest`

- [ ] 界面优化
- [ ] 多语言
- [ ] 分词优化

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
