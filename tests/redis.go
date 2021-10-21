package tests

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/job"
)

func Redis() {
	// parse redis URI
	client, err := job.ParseRedisURI(global.Conf.Redis.Uri)
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize redis failed: %v", err))
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize redis failed: %v", err))
	}
	global.Redis = client
	global.Log.Debug(ctx, "[unit test]initialize redis success")
}
