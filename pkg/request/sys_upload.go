package request

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/golang-module/carbon"
	"regexp"
)

const (
	// chunk tmp path
	ChunkTmpPath = "chunks"
)

type FilePartInfoReq struct {
	CurrentSize *uint `json:"-"`
	CurrentCheckChunkNumber uint `json:"-"`
	// uploaded block numbers
	Uploaded []uint `json:"uploaded"`
	// whether transfer complete
	Complete bool `json:"complete"`
	ChunkNumber uint `json:"chunkNumber" form:"chunkNumber"`
	ChunkSize uint `json:"chunkSize" form:"chunkSize"`
	TotalSize uint `json:"totalSize" form:"totalSize"`
	Identifier string `json:"identifier" form:"identifier"`
	Filename string `json:"filename" form:"filename"`
}

// Remove special characters
func (s *FilePartInfoReq) CleanIdentifier() string {
	re, _ := regexp.Compile("[^0-9A-Za-z_-]")
	return re.ReplaceAllString(s.Identifier, "")
}

func (s *FilePartInfoReq) GetTotalChunk() uint {
	// The remainder will be merged with the last block instead of + 1
	// 105 / 25 => 4 chunk
	// 100 / 25 => 4 chunk
	// 99 / 25 => 3 chunk
	// 24 / 25 => 1 chunk
	if s.ChunkSize > 0 && s.TotalSize > s.ChunkSize {
		return s.TotalSize / s.ChunkSize
	}
	return 1
}

func (s *FilePartInfoReq) GetChunkFilename(chunkNumber uint) string {
	identifier := s.CleanIdentifier()
	return fmt.Sprintf(
		"%s/%s/%s/uploader-%s/chunk%d",
		global.Conf.Upload.SaveDir,
		carbon.Now().ToDateString(),
		ChunkTmpPath,
		identifier,
		chunkNumber,
	)
}

func (s *FilePartInfoReq) GetChunkFilenameWithoutChunkNumber() string {
	identifier := s.CleanIdentifier()
	return fmt.Sprintf(
		"%s/%s/%s/uploader-%s/chunk",
		global.Conf.Upload.SaveDir,
		carbon.Now().ToDateString(),
		ChunkTmpPath,
		identifier,
	)
}

func (s *FilePartInfoReq) GetUploadRootPath() string {
	return fmt.Sprintf(
		"%s/%s",
		global.Conf.Upload.SaveDir,
		carbon.Now().ToDateString(),
	)
}

func (s *FilePartInfoReq) GetChunkRootPath() string {
	identifier := s.CleanIdentifier()
	return fmt.Sprintf(
		"%s/%s/%s/uploader-%s",
		global.Conf.Upload.SaveDir,
		carbon.Now().ToDateString(),
		ChunkTmpPath,
		identifier,
	)
}

func (s *FilePartInfoReq) ValidateReq() error {
	filePart := s
	if filePart == nil {
		return fmt.Errorf("file params invalid")
	}
	if filePart.ChunkNumber == 0 ||
		filePart.ChunkSize == 0 ||
		filePart.TotalSize == 0 ||
		filePart.Identifier == "" ||
		filePart.Filename == "" {
		return fmt.Errorf("file name or file size invalid")
	}

	totalChunk := filePart.GetTotalChunk()
	if filePart.ChunkNumber > totalChunk {
		return fmt.Errorf("file chunk number invalid")
	}

	if filePart.CurrentSize != nil {
		if int64(*filePart.CurrentSize) > int64(global.Conf.Upload.SingleMaxSize)<<20 {
			return fmt.Errorf("the file size exceeds the maximum: %dMB, current: %dB", global.Conf.Upload.SingleMaxSize, int64(*filePart.CurrentSize))
		}

		if filePart.ChunkNumber < totalChunk && *filePart.CurrentSize != filePart.ChunkSize {
			return fmt.Errorf("inconsistent file block size: [%d:%d]", filePart.CurrentSize, filePart.ChunkSize)
		}

		if totalChunk > 1 &&
			filePart.ChunkNumber == totalChunk &&
			*filePart.CurrentSize != filePart.TotalSize%filePart.ChunkSize+filePart.ChunkSize {
			return fmt.Errorf("inconsistent file last block size: [%d:%d]", filePart.CurrentSize, filePart.TotalSize%filePart.ChunkSize+filePart.ChunkSize)
		}
		if totalChunk == 1 &&
			*filePart.CurrentSize != filePart.TotalSize {
			return fmt.Errorf("inconsistent file first block size: [%d:%d]", filePart.CurrentSize, filePart.TotalSize)
		}
	}
	return nil
}
