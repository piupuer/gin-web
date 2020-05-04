package initialize

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jinzhu/gorm"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
)

// 初始化mysql数据库
func Mysql() {
	db, err := gorm.Open("mysql", "root:root@tcp(localserver:43306)/goshipment?charset=utf8&parseTime=True&loc=Local&timeout=10000ms")
	if err != nil {
		panic(fmt.Sprintf("初始化mysql异常: %v", err))
	}
	global.Mysql = db
	// 表结构
	autoMigrate()
	// 打印所有执行的sql
	db.LogMode(true)
	global.Log.Debug("初始化mysql完成")
}

// 自动迁移表结构
func autoMigrate() {
	global.Mysql.AutoMigrate(
		new(models.SysUser),
		new(models.SysRole),
		new(models.SysMenu),
		new(models.SysApi),
	)
}
