package models

import (
	"database/sql/driver"
	"fmt"
	"gin-web/pkg/global"
	"time"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	Id        uint       `gorm:"primary_key;comment:'自增编号'" json:"id"`
	CreatedAt LocalTime  `gorm:"comment:'创建时间'" json:"createdAt"`
	UpdatedAt LocalTime  `gorm:"comment:'更新时间'" json:"updatedAt"`
	DeletedAt *LocalTime `gorm:"comment:'删除时间(软删除)'" sql:"index" json:"deletedAt"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s%s", global.Conf.Mysql.TablePrefix, name)
}

// 自定义时间json转换
const TimeFormat = "2006-01-02 15:04:05"
const DateFormat = "2006-01-02"

type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime{Time: time.Time{}}
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime{Time: now}
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(output), nil
}

// gorm 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// gorm 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return t.Format(TimeFormat)
}

// 只需要日期
func (t LocalTime) DateString() string {
	return t.Format(DateFormat)
}
