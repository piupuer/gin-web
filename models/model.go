package models

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/golang-module/carbon"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	Id        uint                    `gorm:"primaryKey;comment:'自增编号'" json:"id"`
	CreatedAt carbon.ToDateTimeString `gorm:"comment:'创建时间'" json:"createdAt"`
	UpdatedAt carbon.ToDateTimeString `gorm:"comment:'更新时间'" json:"updatedAt"`
	DeletedAt DeletedAt               `gorm:"index:idx_deleted_at;comment:'删除时间(软删除)'" json:"deletedAt"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, name)
}
