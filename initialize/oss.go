package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/oss"
	"github.com/pkg/errors"
	"time"
)

func Oss() {
	Minio()
}

func Minio() {
	if !global.Conf.Upload.Minio.Enable {
		log.WithContext(ctx).Info("minio is not enabled")
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
		oss.WithMinioEndpoint(global.Conf.Upload.Minio.Endpoint),
		oss.WithMinioAccessId(global.Conf.Upload.Minio.AccessId),
		oss.WithMinioSecret(global.Conf.Upload.Minio.Secret),
		oss.WithMinioHttps(global.Conf.Upload.Minio.UseHttps),
	}

	minio, err := oss.NewMinio(ops...)
	if err != nil {
		panic(errors.Wrap(err, "initialize object storage minio failed"))
	}

	minio.MakeBucket(ctx, global.Conf.Upload.Minio.Bucket)
	init = true
	global.Minio = minio
	log.WithContext(ctx).Info("initialize object storage minio success")
}
