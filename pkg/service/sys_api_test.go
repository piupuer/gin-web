package service

import (
	"fmt"
	"go-shipment-api/tests"
	"testing"
)

func TestGetRoleCategoryApisByRoleId(t *testing.T) {
	tests.InitTestEnv()

	s := New(nil)
	m1, a1, _ := s.GetAllApiGroupByCategoryByRoleId(3)
	m2, a2, _ := s.GetAllApiGroupByCategoryByRoleId(1)
	fmt.Println(m1, a1)
	fmt.Println(m2, a2)

	defer s.tx.Close()
}
