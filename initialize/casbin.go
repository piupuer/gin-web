package initialize

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 初始化casbin
func CasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		panic(fmt.Sprintf("初始化casbin策略管理器: %v", err))
	}
	global.CasbinEnforcer = e
	global.Log.Info(ctx, "初始化casbin策略管理器完成")
}

func mysqlCasbin() (*casbin.Enforcer, error) {
	// 初始化数据库适配器, 添加自定义表前缀, casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	// casbin默认表名casbin_rule, 为了与项目统一改写一下规则
	// 注意: gormadapter.CasbinTableName内部添加了下划线, 这里不再多此一举
	a, err := gormadapter.NewAdapterByDBUseTableName(global.Mysql, global.Conf.Mysql.TablePrefix, "sys_casbin")
	if err != nil {
		return nil, err
	}
	// 加锁避免并发多次初始化cabinModel
	// 读取配置文件
	config, err := global.ConfBox.Find(global.Conf.Casbin.ModelPath)
	cabinModel := model.NewModel()
	// 从字符串中加载casbin配置
	err = cabinModel.LoadModelFromText(string(config))
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(cabinModel, a)
	if err != nil {
		return nil, err
	}
	// 加载策略
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
