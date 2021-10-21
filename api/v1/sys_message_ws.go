package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/middleware"
	"net/http"
)

// message center websocket

var upgrade = websocket.Upgrader{
	// allow origin request
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartMessageHub(idempotenceOps *middleware.IdempotenceOptions) cache_service.MessageHub {
	s := cache_service.New(nil)
	return s.StartMessageHub(idempotenceOps)
}

func MessageWs(c *gin.Context) {
	h := make(http.Header)
	h.Add(constant.MiddlewareRequestIdHeaderName, c.GetString(constant.MiddlewareRequestIdCtxKey))
	conn, err := upgrade.Upgrade(c.Writer, c.Request, h)
	if err != nil {
		global.Log.Error(c, "upgrade websocket failed: %v", err)
		return
	}

	user := GetCurrentUser(c)
	s := cache_service.New(c)
	s.MessageWs(c, conn, c.Request.Header.Get("Sec-WebSocket-Key"), user, c.ClientIP())
}
