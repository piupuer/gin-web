package initialize

import (
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/logger"
)

func Logger() {
	global.Log = logger.NewDefaultWrapper()
	global.Log.Debug("initialize logger success")
}
