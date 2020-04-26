package models

// User
type SysUser struct {
	Model
	Username     string  `gorm:"comment:'用户名'" json:"username"`
	Password     string  `gorm:"comment:'密码'" json:"password"`
	Mobile       string  `gorm:"comment:'手机'" json:"mobile"`
	Avatar       string  `gorm:"comment:'头像'" json:"avatar"`
	Nickname     string  `gorm:"comment:'昵称'" json:"nickname"`
	Introduction string  `gorm:"comment:'自我介绍'" json:"introduction"`
	RoleId       uint    `gorm:"comment:'角色Id外键'" json:"role_id"`
	Role         SysRole `gorm:"foreignkey:RoleId" json:"role"` // 将SysUser.RoleId指定为外键
}

func (m SysUser) TableName() string {
	return m.Model.TableName("sys_user")
}
