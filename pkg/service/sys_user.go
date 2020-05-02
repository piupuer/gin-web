package service

import (
	"errors"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/utils"
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

// 获取所有用户
func GetUsers() ([]models.SysUser, error) {
	users := make([]models.SysUser, 0)
	err := global.Mysql.Find(&users).Error
	return users, err
}
