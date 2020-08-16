package request

import (
	"gin-web/pkg/response"
)

// 获取消息列表结构体
type MessageListRequestStruct struct {
	ToUserId          uint   `json:"toUserId"`
	Title             string `json:"title" form:"title"`
	Content           string `json:"content" form:"content"`
	Type              *uint  `json:"type" form:"type"`
	Status            *uint  `json:"status" form:"status"`
	response.PageInfo        // 分页参数
}
