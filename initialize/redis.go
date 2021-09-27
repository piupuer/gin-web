package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/job"
	"time"
)

// 初始化redis数据库
func Redis() {
	if !global.Conf.System.UseRedis {
		global.Log.Info(ctx, "未使用redis, 无需初始化")
		return
	}
	init := false
	ctx, cancel := context.WithTimeout(ctx, time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
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
	// parse redis URI
	client, err := job.ParseRedisURI(global.Conf.Redis.Uri)
	if err != nil {
		panic(fmt.Sprintf("初始化redis异常: %v", err))
	}
	err = client.Ping().Err()
	if err != nil {
		panic(fmt.Sprintf("初始化redis异常: %v", err))
	}
	global.Redis = client

	init = true
	global.Log.Info(ctx, "初始化redis完成")
}
