/**
 * 日志仓信息
 */
package sysmnt

import (
	"fmt"
	"glc/cmn"
	"glc/conf"
	"strconv"
)

type StorageData struct {
	Name       string
	DataCount  uint32
	IndexCount uint32
	FileCount  uint32
	TotalSize  string
}

func (s *SysmntStorage) GetStorageData() []*StorageData {

	var rs []*StorageData
	names := cmn.GetStorageNames(conf.GetStorageRoot(), ".sysmnt")
	for _, name := range names {
		d := &StorageData{
			Name: name,
		}

		cnt, size, _ := cmn.GetDirInfo(conf.GetStorageRoot() + cmn.PathSeparator() + name)
		sizem, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1024/1024), 64)
		d.TotalSize = fmt.Sprintf("%fM", sizem)
		d.FileCount = cnt
		rs = append(rs, d)
	}
	return rs
}
