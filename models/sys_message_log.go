package models

const (
	// 消息状态
	SysMessageLogStatusUnRead     uint   = 0 // 未读
	SysMessageLogStatusRead       uint   = 1 // 已读
	SysMessageLogStatusDeleted    uint   = 2 // 删除
	SysMessageLogStatusUnReadStr  string = "未读"
	SysMessageLogStatusReadStr    string = "已读"
	SysMessageLogStatusDeletedStr string = "已删除"
)

// 定义map方便取值
var SysMessageLogStatusConst = map[uint]string{
	SysMessageLogStatusUnRead:  SysMessageLogStatusUnReadStr,
	SysMessageLogStatusRead:    SysMessageLogStatusReadStr,
	SysMessageLogStatusDeleted: SysMessageLogStatusDeletedStr,
}

// 系统消息日志, 主要记录消息接收人以及消息状态
type SysMessageLog struct {
	Model
	ToUserId  uint       `gorm:"comment:'消息接收者'" json:"toUserId"`
	ToUser    SysUser    `gorm:"foreignkey:ToUserId" json:"toUser"`
	MessageId uint       `gorm:"comment:'消息编号'" json:"messageId"`
	Message   SysMessage `gorm:"foreignkey:MessageId" json:"message"`
	Status    uint       `gorm:"type:tinyint(1);default:0;comment:'消息状态(0: 未读, 1: 已读, 2: 删除)'" json:"status"`
}

func (m SysMessageLog) TableName() string {
	return m.Model.TableName("sys_message_log")
}
