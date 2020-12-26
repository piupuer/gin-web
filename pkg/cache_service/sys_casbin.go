package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	redisadapter "github.com/casbin/redis-adapter/v2"
)

var cabinAdapter *redisadapter.Adapter

// 获取casbin策略管理器
func (s *RedisService) Casbin() (*casbin.Enforcer, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.Casbin()
	}
	if cabinAdapter == nil {
		// 这里使用redis适配器
		cabinAdapter = redisadapter.NewAdapterWithKey(
			// 使用tcp连接redis
			"tcp",
			// 主机地址+端口
			fmt.Sprintf("%s:%d", global.Conf.Redis.Host, global.Conf.Redis.Port),
			// 缓存key由数据库名+表名组成, 见redis.RowChange方法cacheKey
			fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, new(models.SysCasbin).TableName()),
		)
	}
	// 读取配置文件
	config, err := global.ConfBox.Find(global.Conf.Casbin.ModelPath)
	cabinModel := model.NewModel()
	// 从字符串中加载casbin配置
	err = cabinModel.LoadModelFromText(string(config))
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(cabinModel, cabinAdapter)
	if err != nil {
		return nil, err
	}
	// 加载策略
	err = e.LoadPolicy()
	return e, err
}

// 获取符合条件的casbin规则, 按角色
func (s *RedisService) GetRoleCasbins(c models.SysRoleCasbin) []models.SysRoleCasbin {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetRoleCasbins(c)
	}
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

// 根据权限编号读取casbin规则(如果roleId为0表示读取全部)
func (s *RedisService) GetCasbinListByRoleId(roleId uint) ([]models.SysCasbin, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetCasbinListByRoleId(roleId)
	}
	list := make([][]string, 0)
	casbins := make([]models.SysCasbin, 0)
	e, _ := s.Casbin()
	if roleId > 0 {
		// 读取角色缓存
		var role models.SysRole
		err := s.redis.Table(new(models.SysRole).TableName()).Where("id", "=", roleId).First(&role).Error
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
