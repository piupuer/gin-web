package v1

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/utils"
)

// CacheGetUserInfo get user info from cache by uid
func CacheGetUserInfo(ctx context.Context, uid uint) (u response.UserInfo, exists bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(ctx, getUserInfoCacheKeyPrefix(), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			utils.Json2Struct(res, &u)
			if u.Id > constant.Zero {
				exists = true
			}
			return
		}
	}
	return
}

// CacheSetUserInfo set user info to cache by uid
func CacheSetUserInfo(ctx context.Context, uid uint, data response.UserInfo) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(ctx, getUserInfoCacheKeyPrefix(), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// CacheDeleteUserInfo delete user info
func CacheDeleteUserInfo(ctx context.Context, uid uint) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(ctx, getUserInfoCacheKeyPrefix(), fmt.Sprintf("%d", uid))
	}
}

// CacheFlushUserInfo clear user info cache
func CacheFlushUserInfo(ctx context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(ctx, getUserInfoCacheKeyPrefix())
	}
}

// CacheGetUser get user from cache by uid
func CacheGetUser(ctx context.Context, uid uint) (u models.SysUser, exists bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(ctx, getUserCacheKeyPrefix(), fmt.Sprintf("%d", uid)).Result()
		if err == nil && res != "" {
			utils.Json2Struct(res, &u)
			if u.Id > constant.Zero {
				exists = true
			}
			return
		}
	}
	return
}

// CacheSetUser set user to cache by uid
func CacheSetUser(ctx context.Context, uid uint, data models.SysUser) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(ctx, getUserCacheKeyPrefix(), fmt.Sprintf("%d", uid), utils.Struct2Json(data))
	}
}

// CacheDeleteUser delete user
func CacheDeleteUser(ctx context.Context, uid uint) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(ctx, getUserCacheKeyPrefix(), fmt.Sprintf("%d", uid))
	}
}

// CacheFlushUser clear user cache
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
