package models

import (
	"github.com/jinzhu/gorm"
)

// User
type SysUser struct {
	gorm.Model        // gorm提供了基础字段CreatedAt/UpdatedAt/DeletedAt可直接继承
	Username   string `gorm:"comment:'用户名'" json:"username"`
	Password   string `gorm:"comment:'密码'" json:"password"`
}
