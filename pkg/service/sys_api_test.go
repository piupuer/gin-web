package service

import (
	"fmt"
	"go-shipment-api/pkg/global"
	"go-shipment-api/tests"
	"testing"
)

func TestGetRoleCategoryApisByRoleId(t *testing.T) {
	tests.InitTestEnv()

	m1, a1, _ := GetAllApiGroupByCategoryByRoleId(3)
	m2, a2, _ := GetAllApiGroupByCategoryByRoleId(1)
	fmt.Println(m1, a1)
	fmt.Println(m2, a2)

	defer s.tx.Close()
}
