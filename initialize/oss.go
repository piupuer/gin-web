package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/oss"
	"time"
)

// 初始化对象存储
func Oss() {
	// 初始化minio
	Minio()
	// 这里预留其他对象存储，如阿里云/七牛云等
	// global.Log.Info("初始化对象存储完成")
}

// 初始化minio对象存储
func Minio() {
	if !global.Conf.Upload.Minio.Enable {
		global.Log.Info("未开启minio, 无需初始化")
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
					panic(fmt.Sprintf("初始化minio异常: 连接超时(%ds)", global.Conf.System.ConnectTimeout))
				}
				// 此处需return避免协程空跑
				return
			}
		}
	}()
	minio := oss.GetMinio(
		global.Conf.Upload.Minio.Endpoint,
		global.Conf.Upload.Minio.AccessId,
		global.Conf.Upload.Minio.Secret,
		global.Conf.Upload.Minio.UseHttps,
	)

	// 初始化一个默认存储桶
	minio.MakeBucket(global.Conf.Upload.Minio.Bucket)
	init = true
	global.Minio = minio
	global.Log.Info("初始化对象存储: minio完成")
}
