package request

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/utils"
	"strings"
)

// 适用于大多数场景的请求参数绑定
type Req struct {
	Ids string `json:"ids" form:"ids"` // 传多个id
}

// 获取
func (s *Req) GetUintIds() []uint {
	return utils.Str2UintArr(s.Ids)
}

// 增量更新id集合结构体
type UpdateIncrementalIdsRequestStruct struct {
	Create []uint `json:"create"` // 需要新增的编号集合
	Delete []uint `json:"delete"` // 需要删除的编号集合
}

// 获取增量, 可直接更新的结果
func (s *UpdateIncrementalIdsRequestStruct) GetIncremental(oldMenuIds []uint, allMenu []models.SysMenu) []uint {
	// 保留选中流水线
	s.Create = models.GetCheckedMenuIds(s.Create, allMenu)
	s.Delete = models.GetCheckedMenuIds(s.Delete, allMenu)
	newList := make([]uint, 0)
	for _, oldItem := range oldMenuIds {
		// 已删除数据不加入列表
		if !utils.Contains(s.Delete, oldItem) {
			newList = append(newList, oldItem)
		}
	}
	// 将需要新增的数据合并
	return append(newList, s.Create...)
}

// 请求数字类型, 先用string接收(类似json.Number, 相比之下多扩展几种类型)
type ReqNumber string

// 获取字符串
func (n ReqNumber) String() string {
	return string(n)
}

// 获取int
func (n ReqNumber) Int() (int, bool) {
	str := n.String()
	if utils.StrIsEmpty(str) {
		return 0, false
	}
	return utils.Str2Int(str), true
}

// 获取uint
func (n ReqNumber) Uint() (uint, bool) {
	str := n.String()
	if utils.StrIsEmpty(str) {
		return 0, false
	}
	return utils.Str2Uint(str), true
}

// 获取uint32
func (n ReqNumber) Uint32() (uint32, bool) {
	str := n.String()
	if utils.StrIsEmpty(str) {
		return 0, false
	}
	return utils.Str2Uint32(str), true
}

// 获取float64
func (n ReqNumber) Float64() (float64, bool) {
	str := n.String()
	if utils.StrIsEmpty(str) {
		return 0, false
	}
	return utils.Str2Float64(str), true
}

// 获取bool
func (n ReqNumber) Bool() (bool, bool) {
	str := n.String()
	if utils.StrIsEmpty(str) {
		return false, false
	}
	return utils.Str2Bool(str), true
}

func (n *ReqNumber) UnmarshalJSON(data []byte) (err error) {
	str := strings.Trim(string(data), "\"")
	// ""空值不进行解析
	if utils.StrIsEmpty(str) {
		*n = ReqNumber(0)
		return
	}

	// 指定解析的格式(设置转为本地格式)
	*n = ReqNumber(str)
	return
}

func (n ReqNumber) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", n.String())
	return []byte(output), nil
}
