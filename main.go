package main

import (
	"context"
	"fmt"
	"go-shipment-api/initialize"
	"go-shipment-api/pkg/global"
	"net/http"
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
			global.Log.Error(fmt.Sprintf("项目启动失败: %v\n堆栈信息: %v", err, string(string(debug.Stack()))))
		}
	}()

	// 初始化日志
	initialize.Logger()

	// 初始化数据库
	initialize.Mysql()

	// 结束后关闭数据库
	defer global.Mysql.Close()

	// 初始化路由
	r := initialize.Routers()

	// 初始化数据
	initialize.InitData()

	host := "127.0.0.1"
	port := 9999
	// 服务器启动以及优雅的关闭
	// 参考地址https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error("listen error: ", err)
		}
	}()

	global.Log.Info(fmt.Sprintf("Server is running at %s:%d", host, port))

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
