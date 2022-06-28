package gweb

import (
	"context"
	"fmt"
	"glc/ldb/conf"
	"glc/onexit"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	ginEngine := gin.Default()

	ginEngine.NoRoute(func(c *gin.Context) {
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

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.GetServerPort()), // :8080
		Handler: ginEngine,
	}

	// 优雅退出
	onexit.RegisterExitHandle(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println("退出Web服务:", err)
		} else {
			log.Println("退出Web服务")
		}
	})

	// 启动Web服务
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("%s", err) // 启动失败的话打印错误信息后退出
	}
}
