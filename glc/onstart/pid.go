package onstart

import (
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
func NewPid(pathfile string, pid string) *PidFile {

	// 运行中时直接返回
	if opid := checkPidFile(pathfile); opid != nil {
		return opid
	}

	// 保存PID
	if err := savePid(pathfile, pid); err != nil {
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
	pid := readPid(path)
	if pid == "" {
		return nil
	}
	if _, err := os.Stat(filepath.Join("/proc", pid)); err == nil {
		return &PidFile{
			Path:  path,
			Pid:   pid,
			IsNew: false,
		}
	}
	return nil
}

func readPid(path string) string {
	if pidByte, err := os.ReadFile(path); err == nil {
		return strings.TrimSpace(string(pidByte))
	}
	return ""
}

func savePid(path string, pid string) error {

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		log.Println("create pid file failed", path)
		return err
	}

	if err := os.WriteFile(path, []byte(pid), 0644); err != nil {
		log.Println("save pid file failed", path)
		return err
	}
	return nil
}
