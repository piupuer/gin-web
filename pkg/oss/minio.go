package oss

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"time"
)

// minio对象存储
type MinioOss struct {
	// minio客户端实例
	Client *minio.Client
	// 当前会话
	Ctx context.Context
}

// 获取minio实例
func GetMinio(endpoint, accessId, secret string, useHttps bool) *MinioOss {
	ctx := context.Background()
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessId, secret, ""),
		Secure: useHttps,
	})

	if err != nil {
		panic(fmt.Sprintf("获取minio实例失败: %v", err))
	}
	return &MinioOss{
		Client: minioClient,
		Ctx:    ctx,
	}
}

// 创建存储桶bucketName(系统未配置可用区)
func (s *MinioOss) MakeBucket(bucketName string) {
	s.MakeBucketWithLocation(bucketName, "")
}

// 创建存储桶bucketName(在可用区location中)
func (s *MinioOss) MakeBucketWithLocation(bucketName, location string) {
	err := s.Client.MakeBucket(s.Ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.Client.BucketExists(s.Ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Println(fmt.Errorf("存储桶%s已经存在(可用区%s)", bucketName, location))
		} else {
			fmt.Println("创建存储桶失败: ", err)
		}
	}
}

// 查找符合条件的对象
func (s *MinioOss) ListObjects(bucketName, prefix string, recursive bool) <-chan minio.ObjectInfo {
	return s.Client.ListObjects(s.Ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})
}

// 上传一个对象(本地)
func (s *MinioOss) PutLocalObject(bucketName, objectName, filePath string) error {
	_, err := s.Client.FPutObject(s.Ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}

// 上传一个对象(文件流)
func (s *MinioOss) PutObject(bucketName, objectName string, file io.Reader, fileSize int64) error {
	_, err := s.Client.PutObject(s.Ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{})
	return err
}

// 批量删除对象
func (s *MinioOss) RemoveObjects(bucketName string, objectNames []string) error {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, name := range objectNames {
			objectsCh <- minio.ObjectInfo{
				Key: name,
			}
		}
	}()

	for rErr := range s.Client.RemoveObjects(s.Ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		if rErr.Err != nil {
			return rErr.Err
		}
	}
	return nil
}

// 获取一个对象的预览地址
func (s *MinioOss) GetObjectPreviewUrl(bucketName, objectName string) string {
	u, err := s.Client.PresignedGetObject(s.Ctx, bucketName, objectName, time.Second*24*60*60, url.Values{})
	if err != nil {
		return ""
	}
	return u.String()
}

// 判断一个对象是否存在
func (s *MinioOss) ObjectExists(bucketName, objectName string) bool {
	_, err := s.Client.StatObject(s.Ctx, bucketName, objectName, minio.StatObjectOptions{})
	return err == nil
}
