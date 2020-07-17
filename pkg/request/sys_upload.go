package request

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"regexp"
	"time"
)

// 文件分片信息结构体
type FilePartInfo struct {
	// 当前文件大小(post传输时值非空)
	CurrentSize *uint `json:"-"`
	// 当前块编号(检查文件完整性会用到)
	CurrentCheckChunkNumber uint `json:"-"`
	// 已上传的块编号
	Uploaded []uint `json:"uploaded"`
	// 是否传输完成
	Complete bool `json:"complete"`
	// 块编号
	ChunkNumber uint `json:"chunkNumber" form:"chunkNumber"`
	// 每块大小
	ChunkSize uint `json:"chunkSize" form:"chunkSize"`
	// 总大小
	TotalSize uint `json:"totalSize" form:"totalSize"`
	// 唯一标识(其中包含文件类型)
	Identifier string `json:"identifier" form:"identifier"`
	// 文件名
	Filename string `json:"filename" form:"filename"`
}

// 从文件唯一标识符中去除特殊字符, 只保留数字字母
func (s *FilePartInfo) CleanIdentifier() string {
	// 从文件标识获取文件类型
	re, _ := regexp.Compile("[^0-9A-Za-z_-]")
	return re.ReplaceAllString(s.Identifier, "")
}

// 获取文件块总数
func (s *FilePartInfo) GetTotalChunk() uint {
	// 余数部分将与最后一块合并, 而不是+1块
	// 105 / 25 => 4块
	// 100 / 25 => 4块
	// 99 / 25 => 3块
	// 24 / 25 => 1块
	if s.ChunkSize > 0 && s.TotalSize > s.ChunkSize {
		return s.TotalSize / s.ChunkSize
	}
	return 1
}

// 获取块文件名
func (s *FilePartInfo) GetChunkFilename(chunkNumber uint) string {
	// 清理特殊字符
	identifier := s.CleanIdentifier()
	// 定义块文件名
	return fmt.Sprintf(
		"%s/%s/uploader-%s.chunk%d",
		global.Conf.Upload.SaveDir,
		models.LocalTime{
			Time: time.Now(),
		}.DateString(),
		identifier,
		chunkNumber,
	)
}
