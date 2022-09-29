package service

// // 添加日志
// func AddTextLog(md *logdata.LogDataModel) {
// 	engine := ldb.NewDefaultEngine()
// 	engine.AddTextLog(md.Date, md.Text, md.System)
// }

// // 转发其他GLC服务
// func TransferGlc(jsonlog string) {
// 	kv, err := GetSysmntItem(cluster.KEY_CLUSTER)
// 	if kv == nil || err != nil {
// 		log.Panicln(err)
// 		return
// 	}

// 	ci := &cluster.ClusterInfo{}
// 	ci.LoadJson(kv.ToJson())

// 	hosts := strings.Split(ci.NodeUrls, ";")
// 	for i := 0; i < len(hosts); i++ {
// 		if hosts[i] != cmn.GetLocalGlcUrl() {
// 			go httpPostJson(hosts[i]+conf.GetContextPath()+"/v1/log/transferAdd", jsonlog) // TODO 失败处理
// 		}
// 	}
// }

// func httpPostJson(url string, jsondata string) ([]byte, error) {

// 	// 请求
// 	req, err := http.NewRequest("POST", url, strings.NewReader(jsondata))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 请求头
// 	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
// 	req.Header.Set(conf.GetHeaderSecurityKey(), conf.GetSecurityKey())

// 	// 读取响应内容
// 	client := http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	return io.ReadAll(res.Body)
// }
