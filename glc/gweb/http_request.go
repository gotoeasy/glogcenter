package gweb

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpRequest struct {
	ginCtx  *gin.Context
	mapHead map[string][]string
}

func NewHttpRequest(c *gin.Context) *HttpRequest {

	// header整理，键忽略大小写
	mapHead := make(map[string][]string)
	for k, v := range c.Request.Header {
		key := strings.ToLower(k)
		val := mapHead[key]
		if val == nil {
			val = []string{}
		}
		val = append(val, v...)
		mapHead[key] = val
	}

	return &HttpRequest{
		ginCtx:  c,
		mapHead: mapHead,
	}
}

func (r *HttpRequest) GetHeader(name string) string {
	ary := r.mapHead[strings.ToLower(name)]
	if ary == nil {
		return ""
	}
	return ary[0]
}

func (r *HttpRequest) GetHeaders(name string) []string {
	ary := r.mapHead[strings.ToLower(name)]
	if ary == nil {
		return []string{}
	}
	return ary
}

func (r *HttpRequest) GetUrlParameter(name string) string {
	return r.ginCtx.Query(name)
}

func (r *HttpRequest) GetFormParameter(name string) string {
	return r.ginCtx.Request.PostFormValue(name)
}

func (r *HttpRequest) Redirect(url string) {
	r.ginCtx.Redirect(http.StatusMovedPermanently, url)
}
