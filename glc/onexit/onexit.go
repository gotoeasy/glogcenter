package onexit

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var fnExits []func()

func init() {
	go func() {
		osc := make(chan os.Signal, 1)
		signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
		<-osc
		log.Println("收到退出信号准备退出")
		for _, fnExit := range fnExits {
			fnExit()
		}
	}()
}

func RegisterExitHandle(fnExit func()) {
	fnExits = append(fnExits, fnExit)
}
