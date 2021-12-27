package initialize

import (
	"context"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/job"
	"github.com/piupuer/go-helper/pkg/query"
	"os"
)

func Cron() {
	j, err := job.New(
		job.Config{
			RedisClient: global.Redis,
		},
		job.WithLogger(global.Log),
		job.WithCtx(ctx),
		job.WithAutoRequestId(true),
	)
	if err != nil {
		panic(err)
	}
	expr := os.Getenv("CRON_RESET")
	if expr != "" {
		j.AddTask(job.GoodTask{
			Name: "reset",
			Expr: expr,
			Func: reset,
		}).Start()
	}
	global.Log.Debug(ctx, "initialize cron job success")
}

func reset(c context.Context) error {
	ctx := query.NewRequestId(c, constant.MiddlewareRequestIdCtxKey)
	global.Log.Info(ctx, "[cron job][reset]starting...")

	if global.Conf.Redis.EnableBinlog {
		global.Redis.FlushAll(ctx)
	}
	tables := make([]string, 0)
	global.Mysql.Raw("show tables").Scan(&tables)
	for _, item := range tables {
		if item == "tb_sys_operation_log" {
			continue
		}
		global.Mysql.Exec("TRUNCATE TABLE " + item)
	}
	Data()

	global.Log.Info(ctx, "[cron job][reset]ended")
	return nil
}
