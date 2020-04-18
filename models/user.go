package models

import (
	"github.com/jinzhu/gorm"
)

// User
type User struct {
	gorm.Model        // gorm提供了基础字段CreatedAt/UpdatedAt/DeletedAt可直接继承
	Username   string `gorm:"comment:'用户名'" json:"username"` // 用户名
	Sex        int    `gorm:"comment:'性别'" json:"sex"`       // 性别
}

// 修改默认表名
func (User) TableName() string {
	return "user"
}
