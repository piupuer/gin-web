package models

// 角色与菜单关联关系
type RelationRoleMenu struct {
	SysRoleId uint `json:"sysRoleId"`
	SysMenuId uint `json:"sysMenuId"`
}

func (m RelationRoleMenu) TableName() string {
	// 多对多关系表在tag中写死, 不能加自定义表前缀
	return "relation_role_menu"
}
