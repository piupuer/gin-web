package middleware

import (
	"github.com/gin-gonic/gin"
	v1 "go-shipment-api/api/v1"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
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
	// 获取casbin策略管理器
	e, err := service.Casbin()
	if err != nil {
		response.FailWithMsg(c, "获取资源访问策略失败")
		c.Abort()
		return
	}
	// 检查策略
	pass, _ := e.Enforce(sub, obj, act)
	if pass {
		c.Next()
	} else {
		response.FailWithCode(c, response.Forbidden)
		c.Abort()
		return
	}
}
