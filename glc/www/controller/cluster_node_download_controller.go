package controller

import (
	"fmt"
	"glc/conf"
	"glc/gweb"
	"path/filepath"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

// 打包下载指定日志仓数据
func ClusterDownloadStoreDataController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	storeName := req.GetUrlParameter("storeName")
	// 打包tar
	dir := filepath.Join(conf.GetStorageRoot(), storeName)
	if !cmn.IsExistDir(dir) {
		return nil // 目录不存在
	}

	tarfile := storeName + "-" + fmt.Sprintf("%d", time.Now().UnixNano()) + ".tar" // logdata-20221030-1491888244752784461.tar
	tarfilename := filepath.Join(conf.GetStorageRoot(), ".tmp", tarfile)           // %ROOT%/.tmp/logdata-20221030-1491888244752784461.tar，会自动建目录

	err := cmn.TarDir(dir, tarfilename)
	if err != nil {
		cmn.Error(err)
		return nil // 打包失败
	}

	// 下载
	ctx := req.GinCtx
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+storeName+".tar")
	ctx.Header("Content-Length", "-1")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.File(tarfilename)

	return nil
}
