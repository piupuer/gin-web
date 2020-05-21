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
