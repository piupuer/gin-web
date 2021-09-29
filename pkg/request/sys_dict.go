package request

import (
	"gin-web/pkg/response"
)

type DictReq struct {
	Name              string   `json:"name" form:"name"`
	Desc              string   `json:"desc" form:"desc"`
	Status            *ReqUint `json:"status" form:"status"`
	response.PageInfo          // 分页参数
}

type CreateDictReq struct {
	Name   string   `json:"name" validate:"required"`
	Desc   string   `json:"desc" validate:"required"`
	Status *ReqUint `json:"status"`
	Remark string   `json:"remark"`
}

func (s CreateDictReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "字典名称"
	m["Desc"] = "字典描述"
	return m
}

type UpdateDictReq struct {
	Name   *string  `json:"name"`
	Desc   *string  `json:"desc"`
	Status *ReqUint `json:"status"`
	Remark *string  `json:"remark"`
}

type DictDataReq struct {
	DictId            *ReqUint `json:"dictId" form:"dictId"`
	Key               string   `json:"key" form:"key"`
	Attr              string   `json:"attr" form:"attr"`
	Val               string   `json:"val" form:"val"`
	Status            *ReqUint `json:"status" form:"sort"`
	response.PageInfo          // 分页参数
}

type CreateDictDataReq struct {
	Key      string `json:"key" validate:"required"`
	Val      string `json:"val" validate:"required"`
	Attr     string `json:"attr"`
	Addition string `json:"addition"`
	Sort     *uint  `json:"sort"`
	Status   *uint  `json:"status"`
	Remark   string `json:"remark"`
	DictId   uint   `json:"dictId" validate:"required"`
}

func (s CreateDictDataReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Key"] = "数据键"
	m["Val"] = "数据值"
	m["DictId"] = "字典编号"
	return m
}

type UpdateDictDataReq struct {
	Key      *string  `json:"key"`
	Val      *string  `json:"val"`
	Attr     *string  `json:"attr"`
	Addition *string  `json:"addition"`
	Sort     *ReqUint `json:"sort"`
	Status   *ReqUint `json:"status"`
	Remark   *string  `json:"remark"`
	DictId   *ReqUint `json:"dictId"`
}
