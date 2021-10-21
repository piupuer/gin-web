package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
)

// 获取符合条件的casbin规则, 按角色
func (my MysqlService) GetRoleCasbins(c models.SysRoleCasbin) []models.SysRoleCasbin {
	policies := global.CasbinEnforcer.GetFilteredPolicy(0, c.Keyword, c.Path, c.Method)
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
func (my MysqlService) CreateRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	return global.CasbinEnforcer.AddPolicy(c.Keyword, c.Path, c.Method)
}

// 创建多条casbin规则, 按角色
func (my MysqlService) CreateRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return global.CasbinEnforcer.AddPolicies(rules)
}

// 批量创建多条casbin规则, 按角色
func (my MysqlService) BatchCreateRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
	// 按角色构建
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return global.CasbinEnforcer.AddPolicies(rules)
}

// 删除一条casbin规则, 按角色
func (my MysqlService) DeleteRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	return global.CasbinEnforcer.RemovePolicy(c.Keyword, c.Path, c.Method)
}

// 批量删除多条casbin规则, 按角色
func (my MysqlService) BatchDeleteRoleCasbins(cs []models.SysRoleCasbin) (bool, error) {
	// 按角色构建
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return global.CasbinEnforcer.RemovePolicies(rules)
}

// 根据权限编号读取casbin规则
func (my MysqlService) FindCasbinByRoleId(roleId uint) ([]models.SysCasbin, error) {
	list := make([][]string, 0)
	casbins := make([]models.SysCasbin, 0)
	if roleId > 0 {
		// 读取角色缓存
		var role models.SysRole
		err := my.Q.Tx.Where("id = ?", roleId).First(&role).Error
		if err != nil {
			return casbins, err
		}
		// 查询符合字段v0=role.Keyword所有casbin规则
		list = global.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
	} else {
		list = global.CasbinEnforcer.GetFilteredPolicy(0)
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
