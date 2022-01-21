package initialize

import (
	"context"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/job"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/query"
	"os"
)

func Cron() {
	j, err := job.New(
		job.Config{
			RedisClient: global.Redis,
		},
		job.WithCtx(ctx),
		job.WithAutoRequestId(true),
	)
	if err != nil {
		panic(err)
	}
	expr := os.Getenv("CRON_RESET")
	// expr := "@every 600s"
	if expr != "" {
		j.AddTask(job.GoodTask{
			Name: "reset",
			Expr: expr,
			Func: reset,
		}).Start()
	}
	log.WithRequestId(ctx).Debug("initialize cron job success")
}

func reset(ctx context.Context) error {
	ctx = query.NewRequestId(ctx, constant.MiddlewareRequestIdCtxKey)
	log.WithRequestId(ctx).Info("[cron job][reset]starting...")

	if global.Conf.Redis.EnableBinlog {
		global.Redis.Del(ctx, []string{
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_event",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_event_item",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_event_role_relation",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_event_src_item_relation",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_event_user_relation",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_log",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_log_approval_role_relation",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_log_approval_user_relation",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_machine",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_role",
			global.Conf.Mysql.DSN.DBName + "_tb_fsm_user",
			global.Conf.Mysql.DSN.DBName + "_tb_leave",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_api",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_dict",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_dict_data",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_machine",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_menu",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_menu_role_relation",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_message",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_message_log",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_role",
			global.Conf.Mysql.DSN.DBName + "_tb_sys_user",
		}...)
	}
	tables := make([]string, 0)
	global.Mysql.Raw("show tables").Scan(&tables)
	for _, item := range tables {
		if item == "tb_sys_operation_log" || item == "tb_sys_casbin" {
			continue
		}
		global.Mysql.Exec("TRUNCATE TABLE " + item)
	}
	Data()

	log.WithRequestId(ctx).Info("[cron job][reset]ended")
	return nil
}
