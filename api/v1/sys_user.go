package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

var (
	// 定期缓存, 避免每次频繁查询数据库
	userInfoCache = cache.New(24*time.Hour, 48*time.Hour)
	userByIdCache = cache.New(24*time.Hour, 48*time.Hour)
)

// 获取当前用户信息
func GetUserInfo(c *gin.Context) {
	user := GetCurrentUser(c)
	oldCache, ok := userInfoCache.Get(fmt.Sprintf("%d", user.Id))
	if ok {
		resp, _ := oldCache.(response.UserInfoResponseStruct)
		response.SuccessWithData(resp)
		return
	}

	// 转为UserInfoResponseStruct, 隐藏部分字段
	var resp response.UserInfoResponseStruct
	utils.Struct2StructByJson(user, &resp)
	resp.Roles = []string{
		"admin",
	}
	resp.Keyword = user.Role.Keyword
	resp.RoleSort = *user.Role.Sort
	// 写入缓存
	userInfoCache.Set(fmt.Sprintf("%d", user.Id), resp, cache.DefaultExpiration)
	response.SuccessWithData(resp)
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	var req request.UserRequestStruct
	request.ShouldBind(c, &req)
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	req.CurrentRole = user.Role
	s := cache_service.New(c)
	users, err := s.GetUsers(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.UserListResponseStruct
	utils.Struct2StructByJson(users, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 修改密码
func ChangePwd(c *gin.Context) {
	// 请求json绑定
	var req request.ChangePwdRequestStruct
	request.ShouldBind(c, &req)
	// 获取当前用户
	user := GetCurrentUser(c)
	query := global.Mysql.Where("username = ?", user.Username).First(&user)
	// 查询用户
	err := query.Error
	response.CheckErr(err)
	// 校验密码
	if ok := utils.ComparePwd(req.OldPassword, user.Password); !ok {
		response.CheckErr("原密码错误")
	}
	// 更新密码
	err = query.Update("password", utils.GenPwd(req.NewPassword)).Error
	response.CheckErr(err)
	response.Success()
}

// 获取当前请求用户信息
func GetCurrentUser(c *gin.Context) models.SysUser {
	userId, exists := c.Get("user")
	var newUser models.SysUser
	if !exists {
		return newUser
	}
	uid, _ := userId.(uint)
	oldCache, ok := userByIdCache.Get(fmt.Sprintf("%d", uid))
	if ok {
		u, _ := oldCache.(models.SysUser)
		return u
	}
	s := service.New(c)
	newUser, _ = s.GetUserById(uid)
	// 写入缓存
	userByIdCache.Set(fmt.Sprintf("%d", uid), newUser, cache.DefaultExpiration)
	return newUser
}

// 创建用户
func CreateUser(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateUserRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	s := service.New(c)
	// 将初始密码转为密文
	req.Password = utils.GenPwd(req.InitPassword)
	err := s.Create(req, new(models.SysUser))
	response.CheckErr(err)
	response.Success()
}

// 更新用户
func UpdateUserById(c *gin.Context) {
	var req request.UpdateUserRequestStruct
	request.ShouldBind(c, &req)

	// 获取path中的userId
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		response.CheckErr("用户编号不正确")
	}

	// 填写了新密码
	if req.NewPassword != nil && strings.TrimSpace(*req.NewPassword) != "" {
		password := utils.GenPwd(*req.NewPassword)
		req.Password = &password
	}

	user := GetCurrentUser(c)
	if userId == user.Id {
		if req.Status != nil && uint(*req.Status) == models.SysUserStatusDisabled {
			response.CheckErr("不能禁用自己")
		}
		if req.RoleId != nil && user.RoleId != *req.RoleId {
			if *user.Role.Sort != models.SysRoleSuperAdminSort {
				response.CheckErr("无法更改自己的角色, 如需更改请联系上级领导")
			} else {
				response.CheckErr("无法更改超级管理员的角色")
			}
		}
	}

	s := service.New(c)
	err := s.UpdateById(userId, req, new(models.SysUser))
	response.CheckErr(err)
	userInfoCache.Delete(fmt.Sprintf("%d", user.Id))
	userByIdCache.Delete(fmt.Sprintf("%d", user.Id))
	response.Success()
}

// 批量删除用户
func BatchDeleteUserByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	user := GetCurrentUser(c)
	if utils.ContainsUint(req.GetUintIds(), user.Id) {
		response.CheckErr("不能删除自己")
	}

	s := service.New(c)
	err := s.DeleteByIds(req.GetUintIds(), new(models.SysUser))
	response.CheckErr(err)
	userInfoCache.Delete(fmt.Sprintf("%d", user.Id))
	userByIdCache.Delete(fmt.Sprintf("%d", user.Id))
	response.Success()
}
