package middleware

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/tests"
	"sync"
	"testing"
)

func TestTransaction(t *testing.T) {
	tests.InitTestEnv()
	// 测试并发访问事务会不会有问题
	var wg sync.WaitGroup
	// 这里10000按需设置, 可能会导致mysql崩溃, 建议使用本地数据库
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
