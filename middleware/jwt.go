package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"strconv"
	"time"
)

var jwtSecret = "jwt-secret"

func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test realm",
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     jwtSecret,                                          // jwt密钥
		PayloadFunc:     payloadFunc,                                        // 有效载荷处理
		IdentityHandler: identityHandler,                                    // 解析Claims
		Authenticator:   authenticator,                                      // 校验token的正确性
		Authorizator:    authorizator,                                       // 校验用户的正确性
		Unauthorized:    unauthorized,                                       // 校验失败处理
		LoginResponse:   loginResponse,                                      // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                     // 登出后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt", // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                           // header名称
		TimeFunc:        time.Now,
	})
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		user, _ := v["user"].(models.SysUser)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          user,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

func authenticator(c *gin.Context) (interface{}, error) {
	var req request.RegisterAndLoginStruct
	// 请求json绑定
	_ = c.ShouldBindJSON(&req)

	u := &models.SysUser{
		Username: req.Username,
		Password: req.Password,
	}

	// 密码校验
	user, err := service.LoginCheck(u)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	// 写入用户, payloadFunc/authorizator会使用到
	return map[string]interface{}{
		"user": user,
	}, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userMap, _ := v["user"].(map[string]interface{})
		// id需要从float64转为uint
		s := strconv.FormatFloat(userMap["ID"].(float64), 'f', -1, 64)
		i, _ := strconv.Atoi(s)
		user := &models.SysUser{
			Model: gorm.Model{
				ID: uint(i),
			},
			Username: userMap["username"].(string),
		}
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	global.Log.Debug(message)
	response.FailWithMsg(c, "jwt校验失败")
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.SuccessWithData(c, map[string]interface{}{
		"token":   token,
		"expires": expires,
	})
}

func logoutResponse(c *gin.Context, code int) {
	response.Success(c)
}
