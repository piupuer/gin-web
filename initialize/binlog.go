package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
	"gin-web/pkg/utils"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/golang-module/carbon"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// 使用siddontang/go-mysql监听mysql binlog
func MysqlBinlog(ignoreTables []string, tableModels ...interface{}) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		global.Log.Info(ctx, "未使用redis或未开启binlog, 无需初始化mysql binlog监听器")
		return
	}
	l := len(tableModels)
	tableNames := make([]string, l)
	for i := 0; i < l; i++ {
		t := reflect.ValueOf(tableModels[i]).Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		tableNames[i] = global.Mysql.NamingStrategy.TableName(reflect.New(t).Elem().Type().Name())
	}
	// 监听器配置
	cfg := canal.NewDefaultConfig()

	cfg.Addr = global.Conf.Mysql.DSN.Addr
	cfg.User = global.Conf.Mysql.DSN.User
	cfg.Password = global.Conf.Mysql.DSN.Passwd
	// 数据库类型mysql/mariadb
	cfg.Flavor = "mysql"
	// 集群中的唯一编号, 单机随意设定
	cfg.ServerID = 100
	// dump程序名
	cfg.Dump.ExecutionPath = "mysqldump"
	// 目标数据库
	cfg.Dump.TableDB = global.Conf.Mysql.DSN.DBName
	// 目标表名
	cfg.Dump.Tables = tableNames
	// server编号(从1开始)
	cfg.ServerID = global.Conf.System.MachineId + 1

	// 创建canal实例
	c, err := canal.NewCanal(cfg)
	if err != nil {
		global.Log.Info(ctx, "初始化mysql binlog监听器失败: %v", err)
	}
	// 添加忽略表
	c.AddDumpIgnoreTables(cfg.Dump.TableDB, ignoreTables...)
	// 设置事件处理器
	c.SetEventHandler(&BinlogEventHandler{
		ctx:          global.RequestIdContext(""),
		binlogPos:    fmt.Sprintf("%s_%s", global.Conf.Redis.BinlogPos, global.Conf.Mysql.DSN.DBName),
		IgnoreTables: ignoreTables,
	})
	// 刷新数据
	refresh(tableNames, tableModels)
	// 从最后一个位置开始运行
	pos, _ := c.GetMasterPos()
	go c.RunFrom(pos)
	global.Log.Info(ctx, "初始化mysql binlog监听器完成")
}

// 自定义事件处理器
type BinlogEventHandler struct {
	ctx       context.Context
	binlogPos string
	canal.DummyEventHandler
	IgnoreTables []string
}

// 同步器启动时, 刷新redis数据
func refresh(tableNames []string, tableModels []interface{}) {
	database := global.Conf.Mysql.DSN.DBName
	for i, table := range tableNames {
		// 缓存键由数据库名与表名组成
		cacheKey := fmt.Sprintf("%s_%s", database, table)
		// 查询mysql数据
		oldRows, err := getRows(table, tableModels[i])
		if err != nil {
			continue
		}
		newRows := make([]map[string]interface{}, 0)
		for _, oldRow := range oldRows {
			row := make(map[string]interface{}, 0)
			for key, item := range oldRow {
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
		err = global.Redis.Set(ctx, cacheKey, *compress, 0).Err()
		if err != nil {
			panic(fmt.Sprintf("刷新redis数据失败, %v", err))
		}
	}
}

func getRows(table string, model interface{}) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0)
	// 读取数据
	rows, err := global.Mysql.Table(table).Rows()
	if err != nil {
		return nil, err
	}
	// 结束时关闭
	defer rows.Close()
	// 列名
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	mt := reflect.TypeOf(model).Elem()

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		item := make(map[string]interface{}, 0)
		for i, colName := range cols {
			val := columns[i]
			var s interface{}
			// 反射model中的字段
			field, exists := mt.FieldByName(utils.CamelCase(colName))
			if exists && val != nil {
				switch val.(type) {
				case time.Time:
					local := carbon.Time2Carbon(val.(time.Time))
					s = local.String()
				case []uint8:
					vs := string(val.([]uint8))
					k := field.Type.Kind()
					if field.Type.Kind() == reflect.Ptr {
						// 指针变量继续下一层
						k = field.Type.Elem().Kind()
					}
					switch k {
					case reflect.Uint:
						s = utils.Str2Uint(vs)
					case reflect.Int:
						f, _ := strconv.Atoi(vs)
						s = f
					case reflect.Float64:
						f, _ := strconv.ParseFloat(vs, 64)
						s = f
					case reflect.Float32:
						f, _ := strconv.ParseFloat(vs, 32)
						s = f
					default:
						s = vs
					}
				}
			} else {
				s = nil
			}
			item[colName] = s
		}
		list = append(list, item)
	}
	return list, nil
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
			global.Log.Error(s.ctx, "[OnRow]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack()))
			return
		}
	}()
	global.Log.Debug(s.ctx, "行变化: %s %v", e.Action, e.Rows)
	// 同步数据到redis
	redis.RowChange(s.ctx, e)
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
			err := global.Redis.Del(ctx, cacheKey).Err()
			if err != nil {
				global.Log.Error(s.ctx, "删除表%s, 同步binlog增量数据到redis失败: %v", table, err)
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
			err := global.Redis.Del(ctx, cacheKey).Err()
			if err != nil {
				global.Log.Error(s.ctx, "清空表%s, 同步binlog增量数据到redis失败: %v", table, err)
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
			global.Log.Error(s.ctx, "[OnPosSynced]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack()))
			return
		}
	}()
	global.Log.Debug(s.ctx, "日志位置变化: %s %v %t", pos, set, force)
	redis.PosChange(s.ctx, s.binlogPos, pos)
	return nil
}

// 处理器名称
func (s *BinlogEventHandler) String() string {
	return "BinlogEventHandler"
}
