package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/query"
	"github.com/pkg/errors"
	"time"
)

func Redis() {
	if !global.Conf.Redis.Enable {
		log.WithContext(ctx).Info("redis is not enabled")
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
	client, err := query.ParseRedisURI(global.Conf.Redis.Uri)
	if err != nil {
		panic(errors.Wrap(err, "initialize redis failed"))
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		panic(errors.Wrap(err, "initialize redis failed"))
	}
	global.Redis = client

	init = true
	log.WithContext(ctx).Info("initialize redis success")
}
