package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/oss"
	"time"
)

func Oss() {
	Minio()
}

func Minio() {
	if !global.Conf.Upload.Minio.Enable {
		global.Log.Info(ctx, "if minio is not enabled, there is no need to initialize minio")
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
					panic(fmt.Sprintf("initialize object storage minio failed: connect timeout(%ds)", global.Conf.System.ConnectTimeout))
				}
				return
			}
		}
	}()
	ops := []func(*oss.MinioOptions){
		oss.WithMinioLogger(global.Log),
		oss.WithMinioEndpoint(global.Conf.Upload.Minio.Endpoint),
		oss.WithMinioAccessId(global.Conf.Upload.Minio.AccessId),
		oss.WithMinioSecret(global.Conf.Upload.Minio.Secret),
		oss.WithMinioHttps(global.Conf.Upload.Minio.UseHttps),
	}

	minio, err := oss.NewMinio(ops...)
	if err != nil {
		panic(fmt.Sprintf("initialize object storage minio failed: %v", err))
	}

	minio.MakeBucket(ctx, global.Conf.Upload.Minio.Bucket)
	init = true
	global.Minio = minio
	global.Log.Info(ctx, "initialize object storage minio success")
}
