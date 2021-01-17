package response

// 上传文件合并响应, 字段含义见models
type UploadMergeResponseStruct struct {
	Filename   string `json:"filename"`
	PreviewUrl string `json:"previewUrl"`
}

// 上传解压zip响应, 字段含义见models
type UploadUnZipResponseStruct struct {
	Files []string `json:"files"`
}
