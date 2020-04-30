package service

import (
	"fmt"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/utils"
	"strings"
)

// 获取所有角色
func GetRoles(req *request.RoleListRequestStruct) (list []models.SysRole, err error) {
	db := global.Mysql
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
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

// 创建角色
func CreateRole(req *request.CreateRoleRequestStruct) (err error) {
	var role models.SysRole
	utils.Struct2StructByJson(req, &role)
	// 查询条数
	err = global.Mysql.Create(&role).Error
	return
}
