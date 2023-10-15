package main

import (
	"glc/conf"
	"glc/onstart"
	"runtime"

	"github.com/gotoeasy/glang/cmn"
)

func main() {
	cmn.SetGlcClient(cmn.NewGlcClient(&cmn.GlcOptions{Enable: false, LogLevel: "INFO"})) // 控制台INFO日志级别输出

	runtime.GOMAXPROCS(conf.GetGoMaxProcess()) // 使用最大CPU数量
	onstart.Run()
}
