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
