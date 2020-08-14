package models

const (
	// 消息类型
	SysMessageTypeOneToOne     uint   = 0 // 一对一
	SysMessageTypeOneToMany    uint   = 1 // 一对多
	SysMessageTypeSystem       uint   = 2 // 系统
	SysMessageTypeOneToOneStr  string = "一对一"
	SysMessageTypeOneToManyStr string = "一对多"
	SysMessageTypeSystemStr    string = "系统"
)

// 定义map方便取值
var SysMessageTypeConst = map[uint]string{
	SysMessageTypeOneToOne:  SysMessageTypeOneToOneStr,
	SysMessageTypeOneToMany: SysMessageTypeOneToManyStr,
	SysMessageTypeSystem:    SysMessageTypeSystemStr,
}

// 系统消息表, 主要记录消息内容
type SysMessage struct {
	Model
	FromUserId uint       `gorm:"comment:'消息发送者'" json:"fromUserId"`
	FromUser   SysUser    `gorm:"foreignkey:FromUserId" json:"fromUser"`
	Content    string     `gorm:"comment:'消息内容'" json:"content"`
	Type       uint       `gorm:"type:tinyint(1);default:0;comment:'消息类型(0: 一对一私信, 1: 一对多群发[主要通过角色指定], 2: 系统消息)'" json:"type"`
	RoleId     uint       `gorm:"comment:'一对多角色编号'" json:"roleId"`
	Role       SysRole    `gorm:"foreignkey:RoleId" json:"role"`
	ExpiredAt  *LocalTime `gorm:"comment:'过期时间'" json:"expiredAt"`
}

func (m SysMessage) TableName() string {
	return m.Model.TableName("sys_message")
}
