package response

import "gin-web/models"

// 获取消息列表结构体
type MessageListResponseStruct struct {
	Id           uint             `json:"id"`
	Status       uint             `json:"status"`
	ToUserId     uint             `json:"toUserId"`
	ToUsername   string           `json:"toUsername"`
	Type         uint             `json:"type"`
	Title        string           `json:"title"`
	Content      string           `json:"content"`
	CreatedAt    models.LocalTime `json:"createdAt"`
	FromUserId   uint             `json:"fromUserId"`
	FromUsername string           `json:"fromUsername"`
}
