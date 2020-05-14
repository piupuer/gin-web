package tests

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jinzhu/gorm"
)

// 初始化mysql数据库
func Mysql() {
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		global.Conf.Mysql.Username,
		global.Conf.Mysql.Password,
		global.Conf.Mysql.Host,
		global.Conf.Mysql.Port,
		global.Conf.Mysql.Database,
		global.Conf.Mysql.Query,
	))
	if err != nil {
		panic(fmt.Sprintf("[单元测试]初始化mysql异常: %v", err))
	}
	global.Mysql = db
	// 表结构
	autoMigrate()
	// 打印所有执行的sql
	db.LogMode(global.Conf.Mysql.LogMode)
	global.Log.Debug("[单元测试]初始化mysql完成")
}

// 自动迁移表结构
func autoMigrate() {
	global.Mysql.AutoMigrate(
		new(models.SysUser),
		new(models.SysRole),
		new(models.SysMenu),
		new(models.SysApi),
		new(models.SysCasbin),
	)
}
