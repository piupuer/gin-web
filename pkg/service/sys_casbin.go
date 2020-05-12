package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
)

// 获取casbin策略管理器
func (s *CommonService) Casbin() (*casbin.Enforcer, error) {
	// 初始化数据库适配器, 添加自定义表前缀, casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	a, err := gormadapter.NewAdapterByDB(s.db)
	if err != nil {
		return nil, err
	}
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
	return e, err
}

// 获取符合条件的casbin规则, 按角色
func (s *CommonService) GetRoleCasbins(c models.SysRoleCasbin) []models.SysRoleCasbin {
	e, _ := s.Casbin()
	policies := e.GetFilteredPolicy(0, c.Keyword, c.Path, c.Method)
	cs := make([]models.SysRoleCasbin, 0)
	for _, policy := range policies {
		cs = append(cs, models.SysRoleCasbin{
			Keyword: policy[0],
			Path:    policy[1],
			Method:  policy[2],
		})
	}
	return cs
}

// 创建一条casbin规则, 按角色
func (s *CommonService) CreateRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
	return e.AddPolicy(c.Keyword, c.Path, c.Method)
}

// 批量创建多条casbin规则, 按角色
func (s *CommonService) BatchCreateRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
	// 按角色构建
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return e.AddPolicies(rules)
}

// 删除一条casbin规则, 按角色
func (s *CommonService) DeleteRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
	return e.RemovePolicy(c.Keyword, c.Path, c.Method)
}

// 批量删除多条casbin规则, 按角色
func (s *CommonService) BatchDeleteRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
	// 按角色构建
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return e.RemovePolicies(rules)
}

// 根据权限编号读取casbin规则
func (s *CommonService) GetCasbinListByRoleId(roleId uint) ([]models.SysCasbin, error) {
	casbins := make([]models.SysCasbin, 0)
	var role models.SysRole
	err := s.tx.Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return casbins, err
	}
	e, _ := s.Casbin()
	// 查询符合字段v0=role.Keyword所有casbin规则
	list := e.GetFilteredPolicy(0, role.Keyword)
	for _, v := range list {
		casbins = append(casbins, models.SysCasbin{
			PType: "p",
			V0:    v[0],
			V1:    v[1],
			V2:    v[2],
		})
	}
	return casbins, nil
}
