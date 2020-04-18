package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jinzhu/gorm"
)

func main() {
	var err error
	db, err := gorm.Open("mysql", "root:root@tcp(localserver:43306)/goshipment?charset=utf8&parseTime=True&loc=Local&timeout=1000ms")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)
}
