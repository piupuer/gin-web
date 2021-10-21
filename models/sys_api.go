package models

import "github.com/piupuer/go-helper/models"

type SysApi struct {
	models.M
	Method   string `gorm:"comment:'request method'" json:"method"`
	Path     string `gorm:"comment:'api path'" json:"path"`
	Category string `gorm:"comment:'api group category'" json:"category"`
	Desc     string `gorm:"comment:'api description'" json:"desc"`
}
