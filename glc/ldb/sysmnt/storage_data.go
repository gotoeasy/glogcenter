/**
 * 日志仓信息
 */
package sysmnt

import (
	"fmt"
	"glc/cmn"
	"glc/conf"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/disk"
)

type StorageResult struct {
	Info string          `json:"info,omitempty"`
	Data []*StorageModel `json:"data,omitempty"` // 占用空间
}
type StorageModel struct {
	Name       string `json:"name"`       // 名称
	LogCount   uint32 `json:"logCount"`   // 日志量
	IndexCount uint32 `json:"indexCount"` // 已建索引数量
	FileCount  uint32 `json:"fileCount"`  // 文件数量
	TotalSize  string `json:"totalSize"`  // 占用空间
}

func init() {
	go func() {
		if conf.IsStoreNameAutoAddDate() && conf.GetSaveDays() > 0 {
			removeStorageByDays()
			ticker := time.NewTicker(time.Hour) // 一小时检查一次是否有待删除的日志仓
			for {
				<-ticker.C
				removeStorageByDays()
			}
		}
	}()
}

func GetStorageList() *StorageResult {

	var datas []*StorageModel
	var sum int64
	names := cmn.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	for _, name := range names {
		d := &StorageModel{
			Name: name,
		}

		cnt, size, _ := cmn.GetDirInfo(conf.GetStorageRoot() + cmn.PathSeparator() + name)
		d.TotalSize = cmn.GetSizeInfo(uint64(size))
		d.FileCount = cnt

		sum += size

		if cnt == 0 {
			d.LogCount = 0
			d.IndexCount = 0
		} else {
			sysmntStore := NewSysmntStorage()
			d.LogCount = sysmntStore.GetStorageDataCount(name)
			d.IndexCount = sysmntStore.GetStorageIndexCount(name)
		}

		datas = append(datas, d)
	}

	stat, _ := disk.Usage(conf.GetStorageRoot())

	rs := &StorageResult{
		Info: fmt.Sprintf("共占用空间 " + cmn.GetSizeInfo(uint64(sum)) + "，剩余空间 " + cmn.GetSizeInfo(stat.Free)),
		Data: datas,
	}
	return rs
}

// 删除指定日志仓目录
func DeleteStorage(name string) error {
	// 先尝试目录改名，改成功后再删除
	oldPath := conf.GetStorageRoot() + cmn.PathSeparator() + name
	newPath := conf.GetStorageRoot() + cmn.PathSeparator() + "_x_" + name
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(newPath)
	if err != nil {
		return err
	}
	return NewSysmntStorage().DeleteStorageInfo(name)
}

func removeStorageByDays() {
	// 日志按日期分仓存储时，按保存天数自动删除
	minYmd := cmn.GetYyyymmdd(-1 * conf.GetSaveDays())
	dirs := cmn.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	for _, dir := range dirs {
		ymd := cmn.RightRune(dir, 8)
		if !cmn.StartwithsRune(ymd, "20") {
			continue
		}

		if ymd < minYmd {
			err := DeleteStorage(dir)
			if err != nil {
				log.Println("日志仓最多保存", conf.GetSaveDays(), "天，", "删除", dir, "失败", err)
			} else {
				log.Println("日志仓最多保存", conf.GetSaveDays(), "天，", "已删除", dir)
			}
		}
	}
}
