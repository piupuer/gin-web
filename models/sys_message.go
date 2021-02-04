package models

const (
	// 消息类型
	SysMessageTypeOneToOne     uint   = 0 // 一对一
	SysMessageTypeOneToMany    uint   = 1 // 一对多
	SysMessageTypeSystem       uint   = 2 // 系统
	SysMessageTypeOneToOneStr  string = "一对一"
	SysMessageTypeOneToManyStr string = "一对多"
	SysMessageTypeSystemStr    string = "系统"

	// 消息状态
	SysMessageLogStatusUnRead     uint   = 0 // 未读
	SysMessageLogStatusRead       uint   = 1 // 已读
	SysMessageLogStatusDeleted    uint   = 2 // 删除
	SysMessageLogStatusUnReadStr  string = "未读"
	SysMessageLogStatusReadStr    string = "已读"
	SysMessageLogStatusDeletedStr string = "已删除"
)

// 定义map方便取值
var (
	SysMessageTypeConst = map[uint]string{
		SysMessageTypeOneToOne:  SysMessageTypeOneToOneStr,
		SysMessageTypeOneToMany: SysMessageTypeOneToManyStr,
		SysMessageTypeSystem:    SysMessageTypeSystemStr,
	}
	SysMessageLogStatusConst = map[uint]string{
		SysMessageLogStatusUnRead:  SysMessageLogStatusUnReadStr,
		SysMessageLogStatusRead:    SysMessageLogStatusReadStr,
		SysMessageLogStatusDeleted: SysMessageLogStatusDeletedStr,
	}
)

// 系统消息表, 主要记录消息内容
type SysMessage struct {
	Model
	FromUserId uint       `gorm:"comment:'消息发送者'" json:"fromUserId"`
	FromUser   SysUser    `gorm:"foreignKey:FromUserId" json:"fromUser"`
	Title      string     `gorm:"comment:'消息标题'" json:"title"`
	Content    string     `gorm:"comment:'消息内容'" json:"content"`
	Type       uint       `gorm:"type:tinyint;default:0;comment:'消息类型(0: 一对一私信, 1: 一对多群发[主要通过角色指定], 2: 系统消息)'" json:"type"`
	RoleId     uint       `gorm:"comment:'一对多角色编号'" json:"roleId"`
	Role       SysRole    `gorm:"foreignKey:RoleId" json:"role"`
	ExpiredAt  *LocalTime `gorm:"comment:'过期时间'" json:"expiredAt"`
}

func (m SysMessage) TableName() string {
	return m.Model.TableName("sys_message")
}

// 系统消息日志, 主要记录消息接收人以及消息状态
type SysMessageLog struct {
	Model
	ToUserId  uint       `gorm:"comment:'消息接收者'" json:"toUserId"`
	ToUser    SysUser    `gorm:"foreignKey:ToUserId" json:"toUser"`
	MessageId uint       `gorm:"comment:'消息编号'" json:"messageId"`
	Message   SysMessage `gorm:"foreignKey:MessageId" json:"message"`
	Status    uint       `gorm:"type:tinyint;default:0;comment:'消息状态(0: 未读, 1: 已读, 2: 删除)'" json:"status"`
}

func (m SysMessageLog) TableName() string {
	return m.Model.TableName("sys_message_log")
}
