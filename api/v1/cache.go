package v1

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/piupuer/go-helper/pkg/utils"
)

const (
	CacheSuffixMenuTree = "menu_tree"
	CacheSuffixUserInfo = "user_info"
	CacheSuffixUser     = "user"
)

// get menu tree from cache by uid
func CacheGetMenuTree(c context.Context, uid uint) ([]response.MenuTreeResp, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixMenuTree), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			list := make([]response.MenuTreeResp, 0)
			utils.Json2Struct(res, &list)
			return list, true
		}
	}
	return nil, false
}

// set menu tree to cache by uid
func CacheSetMenuTree(c context.Context, uid uint, data []response.MenuTreeResp) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixMenuTree), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// clear menu tree cache
func CacheFlushMenuTree(c context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixMenuTree))
	}
}

// get user info from cache by uid
func CacheGetUserInfo(c context.Context, uid uint) (*response.UserInfoResp, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUserInfo), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			item := response.UserInfoResp{}
			utils.Json2Struct(res, &item)
			return &item, true
		}
	}
	return nil, false
}

// set user info to cache by uid
func CacheSetUserInfo(c context.Context, uid uint, data response.UserInfoResp) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUserInfo), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// delete user info
func CacheDeleteUserInfo(c context.Context, uid uint) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUserInfo), fmt.Sprintf("%d", uid))
	}
}

// clear user info cache
func CacheFlushUserInfo(c context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUserInfo))
	}
}

// get user from cache by uid
func CacheGetUser(c context.Context, uid uint) (*models.SysUser, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUser), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			item := models.SysUser{}
			utils.Json2Struct(res, &item)
			return &item, true
		}
	}
	return nil, false
}

// set user to cache by uid
func CacheSetUser(c context.Context, uid uint, data models.SysUser) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUser), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// delete user
func CacheDeleteUser(c context.Context, uid uint) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUser), fmt.Sprintf("%d", uid))
	}
}

// clear user cache
func CacheFlushUser(c context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUser))
	}
}
