package request

import (
	"database/sql/driver"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// 适用于大多数场景的请求参数绑定
type Req struct {
	Ids string `json:"ids" form:"ids"` // 传多个id
}

// 获取
func (s Req) GetUintIds() []uint {
	return utils.Str2UintArr(s.Ids)
}

// 增量更新id集合结构体
type UpdateIncrementalIdsRequestStruct struct {
	Create []uint `json:"create"` // 需要新增的编号集合
	Delete []uint `json:"delete"` // 需要删除的编号集合
}

// 获取增量, 可直接更新的结果
func (s UpdateIncrementalIdsRequestStruct) GetIncremental(oldMenuIds []uint, allMenu []models.SysMenu) []uint {
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

// 请求uint类型
type ReqUint uint

func (r *ReqUint) UnmarshalJSON(data []byte) (err error) {
	str := strings.Trim(string(data), "\"")
	// ""空值不进行解析
	if utils.StrIsEmpty(str) {
		*r = ReqUint(0)
		return
	}
	*r = ReqUint(utils.Str2Uint(str))
	return
}

func (r ReqUint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", r)), nil
}

// gorm 写入 mysql 时调用
// driver.Value可取值int64/float64/bool/[]byte/string/time.Time
func (r ReqUint) Value() (driver.Value, error) {
	return int64(r), nil
}

// gorm 检出 mysql 时调用
func (r *ReqUint) Scan(v interface{}) error {
	value, ok := v.(ReqUint)
	if ok {
		*r = value
		return nil
	}
	return fmt.Errorf("can not convert %v to ReqUint", v)
}

// 请求float64类型
type ReqFloat64 float64

func (r *ReqFloat64) UnmarshalJSON(data []byte) (err error) {
	str := strings.Trim(string(data), "\"")
	// ""空值不进行解析
	if utils.StrIsEmpty(str) {
		*r = ReqFloat64(0)
		return
	}
	*r = ReqFloat64(utils.Str2Float64(str))
	return
}

func (r ReqFloat64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%f", r)), nil
}

// gorm 写入 mysql 时调用
// driver.Value可取值int64/float64/bool/[]byte/string/time.Time
func (r ReqFloat64) Value() (driver.Value, error) {
	return float64(r), nil
}

// gorm 检出 mysql 时调用
func (r *ReqFloat64) Scan(v interface{}) error {
	value, ok := v.(ReqFloat64)
	if ok {
		*r = value
		return nil
	}
	return fmt.Errorf("can not convert %v to ReqFloat64", v)
}

// 参数绑定
func ShouldBind(c *gin.Context, req interface{}) {
	err := c.ShouldBind(req)
	if err != nil {
		global.Log.Error(c, "参数绑定失败: %v", err)
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
	}
}

// 参数校验
func Validate(c *gin.Context, req interface{}, trans map[string]string) {
	err := global.NewValidatorError(global.Validate.Struct(req), trans)
	if err != nil {
		global.Log.Error(c, "参数校验失败: %v", err)
		response.FailWithMsg(err)
	}
}
