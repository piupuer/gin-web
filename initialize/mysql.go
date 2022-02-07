package initialize

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	m "github.com/go-sql-driver/mysql"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/binlog"
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/query"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func Mysql() {
	cfg, err := m.ParseDSN(global.Conf.Mysql.Uri)
	if err != nil {
		panic(errors.Wrap(err, "initialize mysql failed"))
	}
	// create database if not exists
	query.MigrateDatabase(*cfg)
	global.Conf.Mysql.DSN = *cfg

	log.WithRequestId(ctx).Info("mysql dsn: %s", cfg.FormatDSN())
	init := false
	ctx, cancel := context.WithTimeout(ctx, time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				if !init {
					panic(fmt.Sprintf("initialize mysql failed: connect timeout(%ds)", global.Conf.System.ConnectTimeout))
				}
				// avoid goroutine dead lock
				return
			}
		}
	}()
	l := log.NewDefaultGormLogger()
	if global.Conf.Mysql.NoSql {
		// not show sql log
		l = l.LogMode(glogger.Silent)
	} else {
		l = l.LogMode(glogger.Info)
	}
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   global.Conf.Mysql.TablePrefix + "_",
			SingularTable: true,
		},
		// select * from xxx => select a,b,c from xxx
		QueryFields: true,
		Logger:      l,
	})
	if err != nil {
		panic(errors.Wrap(err, "initialize mysql failed"))
	}
	init = true
	global.Mysql = db
	autoMigrate()
	binlogListen()
	log.WithRequestId(ctx).Info("initialize mysql success")
}

func autoMigrate() {
	// migrate tables
	global.Mysql.WithContext(ctx).AutoMigrate(
		new(ms.SysMenu),
		new(ms.SysMenuRoleRelation),
		new(ms.SysApi),
		new(ms.SysCasbin),
		new(ms.SysOperationLog),
		new(ms.SysMessage),
		new(ms.SysMessageLog),
		new(ms.SysMachine),
		new(ms.SysDict),
		new(ms.SysDictData),
		new(models.SysUser),
		new(models.SysRole),
		new(models.Leave),
	)
	// auto migrate fsm
	fsm.Migrate(global.Mysql, fsm.WithCtx(ctx))
}

func binlogListen() {
	if !global.Conf.Redis.EnableBinlog {
		log.WithRequestId(ctx).Info("if redis is not used or binlog is not enabled, there is no need to initialize the MySQL binlog listener")
		return
	}
	err := binlog.NewMysqlBinlog(
		binlog.WithCtx(ctx),
		binlog.WithRedis(global.Redis),
		binlog.WithDb(global.Mysql),
		binlog.WithDsn(&global.Conf.Mysql.DSN),
		binlog.WithBinlogPos(global.Conf.Redis.BinlogPos),
		binlog.WithIgnore(
			// The following tables will have more and more data over time
			// It is not suitable to store the entire table JSON in redis
			"sys_operation_log",
		),
		binlog.WithModels(
			// The following tables will be sync to redis
			new(ms.SysMenu),
			new(ms.SysMenuRoleRelation),
			new(ms.SysApi),
			new(ms.SysCasbin),
			new(ms.SysMessage),
			new(ms.SysMessageLog),
			new(ms.SysMachine),
			new(ms.SysDict),
			new(ms.SysDictData),
			new(models.SysUser),
			new(models.SysRole),
			new(models.Leave),
		),
	)
	if err != nil {
		panic(err)
	}
}
