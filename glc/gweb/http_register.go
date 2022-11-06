package gweb

import (
	"glc/gweb/method"

	"github.com/gotoeasy/glang/cmn"
)

type HttpController struct {
	Method     string
	Path       string
	Controller func(*HttpRequest) *HttpResult
}

var mapHandleGet map[string]*HttpController
var mapHandlePost map[string]*HttpController
var filters []func(*HttpRequest) *HttpResult

func init() {
	mapHandleGet = make(map[string]*HttpController)
	mapHandlePost = make(map[string]*HttpController)
}

func getHttpController(methodType string, path string) *HttpController {
	switch methodType {
	case method.GET:
		return mapHandleGet[cmn.ToLower(path)]
	case method.POST:
		return mapHandlePost[cmn.ToLower(path)]
	default:
		return nil
	}
}

func RegisterController(methodType string, path string, fnController func(*HttpRequest) *HttpResult) {
	pathLower := cmn.ToLower(path) // path匹配比较忽略大小写
	if mapHandleGet[pathLower] != nil || mapHandlePost[pathLower] != nil {
		panic("duplicate controller path: " + path)
	}

	r := &HttpController{
		Method:     methodType,
		Path:       pathLower,
		Controller: fnController,
	}

	switch methodType {
	case method.GET:
		mapHandleGet[pathLower] = r
	case method.POST:
		mapHandlePost[pathLower] = r
	default:
		panic("unsuport method: " + methodType)
	}

}

func RegisterFilter(fnFilter func(*HttpRequest) *HttpResult) {
	filters = append(filters, fnFilter)
}

func getFilters() []func(*HttpRequest) *HttpResult {
	return filters
}
