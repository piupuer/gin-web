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

// 登录校验
func LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	var u models.SysUser
	// 查询用户
	err := global.Mysql.Where("username = ?", user.Username).First(&u).Error
	if err != nil {
		return nil, err
	}
	// 校验密码
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, errors.New("用户名或密码错误")
	}
	return &u, err
}

// 获取用户
func GetUsers(req *request.UserListRequestStruct) ([]models.SysUser, error) {
	var err error
	list := make([]models.SysUser, 0)
	db := global.Mysql
	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
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

// 创建用户
func CreateUser(req *request.CreateUserRequestStruct) (err error) {
	var user models.SysUser
	utils.Struct2StructByJson(req, &user)
	// 创建数据
	err = global.Mysql.Create(&user).Error
	return
}

// 更新用户
func UpdateUserById(id uint, req *request.CreateUserRequestStruct) (err error) {
	var oldUser models.SysUser
	if global.Mysql.Where("id = ?", id).First(&oldUser).RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	var user models.SysUser
	utils.CompareDifferenceStructByJson(req, oldUser, &user)

	// 更新指定列
	err = global.Mysql.Model(&oldUser).UpdateColumns(user).Error
	return
}

// 批量删除用户
func DeleteUserByIds(ids []uint) (err error) {
	return global.Mysql.Where("id IN (?)", ids).Delete(models.SysUser{}).Error
}
