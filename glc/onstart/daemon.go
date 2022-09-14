package onstart

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func init() {

	// 仅支持linux
	if runtime.GOOS != "linux" {
		return
	}

	pidfile := "~/.gologcenter/glc.pid"

	// 操作系统是linux时，支持以命令行参数【-d】后台方式启动，【stop】停止
	daemon := false
	stop := false
	for index, arg := range os.Args {
		if index == 0 {
			continue
		}
		if arg == "-d" {
			daemon = true
		}
		if arg == "stop" {
			stop = true
		}
	}

	// 停止
	if stop {
		tmpPid := readPid(pidfile)
		if tmpPid != "" {
			exec.Command("sh", "-c", "kill "+tmpPid)
		}
		os.Exit(0)
	}

	// 启动守护进程
	if daemon {

		// 已启动时忽略
		chk := checkPidFile(pidfile)
		if chk != nil {
			log.Println(chk.Pid)
			os.Exit(0) // 进程已存在，不重复启动
		}

		// 启动失败则退出
		cmd := exec.Command(os.Args[0], flag.Args()...)
		if err := cmd.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}

		// 启动成功则保存pid
		daemonPid := fmt.Sprintf("%d", cmd.Process.Pid)
		err := savePid(pidfile, daemonPid)
		if err != nil {
			// 保存pid失败则退出
			exec.Command("sh", "-c", "kill "+daemonPid)
			os.Exit(1) // 创建或保存pid文件失败
		}

		log.Println(daemonPid)
		os.Exit(0)
	}

	// 普通启动
	// 已启动时忽略
	chk := checkPidFile(pidfile)
	if chk != nil {
		log.Println(chk.Pid)
		os.Exit(0) // 进程已存在，不重复启动
	}

	// 保存pid失败则退出
	pid := fmt.Sprintf("%d", os.Getpid())
	err := savePid(pidfile, pid)
	if err != nil {
		exec.Command("sh", "-c", "kill "+pid)
		os.Exit(1) // 创建或保存pid文件失败
	}

}
