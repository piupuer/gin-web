package redis

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/schema"
)

const (
	idName        = "id"
	deletedAtName = "deletedAt"
)

// mysql数据行发生变化, 同步数据到redis
func RowChange(e *canal.RowsEvent) {
	database := e.Table.Schema
	table := e.Table.Name
	// 默认以id为主键, 查找id的索引位置
	idIndex := -1
	deletedAtIndex := -1
	for i, column := range e.Table.Columns {
		name := utils.CamelCaseLowerFirst(column.Name)
		if name == idName {
			idIndex = i
		}
		if name == deletedAtName {
			deletedAtIndex = i
		}
		if idIndex >= 0 && deletedAtIndex >= 0 {
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
		row := getRow(changeRows[0], e.Table)
		if row[deletedAtName] == nil {
			// 由于gorm默认执行软删除, 当delete_at为空时才加入redis缓存
			newRows = append(newRows, row)
		}
		break
	case canal.UpdateAction:
		// 更新数据
		// 通过历史数据changeRows[0]去匹配需要更新的数据所在索引
		index := getOldRowIndex(newRows, changeRows[0], e.Table)
		if deletedAtIndex >= 0 && changeRows[0][deletedAtIndex] == nil && changeRows[1][deletedAtIndex] != nil {
			// 由于gorm默认执行软删除, 当delete_at发生变化时清理redis缓存
			newRows = append(newRows[:index], newRows[index+1:]...)
		} else {
			// 执行更新
			newRows[index] = getRow(changeRows[1], e.Table)
		}
		break
	case canal.DeleteAction:
		// 删除数据
		// 找到需要删除的索引
		indexes := make([]int, 0)
		for _, changeRow := range changeRows {
			// 找到没有改变的数据所在索引
			index := getOldRowIndex(newRows, changeRow, e.Table)
			if index > -1 {
				indexes = append(indexes, index)
			}
		}
		// 删除对应数据
		for _, index := range indexes {
			newRows = append(newRows[:index], newRows[index+1:]...)
		}
		break
	}
	// 将数据转为json字符串写入redis, expiration=0永不过期
	err = global.Redis.Set(cacheKey, utils.Struct2Json(newRows), 0).Err()
	if err != nil {
		global.Log.Error("同步binlog增量数据到redis失败: ", err, e)
	}
}

// 获取旧数据所在行索引
func getOldRowIndex(oldRows []map[string]interface{}, data []interface{}, table *schema.Table) int {
	for i, row := range oldRows {
		newRow := getRow(data, table)
		// 比对增量字段
		m := make(gin.H, 0)
		utils.CompareDifferenceStructByJson(row, newRow, &m)
		// 字段没有任何变化
		if len(m) == 0 {
			return i
		}
	}
	return -1
}

// 获取一列
func getRow(data []interface{}, table *schema.Table) map[string]interface{} {
	row := make(map[string]interface{}, 0)
	count := len(data)
	for i, column := range table.Columns {
		var item interface{}
		if i < count {
			// canal没有对tinyint(1)做bool转换, 这里自行转换
			if column.RawType == "tinyint(1)" {
				item = false
				switch data[i].(type) {
				// canal中的tinyint(1)为float64格式
				case float64:
					if int(data[i].(float64)) == 1 {
						item = true
					}
					break
				}
			} else {
				item = data[i]
			}
			// 由于gorm以驼峰命名, 这里将蛇形转为驼峰
			row[utils.CamelCaseLowerFirst(column.Name)] = item
		}
	}
	if count != len(table.Columns) {
		global.Log.Warn(fmt.Sprintf("数据字段可能不匹配, columns: %v, data: %v", table.Columns, data))
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
