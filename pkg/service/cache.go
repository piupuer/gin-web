package service

import (
	"context"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/utils"
)

const (
	CacheSuffixDictName       = "dict_name"
	CacheSuffixDictNameAndKey = "dict_name_and_key"
)

// get dict name from cache by uid
func CacheGetDictName(c context.Context, name string) ([]models.SysDictData, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictName), name).Result()
		if err == nil && res != "" {
			list := make([]models.SysDictData, 0)
			utils.Json2Struct(res, &list)
			return list, true
		}
	}
	return nil, false
}

// set dict name to cache by uid
func CacheSetDictName(c context.Context, name string, data []models.SysDictData) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictName), name, utils.Struct2Json(data))
	}
}

// delete dict name
func CacheDeleteDictName(c context.Context, name string) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictName), name)
	}
}

// clear dict name cache
func CacheFlushDictName(c context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictName))
	}
}

// get dict name and key from cache by uid
func CacheGetDictNameAndKey(c context.Context, name, key string) (*models.SysDictData, bool) {
	if global.Conf.Redis.Enable {
		res, err := global.Redis.HGet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictNameAndKey), fmt.Sprintf("%s_%s", name, key)).Result()
		if err == nil && res != "" {
			item := models.SysDictData{}
			utils.Json2Struct(res, &item)
			return &item, true
		}
	}
	return nil, false
}

// set dict name and key to cache by uid
func CacheSetDictNameAndKey(c context.Context, name, key string, data models.SysDictData) {
	if global.Conf.Redis.Enable {
		global.Redis.HSet(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictNameAndKey), fmt.Sprintf("%s_%s", name, key), utils.Struct2Json(data))
	}
}

// delete dict name and key
func CacheDeleteDictNameAndKey(c context.Context, name, key string) {
	if global.Conf.Redis.Enable {
		global.Redis.HDel(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictNameAndKey), fmt.Sprintf("%s_%s", name, key))
	}
}

// clear dict name and key cache
func CacheFlushDictNameAndKey(c context.Context) {
	if global.Conf.Redis.Enable {
		global.Redis.Del(c, fmt.Sprintf("%s_%s", global.ProName, CacheSuffixDictNameAndKey))
	}
}
