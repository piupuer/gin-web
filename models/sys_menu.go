package models

import "github.com/piupuer/go-helper/models"

type SysMenu struct {
	models.M
	Name       string    `gorm:"comment:'name'" json:"name"`
	Title      string    `gorm:"comment:'title'" json:"title"`
	Icon       string    `gorm:"comment:'icon'" json:"icon"`
	Path       string    `gorm:"comment:'url path'" json:"path"`
	Redirect   string    `gorm:"comment:'redirect url'" json:"redirect"`
	Component  string    `gorm:"comment:'ui component name'" json:"component"`
	Permission string    `gorm:"comment:'permission'" json:"permission"`
	Sort       *uint     `gorm:"type:int unsigned;comment:'sort(>=0)'" json:"sort"`
	Status     *uint     `gorm:"type:tinyint(1);default:1;comment:'status(0: disabled 1: enabled)'" json:"status"`
	Visible    *uint     `gorm:"type:tinyint(1);default:1;comment:'visible(0: hidden 1: visible)'" json:"visible"`
	Breadcrumb *uint     `gorm:"type:tinyint(1);default:1;comment:'breadcrumb(0: disabled 1: enabled)'" json:"breadcrumb"`
	ParentId   uint      `gorm:"default:0;comment:'parent menu id'" json:"parentId"`
	Creator    string    `gorm:"comment:'creator'" json:"creator"`
	Children   []SysMenu `gorm:"-" json:"children"`
	Roles      []SysRole `gorm:"many2many:sys_role_menu_relation;" json:"roles"`
}

func FindCheckedMenuId(list []uint, allMenu []SysMenu) []uint {
	checked := make([]uint, 0)
	for _, c := range list {
		parent := SysMenu{
			ParentId: c,
		}
		children := parent.FindChildrenId(allMenu)
		count := 0
		for _, child := range children {
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
			// all checked
			checked = append(checked, c)
		}
	}
	return checked
}

// find children menu ids
func (m SysMenu) FindChildrenId(allMenu []SysMenu) []uint {
	childrenIds := make([]uint, 0)
	for _, menu := range allMenu {
		if menu.ParentId == m.ParentId {
			childrenIds = append(childrenIds, menu.Id)
		}
	}
	return childrenIds
}
