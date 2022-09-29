package cluster

import (
	"glc/cmn"
	"glc/conf"
	"glc/www/service"
	"log"
)

func Start() {
	if !conf.IsClusterMode() {
		return
	}

	log.Println("集群节点启动", cmn.GetLocalGlcUrl())
	kv, err := service.GetSysmntItem(KEY_CLUSTER)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(kv.ToJson())
	}
}
