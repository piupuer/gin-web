package models

import (
	"github.com/jinzhu/gorm"
)

// User
type SysUser struct {
	gorm.Model          // gorm提供了基础字段CreatedAt/UpdatedAt/DeletedAt可直接继承
	Username     string `gorm:"comment:'用户名'" json:"username"`
	Password     string `gorm:"comment:'密码'" json:"password"`
	Mobile       string `gorm:"comment:'手机'" json:"mobile"`
	Avatar       string `gorm:"comment:'头像'" json:"avatar"`
	Nickname     string `gorm:"comment:'昵称'" json:"nickname"`
	Introduction string `gorm:"comment:'自我介绍'" json:"introduction"`
}
