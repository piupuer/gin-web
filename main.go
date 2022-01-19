package main

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/listen"
	"github.com/piupuer/go-helper/pkg/query"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"strings"
)

var ctx = query.NewRequestId(nil, constant.MiddlewareRequestIdCtxKey)

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
			if global.Log != nil {
				global.Log.Error("[%s]project run failed: %v\nstack: %v", global.ProName, err, string(debug.Stack()))
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
	initialize.Data()
	initialize.Cron()
	initialize.Oss()

	// listen http
	listen.Http(
		listen.WithHttpLogger(global.Log),
		listen.WithHttpCtx(ctx),
		listen.WithHttpProName(global.ProName),
		listen.WithHttpPort(global.Conf.System.Port),
		listen.WithHttpPprofPort(global.Conf.System.PprofPort),
		listen.WithHttpHandler(router.RegisterServers(ctx)),
	)
}
