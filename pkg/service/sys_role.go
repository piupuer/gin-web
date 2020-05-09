package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/utils"
	"strings"
)

// 注意golang坑, 返回值有数组时, 不要为了return方便使用命名返回值
// 如func GetRoles(req *request.RoleListRequestStruct) (roles []models.SysRole, err error) {
// 这会导致roles被初始化无穷大的集合, 代码无法断点调试: collecting data..., 几秒后程序异常退出:
// error layer=rpc writing response:write tcp xxx write: broken pipe
// error layer=rpc rpc:read tcp xxx read: connection reset by peer
// exit code 0
// 我曾一度以为调试工具安装配置错误, 使用其他项目代码却能稳定调试, 最终还是定位到代码本身. 踩过的坑希望大家不要再踩
// 获取所有角色
func GetRoles(req *request.RoleListRequestStruct) ([]models.SysRole, error) {
	var err error
	list := make([]models.SysRole, 0)
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
	return list, err
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
func UpdateRoleById(id uint, req gin.H) (err error) {
	var oldRole models.SysRole
	query := global.Mysql.Table(oldRole.TableName()).Where("id = ?", id).First(&oldRole)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(oldRole, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
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
