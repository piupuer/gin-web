package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// 消息中心websocket

// 升级请求头
var upgrade = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 创建消息中心仓库
func StartMessageHub(checkIdempotenceTokenFunc func(token string) bool) cache_service.MessageHub {
	s := cache_service.New(nil)
	return s.StartMessageHub(checkIdempotenceTokenFunc)
}

// 启动消息连接
func MessageWs(c *gin.Context) {
	h := make(http.Header)
	h.Add(global.RequestIdHeader, c.GetString(global.RequestIdContextKey))
	conn, err := upgrade.Upgrade(c.Writer, c.Request, h)
	if err != nil {
		global.Log.Error(c, "创建消息连接失败", err)
		return
	}

	// 获取当前登录用户
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	// 启动连接
	s.MessageWs(conn, c.Request.Header.Get("Sec-WebSocket-Key"), user, c.ClientIP())
}
