package request

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"regexp"
	"time"
)

const (
	// 上传分片文件临时目录
	ChunkTmpPath = "chunks"
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
		"%s/%s/%s/uploader-%s/chunk%d",
		global.Conf.Upload.SaveDir,
		models.LocalTime{
			Time: time.Now(),
		}.DateString(),
		ChunkTmpPath,
		identifier,
		chunkNumber,
	)
}

// 获取块文件名(不带分片编号)
func (s *FilePartInfo) GetChunkFilenameWithoutChunkNumber() string {
	// 清理特殊字符
	identifier := s.CleanIdentifier()
	// 定义块文件名
	return fmt.Sprintf(
		"%s/%s/%s/uploader-%s/chunk",
		global.Conf.Upload.SaveDir,
		models.LocalTime{
			Time: time.Now(),
		}.DateString(),
		ChunkTmpPath,
		identifier,
	)
}

// 获取块文件名(不带临时目录)
func (s *FilePartInfo) GetUploadRootPath() string {
	// 定义块文件名
	return fmt.Sprintf(
		"%s/%s",
		global.Conf.Upload.SaveDir,
		models.LocalTime{
			Time: time.Now(),
		}.DateString(),
	)
}

// 获取块文件名(不带临时目录)
func (s *FilePartInfo) GetChunkRootPath() string {
	// 清理特殊字符
	identifier := s.CleanIdentifier()
	// 定义块文件名
	return fmt.Sprintf(
		"%s/%s/%s/uploader-%s",
		global.Conf.Upload.SaveDir,
		models.LocalTime{
			Time: time.Now(),
		}.DateString(),
		ChunkTmpPath,
		identifier,
	)
}

// 请求校验
func (s *FilePartInfo) ValidateReq() error {
	filePart := s
	if filePart == nil {
		return fmt.Errorf("文件参数不合法")
	}
	// 文件大小不能为0
	if filePart.ChunkNumber == 0 ||
		filePart.ChunkSize == 0 ||
		filePart.TotalSize == 0 ||
		filePart.Identifier == "" ||
		filePart.Filename == "" {
		return fmt.Errorf("文件名称或大小不合法")
	}

	// 块编号不能超出总块数
	totalChunk := filePart.GetTotalChunk()
	if filePart.ChunkNumber > totalChunk {
		return fmt.Errorf("文件块编号不合法")
	}

	// 继续比较当前文件大小
	if filePart.CurrentSize != nil {
		// 不能超出文件大小最大值
		if int64(*filePart.CurrentSize) > int64(global.Conf.Upload.SingleMaxSize)<<20 {
			return fmt.Errorf("文件大小超出最大值%dMB, 当前%dB", global.Conf.Upload.SingleMaxSize, int64(*filePart.CurrentSize))
		}

		// 正常块, 当前文件大小必须等于块大小
		if filePart.ChunkNumber < totalChunk && *filePart.CurrentSize != filePart.ChunkSize {
			return fmt.Errorf("文件块大小不一致[%d:%d]", filePart.CurrentSize, filePart.ChunkSize)
		}

		// 当前块为最后一块
		// 总块数>1
		if totalChunk > 1 &&
			filePart.ChunkNumber == totalChunk &&
			*filePart.CurrentSize != filePart.TotalSize%filePart.ChunkSize+filePart.ChunkSize {
			return fmt.Errorf("文件块大小不一致(末尾块)[%d:%d]", filePart.CurrentSize, filePart.TotalSize%filePart.ChunkSize+filePart.ChunkSize)
		}
		// 总块数=1
		if totalChunk == 1 &&
			*filePart.CurrentSize != filePart.TotalSize {
			return fmt.Errorf("文件块大小不一致(首块)[%d:%d]", filePart.CurrentSize, filePart.TotalSize)
		}
	}
	return nil
}
