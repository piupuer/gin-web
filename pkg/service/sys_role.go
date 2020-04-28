package service

import (
	"fmt"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
)

// 获取所有角色
func GetRoles(req *request.RoleListRequestStruct) (list []models.SysRole, err error) {
	db := global.Mysql
	if req.Name != "" {
		db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Name))
	}
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		// 获取分页参数
		limit, offset := req.GetLimit()
		err = db.Limit(limit).Offset(offset).Find(&list).Error
	}
	return
}
