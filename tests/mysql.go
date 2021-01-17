package tests

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 初始化mysql数据库
func Mysql() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		global.Conf.Mysql.Username,
		global.Conf.Mysql.Password,
		global.Conf.Mysql.Host,
		global.Conf.Mysql.Port,
		global.Conf.Mysql.Database,
		global.Conf.Mysql.Charset,
		global.Conf.Mysql.Collation,
		global.Conf.Mysql.Query,
	)
	global.Log.Debug("[单元测试]数据库连接DSN: ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		// 指定表前缀
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: global.Conf.Mysql.TablePrefix + "_",
		},
	})
	if err != nil {
		panic(fmt.Sprintf("[单元测试]初始化mysql异常: %v", err))
	}
	global.Mysql = db
	// 表结构
	autoMigrate()
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
		new(models.SysWorkflow),
		new(models.SysWorkflowLine),
		new(models.SysWorkflowLog),
		new(models.SysOperationLog),
		new(models.SysMessage),
		new(models.SysMessageLog),
		new(models.SysMachine),
	)
}
