package request

import (
	"go-shipment-api/pkg/utils"
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
func (s *UpdateIncrementalIdsRequestStruct) GetIncremental(oldList []uint) []uint {
	newList := make([]uint, 0)
	for _, oldItem := range oldList {
		// 已删除数据不加入列表
		if !utils.Contains(s.Delete, oldItem) {
			newList = append(newList, oldItem)
		}
	}
	// 将需要新增的数据合并
	return append(newList, s.Create...)
}
