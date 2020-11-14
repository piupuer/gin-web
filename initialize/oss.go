package initialize

import (
	"gin-web/pkg/global"
	"gin-web/pkg/oss"
)

// 初始化对象存储
func Oss() {
	// 初始化minio
	InitMinio()
	// 这里预留其他对象存储，如阿里云/七牛云等
	// global.Log.Info("初始化对象存储完成")
}

// 初始化minio对象存储
func InitMinio() {
	if !global.Conf.Upload.Minio.Enable {
		global.Log.Info("未开启minio, 无需初始化")
		return
	}
	global.Minio = oss.GetMinio(
		global.Conf.Upload.Minio.Endpoint,
		global.Conf.Upload.Minio.AccessId,
		global.Conf.Upload.Minio.Secret,
		global.Conf.Upload.Minio.UseHttps,
	)
	// 初始化一个默认存储桶
	global.Minio.MakeBucket(global.Conf.Upload.Minio.Bucket)
	global.Log.Info("初始化对象存储: minio完成")
}
