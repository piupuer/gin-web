package main

import (
	"context"
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"net/http"
	"runtime"
	"strings"
	
	// 加入pprof性能分析
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

var ctx = global.RequestIdContext("") // 生成启动时request id

func main() {
	defer func() {
		if err := recover(); err != nil {
			if global.Log != nil {
				// 将异常写入日志
				global.Log.Error(ctx, "项目启动失败: %v\n堆栈信息: %v", err, string(debug.Stack()))
			} else {
				fmt.Printf("项目启动失败: %v\n堆栈信息: %v\n", err, string(debug.Stack()))
			}
		}
	}()

	// get runtime root
	_, file, _, _ := runtime.Caller(0)
	global.RuntimeRoot = strings.TrimSuffix(file, "main.go")

	// 初始化配置
	initialize.Config(ctx)

	// 初始化日志
	initialize.Logger()

	// 初始化redis数据库
	initialize.Redis()

	// 初始化mysql数据库
	initialize.Mysql()

	// 初始化casbin策略管理器
	initialize.CasbinEnforcer()

	// 初始校验器
	initialize.Validate()

	// 结束后关闭数据库(gorm2.0升级为连接池模式, 无需手动关闭)

	// 初始化路由
	r := initialize.Routers()

	// 初始化数据
	initialize.Data()

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
		global.Log.Info(ctx, "Debug pprof is running at %s:%d", host, global.Conf.System.PprofPort)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, global.Conf.System.PprofPort), nil); err != nil {
			global.Log.Error(ctx, "listen pprof error: %v", err)
		}
	}()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error(ctx, "listen error: %v", err)
		}
	}()

	global.Log.Info(ctx, "Server is running at %s:%d/%s", host, port, global.Conf.System.UrlPathPrefix)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info(ctx, "Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Error(ctx, "Server forced to shutdown: %v", err)
	}

	global.Log.Info(ctx, "Server exiting")
}
