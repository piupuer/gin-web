package main

import (
	"context"
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/query"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

var ctx = query.NewRequestId(nil, constant.MiddlewareRequestIdCtxKey)

func main() {
	defer func() {
		if err := recover(); err != nil {
			if global.Log != nil {
				global.Log.Error(ctx, "[%s]project run failed: %v\nstack: %v", global.ProName, err, string(debug.Stack()))
			} else {
				fmt.Printf("[%s]project run failed: %v\nstack: %v\n", global.ProName, err, string(debug.Stack()))
			}
		}
	}()

	// get runtime root
	_, file, _, _ := runtime.Caller(0)
	global.RuntimeRoot = strings.TrimSuffix(file, "main.go")

	// initialize components
	initialize.Config(ctx)
	initialize.Logger()
	initialize.Redis()
	initialize.Mysql()
	initialize.CasbinEnforcer()
	r := initialize.Routers()
	initialize.Data()
	initialize.Cron()
	initialize.Oss()

	host := "0.0.0.0"
	port := global.Conf.System.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	go func() {
		// listen pprof port
		global.Log.Info(ctx, "[%s]debug pprof is running at %s:%d", global.ProName, host, global.Conf.System.PprofPort)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, global.Conf.System.PprofPort), nil); err != nil {
			global.Log.Error(ctx, "[%s]listen pprof error: %v", global.ProName, err)
		}
	}()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error(ctx, "[%s]listen error: %v", global.ProName, err)
		}
	}()

	global.Log.Info(ctx, "[%s]server is running at %s:%d/%s", global.ProName, host, port, global.Conf.System.UrlPrefix)

	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info(ctx, "[%s]shutting down server...", global.ProName)

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Error(ctx, "[%s]server forced to shutdown: %v", global.ProName, err)
	}

	global.Log.Info(ctx, "[%s]server exiting", global.ProName)
}
