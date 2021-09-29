package response

type MessageResp struct {
	BaseData
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
	Type   string `json:"type"`   // 消息类型
	Detail Resp   `json:"detail"` // 消息详情
}
