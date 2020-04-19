package global

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var (
	Log *zap.SugaredLogger
	Mysql *gorm.DB
)
