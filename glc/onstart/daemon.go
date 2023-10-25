package onstart

import (
	"fmt"
	"glc/conf"
	"glc/ver"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

func init() {

	// 命令行参数解析，【-d】后台方式启动，【stop】停止，【restart】重启，【-v/version/--version/-version】查看版本
	// 内部用特殊参数，提示docker方式启动【--docker】固定为非后台方式
	daemon := false
	stop := false
	restart := false
	version := false
	docker := false
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
		if arg == "restart" {
			restart = true
		}
		if arg == "--docker" {
			docker = true
		}
		if arg == "-v" || arg == "version" || arg == "--version" || arg == "-version" {
			version = true
		}
	}
	// 查看版本
	if version {
		fmt.Printf("%s\n", "glogcenter "+ver.VERSION)
		os.Exit(0)
	}

	// 【alpine以外未足够测试，暂且默认支持alpine容器及window开发调试，可按需注释掉】
	if !cmn.IsWin() {
		info, _ := cmn.MeasureHost()
		if info == nil || !cmn.ContainsIngoreCase(info.Platform, "alpine") {
			fmt.Printf("%s\n", info.Platform)
			os.Exit(0)
		}
	}

	// 其余参数仅支持linux
	if !cmn.IsLinux() {
		return
	}

	// 端口冲突时退出
	if cmn.IsPortOpening("8080") {
		fmt.Printf("%s\n", "port 8080 conflict, startup failed.")
		os.Exit(0)
	}

	// 自动判断创建目录
	_, err := os.Stat(conf.GetStorageRoot())
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(conf.GetStorageRoot(), 0766)
	}

	// docker方式启动
	if docker {
		return
	}

	// pid 目录、文件
	pidpath := "."
	pidfile := "glc.pid"
	u, err := user.Current()
	if err == nil {
		pidpath = filepath.Join(u.HomeDir, ".glogcenter")
		os.MkdirAll(pidpath, 0766)
	}
	pidpathfile := filepath.Join(pidpath, pidfile)

	rs := checkPidFile(pidpathfile)
	if rs != nil {
		if stop || restart {
			// 退出/重启
			cmd := exec.Command("sh", "-c", "kill "+rs.Pid)
			cmd.Start()
		} else {
			// 禁止重复启动
			fmt.Printf("[PID] %s\n", rs.Pid)
			os.Exit(0)
		}
	}

	if stop {
		os.Exit(0)
	}

	if daemon {
		// cmd := exec.Command(os.Args[0], flag.Args()...)
		cmd := exec.Command(os.Args[0]) // 不再需要启动参数了

		err := cmd.Start()
		for i := 0; i < 60; i++ {
			if err != nil {
				time.Sleep(time.Duration(1) * time.Second) // 原进程没退出的话会导致启动失败，等待1秒后再试
			} else {
				break
			}
		}
		if err != nil {
			// 最多等1分钟，仍旧启动失败就放弃
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}

		fmt.Printf("[PID] %d\n", cmd.Process.Pid)
		os.Exit(0)
	} else {
		npid := fmt.Sprintf("%d", os.Getpid())
		nerr := savePid(pidpathfile, npid)
		if nerr != nil {
			cmd := exec.Command("sh", "-c", "kill "+npid)
			cmd.Start()
			os.Exit(1)
		}
	}

}
