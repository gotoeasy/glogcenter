package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb"
	"glc/ldb/storage/logdata"
	"io"
	"log"
	"net/http"
	"strings"
)

// 添加日志（JSON提交方式）
func JsonLogAddController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() {
		auth := req.GetHeader(conf.GetHeaderSecurityKey())
		if auth != conf.GetSecurityKey() {
			return gweb.Error(403, "未经授权的访问，拒绝服务")
		}
	}

	md := &logdata.LogDataModel{}
	err := req.BindJSON(md)
	if err != nil {
		log.Println("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	if conf.IsEnableSlaveTransfer() {
		transferGlc(md) // 转发其他GLC服务
	}

	engine := ldb.NewDefaultEngine()
	engine.AddTextLog(md.Date, md.Text, md.System)
	return gweb.Ok()
}

// 转发其他GLC服务
func transferGlc(md *logdata.LogDataModel) {
	hosts := conf.GetSlaveHosts()
	for i := 0; i < len(hosts); i++ {
		go httpPostJson(hosts[i]+conf.GetContextPath()+"/v1/log/add", md.ToJson())
	}
}

func httpPostJson(url string, jsondata string) ([]byte, error) {

	// 请求
	req, err := http.NewRequest("POST", url, strings.NewReader(jsondata))
	if err != nil {
		return nil, err
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set(conf.GetHeaderSecurityKey(), conf.GetSecurityKey())

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
