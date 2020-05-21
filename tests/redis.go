package tests

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/go-redis/redis"
)

// 初始化redis数据库
func Redis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Conf.Redis.Host, global.Conf.Redis.Port),
		DB:       global.Conf.Redis.Database,
		Password: global.Conf.Redis.Password,
	})
	global.Redis = client
	global.Log.Debug("[单元测试]初始化redis完成")
}
