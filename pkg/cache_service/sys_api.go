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

func (s RedisService) FindApi(req *request.ApiReq) ([]models.SysApi, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return s.mysql.FindApi(req)
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
	err = s.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}

// find all api group by api category
func (s RedisService) FindAllApiGroupByCategoryByRoleId(currentRole models.SysRole, roleId uint) ([]response.ApiGroupByCategoryResp, []uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return s.mysql.FindAllApiGroupByCategoryByRoleId(currentRole, roleId)
	}
	tree := make([]response.ApiGroupByCategoryResp, 0)
	accessIds := make([]uint, 0)
	allApi := make([]models.SysApi, 0)
	err := s.Q.Table("sys_api").Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	var currentRoleId uint
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		currentRoleId = currentRole.Id
	}
	currentCasbins, err := s.mysql.FindCasbinByRoleId(currentRoleId)
	casbins, err := s.mysql.FindCasbinByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	newApi := make([]models.SysApi, 0)
	for _, api := range allApi {
		path := api.Path
		method := api.Method
		for _, currentCasbin := range currentCasbins {
			if path == currentCasbin.V1 && method == currentCasbin.V2 {
				newApi = append(newApi, api)
				break
			}
		}
	}

	for _, api := range newApi {
		category := api.Category
		path := api.Path
		method := api.Method
		access := false
		for _, casbin := range casbins {
			if path == casbin.V1 && method == casbin.V2 {
				access = true
				break
			}
		}
		if access {
			accessIds = append(accessIds, api.Id)
		}
		existIndex := -1
		children := make([]response.ApiResp, 0)
		for index, leaf := range tree {
			if leaf.Category == category {
				children = leaf.Children
				existIndex = index
				break
			}
		}
		var item response.ApiResp
		utils.Struct2StructByJson(api, &item)
		item.Title = fmt.Sprintf("%s %s[%s]", item.Desc, item.Path, item.Method)
		children = append(children, item)
		if existIndex != -1 {
			tree[existIndex].Children = children
		} else {
			tree = append(tree, response.ApiGroupByCategoryResp{
				Title:    category + " group",
				Category: category,
				Children: children,
			})
		}
	}
	return tree, accessIds, err
}
