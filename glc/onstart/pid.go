package onstart

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

type PidFile struct {
	Path  string // pid文件路径
	Pid   string // pid
	IsNew bool   // 是否新
	Err   error  // error
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
		cmn.Error("create pid file failed", path)
		return err
	}

	if err := os.WriteFile(path, []byte(pid), 0644); err != nil {
		cmn.Error("save pid file failed", path, err)
		return err
	}

	return nil
}
