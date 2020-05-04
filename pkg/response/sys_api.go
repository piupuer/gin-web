package response

import "time"

// 接口信息响应, 字段含义见models.SysRole
type ApiListResponseStruct struct {
	Id        uint      `json:"id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Category  string    `json:"category"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
}
