package service

import (
	"fmt"
	"go-shipment-api/pkg/global"
	"go-shipment-api/tests"
	"testing"
)

func TestGetAllMenuByRoleId(t *testing.T) {
	tests.InitTestEnv()

	m1, a1, _ := GetAllMenuByRoleId(3)
	m2, a2, _ := GetAllMenuByRoleId(1)
	fmt.Println(m1, a1)
	fmt.Println(m2, a2)

	defer global.Mysql.Close()
}
