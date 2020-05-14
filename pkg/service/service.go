package service

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CommonService struct {
	tx *gorm.DB // 事务对象实例
	db *gorm.DB // 无事务对象实例
}

// 初始化服务
func New(c *gin.Context) CommonService {
	// 获取事务对象
	tx := global.GetTx(c)
	return CommonService{
		tx: tx,
		db: global.Mysql,
	}
}
