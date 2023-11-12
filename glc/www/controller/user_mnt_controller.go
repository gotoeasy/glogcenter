package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/ldb/sysmnt"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

type SysUserResult struct {
	Info string            `json:"info,omitempty"`
	Data []*sysmnt.SysUser `json:"data,omitempty"` // 占用空间
}

// 查询用户列表
func UserListController(req *gweb.HttpRequest) *gweb.HttpResult {

	token := req.GetToken()
	if (!InWhiteList(req) && InBlackList(req)) || !conf.IsEnableLogin() || (conf.IsEnableLogin() && GetUsernameByToken(token) != conf.GetUsername()) {
		return gweb.Error403() // 黑名单检查、登录检查、管理员检查
	}
	if conf.IsEnableLogin() {
		catchSession.Set(token, GetUsernameByToken(token)) // 会话重新计时
	}

	var users []*sysmnt.SysUser
	mnt := sysmnt.NewSysmntStorage()
	names := mnt.GetSysUsernames()
	for i := 0; i < len(names); i++ {
		user := mnt.GetSysUser(names[i])
		if user != nil {
			user.Password = "" // 密码不回传前端
			users = append(users, user)
		}
	}

	return gweb.Result(&SysUserResult{Data: users})
}

// 保存用户（JSON提交方式）
func UserSaveController(req *gweb.HttpRequest) *gweb.HttpResult {

	if (!InWhiteList(req) && InBlackList(req)) || !conf.IsEnableLogin() || (conf.IsEnableLogin() && GetUsernameByToken(req.GetToken()) != conf.GetUsername()) {
		return gweb.Error403() // 黑名单检查、登录检查、管理员检查
	}

	user := &sysmnt.SysUser{}
	err := req.BindJSON(user)
	if err != nil {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	if user.Username == "" || user.Username == conf.GetUsername() || cmn.Contains(user.Username, ":") {
		return gweb.Error500("无效的用户名")
	}

	mnt := sysmnt.NewSysmntStorage()
	if user.CreateDate == "" {
		// 新建
		names := mnt.GetSysUsernames()
		for i := 0; i < len(names); i++ {
			if cmn.EqualsIngoreCase(user.Username, names[i]) {
				return gweb.Error500("用户名 " + user.Username + " 已存在") // 重复性检查
			}
		}
		user.CreateDate = time.Now().Format("2006-01-02 15:04:05")
		user.UpdateDate = user.CreateDate
	} else {
		// 更新
		user.UpdateDate = time.Now().Format("2006-01-02 15:04:05")
	}
	if cmn.Trim(user.Systems) == "" {
		user.Systems = "*" // 未设定时可访问全部系统日志
	}

	var ary []string
	var m map[string]bool = make(map[string]bool)
	systems := cmn.Split(cmn.ReplaceAll(user.Systems, "，", ","), ",")
	for i := 0; i < len(systems); i++ {
		s := cmn.Trim(systems[i])
		if s != "" && !m[cmn.ToLower(s)] {
			ary = append(ary, s)
			m[cmn.ToLower(s)] = true
		}
	}
	user.Systems = cmn.Join(ary, ",")

	err = mnt.SaveSysUser(user)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500("处理失败")
	}

	if conf.IsClusterMode() {
		go TransferGlc(conf.SysUserTransferSave, user.ToJson()) // 转发其他GLC服务
	}

	return gweb.Ok()
}

// 修改自己密码（JSON提交方式）
func UserChangePswController(req *gweb.HttpRequest) *gweb.HttpResult {

	token := req.GetToken()
	username := GetUsernameByToken(token)
	if (!InWhiteList(req) && InBlackList(req)) || !conf.IsEnableLogin() || username == "" {
		return gweb.Error403() // 黑名单检查、登录检查、管理员检查
	}

	user := &sysmnt.SysUser{}
	err := req.BindJSON(user)
	if err != nil || user.Username != username || user.OldPassword == "" || user.NewPassword == "" {
		cmn.Error("请求参数有误", err)
		return gweb.Error500("操作失败")
	}

	if user.Username == conf.GetUsername() {
		if user.OldPassword != conf.GetPassword() {
			return gweb.Error500("操作失败")
		}
		conf.SetPassword(user.NewPassword)
		return gweb.Ok()
	}

	mnt := sysmnt.NewSysmntStorage()
	sysuser := mnt.GetSysUser(user.Username)
	if sysuser == nil {
		cmn.Error("找不到用户：", user.Username)
		return gweb.Error500("操作失败")
	}

	if user.OldPassword != sysuser.Password {
		return gweb.Error500("操作失败")
	}

	sysuser.Password = user.NewPassword
	err = mnt.SaveSysUser(sysuser)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500("操作失败")
	}

	if conf.IsClusterMode() {
		go TransferGlc(conf.SysUserTransferChgPsw, sysuser.ToJson()) // 转发其他GLC服务
	}

	return gweb.Ok()
}

// 删除用户（JSON提交方式）
func UserDelController(req *gweb.HttpRequest) *gweb.HttpResult {

	if (!InWhiteList(req) && InBlackList(req)) || !conf.IsEnableLogin() || (conf.IsEnableLogin() && GetUsernameByToken(req.GetToken()) != conf.GetUsername()) {
		return gweb.Error403() // 黑名单检查、登录检查、管理员检查
	}

	user := &sysmnt.SysUser{}
	err := req.BindJSON(user)
	if err != nil {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	err = sysmnt.NewSysmntStorage().DeleteSysUser(user)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500("处理失败")
	}

	// 删除掉的用户清空相关会话
	token, find := catchSession.Get(user.Username)
	if find {
		catchSession.Delete(token.(string))
		catchSession.Delete(user.Username)
	}

	if conf.IsClusterMode() {
		go TransferGlc(conf.UserTransferLogin, user.ToJson()) // 转发其他GLC服务
	}

	return gweb.Ok()
}

// 保存用户（来自数据转发）
func UserTransferSaveController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	user := &sysmnt.SysUser{}
	req.BindJSON(user)
	if cmn.EqualsIngoreCase(user.Username, conf.GetUsername()) {
		return gweb.Error500("无效的用户名") // 禁止和管理员账号冲突
	}

	err := sysmnt.NewSysmntStorage().SaveSysUser(user)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500("处理失败")
	}

	return gweb.Ok()
}

// 修改自己密码（来自数据转发）
func UserTransferChangePswController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	user := &sysmnt.SysUser{}
	err := req.BindJSON(user)
	if err != nil {
		cmn.Error("请求参数有误", err)
		return gweb.Error500("操作失败")
	}

	if user.Username == conf.GetUsername() {
		conf.SetPassword(user.NewPassword)
		return gweb.Ok()
	}

	mnt := sysmnt.NewSysmntStorage()
	sysuser := mnt.GetSysUser(user.Username)
	if sysuser == nil {
		cmn.Error("找不到用户：", user.Username)
		return gweb.Error500("操作失败")
	}

	sysuser.Password = user.NewPassword
	err = mnt.SaveSysUser(sysuser)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500("操作失败")
	}

	return gweb.Ok()
}

// 删除用户（来自数据转发）
func UserTransferDelController(req *gweb.HttpRequest) *gweb.HttpResult {

	// 开启API秘钥校验时才检查
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	user := &sysmnt.SysUser{}
	req.BindJSON(user)

	err := sysmnt.NewSysmntStorage().DeleteSysUser(user)
	if err != nil {
		cmn.Error(err)
		return gweb.Error500("处理失败")
	}

	return gweb.Ok()
}
