package gweb

import (
	"context"
	"fmt"
	"glc/conf"
	"glc/ldb"
	"glc/ldb/storage/logdata"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gotoeasy/glang/cmn"
)

type IgnoreGinStdoutWritter struct{}

func (w *IgnoreGinStdoutWritter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func Run() {

	gin.DisableConsoleColor()                     // 关闭Gin的日志颜色
	gin.DefaultWriter = &IgnoreGinStdoutWritter{} // 关闭Gin的默认日志输出
	gin.SetMode(gin.ReleaseMode)                  // 开启Gin的Release模式

	ginEngine := gin.Default()

	// 允许跨域
	if conf.IsEnableCors() {
		ginEngine.Use(newCors())
	}

	// 按配置判断启用GZIP压缩
	if conf.IsEnableWebGzip() {
		ginEngine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// 请求路径包含system变量，以方便代理转发控制
	ginEngine.POST(conf.GetContextPath()+"/v2/log/add/:system", func(c *gin.Context) {

		req := NewHttpRequest(c)
		if conf.IsEnableSecurityKey() && req.GetHeader(conf.GetHeaderSecurityKey()) != conf.GetSecurityKey() {
			c.JSON(http.StatusForbidden, "未经授权的访问，拒绝服务")
			return
		}

		md := &logdata.LogDataModel{}
		err := c.BindJSON(md)
		if err != nil {
			c.JSON(http.StatusOK, Error500(err.Error()))
			return
		}
		md.System = c.Param("system")

		matched, _ := regexp.MatchString(`^[0-9a-zA-Z]+$`, md.System)
		if !matched {
			cmn.Error("无效的system名： " + md.System + "，仅支持字母数字")
			c.JSON(http.StatusBadRequest, "无效的system名： "+md.System+"，仅支持字母数字")
			return
		}

		ldb.AddTextLog(md)
		c.JSON(http.StatusOK, Ok())
	})

	ginEngine.NoRoute(func(c *gin.Context) {
		req := NewHttpRequest(c)

		// filter
		filters := getFilters()
		for _, fnFilter := range filters {
			rs := fnFilter(req)
			if rs != nil {
				c.JSON(200, rs) // 过滤器返回有内容时直接返回处理结果，结束
				return
			}
		}

		// 静态文件
		path := cmn.ToLower(c.Request.URL.Path)
		if cmn.Endwiths(path, ".html") {
			path = "/**/*.html"
		} else if cmn.Endwiths(path, ".css") {
			path = "/**/*.css"
		} else if cmn.Endwiths(path, ".js") {
			path = "/**/*.js"
		} else if cmn.Endwiths(path, ".txt") {
			path = "/**/*.txt"
		} else if cmn.Endwiths(path, ".png") {
			path = "/**/*.png"
		} else if cmn.Endwiths(path, ".ico") {
			path = "/**/*.ico"
		} else if cmn.Endwiths(path, ".svg") {
			path = "/**/*.svg"
		} else if cmn.Endwiths(path, ".jpg") {
			path = "/**/*.jpg"
		} else if cmn.Endwiths(path, ".jpeg") {
			path = "/**/*.jpeg"
		} else if cmn.Endwiths(path, ".json") {
			path = "/**/*.json"
		} else if cmn.Endwiths(path, ".xml") {
			path = "/**/*.xml"
		}

		// controller
		method := cmn.ToUpper(c.Request.Method)
		handle := getHttpController(method, path)
		if handle == nil {
			if cmn.Contains(path, ".") {
				c.JSON(http.StatusNotFound, Error404()) // 有后缀名的文件请求，找不到则404
			} else {
				req.Redirect(conf.GetContextPath() + "/") // 默认统一跳转大到[/glc/]
			}
			return
		}

		rs := handle.Controller(req)
		if rs != nil {
			c.JSON(http.StatusOK, rs)
		}
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.GetServerPort()), // :8080
		Handler: ginEngine,
	}

	// 优雅退出
	cmn.OnExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmn.Info("退出Web服务")
		if err := httpServer.Shutdown(ctx); err != nil {
			cmn.Error(err)
		}
	})

	// 启动Web服务
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		cmn.Error(err.Error()) // 启动失败的话打印错误信息后退出
	}
}

func newCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	},
	)
}
