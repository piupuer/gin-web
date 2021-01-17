package oss

import (
	"fmt"
	"testing"
)

func TestGetMinio(t *testing.T) {
	m := GetMinio("127.0.0.1", "minio", "minio123", false)
	bucketName := "test"
	m.MakeBucket(bucketName)

	err := m.PutLocalObject(bucketName, "1.jpg", "/Users/piupuer/Downloads/images/1.jpg")
	fmt.Println(err)

	fmt.Println(m.GetObjectPreviewUrl(bucketName, "2.xls"))
	fmt.Println(m.GetObjectPreviewUrl(bucketName, "1.jpg"))
	fmt.Println(m.ObjectExists(bucketName, "2.xls"))
	fmt.Println(m.ObjectExists(bucketName, "3.xls"))

}

func TestOSSMinio_RemoveObjects(t *testing.T) {
	m := GetMinio("127.0.0.1", "minio", "minio123", false)
	bucketName := "test"
	fmt.Println(m.RemoveObjects(bucketName, []string{"1.jpg"}))
}
