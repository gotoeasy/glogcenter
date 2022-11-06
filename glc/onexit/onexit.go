package onexit

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gotoeasy/glang/cmn"
)

var fnExits []func()

func init() {
	go func() {
		osc := make(chan os.Signal, 1)
		signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
		<-osc
		cmn.Info("收到退出信号准备退出")
		for _, fnExit := range fnExits {
			fnExit()
		}
	}()
}

func RegisterExitHandle(fnExit func()) {
	fnExits = append(fnExits, fnExit)
}
