package service

import (
	"errors"
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
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
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
	// 创建数据
	err = global.Mysql.Create(&role).Error
	return
}

// 更新角色
func UpdateRoleById(id uint, req *request.CreateRoleRequestStruct) (err error) {
	var oldRole models.SysRole
	if global.Mysql.Where("id = ?", id).First(&oldRole).RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	var role models.SysRole
	utils.CompareDifferenceStructByJson(req, oldRole, &role)

	// 更新指定列
	err = global.Mysql.Model(&oldRole).UpdateColumns(role).Error
	return
}

// 批量删除角色
func DeleteRoleByIds(ids []uint) (err error) {
	var roles []models.SysRole
	// 查询符合条件的角色, 以及关联的用户
	err = global.Mysql.Preload("Users").Where("id IN (?)", ids).Find(&roles).Error
	if err != nil {
		return
	}
	newIds := make([]uint, 0)
	for _, v := range roles {
		if len(v.Users) > 0 {
			return errors.New(fmt.Sprintf("角色[%s]仍有%d位关联用户, 请先删除用户再删除角色", v.Name, len(v.Users)))
		}
		newIds = append(newIds, v.Id)
	}
	if len(newIds) > 0 {
		// 执行删除
		err = global.Mysql.Where("id IN (?)", newIds).Delete(models.SysRole{}).Error
	}
	return
}
