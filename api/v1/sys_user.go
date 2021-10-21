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
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"strings"
)

func GetUserInfo(c *gin.Context) {
	user := GetCurrentUser(c)
	oldCache, ok := CacheGetUserInfo(c, user.Id)
	if ok {
		resp.SuccessWithData(oldCache)
		return
	}

	var rp response.UserInfoResp
	utils.Struct2StructByJson(user, &rp)
	rp.Roles = []string{
		"admin",
	}
	rp.Keyword = user.Role.Keyword
	rp.RoleSort = *user.Role.Sort
	CacheSetUserInfo(c, user.Id, rp)
	resp.SuccessWithData(rp)
}

func FindUser(c *gin.Context) {
	var r request.UserReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.CurrentRole = user.Role
	s := cache_service.New(c)
	list, err := s.FindUser(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.UserResp{}, r.Page)
}

func ChangePwd(c *gin.Context) {
	var r request.ChangePwdReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	query := global.Mysql.Where("username = ?", user.Username).First(&user)
	err := query.Error
	resp.CheckErr(err)
	if ok := utils.ComparePwd(r.OldPassword, user.Password); !ok {
		resp.CheckErr("the original password is incorrect")
	}
	err = query.Update("password", utils.GenPwd(r.NewPassword)).Error
	resp.CheckErr(err)
	resp.Success()
}

func GetCurrentUser(c *gin.Context) models.SysUser {
	userId, exists := c.Get("user")
	var newUser models.SysUser
	if !exists {
		return newUser
	}
	uid := utils.Str2Uint(fmt.Sprintf("%d", userId))
	oldCache, ok := CacheGetUser(c, uid)
	if ok {
		return *oldCache
	}
	s := service.New(c)
	newUser, _ = s.GetUserById(uid)
	CacheSetUser(c, uid, newUser)
	return newUser
}

func CreateUser(c *gin.Context) {
	var r request.CreateUserReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	// plaintext to ciphertext
	r.Password = utils.GenPwd(r.InitPassword)
	err := s.Q.Create(r, new(models.SysUser))
	resp.CheckErr(err)
	resp.Success()
}

func UpdateUserById(c *gin.Context) {
	var r request.UpdateUserReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)

	// new password is not empty, update password
	if r.NewPassword != nil && strings.TrimSpace(*r.NewPassword) != "" {
		password := utils.GenPwd(*r.NewPassword)
		r.Password = &password
	}

	user := GetCurrentUser(c)
	if id == user.Id {
		if r.Status != nil && uint(*r.Status) == models.SysUserStatusDisabled {
			resp.CheckErr("cannot disable yourself")
		}
		if r.RoleId != nil && user.RoleId != *r.RoleId {
			if *user.Role.Sort != models.SysRoleSuperAdminSort {
				resp.CheckErr("cannot change your role")
			} else {
				resp.CheckErr("cannot change super admin's role")
			}
		}
	}

	s := service.New(c)
	err := s.Q.UpdateById(id, r, new(models.SysUser))
	resp.CheckErr(err)
	CacheDeleteUserInfo(c, user.Id)
	CacheDeleteUser(c, user.Id)
	resp.Success()
}

func BatchDeleteUserByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	if utils.ContainsUint(r.Uints(), user.Id) {
		resp.CheckErr("cannot remove yourself")
	}

	s := service.New(c)
	err := s.Q.DeleteByIds(r.Uints(), new(models.SysUser))
	resp.CheckErr(err)
	CacheFlushUserInfo(c)
	CacheFlushUser(c)
	resp.Success()
}
