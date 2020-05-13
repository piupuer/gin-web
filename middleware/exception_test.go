package middleware

import (
	"fmt"
	"go-shipment-api/pkg/global"
	"go-shipment-api/tests"
	"sync"
	"testing"
)

func TestException(t *testing.T) {
	tests.InitTestEnv()
	// 测试并发访问事务会不会有问题
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			fmt.Println("开启事务")
			tx := global.Mysql.Begin()
			if err := tx.Error; err != nil {
				fmt.Println("开启事务失败", tx.Error)
			} else {
				fmt.Println("开启事务成功")
			}
			defer func() {
				tx.Rollback()
				wg.Done()
			}()
		}()
	}
	wg.Wait()
}
