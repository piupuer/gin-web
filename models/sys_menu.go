package models

import "github.com/piupuer/go-helper/models"

// 系统菜单表
type SysMenu struct {
	models.Model
	Name       string    `gorm:"comment:'菜单名称(英文名, 可用于国际化)'" json:"name"`
	Title      string    `gorm:"comment:'菜单标题(无法国际化时使用)'" json:"title"`
	Icon       string    `gorm:"comment:'菜单图标'" json:"icon"`
	Path       string    `gorm:"comment:'菜单访问路径'" json:"path"`
	Redirect   string    `gorm:"comment:'重定向路径'" json:"redirect"`
	Component  string    `gorm:"comment:'前端组件路径'" json:"component"`
	Permission string    `gorm:"comment:'权限标识'" json:"permission"`
	Sort       *uint     `gorm:"type:int unsigned;comment:'菜单顺序(同级菜单, 从0开始, 越小显示越靠前)'" json:"sort"`
	Status     *uint     `gorm:"type:tinyint(1);default:1;comment:'菜单状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Visible    *uint     `gorm:"type:tinyint(1);default:1;comment:'菜单可见性(可见/隐藏, 默认可见)'" json:"visible"`
	Breadcrumb *uint     `gorm:"type:tinyint(1);default:1;comment:'面包屑可见性(可见/隐藏, 默认可见)'" json:"breadcrumb"`
	ParentId   uint      `gorm:"default:0;comment:'父菜单编号(编号为0时表示根菜单)'" json:"parentId"`
	Creator    string    `gorm:"comment:'创建人'" json:"creator"`
	Children   []SysMenu `gorm:"-" json:"children"`                              // 子菜单集合
	Roles      []SysRole `gorm:"many2many:sys_menu_role_relation;" json:"roles"` // 角色菜单多对多关系
}

// 获取选中列表
func GetCheckedMenuIds(list []uint, allMenu []SysMenu) []uint {
	checked := make([]uint, 0)
	for _, c := range list {
		// 获取子流水线
		parent := SysMenu{
			ParentId: c,
		}
		children := parent.GetChildrenIds(allMenu)
		// 判断子流水线是否全部在create中
		count := 0
		for _, child := range children {
			// 避免环包调用, 不再调用utils
			// if utils.ContainsUint(list, child) {
			// 	count++
			// }
			contains := false
			for _, v := range list {
				if v == child {
					contains = true
				}
			}
			if contains {
				count++
			}
		}
		if len(children) == count {
			// 全部选中
			checked = append(checked, c)
		}
	}
	return checked
}

// 查找子菜单编号
func (m SysMenu) GetChildrenIds(allMenu []SysMenu) []uint {
	childrenIds := make([]uint, 0)
	for _, menu := range allMenu {
		if menu.ParentId == m.ParentId {
			childrenIds = append(childrenIds, menu.Id)
		}
	}
	return childrenIds
}
