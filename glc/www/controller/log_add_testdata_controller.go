package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb/storage/logdata"

	"github.com/gotoeasy/glang/cmn"
)

// 添加2000条测试日志（仅测试模式有效），用于测试或快速体验
func JsonLogAddTestDataController(req *gweb.HttpRequest) *gweb.HttpResult {

	if !conf.IsTestMode() {
		return gweb.Error500("当前不是测试模式，不支持生成测试数据") // 非测试模式时忽略
	}

	cnt := 0
	for {
		cnt++
		traceId := cmn.RandomHashString()
		md := &logdata.LogDataModel{
			Text:       "测试用的日志，字段名为Text，" + "字段Date的格式为YYYY-MM-DD HH:MM:SS.SSS，必须格式一致才能正常使用时间范围检索条件。" + "随机3位字符串：" + cmn.RandomString(3) + "，第" + cmn.IntToString(cnt) + "条",
			Date:       cmn.Now(),
			System:     "demo1",
			ServerName: "default",
			ServerIp:   "127.0.0.1",
			ClientIp:   "127.0.0.1",
			TraceId:    traceId,
			LogLevel:   "INFO",
			User:       "tuser-" + cmn.RandomString(1),
		}
		addDataModelLog(md)

		if conf.IsClusterMode() {
			go TransferGlc(conf.LogTransferAdd, md.ToJson()) // 转发其他GLC服务
		}

		md2 := &logdata.LogDataModel{
			Text:       "几个随机字符串供查询试验：" + cmn.RandomString(1) + "，" + cmn.Right(cmn.ULID(), 2) + "，" + cmn.RandomString(3) + "，" + cmn.Right(cmn.ULID(), 4) + "，" + cmn.Right(cmn.ULID(), 5),
			Date:       cmn.Now(),
			System:     "demo2",
			ServerName: "default",
			ServerIp:   "127.0.0.1",
			ClientIp:   "127.0.0.1",
			TraceId:    traceId,
			LogLevel:   "DEBUG",
			User:       "tuser-" + cmn.RandomString(1),
		}
		addDataModelLog(md2)

		if conf.IsClusterMode() {
			go TransferGlc(conf.LogTransferAdd, md2.ToJson()) // 转发其他GLC服务
		}

		md3 := &logdata.LogDataModel{
			Text:       "几个随机字符串供查询试验：" + cmn.RandomString(1) + "，" + cmn.Right(cmn.ULID(), 2) + "，" + cmn.RandomString(3) + "，" + cmn.Right(cmn.ULID(), 4) + "，" + cmn.Right(cmn.ULID(), 5),
			Date:       cmn.Now(),
			System:     "demo3",
			ServerName: "default",
			ServerIp:   "127.0.0.1",
			ClientIp:   "127.0.0.1",
			TraceId:    traceId,
			LogLevel:   "WARN",
			User:       "tuser-" + cmn.RandomString(1),
		}
		addDataModelLog(md3)

		if conf.IsClusterMode() {
			go TransferGlc(conf.LogTransferAdd, md3.ToJson()) // 转发其他GLC服务
		}

		if cnt >= 1000 {
			break
		}
	}

	return gweb.Ok200("操作成功")
}
