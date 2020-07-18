package response

import (
	"gin-web/models"
	"time"
)

// 接口信息响应, 字段含义见models
type OperationLogListResponseStruct struct {
	Id        uint             `json:"id"`
	ApiDesc    string        `json:"apiDesc"`
	Path       string        `json:"path"`
	Method     string        `json:"method"`
	Params     string        `json:"params"`
	Body       string        `json:"body"`
	Data       string        `json:"data"`
	Status     int           `json:"status"`
	Username   string        `json:"username"`
	RoleName   string        `json:"roleName"`
	Ip         string        `json:"ip"`
	IpLocation string        `json:"ipLocation"`
	Latency    time.Duration `json:"latency"`
	UserAgent  string        `json:"userAgent"`
	CreatedAt models.LocalTime `json:"createdAt"`
}

