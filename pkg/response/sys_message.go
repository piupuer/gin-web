package response

// 获取消息列表结构体
type MessageListResponseStruct struct {
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

// websocket消息响应
type MessageWsResponseStruct struct {
	// 消息类型, 见const
	Type string `json:"type"`
	// 消息详情
	Detail Resp `json:"detail"`
}
