package initialize

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func CasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		panic(fmt.Sprintf("initialize casbin enforcer failed: %v", err))
	}
	global.CasbinEnforcer = e
	global.Log.Info(ctx, "initialize casbin enforcer success")
}

func mysqlCasbin() (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDBUseTableName(
		global.Mysql.WithContext(ctx),
		// add mysql table prefix config
		global.Conf.Mysql.TablePrefix,
		"sys_casbin",
	)
	if err != nil {
		return nil, err
	}
	// read model path
	config := global.ConfBox.Find(global.Conf.System.CasbinModelPath)
	cabinModel := model.NewModel()
	err = cabinModel.LoadModelFromText(string(config))
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(cabinModel, a)
	if err != nil {
		return nil, err
	}
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
