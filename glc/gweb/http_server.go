package gweb

import (
	"context"
	"fmt"
	"glc/cmn"
	"glc/ldb/conf"
	"glc/onexit"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Run() {

	gin.SetMode(gin.ReleaseMode) // 开启Release模式

	ginEngine := gin.Default()

	// 按配置判断启用GZIP压缩
	if conf.IsEnableWebGzip() {
		ginEngine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

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

		// 静态文件
		path := strings.ToLower(c.Request.URL.Path)
		if cmn.EndwithsRune(path, ".html") {
			path = "/**/*.html"
		} else if cmn.EndwithsRune(path, ".css") {
			path = "/**/*.css"
		} else if cmn.EndwithsRune(path, ".js") {
			path = "/**/*.js"
		} else if cmn.EndwithsRune(path, ".png") {
			path = "/**/*.png"
		} else if cmn.EndwithsRune(path, ".ico") {
			path = "/**/*.ico"
		}

		// controller
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
		log.Println("退出Web服务")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	})

	// 启动Web服务
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("%s", err) // 启动失败的话打印错误信息后退出
	}
}
