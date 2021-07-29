package initialize

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
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
	// 隐藏密码
	showDsn := fmt.Sprintf(
		"%s:******@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		global.Conf.Mysql.Username,
		global.Conf.Mysql.Host,
		global.Conf.Mysql.Port,
		global.Conf.Mysql.Database,
		global.Conf.Mysql.Charset,
		global.Conf.Mysql.Collation,
		global.Conf.Mysql.Query,
	)
	global.Log.Info("数据库连接DSN: ", showDsn)
	init := false
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				if !init {
					panic(fmt.Sprintf("初始化mysql异常: 连接超时(%ds)", global.Conf.System.ConnectTimeout))
				}
				// 此处需return避免协程空跑
				return
			}
		}
	}()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		// 指定表前缀
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: global.Conf.Mysql.TablePrefix + "_",
		},
		// 查询全部字段, 某些情况下*不走索引
		QueryFields: true,
	})
	if err != nil {
		panic(fmt.Sprintf("初始化mysql异常: %v", err))
	}
	init = true
	// 开启mysql日志
	if global.Conf.Mysql.LogMode {
		db = db.Debug()
	}
	global.Mysql = db
	// 表结构
	autoMigrate()
	global.Log.Info("初始化mysql完成")
}

// 自动迁移表结构
func autoMigrate() {
	global.Mysql.AutoMigrate(
		new(models.SysUser),
		new(models.SysRole),
		new(models.SysMenu),
		new(models.SysApi),
		new(models.SysCasbin),
		new(models.SysOperationLog),
		new(models.SysMessage),
		new(models.SysMessageLog),
		new(models.SysDict),
		new(models.SysDictData),
	)
}
