package gweb

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Run() {
	server := gin.Default()

	server.NoRoute(func(c *gin.Context) {
		req := NewHttpRequest(c)
		path := strings.ToLower(c.Request.URL.Path)
		method := strings.ToUpper(c.Request.Method)

		handle := getHttpController(method, path)
		if handle == nil {
			c.JSON(http.StatusNotFound, Error404())
			return
		}

		rs := handle.Controller(req)
		c.JSON(rs.Code, rs)
	})

	server.Run(":8080")

}
