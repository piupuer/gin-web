package api

import (
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/resp"
)

func Ping(c *gin.Context) {
	resp.SuccessWithData("pong")
}
