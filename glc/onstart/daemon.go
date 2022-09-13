package onstart

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func init() {

	// 操作系统是linux时，支持以命令行参数【-d】后台方式启动
	daemon := false
	for index, arg := range os.Args {
		if runtime.GOOS == "linux" && index > 0 && arg == "-d" {
			daemon = true
			break
		}
	}

	if daemon {
		cmd := exec.Command(os.Args[0], flag.Args()...)
		if err := cmd.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		fmt.Printf("%s [PID] %d running\n", os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
}
