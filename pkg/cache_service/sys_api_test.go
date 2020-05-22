package cache_service

import (
	"fmt"
	"gin-web/pkg/request"
	"gin-web/tests"
	"testing"
)

func TestRedisService_GetApis(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	req1 := &request.ApiListRequestStruct{}
	fmt.Println(s.GetApis(req1))
	req2 := &request.ApiListRequestStruct{
		Method: "GET",
	}
	fmt.Println(s.GetApis(req2))
}

func TestRedisService_GetAllApiGroupByCategoryByRoleId(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	fmt.Println(s.GetAllApiGroupByCategoryByRoleId(1))
	fmt.Println(s.GetAllApiGroupByCategoryByRoleId(2))
	fmt.Println(s.GetAllApiGroupByCategoryByRoleId(3))
	fmt.Println(s.GetAllApiGroupByCategoryByRoleId(1000))
}
