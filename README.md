<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/logo.png"/>
</div>


# 缘起

日志中心，一直是用传统三件套，打开速度不理想，资源占用太厉害，常闹情绪，替代品总是找不到。终于忍不住用`go`试写一个`logcenter`，结果确实是惊艳到自己了，故起名`glogcenter`，简称`GLC`，开仓建库
<br>


# 特点
- [x] 使用`golang`实现，就是快
- [x] 借助`goleveldb`做数据保存，结合日志写多读少特点稍加设计，真是快
- [x] 日志量虽大，却是真心节省内存资源
- [x] 常用的无条件查询最新日志，快到麻木无感
- [x] 关键词全文检索，支持中文分词，反向索引以空间换时间，快到麻木无感
- [x] 提供`docker`镜像支持容器化部署，方便之极
- [x] 提供`java`项目日志收集包，日志都发来发来发来
- [x] 提供简洁的日志查询界面

<div align=center>
<img src="https://gotoeasy.github.io/screenshots/glogcenter/glogcenter.png"/>
</div>

<br>

# `docker`运行
```
// 简单示例
docker run -d -p 8080:8080 -v /logdata:/glogcenter gotoeasy/glc
```


## 服务接口
- [x] `/glc/add`日志添加，表单提交方式，字段`text`是日志内容，更多参数斟酌中
- [x] `/glc/search`日志查询，表单提交方式，检索条件字段`searchKey`，更多参数斟酌中



# 使用`logback`的`java`项目，支持日志收集
```xml
<!-- pom坐标 -->
<dependency>
    <groupId>top.gotoeasy</groupId>
    <artifactId>glc-logback-appender</artifactId>
    <version>0.2.0</version>
</dependency>
```

```xml
<!-- logback-spring.xml配置例子 -->
<appender name="GLC" class="top.gotoeasy.framework.glc.logback.appender.GlcHttpAppender">
    <glcApiUrl>http://127.0.0.1:8080/glc/add</glcApiUrl>
    <glcApiKey>X-GLC-AUTH:glogcenter</glcApiKey>
    <system>Demo</glcApiKey>
    <layout>
        <pattern><![CDATA[%p %m %n]]></pattern>
    </layout>
</appender>
```


# TODO
- [ ] 为了性能，好像把空间代价有搞大了
- [ ] 还好多



# 更新履历

### 版本`0.2.0`，`latest`

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
