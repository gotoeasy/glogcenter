package controller

import (
	"glc/com"
	"glc/conf"
	"glc/www/cluster"
	"glc/www/service"
	"io"
	"net/http"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

// 转发其他GLC服务
func TransferGlc(uri string, jsonlog string) {
	kv, err := service.GetSysmntItem(cluster.KEY_CLUSTER)
	if kv == nil || err != nil {
		cmn.Error("转发日志失败（取集群信息失败）", err)
		return
	}

	ci := &cluster.ClusterInfo{}
	ci.LoadJson(kv.Value)

	hosts := cmn.Split(ci.NodeUrls, ";")
	for i := 0; i < len(hosts); i++ {
		if hosts[i] != com.GetLocalGlcUrl() {
			_, err := httpPostJson(hosts[i]+conf.GetContextPath()+"/v1/log/transferAdd", jsonlog)
			if err != nil {
				cmn.Error("集群内数据转发失败", hosts[i], err)
			}
		}
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
