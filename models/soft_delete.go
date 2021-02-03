package models

import (
	"database/sql/driver"
	"fmt"
	"gin-web/pkg/global"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
)

// 本地时间
type DeletedAt struct {
	time.Time
}

func (t *DeletedAt) UnmarshalJSON(data []byte) (err error) {
	str := strings.Trim(string(data), "\"")
	// ""空值不进行解析
	// 避免环包调用, 不再调用utils
	if str == "null" || strings.TrimSpace(str) == "" {
		*t = DeletedAt{Time: time.Time{}}
		return
	}
	
	// 设置str
	t.SetString(str)
	return
}

func (t DeletedAt) MarshalJSON() ([]byte, error) {
	s := t.Format(global.SecLocalTimeFormat)
	// 处理时间0值
	if t.IsZero() {
		s = ""
	}
	output := fmt.Sprintf("\"%s\"", s)
	return []byte(output), nil
}

// gorm 写入 mysql 时调用
func (t DeletedAt) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// gorm 检出 mysql 时调用
func (t *DeletedAt) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DeletedAt{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to DeletedAt", v)
}

// 用于 fmt.Println 和后续验证场景
func (t DeletedAt) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(global.SecLocalTimeFormat)
}

// 设置字符串
func (t *DeletedAt) SetString(str string) *DeletedAt {
	if t != nil {
		// 指定解析的格式(设置转为本地格式)
		now, err := time.ParseInLocation(global.SecLocalTimeFormat, str, time.Local)
		if err == nil {
			*t = DeletedAt{Time: now}
		}
	}
	return t
}

func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteQueryClause{Field: f}}
}

type SoftDeleteQueryClause struct {
	Field *schema.Field
}

func (sd SoftDeleteQueryClause) Name() string {
	return ""
}

func (sd SoftDeleteQueryClause) Build(clause.Builder) {
}

func (sd SoftDeleteQueryClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok {
		if c, ok := stmt.Clauses["WHERE"]; ok {
			if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
				for _, expr := range where.Exprs {
					if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
						where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
						c.Expression = where
						stmt.Clauses["WHERE"] = c
						break
					}
				}
			}
		}

		stmt.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
		}})
		stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
	}
}

func (DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteDeleteClause{Field: f}}
}

type SoftDeleteDeleteClause struct {
	Field *schema.Field
}

func (sd SoftDeleteDeleteClause) Name() string {
	return ""
}

func (sd SoftDeleteDeleteClause) Build(clause.Builder) {
}

func (sd SoftDeleteDeleteClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.String() == "" {
		curTime := stmt.DB.NowFunc()
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: curTime}})
		stmt.SetColumn(sd.Field.DBName, curTime, true)

		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
		}

		if _, ok := stmt.Clauses["WHERE"]; !stmt.DB.AllowGlobalUpdate && !ok {
			stmt.DB.AddError(gorm.ErrMissingWhereClause)
		} else {
			SoftDeleteQueryClause{Field: sd.Field}.ModifyStatement(stmt)
		}

		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build("UPDATE", "SET", "WHERE")
	}
}
