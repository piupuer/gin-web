package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/job"
	"time"
)

func Redis() {
	if !global.Conf.System.UseRedis {
		global.Log.Info(ctx, "if redis is not used, there is no need to initialize redis")
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
					panic(fmt.Sprintf("initialize redis failed: connect timeout(%ds)", global.Conf.System.ConnectTimeout))
				}
				return
			}
		}
	}()
	// parse redis URI
	client, err := job.ParseRedisURI(global.Conf.Redis.Uri)
	if err != nil {
		panic(fmt.Sprintf("initialize redis failed: %v", err))
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Sprintf("initialize redis failed: %v", err))
	}
	global.Redis = client

	init = true
	global.Log.Info(ctx, "initialize redis success")
}
