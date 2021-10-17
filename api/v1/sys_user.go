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
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
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
		rp, _ := oldCache.(response.UserInfoResp)
		resp.SuccessWithData(rp)
		return
	}

	// 转为UserInfoResponseStruct, 隐藏部分字段
	var rp response.UserInfoResp
	utils.Struct2StructByJson(user, &rp)
	rp.Roles = []string{
		"admin",
	}
	rp.Keyword = user.Role.Keyword
	rp.RoleSort = *user.Role.Sort
	// 写入缓存
	userInfoCache.Set(fmt.Sprintf("%d", user.Id), rp, cache.DefaultExpiration)
	resp.SuccessWithData(rp)
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	var r request.UserReq
	req.ShouldBind(c, &r)
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	r.CurrentRole = user.Role
	s := cache_service.New(c)
	list, err := s.GetUsers(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.UserResp{}, r.Page)
}

// 修改密码
func ChangePwd(c *gin.Context) {
	// 请求json绑定
	var r request.ChangePwdReq
	req.ShouldBind(c, &r)
	// 获取当前用户
	user := GetCurrentUser(c)
	query := global.Mysql.Where("username = ?", user.Username).First(&user)
	// 查询用户
	err := query.Error
	resp.CheckErr(err)
	// 校验密码
	if ok := utils.ComparePwd(r.OldPassword, user.Password); !ok {
		resp.CheckErr("原密码错误")
	}
	// 更新密码
	err = query.Update("password", utils.GenPwd(r.NewPassword)).Error
	resp.CheckErr(err)
	resp.Success()
}

// 获取当前请求用户信息
func GetCurrentUser(c *gin.Context) models.SysUser {
	userId, exists := c.Get("user")
	var newUser models.SysUser
	if !exists {
		return newUser
	}
	uid := utils.Str2Uint(fmt.Sprintf("%d", userId))
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
	var r request.CreateUserReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	// 记录当前创建人信息
	r.Creator = user.Nickname + user.Username
	s := service.New(c)
	// 将初始密码转为密文
	r.Password = utils.GenPwd(r.InitPassword)
	err := s.Create(r, new(models.SysUser))
	resp.CheckErr(err)
	resp.Success()
}

// 更新用户
func UpdateUserById(c *gin.Context) {
	var r request.UpdateUserReq
	req.ShouldBind(c, &r)

	// 获取path中的userId
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		resp.CheckErr("用户编号不正确")
	}

	// 填写了新密码
	if r.NewPassword != nil && strings.TrimSpace(*r.NewPassword) != "" {
		password := utils.GenPwd(*r.NewPassword)
		r.Password = &password
	}

	user := GetCurrentUser(c)
	if userId == user.Id {
		if r.Status != nil && uint(*r.Status) == models.SysUserStatusDisabled {
			resp.CheckErr("不能禁用自己")
		}
		if r.RoleId != nil && user.RoleId != *r.RoleId {
			if *user.Role.Sort != models.SysRoleSuperAdminSort {
				resp.CheckErr("无法更改自己的角色, 如需更改请联系上级领导")
			} else {
				resp.CheckErr("无法更改超级管理员的角色")
			}
		}
	}

	s := service.New(c)
	err := s.UpdateById(userId, r, new(models.SysUser))
	resp.CheckErr(err)
	userInfoCache.Delete(fmt.Sprintf("%d", user.Id))
	userByIdCache.Delete(fmt.Sprintf("%d", user.Id))
	resp.Success()
}

// 批量删除用户
func BatchDeleteUserByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	if utils.ContainsUint(r.GetUintIds(), user.Id) {
		resp.CheckErr("不能删除自己")
	}

	s := service.New(c)
	err := s.DeleteByIds(r.GetUintIds(), new(models.SysUser))
	resp.CheckErr(err)
	userInfoCache.Delete(fmt.Sprintf("%d", user.Id))
	userByIdCache.Delete(fmt.Sprintf("%d", user.Id))
	resp.Success()
}
