package models

import (
	"database/sql/driver"
	"fmt"
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
	return fmt.Sprintf("%s%s", "tb_prefix_", name)
}

// 自定义时间json转换
const TimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

// gorm 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

// gorm 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
