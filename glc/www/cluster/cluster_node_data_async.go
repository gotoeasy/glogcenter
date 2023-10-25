package cluster

import (
	"bufio"
	"encoding/json"
	"fmt"
	"glc/com"
	"glc/conf"
	"glc/ldb/sysmnt"
	"glc/www/service"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

type httpStoresResult struct {
	Code    int                   `json:"code,omitempty"`
	Message string                `json:"message,omitempty"`
	Success bool                  `json:"success,omitempty"`
	Result  *sysmnt.StorageResult `json:"result,omitempty"`
}

func dataAsync() {
	checkAndCopyDataFromRemote() // 检查更新本地日志仓

	// 定期数据同步（暂且忽略优雅退出）
	ticker := time.NewTicker(time.Hour * 6)
	for {
		<-ticker.C
		checkAndCopyDataFromRemote() // 检查更新本地日志仓
	}
}

func checkAndCopyDataFromRemote() {

	cmn.Info("历史数据检查同步开始")

	// 遍历其他全部节点查询日志仓列表信息，筛选出最完整日志仓信息
	mapStore := make(map[string]*sysmnt.StorageModel, 0)
	urls := getClusterUrls()
	for i := 0; i < len(urls); i++ {
		if com.GetLocalGlcUrl() == urls[i] {
			continue // 跳过本节点
		}

		// 查询日志仓列表信息
		storelist := httpGetStoresInfo(urls[i])

		// 筛选出最完整的日志仓信息
		for j := 0; j < len(storelist); j++ {
			md := storelist[j]
			mstore := mapStore[md.Name]
			if mstore == nil {
				mapStore[md.Name] = md
			} else {
				if mstore.LogCount < md.LogCount {
					mapStore[md.Name] = md
				}
			}
		}
	}

	if len(mapStore) == 0 {
		cmn.Info("历史数据检查同步结束（没有其他节点数据可对比）")
		return
	}

	// 本地
	todayStoreName := com.GeyStoreNameByDate("") // 当天日志仓名
	rs := sysmnt.GetStorageList()
	localStores := rs.Data
	mapLocalStore := make(map[string]*sysmnt.StorageModel, 0)
	for i := 0; i < len(localStores); i++ {
		if localStores[i].Name == todayStoreName {
			continue // 跳过当天的日志仓
		}
		mapLocalStore[localStores[i].Name] = localStores[i]
	}

	// 远程有，本地无，复制
	for k, remoteStore := range mapStore {
		cmn.Info("日志仓有无检查同步", remoteStore.NodeUrl, remoteStore.Name)
		if k == todayStoreName {
			cmn.Info("跳过当天的日志仓", todayStoreName)
			continue // 跳过当天的日志仓
		}

		if mapLocalStore[k] == nil {
			// 下载
			cmn.Info("本地无日志仓", k, "开始从", remoteStore.NodeUrl, "复制", remoteStore.Name)
			tarfile, err := httpDownloadStoreFile(remoteStore.NodeUrl, remoteStore.Name) // 下载
			if err != nil {
				continue
			}

			// 解压
			storeName := cmn.Split(cmn.FileName(tarfile), ".")[1] // download.logdata-20221030.1491888244752784461.tar => logdata-20221030
			distDir := filepath.Join(conf.GetStorageRoot(), storeName)
			cmn.MkdirAll(distDir)
			cmn.UnTar(tarfile, distDir)

			// 保存信息
			sysdb := sysmnt.NewSysmntStorage()
			sysdb.SetStorageDataCount(remoteStore.Name, remoteStore.LogCount)
			sysdb.SetStorageIndexCount(remoteStore.Name, remoteStore.IndexCount)
			cmn.Info("完成从", remoteStore.NodeUrl, "复制日志仓", remoteStore.Name)
		} else {
			cmn.Info("本地有日志仓", k)
		}
	}

	// 远程全，本地缺，覆盖
	for i := 0; i < len(localStores); i++ {
		md := localStores[i]
		cmn.Info("日志仓数据检查同步", md.NodeUrl, md.Name)
		if md.Name == todayStoreName {
			cmn.Info("跳过当天的日志仓", todayStoreName)
			continue // 跳过当天的日志仓
		}

		mstore := mapStore[md.Name]
		if mstore == nil || md.LogCount >= mstore.LogCount {
			cmn.Info("本地完整，跳过", md.Name)
			continue // 本地更完整
		}

		// 下载
		cmn.Info("开始从", mstore.NodeUrl, "复制日志仓", mstore.Name)
		tarfile, err := httpDownloadStoreFile(mstore.NodeUrl, mstore.Name) // 下载
		if err != nil {
			continue
		}

		// 本地先删除（日志仓使用中会删除失败，忽略，待下次同步处理）
		err = sysmnt.DeleteStorage(mstore.Name)
		if err != nil {
			cmn.Error("本地日志仓", mstore.Name, "删除失败", err)
			continue
		}

		// 解压
		storeName := cmn.Split(cmn.FileName(tarfile), ".")[1] // download.logdata-20221030.1491888244752784461.tar => logdata-20221030
		distDir := filepath.Join(conf.GetStorageRoot(), storeName)
		cmn.MkdirAll(distDir)
		cmn.UnTar(tarfile, distDir)

		// 保存信息
		sysdb := sysmnt.NewSysmntStorage()
		sysdb.SetStorageDataCount(mstore.Name, mstore.LogCount)
		sysdb.SetStorageIndexCount(mstore.Name, mstore.IndexCount)

		// 删除下载的临时文件
		os.Remove(tarfile)
		cmn.Info("完成从", mstore.NodeUrl, "复制日志仓", mstore.Name)
	}

	cmn.Info("历史数据检查同步完成")
}

func httpDownloadStoreFile(serverUrl string, storeName string) (string, error) {

	tarfile := "download." + storeName + "." + fmt.Sprintf("%d", time.Now().UnixNano()) + ".tar" // download.logdata-20221030.1491888244752784461.tar
	tarfilename := filepath.Join(conf.GetStorageRoot(), ".tmp", tarfile)                         // %ROOT%/.tmp/download.logdata-20221030.1491888244752784461.tar
	os.MkdirAll(filepath.Dir(tarfilename), 0666)                                                 // 建目录
	file, err := os.Create(tarfilename)
	if err != nil {
		cmn.Error("创建下载文件", tarfilename, "失败", err)
		return "", err
	}

	// 请求
	req, err := http.NewRequest("GET", serverUrl+conf.GetContextPath()+"/sys/cluster/down?storeName="+storeName, strings.NewReader(""))
	if err != nil {
		cmn.Error("从", serverUrl, "下载日志仓", storeName, "数据文件失败", err)
		return "", err
	}

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		cmn.Error("从", serverUrl, "下载日志仓", storeName, "数据文件失败", err)
		return "", err
	}
	defer res.Body.Close()

	reader := bufio.NewReaderSize(res.Body, 64*1024)
	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, reader)
	if err != nil {
		cmn.Error("从", serverUrl, "下载日志仓", storeName, "数据文件失败", err)
		return "", err
	}

	return tarfilename, nil
}

func httpGetStoresInfo(serverUrl string) []*sysmnt.StorageModel {

	// 请求
	req, err := http.NewRequest("POST", serverUrl+conf.GetContextPath()+"/v1/store/list", strings.NewReader(""))
	if err != nil {
		cmn.Error("从", serverUrl, "取日志仓信息失败", err)
		return nil
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set(conf.GetHeaderSecurityKey(), conf.GetSecurityKey())

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		cmn.Error("从", serverUrl, "取日志仓信息失败", err)
		return nil
	}
	defer res.Body.Close()

	by, err := io.ReadAll(res.Body)
	if err != nil {
		cmn.Error("从", serverUrl, "取日志仓信息失败", err)
		return nil
	}

	rs := &httpStoresResult{}
	json.Unmarshal(by, rs)

	if !rs.Success {
		cmn.Error("从", serverUrl, "取日志仓信息失败", rs.Message)
		return nil
	}

	return rs.Result.Data
}

func getClusterUrls() []string {
	kv, err := service.GetSysmntItem(KEY_CLUSTER)
	if err != nil {
		cmn.Error(err)
		return nil
	}

	ci := &ClusterInfo{}
	ci.LoadJson(kv.Value)

	return cmn.Split(ci.NodeUrls, ";")
}
