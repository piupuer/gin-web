package initialize

import (
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/utils"
)

// 初始化数据
func InitData() {
	// 默认头像
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := [3]models.SysUser{
		{
			Username:     "admin",
			Password:     utils.GenPwd("123456"),
			Mobile:       "18888888888",
			Avatar:       avatar,
			Nickname:     "管理员",
			Introduction: "妖怪, 哪里跑",
		},
		{
			Username:     "zhangsan",
			Password:     utils.GenPwd("123456"),
			Mobile:       "15888888888",
			Avatar:       avatar,
			Nickname:     "张三",
			Introduction: "今天是个好日子",
		},
		{
			Username:     "lisi",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13888888888",
			Avatar:       avatar,
			Nickname:     "李四",
			Introduction: "这个人很懒, 什么也没留下",
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
