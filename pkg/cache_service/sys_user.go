package cache_service

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
)

// 登录校验
func (s *RedisService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.LoginCheck(user)
	}
	// 查询用户以及角色
	var u models.SysUser
	var role models.SysRole
	jsonUsers := s.GetListFromCache(nil, u.TableName())
	res1, err := JsonQueryFindOne(s.JsonQuery().FromString(jsonUsers).Where("username", "=", user.Username))
	if err != nil {
		return nil, err
	}
	utils.Struct2StructByJson(res1, &u)
	jsonRoles := s.GetListFromCache(nil, role.TableName())
	res2, err := JsonQueryFindOne(s.JsonQuery().FromString(jsonRoles).Where("id", "=", int(u.RoleId)))
	if err != nil {
		return nil, err
	}
	utils.Struct2StructByJson(res2, &role)
	u.Role = role
	// 校验密码
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	return &u, nil
}

func (s *RedisService) GetUsers(req *request.UserListRequestStruct) ([]models.SysUser, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetUsers(req)
	}
	var err error
	list := make([]models.SysUser, 0)
	// 查询接口表所有缓存
	jsonUsers := s.GetListFromCache(nil, new(models.SysUser).TableName())
	query := s.JsonQuery().FromString(jsonUsers)
	username := strings.TrimSpace(req.Username)
	if username != "" {
		query = query.Where("username", "contains", username)
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		query = query.Where("mobile", "contains", mobile)
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		query = query.Where("nickname", "contains", nickname)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	// 查询条数
	req.PageInfo.Total = uint(query.Count())
	var res interface{}
	if req.PageInfo.NoPagination {
		// 不使用分页
		res = query.Get()
	} else {
		// 获取分页参数
		limit, offset := req.GetLimit()
		res = query.Limit(int(limit)).Offset(int(offset)).Get()
	}
	// 转换为结构体
	utils.Struct2StructByJson(res, &list)
	return list, err
}

// 获取单个用户
func (s *RedisService) GetUserById(id uint) (models.SysUser, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetUserById(id)
	}
	var user models.SysUser
	var role models.SysRole
	// 获取用户
	_ = s.GetItemByIdFromCache(id, &user, user.TableName())
	// 获取用户角色
	jsonRoles := s.GetListFromCache(nil, role.TableName())
	res, err := JsonQueryFindOne(s.JsonQuery().FromString(jsonRoles).Where("id", "=", int(user.RoleId)))
	if err != nil {
		return user, err
	}
	utils.Struct2StructByJson(res, &role)
	// 将角色放到role中
	user.Role = role
	return user, nil
}
