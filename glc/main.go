package main

import (
	"glc/conf"
	"glc/onstart"
	"runtime"

	"github.com/gotoeasy/glang/cmn"
)

func main() {
	cmn.SetGlcClient(cmn.NewGlcClient(&cmn.GlcOptions{
		EnableConsoleLog: cmn.GetEnvStr("GLC_ENABLE_CONSOLE_LOG", "false"), // 关闭控制台日志输出
		LogLevel:         cmn.GetEnvStr("GLC_LOG_LEVEL", "INFO"),           // 控制台INFO日志级别输出
	}))

	runtime.GOMAXPROCS(conf.GetGoMaxProcess()) // 使用最大CPU数量
	onstart.Run()
}
