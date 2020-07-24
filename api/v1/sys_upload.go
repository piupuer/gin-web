package v1

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/siddontang/go/ioutil2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 判断文件块是否存在
func UploadFileChunkExists(c *gin.Context) {
	var filePart request.FilePartInfo
	_ = c.Bind(&filePart)
	// 校验请求
	err := validateReq(filePart)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	filePart.Complete, filePart.Uploaded = getUploadedChunkNumbers(filePart)
	response.SuccessWithData(filePart)
}

// 合并分片文件
func UploadMerge(c *gin.Context) {
	var filePart request.FilePartInfo
	_ = c.Bind(&filePart)
	// 通过文件唯一标识找确定文件
	// 获取块文件名
	chunkName := filePart.GetChunkFilename(filePart.CurrentCheckChunkNumber)
	chunkDir, _ := filepath.Split(chunkName)
	// 创建merge file
	mergeFile, err := os.OpenFile(fmt.Sprintf("%s/%s", chunkDir, filePart.Filename), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	defer mergeFile.Close()

	totalChunk := int(filePart.GetTotalChunk())
	chunkSize := int(filePart.ChunkSize)
	var chunkNumbers []int
	for i := 0; i < totalChunk; i++ {
		chunkNumbers = append(chunkNumbers, i+1)
	}

	// 开启协程并发合并文件
	// 如果文件块总数过大, 性能反而降低, 因此需要配置一个合适的协程数
	var count = int(global.Conf.Upload.MergeConcurrentCount)
	chunkCount := len(chunkNumbers) / count
	lastChunkCount := 0
	if len(chunkNumbers)%count > 0 || count == 1 {
		lastChunkCount = len(chunkNumbers)%count + chunkCount
	}
	// 转为二维数组, 每一组数据分配给一个协程使用
	chunks := make([][]int, count)
	for i := 0; i < count; i++ {
		if i < count-1 {
			chunks[i] = chunkNumbers[i*chunkCount : (i+1)*chunkCount]
		} else {
			chunks[i] = chunkNumbers[i*chunkCount : i*chunkCount+lastChunkCount]
		}
	}
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(arr []int) {
			defer wg.Done()
			for _, item := range arr {
				func() {
					currentChunkName := filePart.GetChunkFilename(uint(item))
					exists := ioutil2.FileExists(currentChunkName)
					if exists {
						// 读取文件分片
						f, err := os.OpenFile(currentChunkName, os.O_RDONLY, os.ModePerm)
						if err != nil {
							response.FailWithMsg(err.Error())
							return
						}
						defer func() {
							// 关闭文件
							f.Close()
							// 删除分片
							os.Remove(currentChunkName)
						}()
						b, err := ioutil.ReadAll(f)
						if err != nil {
							response.FailWithMsg(err.Error())
							return
						}
						// 从指定位置开始写
						mergeFile.WriteAt(b, int64((item-1)*chunkSize))
					}
				}()
			}
		}(chunks[i])
	}
	// 等待协程全部处理结束
	wg.Wait()
	
	// 回写文件信息
	var res response.UploadMergeResponseStruct
	res.Filename = chunkDir + filePart.Filename
	response.SuccessWithData(res)
}

// 上传文件(小文件直接是单个文件, 若是超大文件可能是单个分片)
func UploadFile(c *gin.Context) {
	// 限制文件最大内存(二进制移位xxxMB)
	err := c.Request.ParseMultipartForm(int64(global.Conf.Upload.SingleMaxSize) << 20)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("文件大小超出最大值%dMB", global.Conf.Upload.SingleMaxSize))
		return
	}
	// 读取文件分片
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMsg("无法读取文件")
		return
	}

	// 读取文件分片参数
	var filePart request.FilePartInfo
	// 当前大小
	currentSize := uint(header.Size)
	filePart.CurrentSize = &currentSize
	// 块编号
	filePart.ChunkNumber = utils.Str2Uint(strings.TrimSpace(c.Request.FormValue("chunkNumber")))
	// 块大小
	filePart.ChunkSize = utils.Str2Uint(strings.TrimSpace(c.Request.FormValue("chunkSize")))
	// 总大小
	filePart.TotalSize = utils.Str2Uint(strings.TrimSpace(c.Request.FormValue("totalSize")))
	// 唯一标识
	filePart.Identifier = strings.TrimSpace(c.Request.FormValue("identifier"))
	// 文件名
	filePart.Filename = strings.TrimSpace(c.Request.FormValue("filename"))

	// 校验请求
	err = validateReq(filePart)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// 获取块文件名
	chunkName := filePart.GetChunkFilename(filePart.ChunkNumber)
	// 创建不存在的文件夹
	chunkDir, _ := filepath.Split(chunkName)
	err = os.MkdirAll(chunkDir, os.ModePerm)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// 保存块文件
	out, err := os.Create(chunkName)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// 检查文件块完整性
	filePart.CurrentCheckChunkNumber = 1
	filePart.Complete = checkChunkComplete(filePart)
	// 回写响应数据
	response.SuccessWithData(filePart)
}

// 请求校验
func validateReq(filePart request.FilePartInfo) error {
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

// 检查文件块, 主要用于判断文件完整性
func checkChunkComplete(filePart request.FilePartInfo) bool {
	currentChunkName := filePart.GetChunkFilename(filePart.CurrentCheckChunkNumber)
	exists := ioutil2.FileExists(currentChunkName)
	if exists {
		filePart.CurrentCheckChunkNumber++
		if filePart.CurrentCheckChunkNumber > filePart.GetTotalChunk() {
			// 完成全部传输
			return true
		}
		// 继续
		return checkChunkComplete(filePart)
	}
	// 完成当前块
	return false
}

// 获取已上传完成的块number集合
func getUploadedChunkNumbers(filePart request.FilePartInfo) (bool, []uint) {
	totalChunk := filePart.GetTotalChunk()
	var currentChunkNumber uint = 1
	uploadedChunkNumbers := make([]uint, 0)
	for {
		currentChunkName := filePart.GetChunkFilename(currentChunkNumber)
		exists := ioutil2.FileExists(currentChunkName)
		if exists {
			uploadedChunkNumbers = append(uploadedChunkNumbers, currentChunkNumber)
		}
		currentChunkNumber++
		if currentChunkNumber > totalChunk {
			break
		}
	}
	return len(uploadedChunkNumbers) == int(totalChunk), uploadedChunkNumbers
}
