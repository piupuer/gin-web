package initialize

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	global.Log.Info("数据库连接DSN: ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("初始化mysql异常: %v", err))
	}
	global.Mysql = db
	// 表结构
	autoMigrate()
	global.Log.Info("初始化mysql完成")
	// 初始化数据库日志监听器
	binlog()
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
		new(models.RelationUserWorkflowLine),
		new(models.SysLeave),
		new(models.SysOperationLog),
		new(models.SysMessage),
		new(models.SysMessageLog),
	)
}

func binlog() {
	MysqlBinlog([]string{
		new(models.SysUser).TableName(),
		new(models.SysRole).TableName(),
		new(models.SysMenu).TableName(),
		new(models.SysApi).TableName(),
		new(models.SysCasbin).TableName(),
		new(models.RelationRoleMenu).TableName(),
		new(models.SysWorkflow).TableName(),
		new(models.SysWorkflowLine).TableName(),
		new(models.SysWorkflowLog).TableName(),
		new(models.RelationUserWorkflowLine).TableName(),
		new(models.SysLeave).TableName(),
		new(models.SysOperationLog).TableName(),
		new(models.SysMessage).TableName(),
		new(models.SysMessageLog).TableName(),
	})
}
