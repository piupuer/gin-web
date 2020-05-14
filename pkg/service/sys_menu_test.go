package service

import (
	"fmt"
	"gin-web/tests"
	"testing"
)

func TestGetAllMenuByRoleId(t *testing.T) {
	tests.InitTestEnv()

	s := New(nil)
	m1, a1, _ := s.GetAllMenuByRoleId(3)
	m2, a2, _ := s.GetAllMenuByRoleId(1)
	fmt.Println(m1, a1)
	fmt.Println(m2, a2)

	defer s.tx.Close()
}
