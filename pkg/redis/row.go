package redis

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/golang-module/carbon"
)

const (
	idName        = "id"
	deletedAtName = "deletedAt"
)

// mysql数据行发生变化, 同步数据到redis
func RowChange(ctx context.Context, e *canal.RowsEvent) {
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
	// gorm更新到v2版本后, 某些字段在e.Rows中的表现为[]uint8类型([]byte的别名)
	// 如果遇到json中有字段是uint8类型, golang json包会转为base64字符串, 这里将uint8转为string
	rows := make([][]interface{}, len(e.Rows))
	for i, eRow := range e.Rows {
		row := make([]interface{}, len(eRow))
		for j, eItem := range eRow {
			if eV, ok := eItem.([]uint8); ok {
				row[j] = string(eV)
			} else {
				columnType := e.Table.Columns[j].RawType
				// 转为本地时间
				// grom v2时间类型全部是datetime(3)
				if t, ok := eItem.(string); ok && columnType == "datetime(3)" {
					eItem = carbon.ToDateTimeString{
						Carbon: carbon.Parse(t),
					}
				}
				row[j] = eItem
			}
		}
		rows[i] = row
	}

	// 缓存键由数据库名与表名组成
	cacheKey := fmt.Sprintf("%s_%s", database, table)
	// 读取redis历史数据
	oldRowsStr, err := global.Redis.Get(ctx, cacheKey).Result()
	newRows := make([]map[string]interface{}, 0)
	changeRows := make([][]interface{}, 0)
	if err == nil {
		// 解压缩字符串
		oldRows := utils.DeCompressStrByZlib(oldRowsStr)
		// 将旧数据解析为对象
		utils.Json2Struct(oldRows, &newRows)
	}
	rowCount := len(newRows)
	// 将rows用json解析一下, 否则查找相同元素时可能出现类型不一致
	utils.Struct2StructByJson(rows, &changeRows)

	// gorm更新到v2版本后blob对象不再需要base64解码

	// 选择事件类型
	switch e.Action {
	case canal.InsertAction:
		// 插入数据
		for _, changeRow := range changeRows {
			row := getRow(ctx, changeRow, e.Table)
			if row[deletedAtName] == nil {
				// 由于gorm默认执行软删除, 当delete_at为空时才加入redis缓存
				newRows = append(newRows, row)
			}
		}
		break
	case canal.UpdateAction:
		// 更新数据(新旧数据2个2个一组)
		for i, l := 0, len(changeRows); i < l; i += 2 {
			oldRow := changeRows[i]
			newRow := changeRows[i+1]
			// 通过历史数据changeRows[0]去匹配需要更新的数据所在索引
			index := getOldRowIndex(ctx, newRows, oldRow, e.Table)
			if len(newRows) > 0 && index >= 0 {
				if deletedAtIndex >= 0 && oldRow[deletedAtIndex] == nil && newRow[deletedAtIndex] != nil {
					// 由于gorm默认执行软删除, 当delete_at发生变化时清理redis缓存
					if index < rowCount-1 {
						newRows = append(newRows[:index], newRows[index+1:]...)
					} else {
						newRows = append(newRows[:index])
					}
				} else {
					// 执行更新
					newRows[index] = getRow(ctx, newRow, e.Table)
				}
			} else {
				// 可能是数据反写
				newRows = append(newRows, getRow(ctx, newRow, e.Table))
			}
		}

		break
	case canal.DeleteAction:
		// 找到需要删除的索引
		indexes := make([]int, 0)
		for _, changeRow := range changeRows {
			// 找到没有改变的数据所在索引
			index := getOldRowIndex(ctx, newRows, changeRow, e.Table)
			if index > -1 {
				indexes = append(indexes, index)
			}
		}
		// 记录被删除的元素个数
		deletedCount := 0
		// 删除对应数据
		for _, index := range indexes {
			i := index - deletedCount
			if index < rowCount-1 {
				newRows = append(newRows[:i], newRows[i+1:]...)
				deletedCount++
			} else {
				newRows = append(newRows[:i])
			}
		}
		break
	}
	// 压缩后写入
	compress, err := utils.CompressStrByZlib(utils.Struct2Json(newRows))
	if err != nil {
		global.Log.Error(ctx, "同步binlog增量数据到redis失败: %v, %v", err, e)
		return
	}
	// 将数据转为json字符串写入redis, expiration=0永不过期
	err = global.Redis.Set(ctx, cacheKey, *compress, 0).Err()
	if err != nil {
		global.Log.Error(ctx, "同步binlog增量数据到redis失败: %v, %v", err, e)
	}
}

// 获取旧数据所在行索引
func getOldRowIndex(ctx context.Context, oldRows []map[string]interface{}, data []interface{}, table *schema.Table) int {
	newRow := getRow(ctx, data, table)
	for i, row := range oldRows {
		// 比对增量字段
		m := make(map[string]interface{}, 0)
		utils.CompareDifferenceStructByJson(row, newRow, &m)
		// 字段没有任何变化
		if len(m) == 0 {
			return i
		}
	}
	return -1
}

// 获取一列
func getRow(ctx context.Context, data []interface{}, table *schema.Table) map[string]interface{} {
	row := make(map[string]interface{}, 0)
	count := len(data)
	for i, column := range table.Columns {
		var item interface{}
		if i < count {
			// canal没有对tinyint(1)转换, 这里自行转换为uint
			if column.RawType == "tinyint(1)" {
				switch data[i].(type) {
				// canal中的tinyint(1)为float64格式
				case float64:
					item = uint(data[i].(float64))
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
		global.Log.Warn(ctx, "数据字段可能不匹配, columns: %v, data: %v", table.Columns, data)
	}
	return row
}

// mysql日志位置发生变化
func PosChange(ctx context.Context, pos mysql.Position) {
	err := global.Redis.Set(ctx, global.Conf.Redis.BinlogPos, utils.Struct2Json(pos), 0).Err()
	if err != nil {
		global.Log.Error(ctx, "同步binlog当前位置到redis失败: %v, %v", err, pos)
	}
}
