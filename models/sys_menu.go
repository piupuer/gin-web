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
	Children   []SysMenu `gorm:"-" json:"children"`
	Roles      []SysRole `gorm:"many2many:sys_role_menu_relation;" json:"roles"`
}
