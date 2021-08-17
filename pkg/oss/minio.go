package oss

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm/logger"
	"io"
	"net/url"
	"time"
)

// minio对象存储
type MinioOss struct {
	log logger.Interface
	// minio客户端实例
	Client *minio.Client
}

// 获取minio实例
func GetMinio(log logger.Interface, endpoint, accessId, secret string, useHttps bool) *MinioOss {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessId, secret, ""),
		Secure: useHttps,
	})

	if err != nil {
		panic(fmt.Sprintf("获取minio实例失败: %v", err))
	}
	return &MinioOss{
		log:    log,
		Client: minioClient,
	}
}

// 创建存储桶bucketName(系统未配置可用区)
func (s *MinioOss) MakeBucket(ctx context.Context, bucketName string) {
	s.MakeBucketWithLocation(ctx, bucketName, "")
}

// 创建存储桶bucketName(在可用区location中)
func (s *MinioOss) MakeBucketWithLocation(ctx context.Context, bucketName, location string) {
	err := s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			s.log.Warn(ctx, "存储桶%s已经存在(可用区%s)", bucketName, location)
		} else {
			s.log.Error(ctx, "创建存储桶失败: %v", err)
		}
	}
}

// 查找符合条件的对象
func (s *MinioOss) ListObjects(ctx context.Context, bucketName, prefix string, recursive bool) <-chan minio.ObjectInfo {
	return s.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})
}

// 上传一个对象(本地)
func (s *MinioOss) PutLocalObject(ctx context.Context, bucketName, objectName, filePath string) error {
	_, err := s.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}

// 上传一个对象(文件流)
func (s *MinioOss) PutObject(ctx context.Context, bucketName, objectName string, file io.Reader, fileSize int64) error {
	_, err := s.Client.PutObject(ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{})
	return err
}

// 批量删除对象
func (s *MinioOss) RemoveObjects(ctx context.Context, bucketName string, objectNames []string) error {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, name := range objectNames {
			objectsCh <- minio.ObjectInfo{
				Key: name,
			}
		}
	}()

	for rErr := range s.Client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		if rErr.Err != nil {
			return rErr.Err
		}
	}
	return nil
}

// 获取一个对象的预览地址
func (s *MinioOss) GetObjectPreviewUrl(ctx context.Context, bucketName, objectName string) string {
	u, err := s.Client.PresignedGetObject(ctx, bucketName, objectName, time.Second*24*60*60, url.Values{})
	if err != nil {
		return ""
	}
	return u.String()
}

// 判断一个对象是否存在
func (s *MinioOss) ObjectExists(ctx context.Context, bucketName, objectName string) bool {
	_, err := s.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	return err == nil
}
