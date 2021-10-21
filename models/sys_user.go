package models

import "github.com/piupuer/go-helper/models"

const (
	// 用户状态
	SysUserStatusDisabled    uint   = 0    // 禁用
	SysUserStatusNormal      uint   = 1    // 正常
	SysUserStatusDisabledStr string = "禁用" // 禁用
	SysUserStatusNormalStr   string = "正常" // 正常
)

// 定义map方便取值
var SysUserStatusConst = map[uint]string{
	SysUserStatusDisabled: SysUserStatusDisabledStr,
	SysUserStatusNormal:   SysUserStatusNormalStr,
}

// User
type SysUser struct {
	models.M
	Username     string  `gorm:"idx_username,unique;comment:'用户名'" json:"username"`
	Password     string  `gorm:"comment:'密码'" json:"password"`
	Mobile       string  `gorm:"comment:'手机'" json:"mobile"`
	Avatar       string  `gorm:"comment:'头像'" json:"avatar"`
	Nickname     string  `gorm:"comment:'昵称'" json:"nickname"`
	Introduction string  `gorm:"comment:'自我介绍'" json:"introduction"`
	Status       *uint   `gorm:"type:tinyint(1);default:1;comment:'用户状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	RoleId       uint    `gorm:"comment:'角色Id外键'" json:"roleId"`
	Role         SysRole `gorm:"foreignKey:RoleId" json:"role"` // 将SysUser.RoleId指定为外键
}
