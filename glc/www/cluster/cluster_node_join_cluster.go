package cluster

import (
	"encoding/json"
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"glc/www/service"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

const KEY_CLUSTER string = "$cluster"

type ClusterInfo struct {
	MasterUrl string `json:"masterUrl,omitempty"`
	NodeUrls  string `json:"nodeUrls,omitempty"`
	Version   string `json:"version,omitempty"`
}

func (d *ClusterInfo) ToJson() string {
	bt, _ := json.Marshal(d)
	return cmn.BytesToString(bt)
}

func (d *ClusterInfo) LoadJson(jsonstr string) error {
	if jsonstr == "" {
		return nil
	}
	return json.Unmarshal(cmn.StringToBytes(jsonstr), d)
}

// 启动时加入集群
func joinCluster() {
	if !conf.IsClusterMode() {
		return
	}

	service.DelSysmntItem(KEY_CLUSTER) // 刚启动，默认本机不是Master

	localGlcUrl := com.GetLocalGlcUrl()
	urls := conf.GetClusterUrls()

	// 仅单节点
	if len(urls) == 0 || (len(urls) == 1 && urls[0] == localGlcUrl) {
		cmn.Info("单节点", com.GetLocalGlcUrl())
		ci := &ClusterInfo{
			MasterUrl: localGlcUrl,
			NodeUrls:  localGlcUrl,
			Version:   "1",
		}
		kv := &service.KeyValue{
			Key:     KEY_CLUSTER,
			Value:   ci.ToJson(),
			Version: "1",
		}
		service.SetSysmntItem(kv) // 保存节点信息
		cmn.Info("集群模式，当前为单节点")
		return
	}

	// 集群模式（准备加入集群，找Master，无Master时选举）
	var masters []*ClusterInfo
	var onlines []string
	for i := 0; i < len(urls); i++ {
		if urls[i] == localGlcUrl {
			continue
		}
		master := httpGetMasterFromRemote(urls[i])
		if master != nil {
			onlines = append(onlines, urls[i])
			if master.MasterUrl != "" {
				masters = append(masters, master)
			}
		}
	}

	if len(masters) == 0 {
		// 还没有master， 在线范围倒序选举
		aryUp := []string{}
		aryUp = append(aryUp, localGlcUrl)
		aryUp = append(aryUp, onlines...)
		aryUp = com.Unique(aryUp) // 去重
		sort.Slice(aryUp, func(i, j int) bool {
			return aryUp[i] > aryUp[j] // 倒序
		})

		// 所有节点
		aryall := []string{}
		aryall = append(aryall, localGlcUrl)
		aryall = append(aryall, urls...)
		aryall = com.Unique(aryall) // 去重
		sort.Slice(aryall, func(i, j int) bool {
			return aryall[i] < aryall[j] // 升序
		})

		// 保存节点信息
		masterUrl := aryUp[0]
		ci := &ClusterInfo{
			MasterUrl: masterUrl,
			NodeUrls:  cmn.Join(aryall, ";"),
			Version:   "1",
		}
		kv := &service.KeyValue{
			Key:     KEY_CLUSTER,
			Value:   ci.ToJson(),
			Version: "1",
		}

		cmn.Info("当前无Master，选定Master为", masterUrl)
		_, err := service.SetSysmntItem(kv) // 保存
		if err != nil {
			cmn.Error("保存集群信息失败", err)
		}

		jsonstr := kv.ToJson()
		cmn.Info("本节点已保存集群信息", masterUrl, ci.NodeUrls)
		for i := 0; i < len(aryUp); i++ {
			if aryUp[i] != com.GetLocalGlcUrl() {
				go httpClusterAsyncData(aryUp[i], jsonstr) // 异步发送数据给全部节点保存
			}
		}

	} else {
		// 已有master
		// 倒序(最大版本)
		sort.Slice(masters, func(i, j int) bool {
			return cmn.StringToUint32(masters[i].Version, 0) > cmn.StringToUint32(masters[j].Version, 0)
		})
		// 更新保存集群信息
		nversion := cmn.Uint32ToString(cmn.StringToUint32(masters[0].Version, 0) + 1)
		ourls := masters[0].NodeUrls
		nodes := cmn.Split(ourls, ";")
		nodes = append(nodes, localGlcUrl) // 当前节点加入集群
		nodes = com.Unique(nodes)          // 去重
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i] < nodes[j] // 升序
		})
		nurls := cmn.Join(nodes, ";") // 分号分隔保存

		// 先尝试在原Master上保存集群信息（需触发异步群发同步集群信息）
		master := masters[0]
		master.NodeUrls = nurls // 分号分隔保存
		master.Version = nversion
		kv := &service.KeyValue{
			Key:     KEY_CLUSTER,
			Value:   master.ToJson(),
			Version: master.Version,
		}

		_, err := service.SetSysmntItem(kv) // 保存
		if err != nil {
			cmn.Error("保存集群信息失败", err)
		} else {
			cmn.Info("本节点保存集群信息成功", master.MasterUrl, nurls)
		}

		if master.MasterUrl != com.GetLocalGlcUrl() {
			mkv := httpClusterSaveData(master.MasterUrl, kv) // 发送数据给Master节点保存
			if mkv == nil {
				// Master保存失败则更换Master保存，可能和上一步重复（当做再试一遍）
				for i := 0; i < len(nodes); i++ {
					if nodes[i] == com.GetLocalGlcUrl() {
						continue
					}

					master.MasterUrl = nodes[i]
					kv := &service.KeyValue{
						Key:     KEY_CLUSTER,
						Value:   master.ToJson(),
						Version: master.Version,
					}
					mkv := httpClusterSaveData(master.MasterUrl, kv) // 发送集群信息给Master节点保存
					if mkv != nil {
						cmn.Info("保存集群信息到", master.MasterUrl, "成功")
						break // 保存成功
					}
				}
			} else {
				cmn.Info("保存集群信息到", master.MasterUrl, "成功")
			}
		} else {
			// 当前节点是Master且已保存，群发给其他节点
			jsonstr := kv.ToJson()
			for i := 0; i < len(nodes); i++ {
				if nodes[i] != com.GetLocalGlcUrl() {
					go httpClusterAsyncData(nodes[i], jsonstr) // 异步发送数据给全部节点保存
				}
			}
		}

	}

}

func httpGetMasterFromRemote(url string) *ClusterInfo {
	kv := httpGetClusterInfo(url)
	if kv == nil {
		return nil
	}
	// if cmn.StringToUint32(kv.Version, 0) == 0 {
	// 	return nil
	// }
	cl := &ClusterInfo{}
	cl.LoadJson(kv.Value)
	if cl.MasterUrl == "" {
		return nil
	}
	return cl
}

func httpGetClusterInfo(serverUrl string) *service.KeyValue {

	// 请求
	req, err := http.NewRequest("POST", serverUrl+conf.GetContextPath()+"/sys/cluster/info", strings.NewReader(""))
	if err != nil {
		cmn.Error("从", serverUrl, "取集群信息失败", err)
		return nil
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set(conf.GetHeaderSecurityKey(), conf.GetSecurityKey())

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		cmn.Error("从", serverUrl, "取集群信息失败", err)
		return nil
	}
	defer res.Body.Close()

	by, err := io.ReadAll(res.Body)
	if err != nil {
		cmn.Error("从", serverUrl, "取集群信息失败", err)
		return nil
	}

	rs := &gweb.HttpResult{}
	rs.LoadBytes(by)
	if !rs.Success {
		cmn.Error("从", serverUrl, "取集群信息失败", rs.Message)
		return &service.KeyValue{} // 能通信，只是没数据
	}

	kv := &service.KeyValue{}
	kv.LoadJson(rs.Result.(string))
	cmn.Info("从", serverUrl, "取集群信息成功", kv.ToJson())
	return kv
}

func httpClusterSaveData(serverUrl string, clusterKv *service.KeyValue) *service.KeyValue {

	// 请求
	req, err := http.NewRequest("POST", serverUrl+conf.GetContextPath()+"/sys/cluster/save", strings.NewReader(clusterKv.ToJson()))
	if err != nil {
		cmn.Error("请求", serverUrl, "保存集群信息失败", err)
		return nil
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set(conf.GetHeaderSecurityKey(), conf.GetSecurityKey())

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		cmn.Error("请求", serverUrl, "保存集群信息失败", err)
		return nil
	}
	defer res.Body.Close()

	by, err := io.ReadAll(res.Body)
	if err != nil {
		cmn.Error("请求", serverUrl, "保存集群信息失败", err)
		return nil
	}

	rs := &gweb.HttpResult{}
	rs.LoadBytes(by)
	if !rs.Success {
		cmn.Error("请求", serverUrl, "保存集群信息失败", rs.Message)
		return nil
	}

	kv := &service.KeyValue{}
	kv.LoadJson(rs.Result.(string))
	cmn.Info("请求", serverUrl, "保存集群信息成功", kv.Value)
	return kv
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

	kv := &service.KeyValue{}
	kv.LoadJson(rs.Result.(string))

	cmn.Info("异步发送集群信息到", serverUrl, "成功")
	return kv
}
