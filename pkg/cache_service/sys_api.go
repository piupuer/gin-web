package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
)

func (s RedisService) GetApis(req *request.ApiReq) ([]models.SysApi, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return s.mysql.GetApis(req)
	}
	var err error
	list := make([]models.SysApi, 0)
	query := s.Q.
		Table("sys_api").
		Order("created_at DESC")
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method", "contains", method)
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path", "contains", path)
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category", "contains", category)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	err = s.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}

// 根据权限编号获取以api分类分组的权限接口
func (s RedisService) GetAllApiGroupByCategoryByRoleId(currentRole models.SysRole, roleId uint) ([]response.ApiGroupByCategoryResp, []uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return s.mysql.GetAllApiGroupByCategoryByRoleId(currentRole, roleId)
	}
	// 接口树
	tree := make([]response.ApiGroupByCategoryResp, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allApi := make([]models.SysApi, 0)
	// 查询全部api
	err := s.Q.Table("sys_api").Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	var currentRoleId uint
	// 非超级管理员
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		currentRoleId = currentRole.Id
	}
	// 查询当前角色拥有api访问权限的casbin规则
	currentCasbins, err := s.mysql.GetCasbinListByRoleId(currentRoleId)
	// 查询指定角色拥有api访问权限的casbin规则(当前角色只能在自己权限范围内操作, 不得越权)
	casbins, err := s.mysql.GetCasbinListByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	// 找到当前角色的全部api
	newApi := make([]models.SysApi, 0)
	for _, api := range allApi {
		path := api.Path
		method := api.Method
		for _, currentCasbin := range currentCasbins {
			// 该api有权限
			if path == currentCasbin.V1 && method == currentCasbin.V2 {
				newApi = append(newApi, api)
				break
			}
		}
	}

	// 通过分类进行分组归纳
	for _, api := range newApi {
		category := api.Category
		path := api.Path
		method := api.Method
		access := false
		for _, casbin := range casbins {
			// 该api有权限
			if path == casbin.V1 && method == casbin.V2 {
				access = true
				break
			}
		}
		// 加入权限集合
		if access {
			accessIds = append(accessIds, api.Id)
		}
		// 生成接口树
		existIndex := -1
		children := make([]response.ApiResp, 0)
		for index, leaf := range tree {
			if leaf.Category == category {
				children = leaf.Children
				existIndex = index
				break
			}
		}
		// api结构转换
		var item response.ApiResp
		utils.Struct2StructByJson(api, &item)
		item.Title = fmt.Sprintf("%s %s[%s]", item.Desc, item.Path, item.Method)
		children = append(children, item)
		if existIndex != -1 {
			// 更新元素
			tree[existIndex].Children = children
		} else {
			// 新增元素
			tree = append(tree, response.ApiGroupByCategoryResp{
				Title:    category + "分组",
				Category: category,
				Children: children,
			})
		}
	}
	return tree, accessIds, err
}
