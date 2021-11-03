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
	CacheSuffixUserInfo = "user_info"
	CacheSuffixUser     = "user"
)

// get user info from cache by uid
func CacheGetUserInfo(c context.Context, uid uint) (*response.UserInfo, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixUserInfo), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			item := response.UserInfo{}
			utils.Json2Struct(res, &item)
			return &item, true
		}
	}
	return nil, false
}

// set user info to cache by uid
func CacheSetUserInfo(c context.Context, uid uint, data response.UserInfo) {
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
