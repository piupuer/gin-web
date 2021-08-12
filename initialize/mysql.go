package initialize

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	colorful := false
	if global.Conf.Logs.Level <= zapcore.DebugLevel {
		colorful = true
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		// 指定表前缀
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: global.Conf.Mysql.TablePrefix + "_",
		},
		// 查询全部字段, 某些情况下*不走索引
		QueryFields: true,
		Logger:      global.NewGormZapLogger(
			global.Log.Desugar(),
			logger.Config{
				Colorful: colorful,
			},
		),
	})
	if err != nil {
		panic(fmt.Sprintf("初始化mysql异常: %v", err))
	}
	init = true
	global.Mysql = db
	// 表结构
	autoMigrate()
	global.Log.Info("初始化mysql完成")
	// 初始化数据库日志监听器
	binlog()
}

// 自动迁移表结构
func autoMigrate() {
	s := service.New(nil)
	global.Mysql.WithContext(s.RequestIdContext(requestId)).AutoMigrate(
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
		new(models.SysMachine),
		new(models.SysDict),
		new(models.SysDictData),
	)
}

func binlog() {
	MysqlBinlog(
		[]string{
			// 下列表会随着使用时间数据量越来越大, 不适合将整个表json存入redis
			new(models.SysOperationLog).TableName(),
		},
		new(models.SysUser),
		new(models.SysRole),
		new(models.SysMenu),
		new(models.RelationMenuRole),
		new(models.SysApi),
		new(models.SysCasbin),
		new(models.SysWorkflow),
		new(models.SysWorkflowLine),
		new(models.SysWorkflowLog),
		new(models.RelationUserWorkflowLine),
		new(models.SysLeave),
		new(models.SysMessage),
		new(models.SysMessageLog),
		new(models.SysMachine),
		new(models.SysDict),
		new(models.SysDictData),
	)
}
