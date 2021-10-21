package models

import "github.com/piupuer/go-helper/models"

const (
	SysRoleStatusDisabled    uint   = 0
	SysRoleStatusNormal      uint   = 1
	SysRoleStatusDisabledStr string = "disabled"
	SysRoleStatusEnableStr   string = "enable"

	SysRoleSuperAdminSort uint = 0
)

var SysRoleStatusConst = map[uint]string{
	SysRoleStatusDisabled: SysRoleStatusDisabledStr,
	SysRoleStatusNormal:   SysRoleStatusEnableStr,
}

type SysRole struct {
	models.M
	Name    string    `gorm:"comment:'name'" json:"name"`
	Keyword string    `gorm:"index:idx_keyword,unique;comment:'keyword(unique str)'" json:"keyword"`
	Desc    string    `gorm:"comment:'description'" json:"desc"`
	Status  *uint     `gorm:"type:tinyint(1);default:1;comment:'status(0: disabled, 1: enable)'" json:"status"`
	Sort    *uint     `gorm:"default:1;comment:'sort(>=0, the smaller the sort, the greater the permissio, sort=0 is a super admin)'" json:"sort"`
	Menus   []SysMenu `gorm:"many2many:sys_role_menu_relation;" json:"menus"`
	Users   []SysUser `gorm:"foreignKey:RoleId"`
}

type SysRoleMenuRelation struct {
	SysMenuId uint `json:"sysMenuId"`
	SysRoleId uint `json:"sysRoleId"`
}
