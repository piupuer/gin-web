package v1

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/utils"
)

// operation log save callback
func OperationLogSave(c *gin.Context, list []middleware.OperationRecord) {
	arr := make([]ms.SysOperationLog, len(list))
	utils.Struct2StructByJson(list, &arr)
	global.Mysql.Create(arr)
}

// operation log find api callback
func OperationLogFindApi(c *gin.Context) []middleware.OperationApi {
	list := make([]ms.SysApi, 0)
	global.Mysql.
		Model(&ms.SysApi{}).
		Find(&list)
	r := make([]middleware.OperationApi, 0)
	utils.Struct2StructByJson(list, &r)
	return r
}
