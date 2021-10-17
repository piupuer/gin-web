package response

import "github.com/piupuer/go-helper/pkg/resp"

type MessageResp struct {
	resp.Base
	Status       uint   `json:"status"`
	ToUserId     uint   `json:"toUserId"`
	ToUsername   string `json:"toUsername"`
	Type         uint   `json:"type"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	FromUserId   uint   `json:"fromUserId"`
	FromUsername string `json:"fromUsername"`
}

type MessageWsResp struct {
	Type   string    `json:"type"`   // 消息类型
	Detail resp.Resp `json:"detail"` // 消息详情
}
