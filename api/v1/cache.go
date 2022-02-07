package v1

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/piupuer/go-helper/pkg/utils"
)

// get user info from cache by uid
func CacheGetUserInfo(ctx context.Context, uid uint) (*response.UserInfo, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(ctx, getUserInfoCacheKeyPrefix(), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			item := response.UserInfo{}
			utils.Json2Struct(res, &item)
			return &item, true
		}
	}
	return nil, false
}

// set user info to cache by uid
func CacheSetUserInfo(ctx context.Context, uid uint, data response.UserInfo) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(ctx, getUserInfoCacheKeyPrefix(), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// delete user info
func CacheDeleteUserInfo(ctx context.Context, uid uint) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(ctx, getUserInfoCacheKeyPrefix(), fmt.Sprintf("%d", uid))
	}
}

// clear user info cache
func CacheFlushUserInfo(ctx context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(ctx, getUserInfoCacheKeyPrefix())
	}
}

// get user from cache by uid
func CacheGetUser(ctx context.Context, uid uint) (*models.SysUser, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(ctx, getUserCacheKeyPrefix(), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			item := models.SysUser{}
			utils.Json2Struct(res, &item)
			return &item, true
		}
	}
	return nil, false
}

// set user to cache by uid
func CacheSetUser(ctx context.Context, uid uint, data models.SysUser) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(ctx, getUserCacheKeyPrefix(), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// delete user
func CacheDeleteUser(ctx context.Context, uid uint) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(ctx, getUserCacheKeyPrefix(), fmt.Sprintf("%d", uid))
	}
}

// clear user cache
func CacheFlushUser(ctx context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(ctx, getUserCacheKeyPrefix())
	}
}

func getUserCacheKeyPrefix() string {
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, "v1_user")
}

func getUserInfoCacheKeyPrefix() string {
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, "v1_user_info")
}
