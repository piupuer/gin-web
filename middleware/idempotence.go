package middleware

import (
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"strings"
	"sync"
)

// 幂等性中间件

var (
	// 记录是否加锁
	idempotenceLock sync.Mutex
	// 存取token(redis未开启的情况下)
	idempotenceMap = make(map[string]bool)
)

// redis lua脚本(先读取key, 直接删除, 获取删除标志)
const lua string = `
local current = redis.call('GET', KEYS[1])
if current == false then
    return '-1';
end
local del = redis.call('DEL', KEYS[1])
if del == 1 then
     return '1';
else
     return '0';
end
`

// 全局异常处理中间件
func Idempotence(c *gin.Context) {
	// 优先从header提取
	token := c.Request.Header.Get(global.Conf.System.IdempotenceTokenName)
	if token == "" {
		token, _ = c.Cookie(global.Conf.System.IdempotenceTokenName)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		response.FailWithMsg(response.IdempotenceTokenEmptyMsg)
	}
	// token校验
	if !CheckIdempotenceToken(token) {
		response.FailWithMsg(response.IdempotenceTokenInvalidMsg)
	}
	c.Next()
}

// 全局异常处理中间件
func GetIdempotenceToken(c *gin.Context) {
	response.SuccessWithData(GenIdempotenceToken())
}

// 生成一个幂等性token
func GenIdempotenceToken() string {
	token := uuid.NewV4().String()
	// 写入redis或map
	if global.Conf.System.UseRedis {
		global.Redis.Set(token, true, 0)
	} else {
		idempotenceLock.Lock()
		defer idempotenceLock.Unlock()
		idempotenceMap[token] = true
	}
	return token
}

// 校验幂等性token
func CheckIdempotenceToken(token string) bool {
	if global.Conf.System.UseRedis {
		// 执行lua脚本
		res, err := global.Redis.Eval(lua, []string{token}).String()
		if err != nil || res != "1" {
			return false
		}
	} else {
		idempotenceLock.Lock()
		defer idempotenceLock.Unlock()
		// 这里只是模拟单机map缓存，如果是分布式系统多个gin-web应用会导致token不一致，此时建议使用redis或其他分布式唯一token组件
		_, ok := idempotenceMap[token]
		if !ok {
			return false
		}
		// 删除map中的值
		delete(idempotenceMap, token)
	}
	return true
}
