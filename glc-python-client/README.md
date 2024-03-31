日志中心的 python 客户端


## 使用`python`语言的项目，提供工具包，开箱即用
```shell
# 支持以下环境变量配置
export GLC_ENABLE=true # 默认false，要发送日志中心必须配置为true
export GLC_ENABLE_CONSOLE_LOG=true # 默认true，控制台不打印时配置为false
export GLC_API_URL='http://127.0.0.1:8080/glc/v1/log/add' # 未配置时将取消发送
export GLC_API_KEY='X-GLC-AUTH:glogcenter' # 这是默认值，按需修改
export GLC_SYSTEM=default  # 默认default，按需修改
export GLC_LOG_LEVEL=debug # 日志级别（debug/info/warn/error），默认debug
export GLC_TRACE_ID=12345  # 默认空
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
