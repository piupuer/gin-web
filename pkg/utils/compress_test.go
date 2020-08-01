package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCompressJpg(t *testing.T) {
	src := "/Users/piupuer/Downloads/images"
	filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}
		// 压缩图片
		err := CompressImage(path)
		if err != nil {
			fmt.Printf("压缩图片出错: %v\n", err)
		}
		return nil
	})
}

func TestCompressImageSaveOriginal(t *testing.T) {
	src := "/Users/piupuer/Downloads/images"
	filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}
		// 压缩图片
		err := CompressImageSaveOriginal(path, ".before")
		if err != nil {
			fmt.Printf("压缩图片出错: %v\n", err)
		}
		return nil
	})
}
