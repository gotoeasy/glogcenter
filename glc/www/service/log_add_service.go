package service

import (
	"glc/conf"
	"glc/ldb"
	"glc/ldb/storage/logdata"
	"io"
	"net/http"
	"strings"
)

// 添加日志
func AddTextLog(md *logdata.LogDataModel) {
	engine := ldb.NewDefaultEngine()
	engine.AddTextLog(md.Date, md.Text, md.System)
}

// 转发其他GLC服务
func TransferGlc(jsonlog string) {
	hosts := conf.GetSlaveHosts()
	for i := 0; i < len(hosts); i++ {
		go httpPostJson(hosts[i]+conf.GetContextPath()+"/v1/log/transferAdd", jsonlog)
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
