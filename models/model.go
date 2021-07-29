package models

import (
	"database/sql/driver"
	"fmt"
	"gin-web/pkg/global"
	"strings"
	"time"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	Id        uint      `gorm:"primaryKey;comment:'自增编号'" json:"id"`
	CreatedAt LocalTime `gorm:"comment:'创建时间'" json:"createdAt"`
	UpdatedAt LocalTime `gorm:"comment:'更新时间'" json:"updatedAt"`
	DeletedAt DeletedAt `gorm:"index:idx_deleted_at;comment:'删除时间(软删除)'" json:"deletedAt"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, name)
}

// 本地时间
type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	str := strings.Trim(string(data), "\"")
	// ""空值不进行解析
	// 避免环包调用, 不再调用utils
	if str == "null" || strings.TrimSpace(str) == "" {
		*t = LocalTime{Time: time.Time{}}
		return
	}

	// 设置str
	t.SetString(str)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	s := t.Format(global.SecLocalTimeFormat)
	// 处理时间0值
	if t.IsZero() {
		s = ""
	}
	output := fmt.Sprintf("\"%s\"", s)
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
	return fmt.Errorf("can not convert %v to LocalTime", v)
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(global.SecLocalTimeFormat)
}

// 只需要日期
func (t LocalTime) DateString() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(global.DateLocalTimeFormat)
}

// 只需要月份
func (t LocalTime) MonthString() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(global.MonthLocalTimeFormat)
}

// 设置字符串
func (t *LocalTime) SetString(str string) *LocalTime {
	if t != nil {
		// 指定解析的格式(设置转为本地格式)
		now, err := time.ParseInLocation(global.SecLocalTimeFormat, str, time.Local)
		if err == nil {
			*t = LocalTime{Time: now}
			return t
		}
		nowDate, err := time.ParseInLocation(global.DateLocalTimeFormat, str, time.Local)
		if err == nil {
			*t = LocalTime{Time: nowDate}
		}
		nowMonth, err := time.ParseInLocation(global.MonthLocalTimeFormat, str, time.Local)
		if err == nil {
			*t = LocalTime{Time: nowMonth}
		}
	}
	return t
}

// 设置时分
func (t *LocalTime) SetHourAndMinuteString(str string) *LocalTime {
	if t != nil {
		if len(str) != 5 {
			str = "00:00"
		}
		dateStart := t.TodayStart().Format(global.DateLocalTimeFormat)
		return t.SetString(fmt.Sprintf("%s %s:00", dateStart, str))
	}
	return t
}

// 获取今日0点
func (t *LocalTime) TodayStart() *LocalTime {
	if t != nil {
		if t.IsZero() {
			t.Time = time.Now()
		}
		now := t.Time
		// 取当日毫秒数
		dateStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		t.Time = dateStart
	}
	return t
}

// 获取明天0点
func (t *LocalTime) TomorrowStart() *LocalTime {
	if t != nil {
		t.Time = t.TodayStart().Time.AddDate(0, 0, 1)
	}
	return t
}

// 获取下个月0点
func (t *LocalTime) NextMonthStart() *LocalTime {
	if t != nil {
		t.Time = t.TodayStart().Time.AddDate(0, 1, 0)
	}
	return t
}

// 获取当前日期与目标日期之间的全部日期
func (t *LocalTime) GetDates(str string) []string {
	res := make([]string, 0)
	endT := new(LocalTime).SetString(str).TodayStart()
	end := *endT
	start := *t.TodayStart()
	// 交换位置
	if start.Time.After(end.Time) {
		tmp := start
		start = end
		end = tmp
	}
	endStr := end.DateString()
	startStr := start.DateString()
	if startStr == endStr {
		res = append(res, startStr)
		return res
	}
	for {
		current := start.AddDate(0, 0, 1)
		currentStr := start.DateString()
		start.Time = current
		res = append(res, currentStr)
		if currentStr == endStr {
			break
		}
	}
	return res
}

// 获取当前日期与目标日期之间的全部月份
func (t *LocalTime) GetMonths(str string) []string {
	res := make([]string, 0)
	endT := new(LocalTime).SetString(str).TodayStart()
	end := *endT
	start := *t.TodayStart()
	// 交换位置
	if start.Time.After(end.Time) {
		tmp := start
		start = end
		end = tmp
	}
	endStr := end.MonthString()
	startStr := start.MonthString()
	if startStr == endStr {
		res = append(res, startStr)
		return res
	}
	for {
		current := start.AddDate(0, 1, 0)
		currentStr := start.MonthString()
		start.Time = current
		res = append(res, currentStr)
		if currentStr == endStr {
			break
		}
	}
	return res
}
