package backup

import (
	"glc/cmn"
	"glc/conf"
	"glc/onexit"
	"log"
	"math"
	"path/filepath"
	"time"
)

var backupBusy bool = false

func init() {
	if !conf.IsEnableBackup() {
		return
	}
	onexit.RegisterExitHandle(onExit)
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
	ymd := cmn.StringToUint32(cmn.RightRune(storeName, 8), math.MaxUint32)
	today := cmn.StringToUint32(cmn.GetYyyymmdd(0), 0)
	if ymd >= today {
		return false // 仅支持压缩过去日期的日志仓目录
	}

	tarfile := storeName + ".tar"
	tarfilename := filepath.Join(conf.GetStorageRoot(), ".bak", tarfile)
	err := TarDir(dir, tarfilename)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 上传Minio
	if conf.IsEnableUploadMinio() {
		group := conf.GetGlcGroup()
		err := UploadMinio(tarfilename, group+"/"+tarfile)
		if err != nil {
			log.Println(err.Error())
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
