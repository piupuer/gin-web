package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	// 定期缓存, 避免每次频繁查询数据库
	menuTreeCache = cache.New(24*time.Hour, 48*time.Hour)
)

// 查询当前用户菜单树
func GetMenuTree(c *gin.Context) {
	user := GetCurrentUser(c)
	oldCache, ok := menuTreeCache.Get(fmt.Sprintf("%d", user.Id))
	if ok {
		resp, _ := oldCache.([]response.MenuTreeResp)
		response.SuccessWithData(resp)
		return
	}

	s := service.New(c)
	menus, err := s.GetMenuTree(user.RoleId)
	response.CheckErr(err)
	var resp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &resp)
	// 写入缓存
	menuTreeCache.Set(fmt.Sprintf("%d", user.Id), resp, cache.DefaultExpiration)
	response.SuccessWithData(resp)
}

// 查询指定角色的菜单树
func GetAllMenuByRoleId(c *gin.Context) {
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	s := service.New(c)
	menus, ids, err := s.GetAllMenuByRoleId(user.Role, utils.Str2Uint(c.Param("roleId")))
	response.CheckErr(err)
	var resp response.MenuTreeWithAccessResp
	resp.AccessIds = ids
	utils.Struct2StructByJson(menus, &resp.List)
	response.SuccessWithData(resp)
}

// 查询所有菜单
func GetMenus(c *gin.Context) {
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	s := service.New(c)
	menus := s.GetMenus(user.Role)
	var resp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &resp)
	response.SuccessWithData(resp)
}

// 创建菜单
func CreateMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateMenuReq
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.CreateMenu(user.Role, &req)
	response.CheckErr(err)
	response.Success()
}

// 更新菜单
func UpdateMenuById(c *gin.Context) {
	var req request.UpdateMenuReq
	request.ShouldBind(c, &req)
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		response.CheckErr("菜单编号不正确")
	}
	s := service.New(c)
	err := s.UpdateById(menuId, req, new(models.SysMenu))
	response.CheckErr(err)
	response.Success()
}

// 批量删除菜单
func BatchDeleteMenuByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.DeleteByIds(req.GetUintIds(), new(models.SysMenu))
	response.CheckErr(err)
	response.Success()
}
