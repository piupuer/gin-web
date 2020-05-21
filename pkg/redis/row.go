package redis

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/schema"
)

// mysql数据行发生变化, 同步数据到redis
func RowChange(e *canal.RowsEvent) {
	database := e.Table.Schema
	table := e.Table.Name
	// 默认以id为主键, 查找id的索引位置
	idName := "id"
	idIndex := 0
	for i, column := range e.Table.Columns {
		if column.Name == idName {
			idIndex = i
			break
		}
	}
	// 缓存键由数据库名与表名组成
	cacheKey := fmt.Sprintf("%s_%s", database, table)
	// 读取redis历史数据
	oldRows, err := global.Redis.Get(cacheKey).Result()
	newRows := make([]map[string]interface{}, 0)
	changeRows := make([][]interface{}, 0)
	if err == nil {
		// 将旧数据解析为对象
		utils.Json2Struct(oldRows, &newRows)
	}
	// 将rows用json解析一下, 否则查找相同元素时可能出现类型不一致
	utils.Struct2StructByJson(e.Rows, &changeRows)
	// 选择事件类型
	switch e.Action {
	case canal.InsertAction:
		// 插入数据
		newRows = append(newRows, getRow(changeRows[0], e.Table.Columns))
		break
	case canal.UpdateAction:
		// 更新数据
		for i, row := range newRows {
			// 找到相同id
			if row["id"] == changeRows[0][idIndex] {
				// 直接替换原有元素, changeRows[0]表示改变前的元素, changeRows[1]表示改变后的元素
				newRows[i] = getRow(changeRows[1], e.Table.Columns)
				break
			}
		}
		break
	case canal.DeleteAction:
		// 删除数据
		for i, row := range newRows {
			// 找到相同id
			if row[idName] == changeRows[0][idIndex] {
				// 移除当前元素
				newRows = append(newRows[:i], newRows[i+1:]...)
				break
			}
		}
		break
	}
	// 将数据转为json字符串写入redis, expiration=0永不过期
	err = global.Redis.Set(cacheKey, utils.Struct2Json(newRows), 0).Err()
	if err != nil {
		global.Log.Error("同步binlog增量数据到redis失败: ", err, e)
	}
}

// 获取一列
func getRow(data []interface{}, columns []schema.TableColumn) map[string]interface{} {
	row := make(map[string]interface{}, 0)
	for i, column := range columns {
		// 由于gorm以驼峰命名, 这里将蛇形转为驼峰
		row[utils.CamelCaseLowerFirst(column.Name)] = data[i]
	}
	return row
}

// 获取当前日志位置
func GetCurrentPos() mysql.Position {
	var pos mysql.Position
	p, err := global.Redis.Get(global.Conf.Redis.BinlogPos).Result()
	if err == nil {
		// 将旧数据解析为对象
		utils.Json2Struct(p, &pos)
	}
	return pos
}

// mysql日志位置发生变化
func PosChange(pos mysql.Position) {
	err := global.Redis.Set(global.Conf.Redis.BinlogPos, utils.Struct2Json(pos), 0).Err()
	if err != nil {
		global.Log.Error("同步binlog当前位置到redis失败: ", err, pos)
	}
}
