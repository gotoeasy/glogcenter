package html

import (
	"glc/conf"
	"glc/gweb"
	"glc/www"
	"os"

	"github.com/gotoeasy/glang/cmn"
)

// [/]重定向到[/glc/]
func RedirectToHomeController(req *gweb.HttpRequest) *gweb.HttpResult {
	defer req.Redirect(conf.GetContextPath() + "/")
	return nil
}

// 响应请求[/glc/]，读取index.html返回
func HomeIndexHtmlController(req *gweb.HttpRequest) *gweb.HttpResult {
	file, err := www.Static.ReadFile("web/dist/index.html")
	if err != nil && os.IsNotExist(err) {
		req.ResponseData(404, "text/html", cmn.StringToBytes("not found"))
	} else {
		req.ResponseData(200, "text/html", file)
	}
	return nil
}

// 响应 *.html/*.css/*.js/*.png 等静态文件请求
func StaticFileController(req *gweb.HttpRequest) *gweb.HttpResult {

	urlPath := req.RequestUrlPath()
	contentType := getContentType(urlPath)
	file, err := www.Static.ReadFile(getStaticFilePath(urlPath))
	if err != nil && os.IsNotExist(err) {
		cmn.Error("文件找不到", getStaticFilePath(urlPath), err)
		req.ResponseData(404, contentType, cmn.StringToBytes("not found"))
	} else {
		req.ResponseData(200, contentType, file)
	}
	return nil
}

// urlPath如[/glc/assets/index.f0b375ee.js]或[/assets/index.f0b375ee.js]
func getStaticFilePath(urlPath string) string {
	path := urlPath
	if conf.GetContextPath() != "" && cmn.Startwiths(urlPath, conf.GetContextPath()) {
		path = cmn.SubString(urlPath, len(conf.GetContextPath()), len(urlPath))
	}
	return "web/dist" + path
}

func getContentType(urlPath string) string {

	if cmn.Endwiths(urlPath, ".html") {
		return "text/html"
	} else if cmn.Endwiths(urlPath, ".css") {
		return "text/css"
	} else if cmn.Endwiths(urlPath, ".js") {
		return "application/x-javascript"
	} else if cmn.Endwiths(urlPath, ".png") {
		return "image/png"
	} else if cmn.Endwiths(urlPath, ".jpg") || cmn.Endwiths(urlPath, ".jpeg") {
		return "image/jpeg"
	} else if cmn.Endwiths(urlPath, ".gif") {
		return "image/gif"
	} else if cmn.Endwiths(urlPath, ".ico") {
		return "image/x-icon"
	} else if cmn.Endwiths(urlPath, ".svg") {
		return "image/svg+xml"
	} else if cmn.Endwiths(urlPath, ".json") {
		return "application/json"
	} else if cmn.Endwiths(urlPath, ".xml") {
		return "application/xml"
	} else {
		cmn.Info("未识别出ContentType，按text/html处理", urlPath)
		return "text/html"
	}

}
