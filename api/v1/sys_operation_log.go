package v1

import (
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/utils"
)

// operation log save callback
func OperationLogSave(c *gin.Context, list []middleware.OperationRecord) {
	arr := make([]ms.SysOperationLog, len(list))
	utils.Struct2StructByJson(list, &arr)
	my := service.New(c)
	my.Q.Db.Create(arr)
}

// operation log find skip path callback
func OperationLogFindSkipPath(c *gin.Context) []string {
	my := service.New(c)
	return my.Q.FindDictDataValByName(constant.MiddlewareOperationLogSkipPathDict)
}

// operation log find api callback
func OperationLogFindApi(c *gin.Context) []middleware.OperationApi {
	list := make([]ms.SysApi, 0)
	my := service.New(c)
	my.Q.Db.
		Model(&ms.SysApi{}).
		Find(&list)
	r := make([]middleware.OperationApi, 0)
	utils.Struct2StructByJson(list, &r)
	return r
}
