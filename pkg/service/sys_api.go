package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

func (my MysqlService) FindApi(req *request.ApiReq) ([]models.SysApi, error) {
	var err error
	list := make([]models.SysApi, 0)
	query := my.Q.Tx.
		Model(&models.SysApi{}).
		Order("created_at DESC")
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	err = my.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}

// find all api group by api category
func (my MysqlService) FindAllApiGroupByCategoryByRoleId(currentRole models.SysRole, roleId uint) ([]response.ApiGroupByCategoryResp, []uint, error) {
	tree := make([]response.ApiGroupByCategoryResp, 0)
	accessIds := make([]uint, 0)
	allApi := make([]models.SysApi, 0)
	// find all api
	err := my.Q.Tx.Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	var currentRoleId uint
	// not super admin
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		currentRoleId = currentRole.Id
	}
	// find all casbin by current user's role id
	currentCasbins, err := my.FindCasbinByRoleId(currentRoleId)
	// find all casbin by current role id
	casbins, err := my.FindCasbinByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	newApi := make([]models.SysApi, 0)
	for _, api := range allApi {
		path := api.Path
		method := api.Method
		for _, currentCasbin := range currentCasbins {
			// have permission
			if path == currentCasbin.V1 && method == currentCasbin.V2 {
				newApi = append(newApi, api)
				break
			}
		}
	}

	// group by category
	for _, api := range newApi {
		category := api.Category
		path := api.Path
		method := api.Method
		access := false
		for _, casbin := range casbins {
			// have permission
			if path == casbin.V1 && method == casbin.V2 {
				access = true
				break
			}
		}
		// add to access ids
		if access {
			accessIds = append(accessIds, api.Id)
		}
		// generate api tree
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

func (my MysqlService) CreateApi(req *request.CreateApiReq) (err error) {
	api := new(models.SysApi)
	err = my.Q.Create(req, new(models.SysApi))
	if err != nil {
		return err
	}
	if len(req.RoleIds) > 0 {
		var roles []models.SysRole
		err = my.Q.Tx.Where("id IN (?)", req.RoleIds).Find(&roles).Error
		if err != nil {
			return
		}
		// generate casbin rules
		cs := make([]models.SysRoleCasbin, 0)
		for _, role := range roles {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: role.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		_, err = my.BatchCreateRoleCasbin(cs)
	}
	return
}

func (my MysqlService) UpdateApiById(id uint, req request.UpdateApiReq) (err error) {
	var api models.SysApi
	query := my.Q.Tx.Model(&api).Where("id = ?", id).First(&api)
	if query.Error == gorm.ErrRecordNotFound {
		return gorm.ErrRecordNotFound
	}

	m := make(map[string]interface{}, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(api, req, &m)

	oldApi := api
	err = query.Updates(m).Error

	// get diff fields
	diff := make(map[string]interface{}, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(oldApi, api, &diff)

	path, ok1 := diff["path"]
	method, ok2 := diff["method"]
	if (ok1 && path != "") || (ok2 && method != "") {
		// path/method change, the caspin rule needs to be updated
		oldCasbins := my.FindRoleCasbin(models.SysRoleCasbin{
			Path:   oldApi.Path,
			Method: oldApi.Method,
		})
		if len(oldCasbins) > 0 {
			keywords := make([]string, 0)
			for _, oldCasbin := range oldCasbins {
				keywords = append(keywords, oldCasbin.Keyword)
			}
			// delete old rules
			my.BatchDeleteRoleCasbin(oldCasbins)
			// create new rules
			newCasbins := make([]models.SysRoleCasbin, 0)
			for _, keyword := range keywords {
				newCasbins = append(newCasbins, models.SysRoleCasbin{
					Keyword: keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
			_, err = my.BatchCreateRoleCasbin(newCasbins)
		}
	}
	return
}

func (my MysqlService) DeleteApiByIds(ids []uint) (err error) {
	var list []models.SysApi
	query := my.Q.Tx.Where("id IN (?)", ids).Find(&list)
	if query.Error != nil {
		return
	}
	casbins := make([]models.SysRoleCasbin, 0)
	for _, api := range list {
		casbins = append(casbins, my.FindRoleCasbin(models.SysRoleCasbin{
			Path:   api.Path,
			Method: api.Method,
		})...)
	}
	// delete old rules
	my.BatchDeleteRoleCasbin(casbins)
	return query.Delete(&models.SysApi{}).Error
}
