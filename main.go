package main

import (
	"embed"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/piupuer/go-helper/pkg/listen"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/pkg/errors"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"strings"
)

var ctx = tracing.NewId(nil)

//go:embed conf
var conf embed.FS

// @title Gin Web
// @version 1.2.1
// @description A simple RBAC admin system written by golang
// @license.name MIT
// @license.url https://github.com/piupuer/gin-web/blob/dev/LICENSE
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.WithContext(ctx).WithError(errors.Errorf("%v", err)).Error("[%s]project run failed, stack: %s", global.ProName, string(debug.Stack()))
		}
	}()

	// get runtime root
	_, file, _, _ := runtime.Caller(0)
	global.RuntimeRoot = strings.TrimSuffix(file, "main.go")

	// initialize components
	initialize.Config(ctx, conf)
	initialize.Tracer()
	initialize.Redis()
	initialize.Mysql()
	initialize.CasbinEnforcer()
	initialize.Cron()
	initialize.Oss()

	// listen http
	listen.Http(
		listen.WithHttpCtx(ctx),
		listen.WithHttpProName(global.ProName),
		listen.WithHttpPort(global.Conf.System.Port),
		listen.WithHttpPprofPort(global.Conf.System.PprofPort),
		listen.WithHttpHandler(router.RegisterServers(ctx)),
		listen.WithHttpExit(func() {
			global.Tracer.Shutdown(ctx)
		}),
	)
}
