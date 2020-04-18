package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jinzhu/gorm"
	"go-shipment-api/models"
)

func main() {
	var err error
	db, err := gorm.Open("mysql", "root:root@tcp(localserver:43306)/goshipment?charset=utf8&parseTime=True&loc=Local&timeout=1000ms")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// 自动迁移表
	db.AutoMigrate(models.User{})

	var user models.User
	tableName := user.TableName()
	// 查询表
	table := db.Table(tableName)

	var user1 models.User
	err = table.Where("id=?", 1).First(&user1).Error
	if err != nil {
		fmt.Println("查询user id=1的用户张三, 不存在创建新用户张三")
		table.Create(&models.User{
			Model: gorm.Model{
				ID: 1,
			},
			Username: "张三",
			Sex:      0,
		})
	} else {
		fmt.Println(fmt.Printf("查询user id=1的用户: 用户名%s, 性别%d", user1.Username, user1.Sex))
	}

	var user2 models.User
	err = table.Where("id=?", 2).First(&user2).Error
	if err != nil {
		fmt.Println("查询user id=2的用户李四, 不存在创建新用户李四")
		table.Create(&models.User{
			Model: gorm.Model{
				ID: 2,
			},
			Username: "李四",
			Sex:      2,
		})
	} else {
		fmt.Println(fmt.Printf("查询user id=2的用户: 用户名%s, 性别%d", user2.Username, user2.Sex))
	}
	var users1 []models.User
	var count1 int
	table.Find(&users1).Count(&count1)
	fmt.Println(fmt.Printf("第1次查询全部数据: 总条数%d, 集合%v", count1, users1))

	// 更新数据
	table.Where("id=?", 1).Update("username", "张五")
	fmt.Println("更新张三的姓名为张五")

	// 删除数据, 使用硬删除, 不保留记录
	table.Unscoped().Where("id=?", 2).Delete(models.User{})
	fmt.Println("删除李四")

	fmt.Println(fmt.Printf("第2次查询全部数据: 总条数%d, 集合%v", count1, users1))
}
