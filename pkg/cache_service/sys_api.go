package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/piupuer/go-helper/ms"
	"strings"
)

func (rd RedisService) FindApi(req *request.ApiReq) []ms.SysApi {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindApi(req)
	}
	list := make([]ms.SysApi, 0)
	query := rd.Q.
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
	rd.Q.FindWithPage(query, &req.Page, &list)
	return list
}

// find all api group by api category
func (rd RedisService) FindAllApiGroupByCategoryByRoleId(currentRole models.SysRole, roleId uint) ([]response.ApiGroupByCategoryResp, []uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindAllApiGroupByCategoryByRoleId(currentRole, roleId)
	}
	tree := make([]response.ApiGroupByCategoryResp, 0)
	accessIds := make([]uint, 0)
	allApi := make([]ms.SysApi, 0)
	rd.Q.
		Table("sys_api").
		Find(&allApi)
	var currentRoleId uint
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		currentRoleId = currentRole.Id
	}
	currentCasbins, err := rd.mysql.FindCasbinByRoleId(currentRoleId)
	casbins, err := rd.mysql.FindCasbinByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	newApi := make([]ms.SysApi, 0)
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
