package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

// 注意golang坑, 返回值有数组时, 不要为了return方便使用命名返回值
// 如func GetApis(req *request.ApiReq) (apis []models.SysApi, err error) {
// 这会导致apis被初始化无穷大的集合, 代码无法断点调试: collecting data..., 几秒后程序异常退出:
// error layer=rpc writing response:write tcp xxx write: broken pipe
// error layer=rpc rpc:read tcp xxx read: connection reset by peer
// exit code 0
// 我曾一度以为调试工具安装配置错误, 使用其他项目代码却能稳定调试, 最终还是定位到代码本身. 踩过的坑希望大家不要再踩
// 获取所有接口
func (my MysqlService) GetApis(req *request.ApiReq) ([]models.SysApi, error) {
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
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	// 查询列表
	err = my.Q.Find(query, &req.Page, &list)
	return list, err
}

// 根据权限编号获取以api分类分组的权限接口
func (my MysqlService) GetAllApiGroupByCategoryByRoleId(currentRole models.SysRole, roleId uint) ([]response.ApiGroupByCategoryResp, []uint, error) {
	// 接口树
	tree := make([]response.ApiGroupByCategoryResp, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allApi := make([]models.SysApi, 0)
	// 查询全部api
	err := my.Q.Tx.Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	var currentRoleId uint
	// 非超级管理员
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		currentRoleId = currentRole.Id
	}
	// 查询当前角色拥有api访问权限的casbin规则
	currentCasbins, err := my.GetCasbinListByRoleId(currentRoleId)
	// 查询指定角色拥有api访问权限的casbin规则(当前角色只能在自己权限范围内操作, 不得越权)
	casbins, err := my.GetCasbinListByRoleId(roleId)
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

// 创建接口
func (my MysqlService) CreateApi(req *request.CreateApiReq) (err error) {
	api := new(models.SysApi)
	err = my.Q.Create(req, new(models.SysApi))
	if err != nil {
		return err
	}
	// 添加了角色
	if len(req.RoleIds) > 0 {
		// 查询角色关键字
		var roles []models.SysRole
		err = my.Q.Tx.Where("id IN (?)", req.RoleIds).Find(&roles).Error
		if err != nil {
			return
		}
		// 构建casbin规则
		cs := make([]models.SysRoleCasbin, 0)
		for _, role := range roles {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: role.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		// 批量创建
		_, err = my.BatchCreateRoleCasbins(cs)
	}
	return
}

// 更新接口
func (my MysqlService) UpdateApiById(id uint, req request.UpdateApiReq) (err error) {
	var api models.SysApi
	query := my.Q.Tx.Model(&api).Where("id = ?", id).First(&api)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(map[string]interface{}, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(api, req, &m)

	// 记录update前的旧数据, 执行Updates后api会变成新数据
	oldApi := api
	// 更新指定列
	err = query.Updates(m).Error

	// 对比api发生了哪些变化
	diff := make(map[string]interface{}, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(oldApi, api, &diff)

	path, ok1 := diff["path"]
	method, ok2 := diff["method"]
	if (ok1 && path != "") || (ok2 && method != "") {
		// path或method变化, 需要更新casbin规则
		// 查找当前接口都有哪些角色在使用
		oldCasbins := my.GetRoleCasbins(models.SysRoleCasbin{
			Path:   oldApi.Path,
			Method: oldApi.Method,
		})
		if len(oldCasbins) > 0 {
			keywords := make([]string, 0)
			for _, oldCasbin := range oldCasbins {
				keywords = append(keywords, oldCasbin.Keyword)
			}
			// 删除旧规则, 添加新规则
			my.BatchDeleteRoleCasbins(oldCasbins)
			// 构建新casbin规则
			newCasbins := make([]models.SysRoleCasbin, 0)
			for _, keyword := range keywords {
				newCasbins = append(newCasbins, models.SysRoleCasbin{
					Keyword: keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
			// 批量创建
			_, err = my.BatchCreateRoleCasbins(newCasbins)
		}
	}
	return
}

// 批量删除接口
func (my MysqlService) DeleteApiByIds(ids []uint) (err error) {
	var list []models.SysApi
	query := my.Q.Tx.Where("id IN (?)", ids).Find(&list)
	if query.Error != nil {
		return
	}
	// 查找当前接口都有哪些角色在使用
	casbins := make([]models.SysRoleCasbin, 0)
	for _, api := range list {
		casbins = append(casbins, my.GetRoleCasbins(models.SysRoleCasbin{
			Path:   api.Path,
			Method: api.Method,
		})...)
	}
	// 删除所有规则
	my.BatchDeleteRoleCasbins(casbins)
	return query.Delete(&models.SysApi{}).Error
}
