package main

import (
	"context"
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"net/http"
  // 加入pprof性能分析
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("项目启动失败: %v\n堆栈信息: %v", err, string(debug.Stack())))
		}
	}()

	// 初始化配置
	initialize.Config()

	// 初始化日志
	initialize.Logger()

	// 初始化redis数据库
	initialize.Redis()

	// 初始化mysql数据库
	initialize.Mysql()

	// 初始校验器
	initialize.Validate()

	// 结束后关闭数据库(gorm2.0升级为连接池模式, 无需手动关闭)

	// 初始化路由
	r := initialize.Routers()

	// 初始化数据
	initialize.Data()

	// 初始化定时任务
	initialize.Cron()
	
	// 初始化对象存储
	initialize.Oss()

	host := "0.0.0.0"
	port := global.Conf.System.Port
	// 服务器启动以及优雅的关闭
	// 参考地址https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	go func() {
		// 加入pprof性能分析
		if err := http.ListenAndServe(":8005", nil); err != nil {
			global.Log.Error("listen pprof error: ", err)
		}
	}()
	
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error("listen error: ", err)
		}
	}()

	global.Log.Info(fmt.Sprintf("Server is running at %s:%d/%s", host, port, global.Conf.System.UrlPathPrefix))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Error("Server forced to shutdown: ", err)
	}

	global.Log.Info("Server exiting")
}
