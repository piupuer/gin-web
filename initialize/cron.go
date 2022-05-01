package initialize

import (
	"context"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/delay"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/migrate"
	"github.com/pkg/errors"
	"os"
)

func Cron() {
	qu := delay.NewQueue(
		delay.WithQueueName(global.ProName),
		delay.WithQueueRedisUri(global.Conf.Redis.Uri),
		delay.WithQueueHandler(taskHandler),
	)
	if qu.Error != nil {
		panic(errors.Wrap(qu.Error, "initialize cron job failed"))
	}
	expr := os.Getenv("CRON_RESET")
	// expr := "@every 600s"
	if expr != "" {
		err := qu.Cron(
			delay.WithQueueTaskUuid("cron_reset_task"),
			delay.WithQueueTaskName("reset"),
			delay.WithQueueTaskExpr(expr),
		)
		if err != nil {
			panic(errors.Wrap(err, "initialize cron job failed"))
		}
	}
	log.WithContext(ctx).Info("initialize cron job success")
}

func taskHandler(ctx context.Context, t delay.Task) (err error) {
	switch t.Uid {
	case "cron_reset_task":
		log.WithContext(ctx).Info("[cron job][reset]starting...")

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
				global.Conf.Mysql.DSN.DBName + "_tb_sys_casbin",
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
		db := global.Mysql.WithContext(ctx)
		db.Raw("show tables").Scan(&tables)
		for _, item := range tables {
			if item == "tb_sys_operation_log" {
				continue
			}
			db.Exec("TRUNCATE TABLE " + item)
		}
		// reset data
		uri := global.Conf.Mysql.Uri
		migrate.Do(
			migrate.WithCtx(ctx),
			migrate.WithUri(uri),
			migrate.WithFs(sqlFs),
			migrate.WithFsRoot("db"),
		)
		log.WithContext(ctx).Info("[cron job][reset]ended")
	}
	return
}
