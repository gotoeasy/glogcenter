package controller

import (
	"glc/conf"
	"glc/gweb"
	"glc/www/service"

	"github.com/gotoeasy/glang/cmn"
)

func ClusterGetItemController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	pkv := &service.KeyValue{}
	err := req.BindJSON(pkv)
	if err != nil || pkv.Key == "" {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	dkv, err := service.GetSysmntItem(pkv.Key)
	if err != nil || pkv.Key == "" {
		cmn.Error(err.Error())
		return gweb.Error500(err.Error())
	}

	return gweb.Result(dkv)
}

func ClusterSetItemController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	pkv := &service.KeyValue{}
	err := req.BindJSON(pkv)
	if err != nil || pkv.Key == "" {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	_, err = service.SetSysmntItem(pkv)
	if err != nil {
		cmn.Error(err.Error())
		return gweb.Error500(err.Error())
	}

	return gweb.Ok()
}

func ClusterDelItemController(req *gweb.HttpRequest) *gweb.HttpResult {
	if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
		return gweb.Error(403, "未经授权的访问，拒绝服务")
	}

	pkv := &service.KeyValue{}
	err := req.BindJSON(pkv)
	if err != nil || pkv.Key == "" {
		cmn.Error("请求参数有误", err)
		return gweb.Error500(err.Error())
	}

	err = service.DelSysmntItem(pkv.Key)
	if err != nil {
		cmn.Error(err.Error())
		return gweb.Error500(err.Error())
	}

	return gweb.Ok()
}
