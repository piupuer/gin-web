package tests

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/query"
)

func Redis() {
	// parse redis URI
	client, err := query.ParseRedisURI(global.Conf.Redis.Uri)
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize redis failed: %v", err))
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize redis failed: %v", err))
	}
	global.Redis = client
	log.WithContext(ctx).Debug("[unit test]initialize redis success")
}
