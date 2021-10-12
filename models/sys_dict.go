package models

import "github.com/piupuer/go-helper/models"

// 系统字典
type SysDict struct {
	models.Model
	Name      string        `gorm:"index:idx_name_unique,unique;comment:'字典名称(一般是英文)'" json:"name"`
	Desc      string        `gorm:"comment:'字典描述(一般是中文)'" json:"desc"`
	Status    *uint         `gorm:"type:tinyint(1);default:1;comment:'角色状态(正常/禁用, 默认正常)'" json:"status"`
	Remark    string        `gorm:"comment:'备注'" json:"remark"`
	DictDatas []SysDictData `gorm:"foreignKey:DictId;comment:'字典数据'" json:"dictDatas"`
}

// 系统字典数据
type SysDictData struct {
	models.Model
	Key      string  `gorm:"comment:'数据键'" json:"key"`
	Val      string  `gorm:"comment:'数据值'" json:"val"`
	Attr     string  `gorm:"comment:'数据属性(主要是前端区分input框类型)'" json:"attr"`
	Addition string  `gorm:"comment:'附加参数(不同的数据有不同的自定义附加参数, 非必填)'" json:"addition"`
	Sort     *uint   `gorm:"comment:'数据排序'" json:"sort"`
	Status   *uint   `gorm:"type:tinyint(1);default:1;comment:'角色状态(正常/禁用, 默认正常)'" json:"status"`
	Remark   string  `gorm:"comment:'备注'" json:"remark"`
	DictId   uint    `gorm:"comment:'所属字典编号'" json:"dictId"`
	Dict     SysDict `gorm:"foreignKey:DictId" json:"dict"`
}
