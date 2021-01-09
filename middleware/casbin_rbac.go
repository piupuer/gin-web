package middleware

import (
	v1 "gin-web/api/v1"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

var checkLock sync.Mutex

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware(c *gin.Context) {
	// 获取当前登录用户
	user := v1.GetCurrentUser(c)
	// 当前登录用户的角色关键字作为casbin访问实体sub
	sub := user.Role.Keyword
	// 请求URL路径作为casbin访问资源obj(需先清除path前缀)
	obj := strings.Replace(c.Request.URL.Path, "/"+global.Conf.System.UrlPathPrefix, "", 1)
	// 请求方式作为casbin访问动作act
	act := c.Request.Method
	// 创建服务
	s := cache_service.New(c)
	// 校验是否有权限访问资源
	if !check(sub, obj, act, s) {
		response.FailWithCode(response.Forbidden)
		return
	}
	// 处理请求
	c.Next()
}

func check(sub, obj, act string, s cache_service.RedisService) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()
	// 获取casbin策略管理器
	e, err := s.Casbin()
	if err != nil {
		return false
	}
	// 检查策略
	pass, _ := e.Enforce(sub, obj, act)
	return pass
}
