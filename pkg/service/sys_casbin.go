package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
)

func (my MysqlService) FindRoleCasbin(c models.SysRoleCasbin) []models.SysRoleCasbin {
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

func (my MysqlService) CreateRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	return global.CasbinEnforcer.AddPolicy(c.Keyword, c.Path, c.Method)
}

func (my MysqlService) BatchCreateRoleCasbin(cs []models.SysRoleCasbin) (bool, error) {
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

func (my MysqlService) DeleteRoleCasbin(c models.SysRoleCasbin) (bool, error) {
	return global.CasbinEnforcer.RemovePolicy(c.Keyword, c.Path, c.Method)
}

func (my MysqlService) BatchDeleteRoleCasbin(cs []models.SysRoleCasbin) (bool, error) {
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

func (my MysqlService) FindCasbinByRoleId(roleId uint) ([]models.SysCasbin, error) {
	list := make([][]string, 0)
	casbins := make([]models.SysCasbin, 0)
	if roleId > 0 {
		var role models.SysRole
		err := my.Q.Tx.Where("id = ?", roleId).First(&role).Error
		if err != nil {
			return casbins, err
		}
		// filter rules by keyword
		list = global.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
	} else {
		list = global.CasbinEnforcer.GetFilteredPolicy(0)
	}

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
