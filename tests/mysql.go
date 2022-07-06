package tests

import (
	"fmt"
	"gin-web/pkg/global"
	m "github.com/go-sql-driver/mysql"
	"github.com/piupuer/go-helper/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Mysql() {
	cfg, err := m.ParseDSN(global.Conf.Mysql.Uri)
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize mysql failed: %v", err))
	}
	global.Conf.Mysql.DSN = *cfg

	log.WithContext(ctx).Info("[unit test]mysql dsn: %s", cfg.FormatDSN())
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
		panic(fmt.Sprintf("[unit test]initialize mysql failed: %v", err))
	}
	global.Mysql = db
	autoMigrate()
	log.WithContext(ctx).Debug("[unit test]initialize mysql success")
}

func autoMigrate() {

}
