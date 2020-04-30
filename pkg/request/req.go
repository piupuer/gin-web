package request

import "go-shipment-api/pkg/utils"

// 适用于大多数场景的请求参数绑定
type Req struct {
	Ids string `json:"ids" form:"ids"` // 传多个id
}

// 获取
func (s *Req) GetUintIds() []uint {
	return utils.Str2UintArr(s.Ids)
}
