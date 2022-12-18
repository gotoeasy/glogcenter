package backup

import (
	"glc/com"
	"glc/conf"
	"math"
	"path/filepath"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

var backupBusy bool = false

func init() {
	if !conf.IsEnableBackup() {
		return
	}
	cmn.OnExit(onExit)
}

func Start() {
	if !conf.IsEnableBackup() {
		return
	}
	// TODO BackupStorage("sssssssss")
}

// 指定的日志仓打包为tar文件后上传Minio
func BackupStorage(storeName string) bool {
	if !conf.IsEnableBackup() {
		return false
	}

	backupBusy = true
	defer (func() { backupBusy = false })()

	// 打包tar
	dir := filepath.Join(conf.GetStorageRoot(), storeName)
	if cmn.IsExistDir(dir) {
		return false // 目录不存在
	}
	ymd := cmn.StringToUint32(cmn.Right(storeName, 8), math.MaxUint32)
	today := cmn.StringToUint32(com.GetYyyymmdd(0), 0)
	if ymd >= today {
		return false // 仅支持压缩过去日期的日志仓目录
	}

	tarfile := storeName + ".tar"
	tarfilename := filepath.Join(conf.GetStorageRoot(), ".bak", tarfile)
	err := cmn.TarDir(dir, tarfilename)
	if err != nil {
		cmn.Error(err)
		return false
	}

	// 上传Minio
	if conf.IsEnableUploadMinio() {
		group := conf.GetGlcGroup()
		minio := cmn.NewMinio(conf.GetMinioUrl(), conf.GetMinioUser(), conf.GetMinioPassword(), conf.GetMinioBucket())
		err := minio.Upload(tarfilename, group+"/"+tarfile)
		if err != nil {
			cmn.Error(err)
			return false
		}
	}
	return err == nil
}

func onExit() {
	if !backupBusy {
		return
	}

	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if !backupBusy {
			ticker.Stop()
			break
		}
	}
}
