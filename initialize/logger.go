package initialize

import (
	"gin-web/pkg/global"
	"github.com/natefinch/lumberjack"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func Logger() {
	colorful := false
	if global.Conf.Logs.Level <= zapcore.DebugLevel {
		colorful = true
	}
	if global.Mode == constant.Prod {
		colorful = false
	}
	global.Log = logger.New(
		logger.WithLevel(logger.Level(global.Conf.Logs.Level)),
		logger.WithColorful(colorful),
		logger.WithLineNumPrefix(global.RuntimeRoot),
		logger.WithLumberjackOption(logger.LumberjackOption{
			Logger: lumberjack.Logger{
				MaxSize:    global.Conf.Logs.MaxSize,
				MaxBackups: global.Conf.Logs.MaxBackups,
				MaxAge:     global.Conf.Logs.MaxAge,
				Compress:   global.Conf.Logs.Compress,
			},
			LogPath: global.Conf.Logs.Path,
		}),
	)
	global.Log.Debug(ctx, "initialize logger success")
}
