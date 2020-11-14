package models

import (
	"fmt"
	"gin-web/pkg/global"
)

// 角色与菜单关联关系
type RelationRoleMenu struct {
	SysRoleId uint `json:"sysRoleId"`
	SysMenuId uint `json:"sysMenuId"`
}

func (m RelationRoleMenu) TableName() string {
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, "relation_role_menu")
}
