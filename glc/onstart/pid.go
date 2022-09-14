package onstart

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PidFile struct {
	Path  string // pid文件路径
	Pid   string // pid
	IsNew bool   // 是否新
	Err   error  // error
}

// 指定路径下生成pid文件，文件内容为pid，已存在时检查pid有效性
func NewPid(pathfile string) *PidFile {

	// 运行中时直接返回
	if opid := checkPidFile(pathfile); opid != nil {
		return opid
	}

	// 创建文件
	if err := os.MkdirAll(filepath.Dir(pathfile), os.FileMode(0755)); err != nil {
		log.Println("create pid file failed", pathfile)
		return &PidFile{
			Path: pathfile,
			Err:  err,
		}
	}

	// 保存PID
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := os.WriteFile(pathfile, []byte(pid), 0644); err != nil {
		log.Println("save pid file failed", pathfile)
		return &PidFile{
			Path: pathfile,
			Err:  err,
		}
	}

	// 成功创建后返回
	return &PidFile{
		Path:  pathfile,
		Pid:   pid,
		IsNew: false,
		Err:   nil,
	}

}

func checkPidFile(path string) *PidFile {
	if pidByte, err := os.ReadFile(path); err == nil {
		pid := strings.TrimSpace(string(pidByte))
		if _, err := os.Stat(filepath.Join("/proc", pid)); err == nil {
			return &PidFile{
				Path:  path,
				Pid:   pid,
				IsNew: false,
			}
		}
	}
	return nil
}
