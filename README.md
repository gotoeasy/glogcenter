# 缘起

日志中心，一直是忍者用那传统三件套，打开速度不理想，资源占用太厉害，常闹情绪，替代品总是找不到。终于忍不住用`go`试写一个`logcenter`，结果确实是惊艳到自己了，故起名`glogcenter`，简称'GLC'，开仓建库
<br>

# 特点
- [x] 使用`golang`实现，就是快
- [x] 借助`goleveldb`做数据保存，结合日志写多读少特点稍加设计，就是快
- [x] 日志量大，用内存提升性能耗不起，依靠`goleveldb`做`大数组`，真心省内存资源
- [x] 日志接收吞食量大，待压测后展示更多信息
- [x] 支持全文检索，支持中文分词
- [x] 无条件查询日志，默认查最新日志，快到麻木无感
- [x] 单关键词查询，反向索引实现，快到麻木无感
- [x] 多关键词查询，需要求反向索引的交集，常令人颇费脑筋，反复调整尝试，快到欢
- [x] 支持`docker`容器化部署，该有的方式
- [x] 提供`java`项目日志收集包，日志都放心发来发来发来
- [ ] 提供日志查询界面


# `docker`运行
```
// 简单示例
docker run -d -p 8080:8080 gotoeasy/glc
```

## `docker`启动支持的环境变量
- [x] `GLC_STORE_ROOT`数据存储的根目录，默认`/glogcenter`
- [x] `GLC_STORE_CHAN_LENGTH`接收日志的通道长度，`golang`特有的概念，默认值`64`，基本不用修改
- [x] `GLC_MAX_IDLE_TIME`最大闲置时间（秒）,超过闲置时间将自动关闭文件，用时再自动打开，0时表示不关闭，默认`180`
- [x] `GLC_STORE_NAME_AUTO_ADD_DATE`，存储名是否自动添加日期（日志量大通常按日单位区分存储），默认`false`
- [x] `GLC_SERVER_PORT`WEB服务端口，默认`8080`，极少有修改的必要
- [x] `GLC_CONTEXT_PATH`WEB服务的`contextPath`，默认`/glc`
- [x] `GLC_ENABLE_SECURITY_KEY`WEB服务是否开启API秘钥校验，默认`false`
- [x] `GLC_HEADER_SECURITY_KEY`WEB服务API秘钥的header键名，默认`X-GLC-AUTH`
- [x] `GLC_SECURITY_KEY`WEB服务API秘钥，默认`glogcenter`

## 服务接口
- [x] `/glc/add`日志添加，表单提交方式，字段`text`是日志内容，`name`是存储名可省略
- [x] `/glc/search`日志查询，表单提交方式，检索条件字段`searchKey`，更多条件字段先看代码了


# 使用`logback`的`java`项目，支持日志收集
```xml
<!-- pom坐标 -->
<dependency>
    <groupId>top.gotoeasy</groupId>
    <artifactId>glc-logback-appender</artifactId>
    <version>0.1.0</version>
</dependency>
```

```xml
<!-- logback-spring.xml配置例子 -->
<appender name="GLC" class="top.gotoeasy.framework.glc.logback.appender.GlcHttpAppender">
    <glcApiUrl>http://127.0.0.1:8080/glc/add</glcApiUrl>
    <glcApiKey>X-GLC-AUTH:glogcenter</glcApiKey>
    <layout>
        <pattern><![CDATA[ %d %c %t %p %m %n]]></pattern>
    </layout>
</appender>
```

# TODO
- [x] 还好多

# 更新履历

### 初版`0.1.0`，`latest`

- [x] 使用`golang`实现，就是快
- [x] 借助`goleveldb`做数据保存及查询辅助
- [x] 真心省内存资源
- [x] 日志接收吞食量大
- [x] 支持全文检索，支持中文分词
- [x] 支持多关键词查询交集
- [x] 支持`docker`容器化部署
- [x] 提供`java`项目日志收集包
- [x] 服务接口`/glc/add`添加日志
- [x] 服务接口`/glc/search`查询日志
