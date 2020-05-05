package service

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
)

// 获取casbin策略管理器
func Casbin() (*casbin.Enforcer, error) {
	// 初始化数据库适配器, 添加自定义表前缀
	a, err := gormadapter.NewAdapterByDB(global.Mysql)
	if err != nil {
		return nil, err
	}
	// 读取配置文件
	e, err := casbin.NewEnforcer("conf/rbac_model.conf", a)
	if err != nil {
		return nil, err
	}
	// 加载策略
	err = e.LoadPolicy()
	return e, err
}

// 创建一条casbin规则
func CreateCasbin(c models.SysCasbin) (bool, error) {
	e, _ := Casbin()
	return e.AddPolicy(c.V0, c.V1, c.V2)
}
