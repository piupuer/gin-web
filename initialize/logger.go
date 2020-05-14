package initialize

import (
	"gin-web/pkg/global"
)

// 初始化日志
func Logger() {
	global.InitLogger()
}
