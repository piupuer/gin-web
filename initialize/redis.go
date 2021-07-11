package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/go-redis/redis"
	"time"
)

// 初始化redis数据库
func Redis() {
	if !global.Conf.System.UseRedis {
		global.Log.Info("未使用redis, 无需初始化")
		return
	}
	init := false
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				if !init {
					panic(fmt.Sprintf("初始化redis异常: 连接超时(%ds)", global.Conf.System.ConnectTimeout))
				}
				// 此处需return避免协程空跑
				return
			}
		}
	}()
	if !global.Conf.Redis.Sentinel.Enable {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", global.Conf.Redis.Host, global.Conf.Redis.Port),
			DB:       global.Conf.Redis.Database,
			Password: global.Conf.Redis.Password,
		})
		err := client.Ping().Err()
		if err != nil {
			panic(fmt.Sprintf("初始化redis异常: %v", err))
		}
		global.Redis = client
	} else {
		sf := &redis.FailoverOptions{
			MasterName:    global.Conf.Redis.Sentinel.MasterName,
			SentinelAddrs: global.Conf.Redis.Sentinel.AddressArr,
			Password:      global.Conf.Redis.Password,
			DB:            global.Conf.Redis.Database,
		}
		client := redis.NewFailoverClient(sf)
		err := client.Ping().Err()
		if err != nil {
			panic(fmt.Sprintf("初始化redis哨兵异常: %v", err))
		}
		global.Redis = client
	}

	init = true
	global.Log.Info("初始化redis完成")
}
