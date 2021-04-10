package request

import (
	"gin-web/pkg/response"
)

// 获取数据字典列表结构体
type DictRequestStruct struct {
	Name              string   `json:"name" form:"name"`
	Desc              string   `json:"desc" form:"desc"`
	Status            *ReqUint `json:"status" form:"status"`
	response.PageInfo          // 分页参数
}

// 创建数据字典结构体
type CreateDictRequestStruct struct {
	Name   string   `json:"name" validate:"required"`
	Desc   string   `json:"desc" validate:"required"`
	Status *ReqUint `json:"status"`
	Remark string   `json:"remark"`
}

// 翻译需要校验的字段名称
func (s CreateDictRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "字典名称"
	m["Desc"] = "字典描述"
	return m
}

// 更新数据字典结构体
type UpdateDictRequestStruct struct {
	Name   *string  `json:"name"`
	Desc   *string  `json:"desc"`
	Status *ReqUint `json:"status"`
	Remark *string  `json:"remark"`
}

// 获取数据字典数据列表结构体
type DictDataRequestStruct struct {
	DictId            *ReqUint `json:"dictId" form:"dictId"`
	Key               string   `json:"key" form:"key"`
	Attr              string   `json:"attr" form:"attr"`
	Val               string   `json:"val" form:"val"`
	Status            *ReqUint `json:"status" form:"sort"`
	response.PageInfo          // 分页参数
}

// 创建数据字典数据结构体
type CreateDictDataRequestStruct struct {
	Key      string `json:"key" validate:"required"`
	Val      string `json:"val" validate:"required"`
	Attr     string `json:"attr"`
	Addition string `json:"addition"`
	Sort     *uint  `json:"sort"`
	Status   *uint  `json:"status"`
	Remark   string `json:"remark"`
	DictId   uint   `json:"dictId" validate:"required"`
}

// 翻译需要校验的字段名称
func (s CreateDictDataRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Key"] = "数据键"
	m["Val"] = "数据值"
	m["DictId"] = "字典编号"
	return m
}

// 更新数据字典数据结构体
type UpdateDictDataRequestStruct struct {
	Key      *string  `json:"key"`
	Val      *string  `json:"val"`
	Attr     *string  `json:"attr"`
	Addition *string  `json:"addition"`
	Sort     *ReqUint `json:"sort"`
	Status   *ReqUint `json:"status"`
	Remark   *string  `json:"remark"`
	DictId   *ReqUint `json:"dictId"`
}
