package initialize

import (
	"context"
	"embed"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	m "github.com/go-sql-driver/mysql"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/binlog"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/migrate"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

//go:embed db/*.sql
var sqlFs embed.FS

func Mysql() {
	cfg, err := m.ParseDSN(global.Conf.Mysql.Uri)
	if err != nil {
		panic(errors.Wrap(err, "initialize mysql failed"))
	}
	global.Conf.Mysql.DSN = *cfg
	uri := global.Conf.Mysql.Uri
	err = migrate.Do(
		migrate.WithCtx(ctx),
		migrate.WithUri(uri),
		migrate.WithFs(sqlFs),
		migrate.WithFsRoot("db"),
		migrate.WithBefore(beforeMigrate),
	)
	if err != nil {
		panic(errors.Wrap(err, "initialize mysql failed"))
	}
	err = binlogListen()
	if err != nil {
		panic(errors.Wrap(err, "initialize mysql binlog failed"))
	}

	log.WithContext(ctx).Info("initialize mysql success")
}

func beforeMigrate(ctx context.Context) (err error) {
	init := false
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
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
	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(global.Conf.Mysql.DSN.FormatDSN()), &gorm.Config{
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
		return
	}
	init = true
	global.Mysql = db
	autoMigrate()
	return
}

func autoMigrate() {
	// migrate tables change to sql-migrate: initialize/db/***.sql

	// auto migrate fsm
	fsm.Migrate(
		fsm.WithDb(global.Mysql),
		fsm.WithCtx(ctx),
		fsm.WithPrefix(constant.FsmPrefix),
	)
}

func binlogListen() (err error) {
	if !global.Conf.Redis.EnableBinlog {
		log.WithContext(ctx).Info("if redis is not used or binlog is not enabled, there is no need to initialize the MySQL binlog listener")
		return
	}
	err = binlog.NewMysqlBinlog(
		binlog.WithCtx(ctx),
		binlog.WithRedis(global.Redis),
		binlog.WithDb(global.Mysql),
		binlog.WithDsn(&global.Conf.Mysql.DSN),
		binlog.WithBinlogPos(global.Conf.Redis.BinlogPos),
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
	return
}
