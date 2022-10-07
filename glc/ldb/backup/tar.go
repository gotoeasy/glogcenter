package backup

import (
	"archive/tar"
	"errors"
	"fmt"
	"glc/cmn"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func TarDir(directory string, tarfilename string) error {

	if !cmn.IsExistDir(directory) {
		return errors.New("目录不存在")
	}
	dir, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	lenPrefix := cmn.LenRune(filepath.Dir(dir)) // 绝对路径除去末尾目录名后的长度

	// 创建文件
	os.MkdirAll(filepath.Dir(tarfilename), 0777)                      // 建目录确保目录存在
	f, err := os.OpenFile(tarfilename, os.O_WRONLY|os.O_CREATE, 0777) // 建文件
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	// 创建一个Writer
	writer := tar.NewWriter(f)
	defer writer.Close()

	// 遍历需要归档的目录
	return filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		// 如果是目录跳过(空目录？)
		if info.IsDir() {
			return nil
		}

		// 打开文件
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 写入文件头
		name := cmn.SubStringRune(path, lenPrefix+1, cmn.LenRune(path))
		name = strings.ReplaceAll(name, "\\", "/")
		hr := &tar.Header{
			Name:    name,          // 用相对目录名
			Format:  tar.FormatGNU, // 支持中文目录文件名
			Size:    info.Size(),
			Mode:    0777,
			ModTime: info.ModTime(),
		}

		// 将文件头写入文件中
		writer.WriteHeader(hr)
		var buff [1024]byte

		// 不断读取文件中的内容并且写入tar文件中
		for {
			n, err := f.Read(buff[:])
			writer.Write(buff[:n])
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil
			}
		}
		return nil
	})

}

// 把tar文件解压到指定目录中
func UnTar(tarFile string, dist string) error {
	if dist == "" {
		dist, _ = filepath.Abs(".") // 默认解压到当前目录
	}
	distDir, _ := filepath.Abs(dist) // 转绝对路径

	// 打开 tar 包
	fr, err := os.Open(tarFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer fr.Close()

	tr := tar.NewReader(fr)
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			fmt.Println(err)
			return err
		}
		full := filepath.Join(distDir, hdr.Name)
		if runtime.GOOS == "windows" {
			full = filepath.Join(distDir, strings.ReplaceAll(hdr.Name, "/", "\\"))
		}
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(full, 0777)
			continue
		} else {
			os.MkdirAll(filepath.Dir(full), 0777)
		}

		fw, err := os.Create(full) // 创建一个空文件，用来写入解包后的数据
		if err != nil {
			fmt.Println(err)
			return err
		}

		if _, err := io.Copy(fw, tr); err != nil {
			fmt.Println(err)
			fw.Close()
			return err
		}
		os.Chmod(full, fs.FileMode(hdr.Mode))
		fw.Close()
	}

	return nil
}
