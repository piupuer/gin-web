package initialize

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
	"gin-web/pkg/utils"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

// 使用siddontang/go-mysql监听mysql binlog
func MysqlBinlog(tables, ignoreTables []string) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		global.Log.Info("未使用redis或未开启binlog, 无需初始化mysql binlog监听器")
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
		global.Log.Infof("初始化mysql binlog监听器失败: ", err)
	}
	// 添加忽略表
	c.AddDumpIgnoreTables(cfg.Dump.TableDB, ignoreTables...)
	// 设置事件处理器
	c.SetEventHandler(&BinlogEventHandler{
		IgnoreTables: ignoreTables,
	})
	// 刷新数据
	refresh(tables)
	// 从最后一个位置开始运行
	pos, _ := c.GetMasterPos()
	go c.RunFrom(pos)
	global.Log.Info("初始化mysql binlog监听器完成")
}

// 自定义事件处理器
type BinlogEventHandler struct {
	canal.DummyEventHandler
	IgnoreTables []string
}

// 同步器启动时, 刷新redis数据
func refresh(tables []string) {
	database := global.Conf.Mysql.Database
	for _, table := range tables {
		// 缓存键由数据库名与表名组成
		cacheKey := fmt.Sprintf("%s_%s", database, table)
		// 查询mysql数据
		oldRows := make([]map[string]interface{}, 0)
		err := global.Mysql.Table(table).Scan(&oldRows).Error
		if err != nil {
			continue
		}
		newRows := make([]map[string]interface{}, 0)
		for _, oldRow := range oldRows {
			row := make(map[string]interface{}, 0)
			for key, item := range oldRow {
				// 转为本地时间
				if t, ok := item.(time.Time); ok {
					item = models.LocalTime{
						Time: t,
					}
				}
				if t, ok := item.(*time.Time); ok {
					if t != nil {
						item = &models.LocalTime{
							Time: *t,
						}
					}
				}
				// 由于gorm以驼峰命名, 这里将蛇形转为驼峰
				row[utils.CamelCaseLowerFirst(key)] = item
			}
			newRows = append(newRows, row)
		}
		// 压缩后写入
		compress, err := utils.CompressStrByZlib(utils.Struct2Json(newRows))
		if err != nil {
			panic(fmt.Sprintf("刷新redis数据失败, %v", err))
		}
		// 将数据转为json字符串写入redis, expiration=0永不过期
		err = global.Redis.Set(cacheKey, *compress, 0).Err()
		if err != nil {
			panic(fmt.Sprintf("刷新redis数据失败, %v", err))
		}
	}
}

// 数据行发生变化
func (s *BinlogEventHandler) OnRow(e *canal.RowsEvent) error {
	if utils.Contains(s.IgnoreTables, e.Table.Name) {
		return nil
	}
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

// ddl事件
func (s *BinlogEventHandler) OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	database := string(queryEvent.Schema)
	sql := strings.ToLower(string(queryEvent.Query))
	dropReg := regexp.MustCompile("drop table `(.+?)`")
	if dropReg != nil {
		// 提取关键信息
		if m := dropReg.FindAllStringSubmatch(sql, -1); len(m) == 1 {
			table := strings.Trim(m[0][1], "`")
			cacheKey := fmt.Sprintf("%s_%s", database, table)
			// 将数据转为json字符串写入redis, expiration=0永不过期
			err := global.Redis.Del(cacheKey).Err()
			if err != nil {
				global.Log.Errorf("删除表%s, 同步binlog增量数据到redis失败: %v", table, err)
			}
		}
	}
	if strings.Contains(sql, "truncate table") {
		table := ""
		arr := strings.Split(sql, " ")
		l := len(arr)
		for i, item := range arr {
			if item == "table" && i < l {
				table = strings.Trim(arr[i+1], "`")
			}
		}
		if table != "" {
			cacheKey := fmt.Sprintf("%s_%s", database, table)
			// 将数据转为json字符串写入redis, expiration=0永不过期
			err := global.Redis.Del(cacheKey).Err()
			if err != nil {
				global.Log.Errorf("清空表%s, 同步binlog增量数据到redis失败: %v", table, err)
			}
		}
	}
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
