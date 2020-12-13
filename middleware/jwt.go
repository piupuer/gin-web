package middleware

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           global.Conf.Jwt.Realm,                                 // jwt标识
		Key:             []byte(global.Conf.Jwt.Key),                           // 服务端密钥
		Timeout:         time.Hour * time.Duration(global.Conf.Jwt.Timeout),    // token过期时间
		MaxRefresh:      time.Hour * time.Duration(global.Conf.Jwt.MaxRefresh), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // 有效载荷处理
		IdentityHandler: identityHandler,                                       // 解析Claims
		Authenticator:   login,                                                 // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                          // 用户登录校验成功处理
		Unauthorized:    unauthorized,                                          // 用户登录校验失败处理
		LoginResponse:   loginResponse,                                         // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                        // 登出后的响应
		RefreshResponse: refreshResponse,                                       // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                              // header名称
		TimeFunc:        time.Now,
	})
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user models.SysUser
		// 将用户json转为结构体
		utils.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.Id,
			"user":          v["user"],
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

func login(c *gin.Context) (interface{}, error) {
	var req request.RegisterAndLoginRequestStruct
	// 请求json绑定
	_ = c.ShouldBindJSON(&req)

	// 密码通过RSA解密
	decodeData, err := utils.RSADecrypt([]byte(req.Password), global.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, err
	}

	u := &models.SysUser{
		Username: req.Username,
		Password: string(decodeData),
	}

	// 创建服务
	s := cache_service.New(c)
	// 密码校验
	user, err := s.LoginCheck(u)
	if err != nil {
		return nil, err
	}
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	return map[string]interface{}{
		"user": utils.Struct2Json(user),
	}, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		var user models.SysUser
		// 将用户json转为结构体
		utils.JsonI2Struct(v["user"], &user)
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	global.Log.Debug(fmt.Sprintf("JWT认证失败, 错误码%d, 错误信息%s", code, message))
	if message == response.LoginCheckErrorMsg {
		response.FailWithMsg(response.LoginCheckErrorMsg)
		return
	}
	response.FailWithCode(response.Unauthorized)
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.SuccessWithData(map[string]interface{}{
		"token": token,
		"expires": models.LocalTime{
			Time: expires,
		},
	})
}

func logoutResponse(c *gin.Context, code int) {
	response.Success()
}

func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.SuccessWithData(map[string]interface{}{
		"token": token,
		"expires": models.LocalTime{
			Time: expires,
		},
	})
}
