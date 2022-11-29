package main

import (
	"glc/conf"
	"glc/onstart"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(conf.GetGoMaxProcess()) // 使用最大CPU数量
	onstart.Run()
}
