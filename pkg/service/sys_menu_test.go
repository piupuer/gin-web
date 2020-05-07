package service

import (
	"fmt"
	"go-shipment-api/pkg/global"
	"go-shipment-api/tests"
	"testing"
)

func TestGetAllMenuByRoleId(t *testing.T) {
	tests.InitTestEnv()

	m1, _ := GetAllMenuByRoleId(3)
	m2, _ := GetAllMenuByRoleId(1)
	fmt.Println(m1)
	fmt.Println(m2)

	defer global.Mysql.Close()
}
