package gweb

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotoeasy/glang/cmn"
)

type HttpRequest struct {
	GinCtx  *gin.Context
	mapHead map[string][]string
}

func NewHttpRequest(c *gin.Context) *HttpRequest {

	// header整理，键忽略大小写
	mapHead := make(map[string][]string)
	for k, v := range c.Request.Header {
		key := cmn.ToLower(k)
		val := mapHead[key]
		if val == nil {
			val = []string{}
		}
		val = append(val, v...)
		mapHead[key] = val
	}

	return &HttpRequest{
		GinCtx:  c,
		mapHead: mapHead,
	}
}

func (r *HttpRequest) SetHeader(key string, value string) {
	r.GinCtx.Header(key, value)
}

func (r *HttpRequest) GetToken() string {
	token := r.GetFormParameter("token")
	if token == "" {
		token = r.GetHeader("X-Access-Token")
	}
	return token
}

func (r *HttpRequest) GetHeader(name string) string {
	ary := r.mapHead[cmn.ToLower(name)]
	if ary == nil {
		return ""
	}
	return ary[0]
}

func (r *HttpRequest) GetHeaders(name string) []string {
	ary := r.mapHead[cmn.ToLower(name)]
	if ary == nil {
		return []string{}
	}
	return ary
}

func (r *HttpRequest) GetUrlParameter(name string) string {
	return r.GinCtx.Query(name)
}

func (r *HttpRequest) GetFormParameter(name string) string {
	return r.GinCtx.Request.PostFormValue(name)
}

func (r *HttpRequest) Redirect(url string) {
	r.GinCtx.Redirect(http.StatusMovedPermanently, url)
}

func (r *HttpRequest) ResponseData(code int, contentType string, bytes []byte) {
	r.GinCtx.Data(code, contentType, bytes)
}

func (r *HttpRequest) GetMethod() string {
	return r.GinCtx.Request.Method
}

func (r *HttpRequest) AbortWithStatus(code int) {
	r.GinCtx.AbortWithStatus(code)
}

func (r *HttpRequest) RequestURI() string {
	return r.GinCtx.Request.RequestURI
}

func (r *HttpRequest) RequestUrlPath() string {
	return r.GinCtx.Request.URL.Path
}

func (r *HttpRequest) BindJSON(obj any) error {
	return r.GinCtx.BindJSON(obj)
}
