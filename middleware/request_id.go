package middleware

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func RequestId(c *gin.Context) {
	// get from request header
	requestId := c.Request.Header.Get(global.RequestIdHeader)

	if requestId == "" {
		uuid4 := uuid.NewV4()
		requestId = uuid4.String()
	}

	// set to context
	c.Set(global.RequestIdContextKey, requestId)

	// set to header
	c.Writer.Header().Set(global.RequestIdHeader, requestId)
	c.Next()
}
