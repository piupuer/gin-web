package response

type UploadMergeResp struct {
	Filename   string `json:"filename"`
	PreviewUrl string `json:"previewUrl"`
}

type UploadUnZipResp struct {
	Files []string `json:"files"`
}
