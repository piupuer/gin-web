package service

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/tests"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestMysqlService_DlDl(t *testing.T) {
	tests.InitTestEnv()
	tableName := "login"
	type Entity struct {
		Id     int    `gorm:"column:id" json:"id"`
		Mobile string `gorm:"column:username" json:"mobile"`
	}

	cache := global.Redis
	mysql := global.Mysql
	list := make([]Entity, 0)
	query := mysql.Table(tableName)
	_ = query.Find(&list).Error
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, item := range list {
		if item.Mobile == "" {
			continue
		}
		oldMobileRandom, _ := cache.Get(item.Mobile).Result()
		var random int
		if oldMobileRandom != "" {
			// 缓存存在
			random, _ = strconv.Atoi(oldMobileRandom)
		} else {
			random = r.Intn(1000000)
			// 写入缓存
			global.Redis.Set(item.Mobile, random, 0)
		}
		mobile, _ := strconv.Atoi(item.Mobile)
		var entity Entity
		entity.Mobile = strconv.Itoa(mobile+random)
		err := query.Where("id = ?", item.Id).Updates(&entity).Error
		fmt.Println(err)
	}
}
