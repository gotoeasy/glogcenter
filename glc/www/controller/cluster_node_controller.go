package controller

import (
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"glc/www/cluster"
	"glc/www/service"
	"io"
	"net/http"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

// 从本节点取集群信息
func ClusterGetClusterInfoController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	kv, err := service.GetSysmntItem(cluster.KEY_CLUSTER)
	if kv == nil || err != nil {
		return gweb.Error500("没找到集群信息 " + err.Error())
	}

	return gweb.Result(kv.ToJson())
}

// 作为Master接收保存集群信息
func ClusterMasterSaveKvDataController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	// 参数
	pkv := &service.KeyValue{}
	err := req.BindJSON(pkv)
	if err != nil || pkv.Key == "" {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	// 保存
	dkv, err := service.SetSysmntItem(pkv)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500(err.Error())
	}

	// 转发
	cl := &cluster.ClusterInfo{}
	cl.LoadJson(dkv.Value)
	urls := cmn.Split(cl.NodeUrls, ";")
	jsonstr := dkv.ToJson()
	for i := 0; i < len(urls); i++ {
		if urls[i] != com.GetLocalGlcUrl() {
			go httpClusterAsyncData(urls[i], jsonstr) // 异步转发其他节点保存，暂且忽略失败
		}
	}

	return gweb.Result(dkv.ToJson())
}

// 接收其他节点发来的同步信息
func ClusterMasterAsyncDataController(req *gweb.HttpRequest) *gweb.HttpResult {
	cmn.Info("接收其他节点发来的同步信息")
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		cmn.Info("接收其他节点发来的同步信息", 403, "未经授权的访问")
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	pkv := &service.KeyValue{}
	err := req.BindJSON(pkv)
	if err != nil || pkv.Key == "" {
		cmn.Error("接收其他节点发来的同步信息", "请求参数有误", err)
		return gweb.Error500("请求参数有误" + err.Error())
	}

	dkv, err := service.SetSysmntItem(pkv)
	if err != nil {
		cmn.Error("接收其他节点发来的同步信息", "保存失败", err)
	} else {
		cmn.Info("接收其他节点发来的同步信息", dkv.ToJson())
	}
	return gweb.Result(dkv.ToJson())
}

func httpClusterAsyncData(serverUrl string, kvJson string) *service.KeyValue {

	// 请求
	req, err := http.NewRequest("POST", serverUrl+conf.GetContextPath()+"/sys/cluster/async", strings.NewReader(kvJson))
	if err != nil {
		cmn.Error("异步发送集群信息到", serverUrl, "失败1", err)
		return nil
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set(conf.GetHeaderSecurityKey(), conf.GetSecurityKey())

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		cmn.Error("异步发送集群信息到", serverUrl, "失败2", err)
		return nil
	}
	defer res.Body.Close()

	by, err := io.ReadAll(res.Body)
	if err != nil {
		cmn.Error("异步发送集群信息到", serverUrl, "失败3", err)
		return nil
	}

	rs := &gweb.HttpResult{}
	rs.LoadBytes(by)
	if !rs.Success {
		cmn.Error("异步发送集群信息到", serverUrl, "失败4", rs.Message)
		return nil
	}

	cmn.Info("异步发送集群信息到", serverUrl, "成功")

	kv := &service.KeyValue{}
	kv.LoadJson(rs.Result.(string))
	return kv
}
