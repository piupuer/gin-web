package tests

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/job"
)

// 初始化redis数据库
func Redis() {
	// parse redis URI
	client, err := job.ParseRedisURI(global.Conf.Redis.Uri)
	if err != nil {
		panic(fmt.Sprintf("[单元测试]初始化redis异常: %v", err))
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Sprintf("[单元测试]初始化redis异常: %v", err))
	}
	global.Redis = client
	global.Log.Debug(ctx, "[单元测试]初始化redis完成")
}
