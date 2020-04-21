package initialize

import (
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/utils"
)

// 初始化数据
func InitData() {
	users := [3]models.SysUser{
		{
			Username: "admin",
			Password: utils.GenPwd("123456"),
		},
		{
			Username: "zhangsan",
			Password: utils.GenPwd("123456"),
		},
		{
			Username: "lisi",
			Password: utils.GenPwd("123456"),
		},
	}
	for _, user := range users {
		oldUser := models.SysUser{}
		err := global.Mysql.Where("username = ?", user.Username).First(&oldUser).Error
		if err != nil && err.Error() == "record not found" {
			global.Mysql.Create(&user)
		}
	}
}
