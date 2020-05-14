package middleware

import (
	v1 "gin-web/api/v1"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
)

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware(c *gin.Context) {
	// 获取当前登录用户
	user := v1.GetCurrentUser(c)
	// 当前登录用户的角色关键字作为casbin访问实体sub
	sub := user.Role.Keyword
	// 请求URL路径作为casbin访问资源obj
	obj := c.Request.URL.Path
	// 请求方式作为casbin访问动作act
	act := c.Request.Method
	// 创建服务
	s := service.New(c)
	// 获取casbin策略管理器
	e, err := s.Casbin()
	if err != nil {
		response.FailWithMsg("获取资源访问策略失败")
		return
	}
	// 检查策略
	pass, _ := e.Enforce(sub, obj, act)
	if !pass {
		response.FailWithCode(response.Forbidden)
	}
	// 处理请求
	c.Next()
}
