/**
 * 日志仓信息
 */
package sysmnt

import (
	"fmt"
	"glc/cmn"
	"glc/conf"
	"glc/ldb/storage"
	"glc/ldb/storage/indexword"
	"os"
)

type StorageResult struct {
	Total string          `json:"total,omitempty"` // 名称
	Free  uint32          `json:"free,omitempty"`  // 日志量
	Data  []*StorageModel `json:"data,omitempty"`  // 占用空间
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
	names := cmn.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	for _, name := range names {
		d := &StorageModel{
			Name: name,
		}

		cnt, size, _ := cmn.GetDirInfo(conf.GetStorageRoot() + cmn.PathSeparator() + name)
		d.TotalSize = fmt.Sprintf("%.0fM", float64(size)/1024/1024)
		d.FileCount = cnt

		if cnt == 0 {
			d.LogCount = 0
			d.IndexCount = 0
		} else {
			store := storage.NewLogDataStorageHandle(name)
			d.LogCount = store.TotalCount()
			widx := indexword.NewWordIndexStorage(name)
			d.IndexCount = widx.GetIndexedCount()
		}

		datas = append(datas, d)
	}

	rs := &StorageResult{
		Data: datas,
	}
	return rs
}

// 删除指定日志仓目录
func DeleteStorage(name string) error {
	return os.RemoveAll(conf.GetStorageRoot() + cmn.PathSeparator() + name)
}
