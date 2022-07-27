package initialize

import (
	"gin-web/pkg/global"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/pkg/errors"
)

func CasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		panic(errors.Wrap(err, "initialize casbin enforcer failed"))
	}
	global.CasbinEnforcer = &e
	log.WithContext(ctx).Info("initialize casbin enforcer success")
}

func mysqlCasbin() (en casbin.Enforcer, err error) {
	a, err := gormadapter.NewAdapterByDBUseTableName(
		global.Mysql.WithContext(ctx),
		// add mysql table prefix config
		global.Conf.Mysql.TablePrefix,
		"sys_casbin",
	)
	if err != nil {
		return
	}
	// read model path
	config := global.ConfBox.Get(global.Conf.System.CasbinModelPath)
	cabinModel := model.NewModel()
	err = cabinModel.LoadModelFromText(string(config))
	if err != nil {
		return
	}
	var e *casbin.Enforcer
	e, err = casbin.NewEnforcer(cabinModel, a)
	if err != nil {
		return
	}
	err = e.LoadPolicy()
	if err != nil {
		return
	}
	en = *e
	return
}
