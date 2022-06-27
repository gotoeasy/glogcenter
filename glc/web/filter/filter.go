package filter

import (
	"glc/gweb"
	"glc/ldb/conf"
)

// 校验HEADER的API秘钥
func ApiKeyFilter(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if !conf.IsEnableSecurityKey() {
		return nil
	}

	auth := req.GetHeader(conf.GetHeaderSecurityKey())
	if auth != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}
	return nil // 返回nil表示正常过滤成功
}
