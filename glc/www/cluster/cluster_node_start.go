package cluster

import (
	"glc/com"
	"glc/conf"
	"glc/www/service"

	"github.com/gotoeasy/glang/cmn"
)

func Start() {
	if !conf.IsClusterMode() {
		return
	}

	cmn.Info("集群节点启动", com.GetLocalGlcUrl())
	joinCluster()
	kv, err := service.GetSysmntItem(KEY_CLUSTER)
	if err != nil {
		cmn.Error(err)
	} else {
		cmn.Info(kv.ToJson())
	}

	// 异步检查更新数据
	go dataAsync()
}
