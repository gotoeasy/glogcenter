package gweb

import (
	"fmt"
	"glc/ldb/conf"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Run() {
	server := gin.Default()

	server.NoRoute(func(c *gin.Context) {
		req := NewHttpRequest(c)

		// filter
		filters := getFilters()
		for _, fnFilter := range filters {
			rs := fnFilter(req)
			if rs != nil {
				c.JSON(rs.Code, rs) // 过滤器返回有内容时直接返回处理结果，结束
				return
			}
		}

		// controller
		path := strings.ToLower(c.Request.URL.Path)
		method := strings.ToUpper(c.Request.Method)
		handle := getHttpController(method, path)
		if handle == nil {
			c.JSON(http.StatusNotFound, Error404())
			return
		}

		rs := handle.Controller(req)
		if rs != nil {
			c.JSON(rs.Code, rs)
		}
	})

	server.Run(fmt.Sprintf(":%d", conf.GetServerPort()))

}
