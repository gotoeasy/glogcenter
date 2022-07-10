/**
 * 日志仓信息
 */
package sysmnt

import (
	"fmt"
	"glc/cmn"
	"glc/conf"
	"os"

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
		Info: fmt.Sprintf("合计占用空间 " + cmn.GetSizeInfo(uint64(sum)) + "，剩余空间 " + cmn.GetSizeInfo(stat.Free)),
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
