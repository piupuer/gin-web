package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/utils"
	"strings"
)

// 注意golang坑, 返回值有数组时, 不要为了return方便使用命名返回值
// 如func GetApis(req *request.ApiListRequestStruct) (apis []models.SysApi, err error) {
// 这会导致apis被初始化无穷大的集合, 代码无法断点调试: collecting data..., 几秒后程序异常退出:
// error layer=rpc writing response:write tcp xxx write: broken pipe
// error layer=rpc rpc:read tcp xxx read: connection reset by peer
// exit code 0
// 我曾一度以为调试工具安装配置错误, 使用其他项目代码却能稳定调试, 最终还是定位到代码本身. 踩过的坑希望大家不要再踩
// 获取所有接口
func GetApis(req *request.ApiListRequestStruct) ([]models.SysApi, error) {
	var err error
	list := make([]models.SysApi, 0)
	db := global.Mysql
	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		// 获取分页参数
		limit, offset := req.GetLimit()
		err = db.Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, err
}

// 根据权限编号获取以api分类分组的权限接口
func GetAllApiGroupByCategoryByRoleId(roleId uint) ([]response.ApiGroupByCategoryResponseStruct, []uint, error) {
	// 接口树
	tree := make([]response.ApiGroupByCategoryResponseStruct, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allApi := make([]models.SysApi, 0)
	// 查询全部api
	err := global.Mysql.Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	// 查询当前角色拥有api访问权限的casbin规则
	casbins, err := GetCasbinListByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	// 通过分类进行分组归纳
	for _, api := range allApi {
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
		children := make([]response.ApiListResponseStruct, 0)
		for index, leaf := range tree {
			if leaf.Category == category {
				children = leaf.Children
				existIndex = index
				break
			}
		}
		// api结构转换
		var item response.ApiListResponseStruct
		utils.Struct2StructByJson(api, &item)
		item.Title = fmt.Sprintf("%s %s[%s]", item.Desc, item.Path, item.Method)
		children = append(children, item)
		if existIndex != -1 {
			// 更新元素
			tree[existIndex].Children = children
		} else {
			// 新增元素
			tree = append(tree, response.ApiGroupByCategoryResponseStruct{
				Title:    category + "分组",
				Category: category,
				Children: children,
			})
		}
	}
	return tree, accessIds, err
}

// 创建接口
func CreateApi(req *request.CreateApiRequestStruct) (err error) {
	var api models.SysApi
	utils.Struct2StructByJson(req, &api)
	// 创建数据
	err = global.Mysql.Create(&api).Error
	return
}

// 更新接口
func UpdateApiById(id uint, req gin.H) (err error) {
	var oldApi models.SysApi
	query := global.Mysql.Table(oldApi.TableName()).Where("id = ?", id).First(&oldApi)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(oldApi, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除接口
func DeleteApiByIds(ids []uint) (err error) {
	return global.Mysql.Where("id IN (?)", ids).Delete(models.SysApi{}).Error
}
