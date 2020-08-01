package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Zip(src, dst string) error {
	// 初始化给定的目录
	baseDir := CreateDirIfNotExists(src)
	CreateDirIfNotExists(dst)
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			fmt.Printf("[Zip]关闭文件失败: %v", err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名(去除baseDir)
		fh.Name = strings.TrimPrefix(path, baseDir)

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return err
		}
		// 输出压缩的内容
		fmt.Printf("[Zip]成功压缩文件: %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}

func UnZip(src, dst string) ([]string, error) {
	// 记录全部被解压的文件
	files := make([]string, 0)
	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(src)
	if err != nil {
		return files, err
	}
	defer zr.Close()

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dst != "" {
		if err := os.MkdirAll(dst, os.ModePerm); err != nil {
			return files, nil
		}
	}

	// 遍历 zr，将文件写入到磁盘
	for _, file := range zr.File {
		var decodeName string
		if file.Flags == 0 {
			// 如果标致位是0, 则是默认的本地编码gbk
			i := bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			// 如果标志为是 1 << 11也就是2048, 则是utf-8编码
			decodeName = file.Name
		}
		path := filepath.Join(dst, decodeName)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return files, nil
			}
			// 因为是目录，跳过当前循环，因为后面都是文件的处理
			continue
		}

		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			return files, nil
		}

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return files, nil
		}

		n, err := io.Copy(fw, fr)
		if err != nil {
			return files, nil
		}

		// 将解压的结果输出
		fmt.Printf("[UnZip]成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)
		// 记录解压文件名
		files = append(files, path)
		// 因为是在循环中，无法使用 defer ，直接放在最后
		// 不过这样也有问题，当出现 err 的时候就不会执行这个了，
		// 可以把它单独放在一个函数中，这里是个实验，就这样了
		fw.Close()
		fr.Close()
	}
	return files, nil
}

// 创建不存在的目录(给定文件名或文件夹名, 判断对应目录是否存在)
func CreateDirIfNotExists(name string) string {
	info, err := os.Stat(name)
	// 存在
	if err == nil {
		// 是目录
		if info.IsDir() {
			return name
		}
		// 是文件
		dir, _ := filepath.Split(name)
		return dir
	} else {
		dir, filename := filepath.Split(name)
		// 在unix系统下以.开头的为隐藏目录
		if strings.HasPrefix(filename, ".") {
			if !strings.Contains(strings.TrimPrefix(filename, "."), ".") {
				dir = name
			}
		} else {
			// 最后一级只要包含"."则认为是文件, 否则是目录
			if !strings.Contains(filename, ".") {
				dir = name
			}
		}
		// 创建不存在的目录
		os.MkdirAll(dir, os.ModePerm)
		return dir
	}
}

// 获取程序运行目录
func GetWorkDir() string {
	pwd, _ := os.Getwd()
	return pwd
}
