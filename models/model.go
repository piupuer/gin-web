package models

import (
	"fmt"
	"time"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	Id        uint       `gorm:"primary_key;comment:'自增编号'" json:"id"`
	CreatedAt time.Time  `gorm:"comment:'创建时间'" json:"created_at"`
	UpdatedAt time.Time  `gorm:"comment:'更新时间'" json:"updated_at"`
	DeletedAt *time.Time `gorm:"comment:'删除时间(软删除)'" sql:"index" json:"deleted_at"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s_%s", "tb_prefix", name)
}
