package initialize

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"runtime/debug"
)

// 使用siddontang/go-mysql监听mysql binlog
func MysqlBinlog(tables []string) {
	if !global.Conf.System.UseRedis {
		global.Log.Debug("未使用redis, 无需初始化mysql binlog监听器")
		return
	}
	// 监听器配置
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf(fmt.Sprintf("%s:%d", global.Conf.Mysql.Host, global.Conf.Mysql.Port))
	cfg.User = global.Conf.Mysql.Username
	cfg.Password = global.Conf.Mysql.Password
	// 数据库类型mysql/mariadb
	cfg.Flavor = "mysql"
	// 集群中的唯一编号, 单机随意设定
	cfg.ServerID = 100
	// dump程序名
	cfg.Dump.ExecutionPath = "mysqldump"
	// 目标数据库
	cfg.Dump.TableDB = global.Conf.Mysql.Database
	// 目标表名
	cfg.Dump.Tables = tables

	// 创建canal实例
	c, err := canal.NewCanal(cfg)
	if err != nil {
		global.Log.Debug("初始化mysql binlog监听器失败: ", err)
	}
	// 设置事件处理器
	c.SetEventHandler(&BinlogEventHandler{})
	// 从指定位置开始加载(go 后台运行)
	go c.RunFrom(redis.GetCurrentPos())
	global.Log.Debug("初始化mysql binlog监听器完成")
}

// 自定义事件处理器
type BinlogEventHandler struct {
	canal.DummyEventHandler
}

// 数据行发生变化
func (s *BinlogEventHandler) OnRow(e *canal.RowsEvent) error {
	// 避免监听器发生未知异常导致程序退出, 这里加defer
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("[OnRow]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
			return
		}
	}()
	global.Log.Debug(fmt.Sprintf("行变化: %s %v", e.Action, e.Rows))
	// 同步数据到redis
	redis.RowChange(e)
	return nil
}

// 日志位置发生变化
func (s *BinlogEventHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	// 避免监听器发生未知异常导致程序退出, 这里加defer
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("[OnPosSynced]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
			return
		}
	}()
	global.Log.Debug(fmt.Sprintf("日志位置变化: %s %v %t", pos, set, force))
	redis.PosChange(pos)
	return nil
}

// 处理器名称
func (s *BinlogEventHandler) String() string {
	return "BinlogEventHandler"
}
