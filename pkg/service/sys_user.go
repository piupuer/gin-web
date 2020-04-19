package service

import (
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
)

// 登录校验
func LoginCheck(user *models.SysUser) (models.SysUser, error) {
	var u models.SysUser
	// 查询用户名密码是否匹配
	err := global.Mysql.Where("username = ? AND password = ?", user.Username, user.Password).First(&u).Error
	return u, err
}

// 获取所有用户
func GetUsers() ([]models.SysUser, error) {
	var users []models.SysUser
	err := global.Mysql.Find(&users).Error
	return users, err
}
