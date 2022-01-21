package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/job"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/pkg/errors"
	"time"
)

func Redis() {
	if !global.Conf.Redis.Enable {
		log.WithRequestId(ctx).Info("if redis is not used, there is no need to initialize redis")
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
		panic(errors.Wrap(err, "initialize redis failed"))
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		panic(errors.Wrap(err, "initialize redis failed"))
	}
	global.Redis = client

	init = true
	log.WithRequestId(ctx).Info("initialize redis success")
}
