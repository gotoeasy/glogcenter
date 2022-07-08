package status

var mapStorageStatus map[string]string // 日志仓是否正在使用

func init() {
	mapStorageStatus = make(map[string]string)
}

func UpdateStorageStatus(name string, open bool) {
	if open {
		mapStorageStatus[name] = "1"
	} else {
		delete(mapStorageStatus, name)
	}
}

func IsStorageOpening(name string) bool {
	return mapStorageStatus[name] == "1"
}
