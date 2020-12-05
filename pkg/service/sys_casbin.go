package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 获取casbin策略管理器
func (s *MysqlService) Casbin() (*casbin.Enforcer, error) {
	// 初始化数据库适配器, 添加自定义表前缀, casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	// casbin默认表名casbin_rule, 为了与项目统一改写一下规则
	// 注意: gormadapter.CasbinTableName内部添加了下划线, 这里不再多此一举
	a, err := gormadapter.NewAdapterByDBUseTableName(s.db, global.Conf.Mysql.TablePrefix, "sys_casbin")
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
func (s *MysqlService) GetRoleCasbins(c models.SysRoleCasbin) []models.SysRoleCasbin {
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
func (s *MysqlService) CreateRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
	return e.AddPolicy(c.Keyword, c.Path, c.Method)
}

// 创建多条casbin规则, 按角色
func (s *MysqlService) CreateRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
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

// 批量创建多条casbin规则, 按角色
func (s *MysqlService) BatchCreateRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
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
func (s *MysqlService) DeleteRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	e, _ := s.Casbin()
	return e.RemovePolicy(c.Keyword, c.Path, c.Method)
}

// 批量删除多条casbin规则, 按角色
func (s *MysqlService) BatchDeleteRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
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
func (s *MysqlService) GetCasbinListByRoleId(roleId uint) ([]models.SysCasbin, error) {
	list := make([][]string, 0)
	casbins := make([]models.SysCasbin, 0)
	e, _ := s.Casbin()
	if roleId > 0 {
		// 读取角色缓存
		var role models.SysRole
		err := s.tx.Where("id = ?", roleId).First(&role).Error
		if err != nil {
			return casbins, err
		}
		// 查询符合字段v0=role.Keyword所有casbin规则
		list = e.GetFilteredPolicy(0, role.Keyword)
	} else {
		list = e.GetFilteredPolicy(0)
	}

	// 避免重复, 记录添加历史
	var added []string
	for _, v := range list {
		if !utils.Contains(added, v[1]+v[2]) {
			casbins = append(casbins, models.SysCasbin{
				PType: "p",
				V1:    v[1],
				V2:    v[2],
			})
			added = append(added, v[1]+v[2])
		}
	}
	return casbins, nil
}
